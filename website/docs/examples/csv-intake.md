# CSV Intake for Local Processing

Use this pattern when a partner sends CSV exports that need cleaning, filtering to a market, column reduction, and emission as JSONL for downstream processing.

**Input:** `csv` → **Output:** `jsonl`

**Steps used:** `log`, `default`, `normalize`, `filter`, `redact`, `select`, `rename`, `assert`

## Pipeline

```yaml
id: partner-customer-import

input:
  format: csv

steps:
  - log:
      message: Raw partner file received
      sample: 3

  - default:
      fields:
        country: AU
        source: partner-upload

  - normalize:
      fields:
        first_name: capitalize
        last_name: capitalize
        email: lower
        country: upper

  - filter:
      field: country
      equals: AU

  - redact:
      fields: [credit_card, password]
      strategy: mask

  - select:
      fields: [first_name, last_name, email, country, source]

  - rename:
      fields:
        first_name: firstName
        last_name: lastName

  - assert:
      min-records: 1
      field-exists: email

output:
  format: jsonl
```

## Sample input

```csv
first_name,last_name,email,country,credit_card,password
alice,smith,ALICE@EXAMPLE.COM,au,4111111111111111,secret123
bob,jones,BOB@EXAMPLE.COM,nz,5555444433332222,secret456
carol,lee,CAROL@EXAMPLE.COM,,4444333322221111,secret789
david,ng,DAVID@EXAMPLE.COM,au,4012888888881881,secret234
```

## What happens

1. **log** — prints a sample of the raw file to help with debugging
2. **default** — fills in `country: AU` where missing; tags every record with `source: partner-upload`
3. **normalize** — capitalises names, lowercases emails, uppercases country codes
4. **filter** — keeps only Australian customers
5. **redact** — masks credit card and password fields
6. **select** — drops columns not needed downstream
7. **rename** — converts snake_case to camelCase for the downstream system
8. **assert** — guarantees at least one record with an email field made it through

## Run it

```bash
pipectl run partner-customer-import.yaml < partner-customers.csv
```
