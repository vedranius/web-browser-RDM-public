package main

import (
	"net/url"
	"strings"
)

// ─── STORAGES: WHAT EACH PROVIDER BILLS FOR THE REQUESTS WRM SENDS (v14.4.0) ──────────────
//
// WRM sends seven S3 operations: HeadBucket, ListObjectsV2, GetObject, HeadObject, PutObject,
// CopyObject and DeleteObject. For each provider this table says whether the provider bills
// them, in which class and for how much, quoted from the provider's own price page (the source
// URLs are shown with the table). Checked on storagePricingChecked; prices change, so the UI
// always links the source. Storage and egress are billed separately and are not counted here.

const storagePricingChecked = "2026-11-06"

type storagePriceClass struct {
	Name     string  `json:"name"`
	Billed   bool    `json:"billed"`
	Per      int     `json:"per,omitempty"`   // requests per price unit (1000 or 1000000)
	Price    float64 `json:"price,omitempty"` // in Currency, per Per requests
	FreeNote string  `json:"free_note,omitempty"`
}

type storageProvider struct {
	ID       string                       `json:"id"`
	Name     string                       `json:"name"`
	Currency string                       `json:"currency,omitempty"`
	Billed   string                       `json:"billed"` // no | yes | not_listed | self_hosted | unknown
	Classes  map[string]storagePriceClass `json:"classes,omitempty"`
	Ops      map[string]string            `json:"ops,omitempty"` // operation → class name
	Summary  map[string]string            `json:"summary"`       // en / hr
	Quote    string                       `json:"quote,omitempty"`
	Sources  []string                     `json:"sources,omitempty"`
	Checked  string                       `json:"checked"`
	Endpoint string                       `json:"endpoint_hint,omitempty"`
	Region   string                       `json:"region_hint,omitempty"`
	PathOK   bool                         `json:"path_style"`
}

var storageWRMOps = []string{"HeadBucket", "ListObjectsV2", "GetObject", "HeadObject", "PutObject", "CopyObject", "DeleteObject"}

