<p align="center">
  <img src="website/public/logo-light.svg" alt="pipectl" height="60" />
</p>

<p align="center">
  Run YAML-defined data pipelines from the command line.<br/>
  Read from <code>stdin</code>, transform through ordered steps, write to <code>stdout</code>.
</p>

<p align="center">
  <a href="https://github.com/pipectl/pipectl/actions/workflows/ci.yml"><img src="https://github.com/pipectl/pipectl/actions/workflows/ci.yml/badge.svg" alt="CI"/></a>
  <a href="https://goreportcard.com/report/github.com/pipectl/pipectl"><img src="https://goreportcard.com/badge/github.com/pipectl/pipectl" alt="Go Report Card"/></a>
  <a href="https://pkg.go.dev/github.com/pipectl/pipectl"><img src="https://pkg.go.dev/badge/github.com/pipectl/pipectl.svg" alt="Go Reference"/></a>
  <img src="https://img.shields.io/github/license/pipectl/pipectl" alt="License"/>
</p>

---

## What is pipectl?

pipectl is a CLI for running YAML-defined data pipelines. You define an ordered list of steps in a YAML file, pipe your data in, and get transformed data out.

It supports JSON, JSONL, and CSV payloads — including conversions between them — and ships with 16 built-in steps for filtering, normalising, redacting, casting, sorting, validating, and more.

## Quick example

```yaml
# customer-intake.yaml
id: customer-intake
input:
  format: csv
steps:
  - normalize:
      fields:
        first_name: capitalize
        last_name: capitalize
        email: lower
  - filter:
      field: country
      equals: AU
  - redact:
      fields: [credit_card]
      strategy: mask
  - select:
      fields: [first_name, last_name, email, plan]
output:
  format: jsonl
```

```bash
pipectl run customer-intake.yaml < customers.csv
```

```jsonl
{"first_name":"Alice","last_name":"Smith","email":"alice@example.com","plan":"pro"}
{"first_name":"Bob","last_name":"Jones","email":"bob@example.com","plan":"free"}
```

## Features

- **16 built-in steps** — filter, normalize, redact, cast, sort, dedupe, validate, convert, and more
- **Three payload formats** — JSON, JSONL, and CSV with automatic conversion between them
- **Composable pipelines** — chain any number of steps; output of one feeds the next
- **JSON Schema validation** — validate records against a schema at any point in the pipeline
- **HTTP transforms** — POST/PUT your payload to an HTTP endpoint and continue with the response
- **Dry-run mode** — validate and preview a pipeline without reading stdin or executing steps
- **Minimal dependencies** — a single binary with no runtime requirements

## Install

```bash
go install github.com/pipectl/pipectl/cmd/pipectl@latest
```

Or build from source:

```bash
git clone https://github.com/pipectl/pipectl.git
cd pipectl
go build ./cmd/pipectl
```

## Usage

```bash
pipectl run <pipeline.yaml> [flags] < input

Flags:
  -o, --output <path>   Write output to file instead of stdout
  -v, --verbose         Enable verbose logging
      --dry-run         Validate pipeline and print plan without executing
```

## Core concepts

| Concept | Description |
|---------|-------------|
| **Pipeline** | A YAML file with an `id`, `input` format, ordered `steps`, and `output` format |
| **Step** | A single operation that reads the payload, transforms it, and passes it on |
| **Payload** | The data flowing through the pipeline — a JSON object/array, JSONL stream, or CSV |

## Steps

| Step | What it does |
|------|-------------|
| `normalize` | Normalise string fields (lower, upper, trim, capitalize, collapse-spaces) |
| `filter` | Keep records matching a condition or nested `all`/`any` group |
| `redact` | Replace field values with `mask`, `sha256`, or `REDACTED` |
| `cast` | Convert field types — int, float, bool, string, time |
| `default` | Fill missing or empty fields with a default value |
| `select` | Keep only the specified fields |
| `rename` | Rename fields |
| `sort` | Sort records by a field, ascending or descending |
| `limit` | Truncate to N records |
| `dedupe` | Remove duplicate records by key fields |
| `convert` | Convert payload format — json ↔ jsonl ↔ csv |
| `validate-json` | Validate records against a JSON Schema |
| `assert` | Assert record count or field existence |
| `http-transform` | POST/PUT/PATCH payload to an HTTP endpoint, continue with the response |
| `log` | Print current record count and samples to stderr |
| `count` | Print current record count to stderr |

## Documentation

Full documentation, step reference, and example pipelines are at **[pipectl.github.io](https://pipectl.github.io)**.

The `examples/` directory contains ready-to-run pipelines with sample data.

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](website/docs/contributing.md) or the [contributing guide](https://pipectl.github.io/contributing) for architecture details and a step-by-step guide to adding new steps.

## License

MIT — see [LICENSE](LICENSE).
