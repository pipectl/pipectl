# Audit Export with Stable Hashing

Use this pattern when you need to prepare an export for analytics or support teams without exposing raw personal data. SHA-256 hashing is stable — the same input always produces the same hash — so records can be correlated across exports without revealing the original values.

**Input:** `jsonl` → **Output:** `jsonl`

**Steps used:** `assert`, `normalize`, `redact`, `default`, `count`

## Pipeline

```yaml
id: audit-export

input:
  format: jsonl

steps:
  - assert:
      min-records: 1

  - normalize:
      fields:
        email: lower
        status: upper

  - redact:
      fields: [email, phone]
      strategy: sha256

  - default:
      fields:
        exported_by: pipectl

  - count:
      message: Records prepared for audit export

output:
  format: jsonl
```

## Sample input

```jsonl
{"email":"Alice@example.com","phone":"0400000001","status":"active"}
{"email":"Bob@example.com","phone":"0400000002","status":"inactive"}
{"email":"Carol@example.com","phone":"0400000003","status":"active"}
```

## What happens

1. **assert** — fails immediately if the input is empty, preventing a silent empty export
2. **normalize** — lowercases emails (so hashes are consistent) and uppercases status values
3. **redact** — replaces email and phone with their SHA-256 digests
4. **default** — stamps each record with `exported_by: pipectl` for auditability
5. **count** — prints the final record count before writing

## Run it

```bash
pipectl run audit-export.yaml < audit-events.jsonl
```