var storageProviders = []storageProvider{
	{
		ID: "hetzner", Name: "Hetzner Object Storage", Billed: "no", PathOK: true,
		Endpoint: "https://fsn1.your-objectstorage.com", Region: "fsn1",
		Summary: map[string]string{
			"en": "Hetzner does not bill S3 requests: every operation WRM sends is free. Billed: the base price (EUR 6.49/month incl. 1 TB storage and 1 TB egress), more storage, more egress. Objects under 64 KB are billed as 64 KB.",
			"hr": "Hetzner ne naplaćuje S3 zahtjeve: sve operacije koje WRM šalje su besplatne. Naplaćuje se osnovna cijena (6,49 EUR/mj. uz 1 TB prostora i 1 TB odlaznog prometa), dodatni prostor i dodatni odlazni promet. Objekti manji od 64 KB naplaćuju se kao 64 KB.",
		},
		Quote:   `"The following types of traffic are free of charge: … S3 API calls such as PUT, GET, or DELETE requests." / "The following is free of charge: … S3 operations such as PUT, GET, or DELETE"`,
		Sources: []string{"https://www.hetzner.com/storage/object-storage/pricing/", "https://docs.hetzner.com/storage/object-storage/overview"},
	},
	{
		ID: "digitalocean", Name: "DigitalOcean Spaces", Billed: "not_listed", PathOK: true,
		Endpoint: "https://fra1.digitaloceanspaces.com", Region: "fra1",
		Summary: map[string]string{
			"en": "DigitalOcean's Spaces price list has no request fees: it bills the subscription (USD 5/month incl. 250 GiB storage and 1,024 GiB outbound transfer), more storage (USD 0.02/GiB) and more outbound transfer (USD 0.01/GiB). The documentation does not say in so many words that requests are free. New buckets are limited to 800 operations per second, and LIST may be throttled further under load.",
			"hr": "Cjenik DigitalOcean Spacesa nema naknadu po zahtjevu: naplaćuje pretplatu (5 USD/mj. uz 250 GiB prostora i 1.024 GiB odlaznog prometa), dodatni prostor (0,02 USD/GiB) i dodatni odlazni promet (0,01 USD/GiB). Dokumentacija ne kaže izričito da su zahtjevi besplatni. Novi bucketi imaju ograničenje od 800 operacija u sekundi, a LIST se pod opterećenjem može dodatno usporiti.",
		},
		Quote:   `"The base price of a Spaces Standard Storage subscription is $5.00 per month." … "All new Spaces buckets support a limit of 800 total operations per second. During periods of very high load, LIST requests may be rate-limited further."`,
		Sources: []string{"https://docs.digitalocean.com/products/spaces/details/pricing/", "https://docs.digitalocean.com/products/spaces/details/limits/"},
	},
	{
		ID: "aws", Name: "Amazon S3", Billed: "yes", Currency: "USD", PathOK: false,
		Endpoint: "https://s3.eu-central-1.amazonaws.com", Region: "eu-central-1",
		Classes: map[string]storagePriceClass{
			"Tier1": {Name: "PUT, COPY, POST, LIST", Billed: true, Per: 1000, Price: 0.005},
			"Tier2": {Name: "GET and all other requests", Billed: true, Per: 1000, Price: 0.0004},
			"Free":  {Name: "DELETE", Billed: false},
		},
		Ops: map[string]string{"HeadBucket": "Tier2", "ListObjectsV2": "Tier1", "GetObject": "Tier2", "HeadObject": "Tier2", "PutObject": "Tier1", "CopyObject": "Tier1", "DeleteObject": "Free"},
		Summary: map[string]string{
			"en": "Amazon S3 bills every request except DELETE. S3 Standard in US East (N. Virginia): PUT, COPY and LIST USD 0.005 per 1,000; GET and HEAD USD 0.0004 per 1,000. Other regions cost more. Opening a folder (LIST) is the most expensive call WRM makes. Egress is billed separately.",
			"hr": "Amazon S3 naplaćuje svaki zahtjev osim DELETE. S3 Standard u US East (N. Virginia): PUT, COPY i LIST 0,005 USD na 1.000; GET i HEAD 0,0004 USD na 1.000. Druge regije su skuplje. Otvaranje mape (LIST) je najskuplji poziv koji WRM radi. Odlazni promet naplaćuje se posebno.",
		},
		Quote:   `"You pay for requests made against your S3 buckets and objects." … "LIST requests for any storage class are charged at the same rate as S3 Standard PUT, COPY, and POST requests. DELETE and CANCEL requests are free."`,
		Sources: []string{"https://aws.amazon.com/s3/pricing/", "https://docs.aws.amazon.com/AmazonS3/latest/userguide/aws-usage-report-understand.html"},
	},
	{
		ID: "backblaze", Name: "Backblaze B2", Billed: "no", PathOK: true,
		Endpoint: "https://s3.eu-central-003.backblazeb2.com", Region: "eu-central-003",
		Classes: map[string]storagePriceClass{
			"A": {Name: "Class A", Billed: false}, "B": {Name: "Class B", Billed: false}, "C": {Name: "Class C", Billed: false},
		},
		Ops: map[string]string{"HeadBucket": "C", "ListObjectsV2": "C", "GetObject": "B", "HeadObject": "B", "PutObject": "A", "CopyObject": "C", "DeleteObject": "A"},
		Summary: map[string]string{
			"en": "Backblaze B2 lists all three classes of S3 calls WRM uses (A: put, delete; B: get, head object; C: list, head bucket, copy) as free. Billed: storage after the first 10 GB (USD 0.005/GB/month) and egress beyond 3x the average stored data (USD 0.01/GB).",
			"hr": "Backblaze B2 navodi sve tri klase S3 poziva koje WRM koristi kao besplatne (A: put, delete; B: get, head object; C: list, head bucket, copy). Naplaćuje se prostor nakon prvih 10 GB (0,005 USD/GB/mj.) i odlazni promet iznad trostrukog prosjeka pohranjenih podataka (0,01 USD/GB).",
		},
		Quote:   `Class A "Cost: Free", Class B "Cost: Free", Class C "Cost: Free"`,
		Sources: []string{"https://www.backblaze.com/cloud-storage/transaction-pricing"},
	},
	{
		ID: "r2", Name: "Cloudflare R2", Billed: "yes", Currency: "USD", PathOK: true,
		Endpoint: "https://<account-id>.r2.cloudflarestorage.com", Region: "auto",
		Classes: map[string]storagePriceClass{
			"A":    {Name: "Class A", Billed: true, Per: 1000000, Price: 4.50, FreeNote: "1 million / month free"},
			"B":    {Name: "Class B", Billed: true, Per: 1000000, Price: 0.36, FreeNote: "10 million / month free"},
			"Free": {Name: "Free", Billed: false},
		},
		Ops: map[string]string{"HeadBucket": "B", "ListObjectsV2": "A", "GetObject": "B", "HeadObject": "B", "PutObject": "A", "CopyObject": "A", "DeleteObject": "Free"},
		Summary: map[string]string{
			"en": "Cloudflare R2 bills requests in two classes, with a monthly free allowance. Class A (list, put, copy): USD 4.50 per million, first 1 million a month free. Class B (get, head): USD 0.36 per million, first 10 million free. DELETE is free, and so is egress.",
			"hr": "Cloudflare R2 naplaćuje zahtjeve u dvije klase, uz besplatnu mjesečnu kvotu. Klasa A (list, put, copy): 4,50 USD na milijun, prvih milijun mjesečno besplatno. Klasa B (get, head): 0,36 USD na milijun, prvih 10 milijuna besplatno. DELETE je besplatan, kao i odlazni promet.",
		},
		Quote:   `"Class A Operations include … ListObjects, PutObject, CopyObject …" "Class B Operations include HeadBucket, HeadObject, GetObject …" "Free operations include DeleteObject …"`,
		Sources: []string{"https://developers.cloudflare.com/r2/pricing/"},
	},
	{
		ID: "minio", Name: "MinIO (self-hosted)", Billed: "self_hosted", PathOK: true,
		Endpoint: "https://minio.example.com:9000", Region: "us-east-1",
		Summary: map[string]string{
			"en": "Your own server: requests cost nothing, only your own hardware and traffic.",
			"hr": "Vaš vlastiti server: zahtjevi ne koštaju ništa, samo vaš hardver i promet.",
		},
	},
	{
		ID: "other", Name: "Other S3-compatible", Billed: "unknown", PathOK: true,
		Summary: map[string]string{
			"en": "WRM does not know this provider's prices. Check whether it bills requests; the activity panel counts every request WRM sends.",
			"hr": "WRM ne zna cijene ovog providera. Provjerite naplaćuje li zahtjeve; panel aktivnosti broji svaki zahtjev koji WRM pošalje.",
		},
	},
}

