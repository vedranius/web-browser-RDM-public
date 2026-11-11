package main

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"strings"
)

// ─── WEBAUTHN (passkeys), v14.13.0 ───────────────────
//
// A small, dependency-free verifier for WebAuthn Level 2 registrations and assertions:
// enough for passkeys on Windows Hello, Touch ID / Face ID, Android and iOS (also from a phone
// through the browser's QR + Bluetooth "hybrid" flow) and FIDO2 security keys.
//
// Attestation statements are not checked (WRM asks for attestation "none"): a passkey proves
// possession of its private key, which is what sign-in needs, not the authenticator's make.
// Supported public keys: ES256 (P-256), EdDSA (Ed25519) and RS256.

const (
	coseAlgES256 = -7
	coseAlgEdDSA = -8
	coseAlgRS256 = -257

	waFlagUP = 0x01 // user present
	waFlagUV = 0x04 // user verified (biometric or device PIN)
	waFlagBE = 0x08 // backup eligible (synced passkey)
	waFlagBS = 0x10 // backed up
	waFlagAT = 0x40 // attested credential data present
	waFlagED = 0x80 // extension data present
)

var waB64 = base64.RawURLEncoding

// b64urlDecode accepts base64url with or without padding.
func b64urlDecode(s string) ([]byte, error) { return b64url(s) }

// ─── minimal CBOR (RFC 8949) ─────────────────────────
//
// Only what authenticators send: unsigned/negative integers, byte and text strings, arrays,
// maps, simple values and floats (skipped). Indefinite lengths and tags are refused, sizes
// and nesting are bounded.

const cborMaxDepth = 16

type cborReader struct {
	b   []byte
	pos int
}

var errCBOR = errors.New("malformed CBOR")

func (c *cborReader) head() (major byte, arg uint64, err error) {
	if c.pos >= len(c.b) {
		return 0, 0, errCBOR
	}
	ib := c.b[c.pos]
	c.pos++
	major, info := ib>>5, ib&0x1f
	switch {
	case info < 24:
		return major, uint64(info), nil
	case info == 24:
		if c.pos+1 > len(c.b) {
			return 0, 0, errCBOR
		}
		arg = uint64(c.b[c.pos])
		c.pos++
	case info == 25:
		if c.pos+2 > len(c.b) {
			return 0, 0, errCBOR
		}
		arg = uint64(binary.BigEndian.Uint16(c.b[c.pos:]))
		c.pos += 2
	case info == 26:
		if c.pos+4 > len(c.b) {
			return 0, 0, errCBOR
		}
		arg = uint64(binary.BigEndian.Uint32(c.b[c.pos:]))
		c.pos += 4
	case info == 27:
		if c.pos+8 > len(c.b) {
			return 0, 0, errCBOR
		}
		arg = binary.BigEndian.Uint64(c.b[c.pos:])
		c.pos += 8
	default:
		return 0, 0, errCBOR // indefinite lengths and reserved values
	}
	return major, arg, nil
}

// value decodes one item: int64, []byte, string, []interface{}, map[interface{}]interface{},
// bool or nil.
func (c *cborReader) value(depth int) (interface{}, error) {
	if depth > cborMaxDepth {
		return nil, errCBOR
	}
	major, arg, err := c.head()
	if err != nil {
		return nil, err
	}
	switch major {
	case 0:
		if arg > 1<<62 {
			return nil, errCBOR
		}
		return int64(arg), nil
	case 1:
		if arg > 1<<62 {
			return nil, errCBOR
		}
		return -1 - int64(arg), nil
	case 2, 3:
		if arg > uint64(len(c.b)-c.pos) {
			return nil, errCBOR
		}
		v := c.b[c.pos : c.pos+int(arg)]
		c.pos += int(arg)
		if major == 3 {
			return string(v), nil
		}
		return append([]byte(nil), v...), nil
	case 4:
		if arg > uint64(len(c.b)-c.pos) { // every item takes at least one byte
			return nil, errCBOR
		}
		out := make([]interface{}, 0, int(arg))
		for i := uint64(0); i < arg; i++ {
			v, err := c.value(depth + 1)
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		}
		return out, nil
	case 5:
		if arg > uint64(len(c.b)-c.pos)/2 {
			return nil, errCBOR
		}
		out := make(map[interface{}]interface{}, int(arg))
		for i := uint64(0); i < arg; i++ {
			k, err := c.value(depth + 1)
			if err != nil {
				return nil, err
			}
			switch k.(type) {
			case int64, string:
			default:
				return nil, errCBOR
			}
			v, err := c.value(depth + 1)
			if err != nil {
				return nil, err
			}
			out[k] = v
		}
		return out, nil
	case 7:
		switch arg {
		case 20:
			return false, nil
		case 21:
			return true, nil
		case 22, 23:
			return nil, nil
		}
		return nil, nil // floats (their bytes were consumed by head) and other simple values
	}
	return nil, errCBOR // tags (6) are not used by WebAuthn
}

