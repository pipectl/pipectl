---
layout: home

hero:
  name: pipectl
  text: YAML-defined data pipelines
  tagline: Read from stdin. Transform through ordered steps. Write to stdout.
  image:
    light: /logo-light.svg
    dark: /logo-dark.svg
    alt: pipectl
  actions:
    - theme: brand
      text: Get Started
      link: /getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/pipectl/pipectl

features:
  - icon: 🔗
    title: Composable
    details: Chain any number of steps in any order. The output of each step is the input to the next.
  - icon: 📄
    title: Multiple payload formats
    details: JSON, JSONL, and CSV — with automatic conversion between them at any point in the pipeline.
  - icon: ⚡
    title: Built-in steps
    details: Filter, normalize, redact, cast, sort, dedupe, validate, convert, make HTTP calls, and more.
  - icon: 🔒
    title: Built-in redaction
    details: Mask or SHA-256 hash sensitive fields before they leave your pipeline.
  - icon: ✅
    title: JSON Schema validation
    details: Validate records against a JSON Schema and fail fast if the data doesn't conform.
  - icon: 🪶
    title: Minimal by design
    details: A single binary, no runtime dependencies, no framework overhead. Just YAML and your data.
---

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
