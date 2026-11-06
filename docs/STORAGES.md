# Storages (S3) - operator guide

WRM can browse and edit S3-compatible object storage next to your servers. It is a section of its
own: a storage is not a Connection and never reaches the terminal, tunnels, jump hosts or the
deploy code.

Released in v14.0.0; the changes are in [CHANGELOG.md](../CHANGELOG.md).

## Add a storage

**Storages → Add storage.** What to enter depends on the provider:

| Provider | Endpoint | Region | Path-style |
|---|---|---|---|
| Hetzner Object Storage | `https://<location>.your-objectstorage.com` (e.g. `fsn1`) | `fsn1`, `nbg1`, `hel1` | on |
| DigitalOcean Spaces | `https://<region>.digitaloceanspaces.com` | `nyc3`, `ams3`, `fra1`, `sgp1` | on |
| Amazon S3 | `https://s3.<region>.amazonaws.com` | e.g. `eu-central-1` | off |
| MinIO | `https://minio.example.lan:9000` | what the server uses (`us-east-1` by default) | on |
| Backblaze B2 (S3 API) | `https://s3.<region>.backblazeb2.com` | e.g. `eu-central-003` | off |

- **Prefix** is mandatory in practice: set it to the folder the storage may touch (e.g.
  `team-a/backups`). Every key is resolved inside it, so the rest of the bucket is unreachable,
  even with the API.
- **Access key / secret key:** use a key scoped to that bucket and prefix. WRM encrypts the
  secret with the installation key; it is never returned to the browser, written to the audit log
  or included in a configuration export (unless you explicitly export secrets).
- **Test connection** checks the endpoint, the bucket and the key permissions before saving.
- A self-signed endpoint (MinIO) needs its **SHA-256 certificate fingerprint** in the dialog; there
  is no global "ignore certificate" switch.
- A LAN endpoint must pass the egress policy: it is allowed by default (`egress_policy=allow_private`);
  with `public_only`, add the host to `egress_allow_hosts`.

## Edit or remove a storage

A storage in the left list has a pencil and a trash button; the selected storage also shows
*Edit* in the toolbar. The dialog changes the endpoint, prefix, key or permissions (a changed
endpoint needs the secret again). **Remove storage** deletes only the WRM entry; the bucket,
its objects and any share of the provider stay where they are.

## Permissions inside a storage

| Switch | Default | What it allows |
|---|---|---|
| (always) | on | list, preview, download, play video |
| Allow editing files | on | save an existing object from the editor |
| Allow adding new files | **off** | upload a new small file |
| Allow deleting files | **off** | delete an object (also asks for a confirmation, and for the step-up password/2FA when the installation enforces it) |

Actions a storage does not allow are hidden in the UI and refused by the API. The switches are
per storage; an administrator can turn the whole section off with `storages_enabled=0`.

## What the browser does

- Lists prefixes like folders. A folder is loaded **completely**: WRM follows every S3
  continuation token, shows the running count while it loads, and then shows the total object
  and folder count and size with a clear "all objects in this folder are shown" line. Above a
  100-page safety limit it offers *Load more* instead of hiding the rest.
- **Search** filters the complete folder, **sort** orders by name, size or date, **select all**
  and **bulk delete** work on the visible selection, and *Copy path* copies an object key.
  Objects: preview (images), play (video through Range requests), download, edit, rename,
  delete.
- **Editor:** JSON (format, validate), YAML/XML/INI, scripts, and subtitles (SRT, WebVTT, ASS/SSA)
  with a timing check, a shift and a renumber. Saves keep the BOM and the line endings and use
  `If-Match`, so a file changed in the meantime is a conflict, not a silent overwrite.
- **Video:** the browser plays what the bucket stores. H.264/AAC works; HEVC usually does not, and
  WRM says so instead of transcoding. A sibling `x.srt`, `x.vtt` or `x.ass` is offered as a
  subtitle track (SRT/ASS converted to VTT in the browser; the stored file is untouched).
