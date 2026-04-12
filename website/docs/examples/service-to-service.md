# Service-to-Service Transform with CSV Response

Use this pattern when a JSON payload must be sent to an HTTP service that returns CSV — for example, a pricing or reporting service used by a legacy system.

**Input:** `json` → **Output:** `csv`

**Steps used:** `validate-json`, `http-transform`, `select`, `count`

## Pipeline

```yaml
id: pricing-sync

input:
  format: json

steps:
  - validate-json:
      schema: ./schemas/pricing-request.schema.json

  - http-transform:
      url: https://api.example.com/v1/pricing/export
      method: POST
      timeout: 30
      headers:
        Authorization: Bearer ${PRICING_TOKEN}
        X-Source-System: pipectl
      expect-format: csv

  - select:
      fields: [sku, region, price, currency]

  - count:
      message: Rows returned from pricing service

output:
  format: csv
```

## Sample input

```json
{
  "account_id": "acct-123",
  "effective_date": "2026-03-22",
  "items": [
    {"sku": "SKU-001", "region": "AU"},
    {"sku": "SKU-002", "region": "NZ"}
  ]
}
```

## What happens

1. **validate-json** — validates the request payload before sending it anywhere
2. **http-transform** — POSTs the JSON to the pricing service; the response is CSV and becomes the new payload
3. **select** — keeps only the four columns needed downstream, dropping any extra fields the service returns
4. **count** — prints how many rows were returned

## Notes

- `expect-format: csv` tells `http-transform` to parse the response as CSV. The response `Content-Type` must be `text/csv`.
- This example requires a reachable HTTP service at the configured `url`.
- Environment variable substitution in headers (e.g. `${PRICING_TOKEN}`) is not automatic — use your shell or a secrets manager to inject the value before running.

## Run it

```bash
PRICING_TOKEN=mytoken pipectl run pricing-sync.yaml < pricing-request.json
```