// cborDecode decodes the first item of b and returns how many bytes it used.
func cborDecode(b []byte) (interface{}, int, error) {
	c := &cborReader{b: b}
	v, err := c.value(0)
	if err != nil {
		return nil, 0, err
	}
	return v, c.pos, nil
}

// ─── public keys (COSE, RFC 9053) ────────────────────

type waPublicKey struct {
	Alg int64
	ec  *ecdsa.PublicKey
	ed  ed25519.PublicKey
	rsa *rsa.PublicKey
}

func coseInt(m map[interface{}]interface{}, k int64) (int64, bool) {
	v, ok := m[k].(int64)
	return v, ok
}

func coseBytes(m map[interface{}]interface{}, k int64) []byte {
	v, _ := m[k].([]byte)
	return v
}

// parseCOSEKey reads a COSE_Key map (as stored for a credential).
func parseCOSEKey(raw []byte) (*waPublicKey, error) {
	v, n, err := cborDecode(raw)
	if err != nil {
		return nil, err
	}
	if n != len(raw) {
		return nil, errors.New("trailing bytes after the public key")
	}
	m, ok := v.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("the public key is not a COSE map")
	}
	kty, _ := coseInt(m, 1)
	alg, ok := coseInt(m, 3)
	if !ok {
		return nil, errors.New("the public key names no algorithm")
	}
	k := &waPublicKey{Alg: alg}
	switch alg {
	case coseAlgES256:
		crv, _ := coseInt(m, -1)
		x, y := coseBytes(m, -2), coseBytes(m, -3)
		if kty != 2 || crv != 1 || len(x) != 32 || len(y) != 32 {
			return nil, errors.New("unsupported EC2 key")
		}
		// ParseUncompressedPublicKey rejects points that are not on the curve.
		pub, err := ecdsa.ParseUncompressedPublicKey(elliptic.P256(), append(append([]byte{4}, x...), y...))
		if err != nil {
			return nil, errors.New("the EC public key is not on the curve")
		}
		k.ec = pub
	case coseAlgEdDSA:
		crv, _ := coseInt(m, -1)
		x := coseBytes(m, -2)
		if kty != 1 || crv != 6 || len(x) != ed25519.PublicKeySize {
			return nil, errors.New("unsupported OKP key")
		}
		k.ed = ed25519.PublicKey(append([]byte(nil), x...))
	case coseAlgRS256:
		nb, eb := coseBytes(m, -1), coseBytes(m, -2)
		if kty != 3 || len(nb) < 256 || len(nb) > 512 || len(eb) == 0 || len(eb) > 4 {
			return nil, errors.New("unsupported RSA key (2048 to 4096 bits)")
		}
		e := 0
		for _, b := range eb {
			e = e<<8 | int(b)
		}
		if e < 3 || e%2 == 0 {
			return nil, errors.New("bad RSA exponent")
		}
		k.rsa = &rsa.PublicKey{N: new(big.Int).SetBytes(nb), E: e}
	default:
		return nil, fmt.Errorf("unsupported key algorithm %d", alg)
	}
	return k, nil
}

func (k *waPublicKey) verify(data, sig []byte) bool {
	switch {
	case k.ec != nil:
		h := sha256.Sum256(data)
		return ecdsa.VerifyASN1(k.ec, h[:], sig)
	case k.ed != nil:
		return ed25519.Verify(k.ed, data, sig)
	case k.rsa != nil:
		h := sha256.Sum256(data)
		return rsa.VerifyPKCS1v15(k.rsa, crypto.SHA256, h[:], sig) == nil
	}
	return false
}

// ─── authenticator data ──────────────────────────────

type waAuthData struct {
	RPIDHash  []byte
	Flags     byte
	SignCount uint32
	AAGUID    []byte
	CredID    []byte
	PublicKey []byte // COSE_Key bytes
}

