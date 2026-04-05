# pipectl

A CLI tool for running YAML-defined data pipelines. Reads from `stdin`, applies an ordered sequence of steps, and writes to `stdout`.

## Install

```bash
go install github.com/shanebell/pipectl/cmd/pipectl@latest
```

Or build from source:

```bash
go build ./cmd/pipectl
```

## Quick Start

Define a pipeline in YAML:

```yaml
id: my-pipeline

input:
  format: csv

steps:
  - normalize:
      fields:
        email: lower
        country: upper

  - filter:
      field: country
      equals: AU

  - redact:
      fields: [credit_card]
      strategy: mask

output:
  format: jsonl
```

Run it:

```bash
pipectl run my-pipeline.yaml < input.csv
```

Write output to a file:

```bash
pipectl run my-pipeline.yaml -o output.jsonl < input.csv
```

## Documentation

See [DOCS.md](DOCS.md) for the full step reference, all configuration options, and example pipelines.
