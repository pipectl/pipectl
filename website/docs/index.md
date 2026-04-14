---
layout: home

hero:
  name: pipectl
  text: YAML-defined data pipelines
  tagline: Read from files or stdin. Transform through ordered steps. Write to files or stdout.
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

**Input: `customers.csv`**

```csv
first_name,last_name,email,country,plan,credit_card
alice,smith,Alice@Example.com,AU,pro,4111111111111111
BOB,JONES,BOB@EXAMPLE.COM,AU,free,5500005555555559
carol,white,carol@example.com,US,pro,3714496353984312
```

**Pipeline YAML: `customer-intake.yaml`**

```yaml
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
      fields: [first_name, last_name, email, credit_card]
output:
  format: jsonl
```

**Command**

```bash
pipectl run customer-intake.yaml < customers.csv
```

**Output**

```jsonl
{"first_name":"Alice","last_name":"Smith","email":"alice@example.com","credit_card":"****************"}
{"first_name":"Bob","last_name":"Jones","email":"bob@example.com","credit_card":"****************"}
```