func storageProviderByID(id string) (storageProvider, bool) {
	for _, p := range storageProviders {
		if p.ID == id {
			p.Checked = storagePricingChecked
			return p, true
		}
	}
	return storageProvider{}, false
}

// guessStorageProvider recognises the common endpoints, for storages saved before v14.4.0.
func guessStorageProvider(endpoint string) string {
	u, err := url.Parse(strings.TrimSpace(endpoint))
	host := strings.ToLower(endpoint)
	if err == nil && u.Host != "" {
		host = strings.ToLower(u.Hostname())
	}
	switch {
	case strings.HasSuffix(host, "your-objectstorage.com"):
		return "hetzner"
	case strings.HasSuffix(host, "digitaloceanspaces.com"):
		return "digitalocean"
	case strings.HasSuffix(host, "amazonaws.com"):
		return "aws"
	case strings.HasSuffix(host, "backblazeb2.com"):
		return "backblaze"
	case strings.HasSuffix(host, "r2.cloudflarestorage.com"):
		return "r2"
	}
	return "other"
}

func storageProviderFor(s Storage) storageProvider {
	id := s.Provider
	if id == "" {
		id = guessStorageProvider(s.Endpoint)
	}
	if p, ok := storageProviderByID(id); ok {
		return p
	}
	p, _ := storageProviderByID("other")
	return p
}

// storageCostEstimate prices the counted requests with the provider's table: per operation the
// class, whether it is billed, and the cost of the counted requests before any free allowance.
func storageCostEstimate(p storageProvider, counts map[string]storageOpCount) map[string]interface{} {
	rows := []map[string]interface{}{}
	total, today := 0.0, 0.0
	for _, op := range storageWRMOps {
		c := counts[op]
		row := map[string]interface{}{"op": op, "total": c.Total, "today": c.Today}
		switch p.Billed {
		case "no", "self_hosted":
			row["billed"] = false
		case "not_listed":
			row["billed"] = false
			row["note"] = "not_listed"
		case "unknown":
			row["billed"] = nil
		default:
			cls, ok := p.Classes[p.Ops[op]]
			row["class"] = p.Ops[op]
			row["billed"] = ok && cls.Billed
			if ok && cls.Billed && cls.Per > 0 {
				unit := cls.Price / float64(cls.Per)
				row["price_per_request"] = unit
				row["cost_total"] = unit * float64(c.Total)
				row["cost_today"] = unit * float64(c.Today)
				total += unit * float64(c.Total)
				today += unit * float64(c.Today)
			}
		}
		rows = append(rows, row)
	}
	return map[string]interface{}{"rows": rows, "currency": p.Currency, "cost_total": total, "cost_today": today,
		"note": "Before any free allowance; storage and egress are billed separately."}
}
