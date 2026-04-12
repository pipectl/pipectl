# Payload Formats

pipectl supports three payload formats: `json`, `jsonl`, and `csv`. The format is specified in `input.format` and `output.format`, and can be changed mid-pipeline with the [`convert`](./steps/convert) step.

## JSON

A single JSON object or a JSON array of objects.

```json
[
  {"name": "Alice", "email": "alice@example.com"},
  {"name": "Bob",   "email": "bob@example.com"}
]
```

A single object is valid too and is treated as a one-record payload:

```json
{"name": "Alice", "email": "alice@example.com"}
```

**Records:** each element of the array. A single object is wrapped internally as a one-element array.

**Nested fields:** supported via dot notation (`user.address`) and array indexing (`tags[0]`) in the `cast` step.

## JSONL

One JSON object per line — also known as newline-delimited JSON (NDJSON).

```jsonl
{"name":"Alice","email":"alice@example.com"}
{"name":"Bob","email":"bob@example.com"}
```

**Records:** each line.

**Notes:**
- No blank lines between records.
- Zero-record payloads (empty input) are valid.

## CSV

Comma-separated values with a header row.

```csv
name,email,country
Alice,alice@example.com,AU
Bob,bob@example.com,NZ
```

**Records:** each data row. The header row is preserved but not counted as a record.

**Custom delimiter:** set `delimiter` under `input` to use a different separator:

```yaml
input:
  format: csv
  delimiter: ";"
```

**Nested fields:** CSV headers containing dots (e.g. `user.name`) are converted to nested JSON objects when the payload is converted to JSON or JSONL.

## Conversion

Use the [`convert`](./steps/convert) step to change format mid-pipeline. The input and output formats do not need to match.

| From | To | Notes |
|------|----|-------|
| CSV → JSON/JSONL | First row becomes field names. Dot-separated headers become nested objects. |
| JSON/JSONL → CSV | Records are flattened. Nested objects become dot-separated headers. Arrays and objects in field values are JSON-encoded as strings. |
| JSON ↔ JSONL | Direct conversion between array and line-delimited format. |