func parseAuthData(b []byte, wantCredential bool) (*waAuthData, error) {
	if len(b) < 37 {
		return nil, errors.New("authenticator data too short")
	}
	a := &waAuthData{RPIDHash: b[:32], Flags: b[32], SignCount: binary.BigEndian.Uint32(b[33:37])}
	rest := b[37:]
	if a.Flags&waFlagAT != 0 {
		if len(rest) < 18 {
			return nil, errors.New("attested credential data too short")
		}
		a.AAGUID = rest[:16]
		n := int(binary.BigEndian.Uint16(rest[16:18]))
		rest = rest[18:]
		if n == 0 || n > 1023 || len(rest) < n {
			return nil, errors.New("bad credential ID length")
		}
		a.CredID = rest[:n]
		rest = rest[n:]
		_, used, err := cborDecode(rest)
		if err != nil {
			return nil, errors.New("bad credential public key")
		}
		a.PublicKey = rest[:used]
		rest = rest[used:]
	} else if wantCredential {
		return nil, errors.New("the authenticator returned no credential")
	}
	if a.Flags&waFlagED != 0 {
		_, used, err := cborDecode(rest)
		if err != nil {
			return nil, errors.New("bad extension data")
		}
		rest = rest[used:]
	}
	if len(rest) != 0 {
		return nil, errors.New("trailing bytes in authenticator data")
	}
	return a, nil
}

// ─── client data and origin ──────────────────────────

type waClientData struct {
	Type        string `json:"type"`
	Challenge   string `json:"challenge"`
	Origin      string `json:"origin"`
	CrossOrigin bool   `json:"crossOrigin"`
}

// waRelyingParty is who the browser must have talked to: the RP ID the credential is scoped to
// and the origins a ceremony may come from.
type waRelyingParty struct {
	ID      string
	Origins []string // exact origins, e.g. https://wrm.example.com
}

// originAllowed: the origin is one of rp.Origins, or an https origin on rp.ID or a subdomain
// of it, or http://localhost (the only plain-HTTP origin browsers allow WebAuthn on).
func (rp waRelyingParty) originAllowed(origin string) bool {
	for _, o := range rp.Origins {
		if o != "" && strings.EqualFold(o, origin) {
			return true
		}
	}
	u, err := url.Parse(origin)
	if err != nil || u.Host == "" || u.Path != "" || u.User != nil {
		return false
	}
	host := strings.ToLower(u.Hostname())
	if !(host == rp.ID || strings.HasSuffix(host, "."+rp.ID)) {
		return false
	}
	switch u.Scheme {
	case "https":
		return true
	case "http":
		return host == "localhost" || strings.HasSuffix(host, ".localhost")
	}
	return false
}

func checkClientData(raw []byte, wantType string, challenge []byte, rp waRelyingParty) error {
	var cd waClientData
	if err := json.Unmarshal(raw, &cd); err != nil {
		return errors.New("bad client data")
	}
	if cd.Type != wantType {
		return errors.New("wrong ceremony type")
	}
	got, err := b64urlDecode(cd.Challenge)
	if err != nil || subtle.ConstantTimeCompare(got, challenge) != 1 {
		return errors.New("the challenge does not match")
	}
	if cd.CrossOrigin {
		return errors.New("cross-origin ceremonies are not accepted")
	}
	if !rp.originAllowed(cd.Origin) {
		return fmt.Errorf("the browser talked to %q, which is not this WRM", truncateStr(cd.Origin, 120))
	}
	return nil
}

func rpIDHashOK(a *waAuthData, rpID string) bool {
	h := sha256.Sum256([]byte(rpID))
	return subtle.ConstantTimeCompare(a.RPIDHash, h[:]) == 1
}

// validRPID: a registrable host name (not an IP address: browsers refuse those).
func validRPID(s string) bool {
	if s == "" || len(s) > 253 || net.ParseIP(s) != nil || strings.HasPrefix(s, ".") || strings.HasSuffix(s, ".") {
		return false
	}
	for _, label := range strings.Split(s, ".") {
		if label == "" || len(label) > 63 {
			return false
		}
		for _, c := range label {
			if !(c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-') {
				return false
			}
		}
	}
	return true
}

// ─── ceremonies ──────────────────────────────────────

type waAttestationResponse struct {
	ClientDataJSON    string   `json:"clientDataJSON"`
	AttestationObject string   `json:"attestationObject"`
	Transports        []string `json:"transports"`
}

