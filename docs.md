# pipectl Documentation

`pipectl` runs a YAML-defined pipeline against data read from `stdin`.

The runtime flow is:

1. Load `pipeline.yaml`
2. Parse `stdin` using `input.format`
3. Run each step in order
4. Write the final payload to `stdout` using `output.format`

## Command Line Usage

Run a pipeline by passing the pipeline file as the single argument and piping input into `stdin`:

```bash
pipectl run pipeline.yaml < input.json
```

With `go run`:

```bash
go run ./cmd/pipectl run pipeline.yaml < input.json
```

Examples:

```bash
go run ./cmd/pipectl run examples/json/customer-signup-json.yaml < examples/json/input/customer-signup.json
go run ./cmd/pipectl run examples/csv/customer-signup-csv.yaml < examples/csv/customer-signup.csv
```

Write output to a file instead of `stdout`:

```bash
pipectl run pipeline.yaml -o output.json < input.json
```

Notes:

- `run` requires exactly one argument: the pipeline file path.
- Input is read from `stdin`. If nothing is piped in, the runtime still executes, but most pipelines will fail once the configured input format is parsed.
- Step logs are always written to `stdout`. Only the final payload output is affected by `-o`.
- `-o` / `--output`: optional path to write the pipeline output to a file.

## pipeline.yaml Format

Minimal structure:

```yaml
id: example-pipeline

input:
  format: json

steps:
  - log:
      message: Starting pipeline
  - convert:
      format: csv

output:
  format: csv
```

Top-level fields:

- `id`: pipeline identifier used only for console output.
- `input`: input configuration.
- `steps`: ordered list of steps. Each list item must contain exactly one step type.
- `output`: output configuration.

### `input`

Supported runtime field:

- `format`: `json`, `jsonl`, or `csv`

Other parsed fields currently exist in the schema but are not used by the runtime:

- `encoding`
- `schema`
- `delimiter`
- `has_header`
- `max_size`

### `output`

Supported field:

- `format`: `json`, `jsonl`, or `csv`

### Step Syntax

Each step is a single-key object:

```yaml
steps:
  - normalize:
      fields:
        email: lower
  - redact:
      fields: [password]
      strategy: mask
```

Supported step types:

- `validate-json`
- `normalize`
- `default`
- `cast`
- `convert`
- `assert`
- `rename`
- `redact`
- `select`
- `filter`
- `limit`
- `log`
- `count`
- `http-transform`

## Step Reference

If a step receives an unsupported payload type, execution fails.

### `validate-json`

Validates JSON or JSONL payloads against a JSON Schema.

Supported payloads:

- `json`
- `jsonl`

Options:

- `schema`: required. Either a path to a schema file or an inline JSON schema string.

Example:

```yaml
- validate-json:
    schema: ./schema.json
```

### `normalize`

Normalizes string fields in JSON, JSONL, or CSV payloads.

Supported payloads:

- `json`
- `jsonl`
- `csv`

Options:

- `fields`: map of field name to strategy

Supported strategies:

- `lower`
- `upper`
- `trim`
- `trim-left`
- `trim-right`
- `collapse-spaces`
- `capitalize`

Example:

```yaml
- normalize:
    fields:
      email: lower
      first_name: capitalize
      name: trim
```

Notes:

- Only string values are changed.
- Unknown strategies are effectively ignored at runtime.

### `default`

Adds missing fields to JSON/JSONL records or fills empty/missing CSV columns.

Supported payloads:

- `json`
- `jsonl`
- `csv`

Options:

- `fields`: map of field name to default value

Example:

```yaml
- default:
    fields:
      country: AU
      currency: AUD
```

Notes:

- For JSON and JSONL, defaults are applied only when the field does not already exist.
- For CSV, a missing column is added to the header and populated for all rows.
- For CSV, an existing column is only filled where the cell is empty.

### `cast`

Converts field values to a different type.

Supported payloads:

- `json`
- `jsonl`

Options:

- `fields`: map of field path to cast configuration

Each field configuration supports:

- `type`: required. One of `int`, `float`, `bool`, `time`, `string`
- `format`: optional. Date/time parse format (Go layout string). Only valid for `type: time`. Defaults to RFC 3339.
- `true_values`: optional list of strings to treat as `true`. Only valid for `type: bool`.
- `false_values`: optional list of strings to treat as `false`. Only valid for `type: bool`.

Example:

```yaml
- cast:
    fields:
      age:
        type: int
      price:
        type: float
      active:
        type: bool
        true_values: [yes, "1"]
        false_values: [no, "0"]
      created_at:
        type: time
        format: "2006-01-02"
```

Notes:

- Field paths support dot notation and array indexing, eg: `user.address`, `tags[0]`.
- For `type: bool`, default true values are `true`, `t`, `1`, `yes`, `y`, `on`. Default false values are `false`, `f`, `0`, `no`, `n`, `off`.
- `true_values` and `false_values` must not overlap.
- For `type: time`, the field value is parsed into a Go `time.Time`. The serialized output format depends on the output encoder.
- Array fields are cast element-by-element.

### `convert`

Converts the payload to another format.

Supported payloads:

- `json`
- `jsonl`
- `csv`

Options:

- `format`: required. One of `json`, `jsonl`, `csv`

Example:

```yaml
- convert:
    format: csv
```

Notes:

- Converting CSV to JSON/JSONL uses the first row as headers.
- Dot-separated CSV headers such as `user.name` become nested JSON objects.

### `assert`

Checks record-count and field-existence conditions.

Supported payloads:

- `json`
- `jsonl`
- `csv`

Options:

- `min-records`: optional integer, must be `>= 0`
- `max-records`: optional integer, must be `>= 0`
- `records-equal`: optional integer, must be `>= 0`
- `field-exists`: optional string

At least one option is required.

Example:

```yaml
- assert:
    min-records: 10
    max-records: 1000
    field-exists: email
```

Notes:

- `min-records` must be `<= max-records` when both are set.
- `records-equal` must fit within the min/max bounds when they are also set.
- For CSV, `field-exists` checks the header row.
- For JSON and JSONL, `field-exists` passes if any record contains the field.

### `rename`

Renames fields in JSON/JSONL records or CSV headers.

Supported payloads:

- `json`
- `jsonl`
- `csv`

Options:

- `fields`: map of current field name to new field name

Example:

```yaml
- rename:
    fields:
      first_name: firstName
      credit_card: creditCard
```

### `redact`

Redacts selected fields.

Supported payloads:

- `json`
- `jsonl`
- `csv`

Options:

- `fields`: list of field names
- `strategy`: optional redaction strategy

Supported strategies:

- `mask`: replace each character with `*`
- `sha256`: replace with SHA-256 hex digest

If `strategy` is omitted or unknown, the replacement value is `REDACTED`.

Example:

```yaml
- redact:
    fields: [credit_card, password]
    strategy: mask
```

Notes:

- JSON/JSONL redaction only handles top-level string fields.
- Non-string JSON values are not redacted.

### `select`

Keeps only selected CSV columns.

Supported payloads:

- `csv`

Options:

- `fields`: list of column names to keep

Example:

```yaml
- select:
    fields: [first_name, email, dob]
```

### `filter`

Keeps only records where one field matches a value.

Supported payloads:

- `csv`
- `json`
- `jsonl`

Options:

- `field`: column name to test
- `equals`: required match value

Example:

```yaml
- filter:
    field: country
    equals: AU
```

### `limit`

Truncates the payload to at most N records.

Supported payloads:

- `json`
- `jsonl`
- `csv`

Options:

- `count`: required integer, must be `>= 1`

Example:

```yaml
- limit:
    count: 100
```

Notes:

- If the payload already has fewer records than `count`, it passes through unchanged.
- For CSV, the header row is always preserved.
- Useful for sampling large inputs, capping output size before an `http-transform`, or testing a pipeline end-to-end with real data.

### `log`

Prints a message, record count, and sample rows/records.

Supported payloads:

- `json`
- `jsonl`
- `csv`

Options:

- `message`: optional text
- `count`: optional boolean, defaults to `true`
- `sample`: optional integer, defaults to `10`

Example:

```yaml
- log:
    message: After normalization
    count: true
    sample: 5
```

Notes:

- `sample: 0` disables sample output.
- Negative sample values are treated like `0`.

### `count`

Prints the current record count.

Supported payloads:

- `json`
- `jsonl`
- `csv`

Options:

- `message`: optional text printed before the count

Example:

```yaml
- count:
    message: Final record count
```

### `http-transform`

Sends the current JSON or JSONL payload to an HTTP endpoint and replaces the payload with the response.

Supported payloads:

- `json`
- `jsonl`

Options:

- `url`: target URL
- `method`: HTTP method
- `proxy`: optional proxy URL
- `headers`: optional string map of request headers
- `timeout`: optional timeout in seconds
- `expect-format`: optional response format: `json`, `jsonl`, or `csv`

Example:

```yaml
- http-transform:
    url: https://example.com/transform
    method: POST
    timeout: 60
    headers:
      Authorization: Bearer token
    expect-format: json
```

Notes:

- If `expect-format` is omitted, it defaults to `json`.
- `timeout` defaults to `60` seconds and must be between `1` and `300` when set.
- Only HTTP `200 OK` responses are accepted.
- Response `Content-Type` must match `expect-format`.
- Request bodies are only sent for `POST`, `PUT`, `PATCH`, and `DELETE`.
- For JSONL requests without an explicit `Content-Type`, the step sends `application/x-ndjson`.

## Example Pipelines

### Customer Signup Validation And API Enrichment

Use this when ingesting signup events, enforcing a schema, normalizing fields, removing secrets, then posting the cleaned records to another service.

```yaml
id: customer-signup

input:
  format: json

steps:
  - validate-json:
      schema: ./examples/json/schema/signup-schema.json

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
      url: https://httpbin.org/post
      method: POST
      expect-format: json

output:
  format: json
```

Example input file `customer-signup-input.json`:

```json
{
  "email": "  ALICE@EXAMPLE.COM ",
  "name": " Alice Smith ",
  "credit_card": "4111111111111111",
  "password": "secret123"
}
```

Run it:

```bash
go run ./cmd/pipectl run customer-signup.yaml < customer-signup-input.json
```

### Partner CSV Intake For Local Processing

Use this when a partner sends CSV exports that need to be cleaned, filtered to a market, reduced to a smaller set of columns, and emitted as JSONL for downstream processing.

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

Example input file `partner-customers.csv`:

```csv
first_name,last_name,email,country,credit_card,password
alice,smith,ALICE@EXAMPLE.COM,au,4111111111111111,secret123
bob,jones,BOB@EXAMPLE.COM,nz,5555444433332222,secret456
carol,lee,CAROL@EXAMPLE.COM,,4444333322221111,secret789
david,ng,DAVID@EXAMPLE.COM,au,4012888888881881,secret234
emma,brown,EMMA@EXAMPLE.COM,us,4222222222222,secret345
frank,wilson,FRANK@EXAMPLE.COM,au,3530111333300000,secret456
grace,taylor,GRACE@EXAMPLE.COM,au,5555555555554444,secret567
henry,anderson,HENRY@EXAMPLE.COM,sg,378282246310005,secret678
ivy,thomas,IVY@EXAMPLE.COM,au,6011111111111117,secret789
jack,white,JACK@EXAMPLE.COM,,30569309025904,secret890
```

Run it:

```bash
go run ./cmd/pipectl run partner-customer-import.yaml < partner-customers.csv
```

### Audit Export With Stable Hashing

Use this when you need to prepare an export for analytics or support teams without exposing raw personal data.

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

Example input file `audit-events.jsonl`:

```jsonl
{"email":"Alice@example.com","phone":"0400000001","status":"active"}
{"email":"Bob@example.com","phone":"0400000002","status":"inactive"}
{"email":"Carol@example.com","phone":"0400000003","status":"active"}
{"email":"David@example.com","phone":"0400000004","status":"pending"}
{"email":"Emma@example.com","phone":"0400000005","status":"active"}
{"email":"Frank@example.com","phone":"0400000006","status":"inactive"}
{"email":"Grace@example.com","phone":"0400000007","status":"active"}
{"email":"Henry@example.com","phone":"0400000008","status":"pending"}
{"email":"Ivy@example.com","phone":"0400000009","status":"inactive"}
{"email":"Jack@example.com","phone":"0400000010","status":"active"}
```

Run it:

```bash
go run ./cmd/pipectl run audit-export.yaml < audit-events.jsonl
```

### JSON To CSV Report Generation

Use this when application events arrive as JSON but an operations or finance team needs a flat CSV report.

```yaml
id: billing-report

input:
  format: json

steps:
  - validate-json:
      schema: ./schemas/invoice-array.schema.json

  - normalize:
      fields:
        customer_email: lower

  - rename:
      fields:
        total_amount: totalAmount
        created_at: createdAt

  - convert:
      format: csv

  - log:
      message: Report ready for export
      count: true
      sample: 5

output:
  format: csv
```

Note:

- `normalize` and `rename` only operate on top-level fields. If your JSON is nested, flatten or reshape it before relying on those steps.

Example input file `billing-report-input.json`:

```json
[
  {
    "invoice_id": "inv-1001",
    "customer_email": " FINANCE@EXAMPLE.COM ",
    "total_amount": 149.95,
    "created_at": "2026-03-20T10:00:00Z"
  },
  {
    "invoice_id": "inv-1002",
    "customer_email": " OPS@EXAMPLE.COM ",
    "total_amount": 89.50,
    "created_at": "2026-03-20T11:00:00Z"
  }
]
```

Run it:

```bash
go run ./cmd/pipectl run billing-report.yaml < billing-report-input.json
```

### Service-To-Service Transform With CSV Response

Use this when a JSON payload must be sent to an HTTP service that returns CSV for a legacy system.

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

Example input file `pricing-request.json`:

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

Run it:

```bash
go run ./cmd/pipectl run pricing-sync.yaml < pricing-request.json
```

Note:

- This example depends on a reachable HTTP service at the configured `url`.

### JSONL Cleanup Before Re-Publishing

Use this when a stream of records already exists in JSONL and needs a lightweight cleanup pass before being published again.

```yaml
id: event-republish

input:
  format: jsonl

steps:
  - log:
      message: Incoming event batch
      sample: 2

  - default:
      fields:
        pipeline_version: v1

  - normalize:
      fields:
        event_type: lower
        source: trim

  - redact:
      fields: [ip_address]
      strategy: mask

  - assert:
      field-exists: event_type

output:
  format: jsonl
```

Example input file `event-republish-input.jsonl`:

```jsonl
{"event_type":" USER.CREATED ","source":" mobile-app ","ip_address":"203.0.113.10"}
{"event_type":"ORDER.PLACED","source":" web ","ip_address":"203.0.113.11"}
{"event_type":" PASSWORD.RESET ","source":" support ","ip_address":"203.0.113.12"}
{"event_type":"PROFILE.UPDATED","source":" web ","ip_address":"203.0.113.13"}
{"event_type":" USER.LOGGED_IN ","source":" mobile-app ","ip_address":"203.0.113.14"}
{"event_type":"ORDER.CANCELLED","source":" api ","ip_address":"203.0.113.15"}
{"event_type":" EMAIL.VERIFIED ","source":" web ","ip_address":"203.0.113.16"}
{"event_type":"SUBSCRIPTION.STARTED","source":" billing ","ip_address":"203.0.113.17"}
{"event_type":" USER.LOGGED_OUT ","source":" mobile-app ","ip_address":"203.0.113.18"}
{"event_type":"ORDER.REFUNDED","source":" operations ","ip_address":"203.0.113.19"}
```

Run it:

```bash
go run ./cmd/pipectl run event-republish.yaml < event-republish-input.jsonl
```