- WRM is not a transfer tool here: uploads are for small files (16 MB per object by default,
  2 MB in the editor). Large media is streamed, never buffered.

## AI / MCP

An AI connection reaches a storage only through its own scopes:

| Scope | Tools |
|---|---|
| `storage_read` | `storage_list`, `storage_read` |
| `storage_write` | `storage_write` (with `allow_write` and an optional etag) |
| `storage_delete` | `storage_delete` (with `allow_delete` and `confirm="delete"`) |

The scopes are opt-in: an administrator adds them to `ai_mcp_scopes`, then the user grants them to
an AI connection. The AI read-deny list (`.env`, keys, certificates) applies to storage reads as
well.

## Status and notifications

WRM checks each storage's bucket in the background (one `HeadBucket` request per storage every
`storage_check_minutes`, default 15 minutes; `0` turns the background check off) and shows a
green or red dot in the list. A storage that stops answering raises a `storage.down` notification
for its owner, and `storage.up` when it comes back (Notifications → events).

## Which S3 requests WRM sends

WRM never polls a folder or an object. Apart from the two background checks below, requests go out only when someone acts:

| Action | Requests |
|---|---|
| Background status check | `HEAD` bucket, once per storage every `storage_check_minutes` |
| A chain's `storage` check (Chains) | `HEAD` bucket, or one `LIST` for a `prefix` check, at most every `chain_storage_check_s` (default 300 s), whatever the chain interval |
| Add a storage, *Test connection*, edit the endpoint/region/bucket/prefix/keys | `HEAD` bucket + one `LIST` with `max-keys=1` |
| Edit only the name or the permissions, remove a storage | none |
| Open a folder | one `LIST` (`ListObjectsV2`, 1000 keys) per page, at most 100 pages before *Load more* |
| Download, preview, play, open in the editor | `GET` (video: `GET` with `Range`) |
| Upload or save | `HEAD` (does the object exist?) + `PUT` |
| Rename | `HEAD` ×2 + `PUT` copy + `DELETE` |
| Delete | `HEAD` (size for the audit log) + `DELETE` |

A retried request (WRM retries a 5xx answer up to twice) is a request of its own and is counted
as such in the activity panel.

## What each provider bills

Checked on 2026-11-06 against the providers' own pages; prices change, so the storage form and
the activity panel always link the source. Storage and egress are billed separately by every
provider and are not part of this table.

| Operation WRM sends | Hetzner | DigitalOcean Spaces | Amazon S3 (Standard, us-east-1) | Backblaze B2 | Cloudflare R2 (Standard) |
|---|---|---|---|---|---|
| HeadBucket | free | not in the price list | GET/other: USD 0.0004 / 1,000 | Class C: free | Class B: USD 0.36 / million |
| ListObjectsV2 | free | not in the price list | LIST: USD 0.005 / 1,000 | Class C: free | Class A: USD 4.50 / million |
| GetObject | free | not in the price list | GET: USD 0.0004 / 1,000 | Class B: free | Class B: USD 0.36 / million |
| HeadObject | free | not in the price list | GET/other: USD 0.0004 / 1,000 | Class B: free | Class B: USD 0.36 / million |
| PutObject | free | not in the price list | PUT: USD 0.005 / 1,000 | Class A: free | Class A: USD 4.50 / million |
| CopyObject | free | not in the price list | COPY: USD 0.005 / 1,000 | Class C: free | Class A: USD 4.50 / million |
| DeleteObject | free | not in the price list | free | Class A: free | free |

