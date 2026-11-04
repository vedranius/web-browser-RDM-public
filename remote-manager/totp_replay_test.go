package main

import (
	"fmt"
	"testing"
	"time"
)

func secondFactorReplayUser(t *testing.T) (*userRow, []string, []byte) {
	t.Helper()
	secret := newTOTPSecret()
	raw, err := decodeTOTPSecret(secret)
	if err != nil {
		t.Fatal(err)
	}
	codes, recovery := newRecoveryCodes()
	res, err := db.Exec(`INSERT INTO users (username, password_hash, created_at, totp_enabled, totp_secret, recovery_codes) VALUES (?,?,?,1,?,?)`,
		fmt.Sprintf("factor-replay-%d", time.Now().UnixNano()), "unused", nowRFC(), encryptValue(secret), recovery)
	if err != nil {
		t.Fatal(err)
	}
	id, _ := res.LastInsertId()
	t.Cleanup(func() { db.Exec(`DELETE FROM users WHERE id=?`, id) })
	u, err := loadUser("id=?", id)
	if err != nil {
		t.Fatal(err)
	}
	return u, codes, raw
}

func consumeFactorConcurrently(t *testing.T, u *userRow, codes []string) int {
	t.Helper()
	type result struct {
		method string
		err    error
	}
	start := make(chan struct{})
	results := make(chan result, len(codes))
	for _, code := range codes {
		go func(code string) {
			<-start
			method, _, err := consumeUserSecondFactor(u, code)
			results <- result{method, err}
		}(code)
	}
	close(start)
	successes := 0
	for range codes {
		r := <-results
		if r.err != nil {
			t.Errorf("consume second factor: %v", r.err)
		}
		if r.method != "" {
			successes++
		}
	}
	return successes
}

func TestSecondFactorTOTPConsumedAtomically(t *testing.T) {
	u, _, secret := secondFactorReplayUser(t)
	code := totpCodeAt(secret, time.Now().Unix()/totpPeriod)
	codes := make([]string, 16)
	for i := range codes {
		codes[i] = code
	}
	if n := consumeFactorConcurrently(t, u, codes); n != 1 {
		t.Fatalf("concurrent requests accepted the same TOTP %d times, want 1", n)
	}
}

func TestSecondFactorRecoveryConsumedAtomically(t *testing.T) {
	u, recovery, _ := secondFactorReplayUser(t)
	codes := make([]string, 16)
	for i := range codes {
		codes[i] = recovery[0]
	}
	if n := consumeFactorConcurrently(t, u, codes); n != 1 {
		t.Fatalf("concurrent requests accepted the same recovery code %d times, want 1", n)
	}
	var remaining string
	if err := db.QueryRow(`SELECT recovery_codes FROM users WHERE id=?`, u.ID).Scan(&remaining); err != nil {
		t.Fatal(err)
	}
	if n := recoveryCodesLeft(remaining); n != 9 {
		t.Fatalf("remaining recovery codes = %d, want 9", n)
	}
}

func TestSecondFactorConcurrentRecoveryCodesDoNotRestoreConsumedCodes(t *testing.T) {
	u, codes, _ := secondFactorReplayUser(t)
	if n := consumeFactorConcurrently(t, u, codes); n != len(codes) {
		t.Fatalf("concurrent distinct recovery codes accepted = %d, want %d", n, len(codes))
	}
	var remaining string
	if err := db.QueryRow(`SELECT recovery_codes FROM users WHERE id=?`, u.ID).Scan(&remaining); err != nil {
		t.Fatal(err)
	}
	if n := recoveryCodesLeft(remaining); n != 0 {
		t.Fatalf("consumed recovery codes were restored: %d codes remain", n)
	}
}

func TestSecondFactorCannotBeReplayedBetweenLoginAndStepUp(t *testing.T) {
	recordLoginResult("127.0.0.1", true)
	t.Cleanup(func() { recordLoginResult("127.0.0.1", true) })
	srv := newTestServer(t)
	name := fmt.Sprintf("factor-api-%d", time.Now().UnixNano())
	client := newTestUser(t, srv, name, false)
	secret := newTOTPSecret()
	raw, err := decodeTOTPSecret(secret)
	if err != nil {
		t.Fatal(err)
	}
	codes, recovery := newRecoveryCodes()
	if _, err := db.Exec(`UPDATE users SET totp_enabled=1, totp_secret=?, totp_last_step=0, recovery_codes=? WHERE id=?`, encryptValue(secret), recovery, client.userID); err != nil {
		t.Fatal(err)
	}

	code := totpCodeAt(raw, time.Now().Unix()/totpPeriod)
	if status := client.jsonDo("POST", "/api/auth/mfa", map[string]string{"mfa_token": newMFAToken(client.userID), "code": code}, nil); status != 200 {
		t.Fatalf("TOTP sign-in: %d", status)
	}
	password := "pw-" + name + "-123456"
	if status := client.jsonDo("POST", "/api/auth/step-up", map[string]string{"password": password, "code": code}, nil); status != 403 {
		t.Fatalf("step-up reused the sign-in TOTP: %d", status)
	}
	if status := client.jsonDo("POST", "/api/auth/step-up", map[string]string{"password": password, "code": codes[0]}, nil); status != 200 {
		t.Fatalf("recovery-code step-up: %d", status)
	}
	if status := client.jsonDo("POST", "/api/auth/mfa", map[string]string{"mfa_token": newMFAToken(client.userID), "code": codes[0]}, nil); status != 401 {
		t.Fatalf("sign-in reused the step-up recovery code: %d", status)
	}
}
