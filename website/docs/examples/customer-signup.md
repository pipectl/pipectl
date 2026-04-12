# Customer Signup Validation and Enrichment

Use this pattern when ingesting signup events: enforce a schema, normalise fields, remove secrets, then post the cleaned record to another service.

**Input:** `json` → **Output:** `json`

**Steps used:** `validate-json`, `default`, `normalize`, `redact`, `http-transform`

## Pipeline

```yaml
id: customer-signup

input:
  format: json

steps:
  - validate-json:
      schema: ./schemas/signup-schema.json

  - default:
      fields:
        country: AU
        currency: AUD

  - normalize:
      fields:
        email: lower
        name: trim

  - redact:
      fields: [credit_card, password]
      strategy: mask

  - http-transform:
      url: https://api.example.com/customers
      method: POST
      timeout: 60
      headers:
        Authorization: Bearer ${API_TOKEN}
      expect-format: json

output:
  format: json
```

## Sample input

```json
{
  "email": "  ALICE@EXAMPLE.COM ",
  "name": " Alice Smith ",
  "credit_card": "4111111111111111",
  "password": "secret123"
}
```

## What happens

1. **validate-json** — fails fast if required fields are missing or types are wrong
2. **default** — adds `country: AU` and `currency: AUD` if not present
3. **normalize** — lowercases `email` and trims whitespace from `name`
4. **redact** — masks `credit_card` and `password` before they leave the pipeline
5. **http-transform** — POSTs the cleaned record; replaces payload with the response

## Run it

```bash
pipectl run customer-signup.yaml < signup.json
```