- **Hetzner Object Storage.** The FAQ on the
  [pricing page](https://www.hetzner.com/storage/object-storage/pricing/) ("What traffic is free
  of charge?") says: "S3 API calls such as PUT, GET, or DELETE requests". The
  [documentation](https://docs.hetzner.com/storage/object-storage/overview) ("Pricing") says:
  "S3 operations such as PUT, GET, or DELETE", with the note that the egress the operations cause
  is billed. Hetzner bills a base price of EUR 6.49 per month, which includes 1 TB of storage and
  1 TB of egress, plus extra storage and egress. Objects under 64 KB are billed as 64 KB.
- **DigitalOcean Spaces.** The [pricing page](https://docs.digitalocean.com/products/spaces/details/pricing/)
  lists only these charges:
  - the subscription (USD 5 per month, including 250 GiB of storage and 1,024 GiB of outbound transfer);
  - extra storage (USD 0.02 per GiB);
  - extra outbound transfer (USD 0.01 per GiB).

  There is no fee per request, but the documentation does not say in so many words that requests
  are free. The [limits page](https://docs.digitalocean.com/products/spaces/details/limits/)
  limits new buckets to 800 operations per second and says LIST "may be rate-limited further"
  under high load. A very large folder in the browser (one LIST per 1,000 objects) is where WRM
  can come close. The minimum billable object size is 4 KiB.
- **Amazon S3.** The [pricing page](https://aws.amazon.com/s3/pricing/) says: "You pay for requests
  made against your S3 buckets and objects", and "LIST requests … are charged at the same rate as
  S3 Standard PUT, COPY, and POST requests. DELETE and CANCEL requests are free." The
  [usage report reference](https://docs.aws.amazon.com/AmazonS3/latest/userguide/aws-usage-report-understand.html)
  puts PUT, COPY, POST and LIST in Requests-Tier1, and GET and "all other non-Tier1 requests"
  (HEAD) in Requests-Tier2. The prices in the table are the usual S3 Standard us-east-1 figures.
  WRM could not read the AWS price table itself when this was written, and other regions cost
  more, so check the table for your region.
- **Backblaze B2.** The [transaction pricing page](https://www.backblaze.com/cloud-storage/transaction-pricing)
  lists Class A, B and C calls as "Cost: Free". Backblaze bills storage after the first 10 GB, and
  egress beyond 3 times the average stored data.
- **Cloudflare R2.** The [pricing page](https://developers.cloudflare.com/r2/pricing/) lists
  ListObjects, PutObject and CopyObject as Class A, HeadBucket, HeadObject and GetObject as
  Class B, and DeleteObject as free. Each month the first 1 million Class A and 10 million Class B
  requests are free. Egress is free.
- **MinIO** and other self-hosted servers have no request fees. For any **other** provider WRM
  does not know the prices; the activity panel still counts every request.

## S3 activity panel

The storage window shows the requests WRM sends to the open storage live, in a panel beside the
file browser. Each line gives:
- the time, the S3 operation and the key (inside the storage prefix);
- the HTTP answer or the error, and how long the request took;
- what caused it: *you* (the file window), *background check*, *chain*, *AI connection* or
  *connection check*; a retry is marked ↻.

Above the list, each operation is counted for today and since WRM started (or the first request).
For the storage's provider each count is marked free or billed, and the panel estimates the cost
of the billed requests before any free allowance. The panel is shown by default; ✕ or the 📡
button hides it, and the browser remembers the choice. The counts live in memory only and hold no
secret and no object content. The checks *Test connection* runs before a storage is saved are not
counted, because the storage does not exist yet.

## Troubleshooting

| Message | Cause |
|---|---|
| "The bucket does not exist (check the name and the region)" | wrong bucket, wrong region, or the endpoint is not the S3 API |
| "Access denied: the key may not allow this operation" | the key policy does not allow list/get/put on that prefix |
| "The secret key or the region does not match the endpoint" | region mismatch (Hetzner uses the location code) |
| "The endpoint is in another region" | the bucket lives in another region than the one entered |
| "The certificate fingerprint does not match" | the endpoint certificate changed; verify it and update the pin |
| "This storage does not allow …" | the capability switch is off for that storage |

## Backup and rollback

The `storages` table is created on start; nothing else changes. A database backup covers it. To
remove the feature, set `storages_enabled=0`; to remove the data as well, drop the table. An old
installation may still have the unused `storage_shares` table from v14.0.0-v14.2.0; the share
feature and its `/s3/` route are gone, so those links no longer work and the table can be dropped.