type waNewCredential struct {
	ID         []byte
	PublicKey  []byte
	Alg        int64
	SignCount  uint32
	AAGUID     []byte
	UV         bool
	Backup     bool // backup eligible: a synced passkey
	Transports []string
}

// verifyRegistration checks a navigator.credentials.create() answer.
func verifyRegistration(resp waAttestationResponse, challenge []byte, rp waRelyingParty, requireUV bool) (*waNewCredential, error) {
	cdj, err := b64urlDecode(resp.ClientDataJSON)
	if err != nil || len(cdj) > 4096 {
		return nil, errors.New("bad client data")
	}
	if err := checkClientData(cdj, "webauthn.create", challenge, rp); err != nil {
		return nil, err
	}
	ao, err := b64urlDecode(resp.AttestationObject)
	if err != nil || len(ao) > 16384 {
		return nil, errors.New("bad attestation object")
	}
	v, _, err := cborDecode(ao)
	if err != nil {
		return nil, errors.New("bad attestation object")
	}
	m, ok := v.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("bad attestation object")
	}
	ad, _ := m["authData"].([]byte)
	a, err := parseAuthData(ad, true)
	if err != nil {
		return nil, err
	}
	if !rpIDHashOK(a, rp.ID) {
		return nil, errors.New("the passkey was made for another site")
	}
	if a.Flags&waFlagUP == 0 {
		return nil, errors.New("the authenticator did not confirm user presence")
	}
	if requireUV && a.Flags&waFlagUV == 0 {
		return nil, errors.New("the authenticator did not verify the user (fingerprint, face or device PIN)")
	}
	key, err := parseCOSEKey(a.PublicKey)
	if err != nil {
		return nil, err
	}
	tr := []string{}
	for _, t := range resp.Transports {
		if len(tr) < 8 && len(t) <= 20 {
			tr = append(tr, t)
		}
	}
	return &waNewCredential{
		ID: append([]byte(nil), a.CredID...), PublicKey: append([]byte(nil), a.PublicKey...), Alg: key.Alg,
		SignCount: a.SignCount, AAGUID: append([]byte(nil), a.AAGUID...), UV: a.Flags&waFlagUV != 0,
		Backup: a.Flags&waFlagBE != 0, Transports: tr,
	}, nil
}

type waAssertionResponse struct {
	ClientDataJSON    string `json:"clientDataJSON"`
	AuthenticatorData string `json:"authenticatorData"`
	Signature         string `json:"signature"`
	UserHandle        string `json:"userHandle"`
}

type waAssertion struct {
	SignCount uint32
	UV        bool
	Backup    bool
}

// verifyAssertion checks a navigator.credentials.get() answer against a stored public key.
func verifyAssertion(resp waAssertionResponse, challenge []byte, rp waRelyingParty, publicKey []byte, requireUV bool) (*waAssertion, error) {
	cdj, err := b64urlDecode(resp.ClientDataJSON)
	if err != nil || len(cdj) > 4096 {
		return nil, errors.New("bad client data")
	}
	if err := checkClientData(cdj, "webauthn.get", challenge, rp); err != nil {
		return nil, err
	}
	ad, err := b64urlDecode(resp.AuthenticatorData)
	if err != nil || len(ad) > 4096 {
		return nil, errors.New("bad authenticator data")
	}
	a, err := parseAuthData(ad, false)
	if err != nil {
		return nil, err
	}
	if !rpIDHashOK(a, rp.ID) {
		return nil, errors.New("the passkey belongs to another site")
	}
	if a.Flags&waFlagUP == 0 {
		return nil, errors.New("the authenticator did not confirm user presence")
	}
	if requireUV && a.Flags&waFlagUV == 0 {
		return nil, errors.New("the authenticator did not verify the user (fingerprint, face or device PIN)")
	}
	sig, err := b64urlDecode(resp.Signature)
	if err != nil || len(sig) == 0 || len(sig) > 1024 {
		return nil, errors.New("bad signature")
	}
	key, err := parseCOSEKey(publicKey)
	if err != nil {
		return nil, err
	}
	h := sha256.Sum256(cdj)
	if !key.verify(bytes.Join([][]byte{ad, h[:]}, nil), sig) {
		return nil, errors.New("the signature is not valid")
	}
	return &waAssertion{SignCount: a.SignCount, UV: a.Flags&waFlagUV != 0, Backup: a.Flags&waFlagBE != 0}, nil
}
