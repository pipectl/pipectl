# Pipectl

Pipe Control - a pipeline execution tool

# Usage

## CSV example

Transform, filter and redact data from a CSV file:

```yaml
id: customer-pipeline

input:
  format: csv

steps:
  - normalize:
      fields:
        country: upper
        first_name: capitalize
        last_name: capitalize

  - filter:
      field: country
      equals: AU

  - redact:
      fields: [credit_card, password]
      strategy: mask

output:
  format: csv
```

## JSON example

Validate and transform JSON data and POST to an API endpoint

```yaml
id: customer-signup

input:
  format: json

steps:
  - validate-json:
      schema: ./schema.json

  - normalize:
      fields:
        email: lower
        name: trim
        string: collapse-spaces

  - redact:
      strategy: mask
      fields: [credit_card, password]

  - http-transform:
      url: https://example.com/register-customer
      method: POST

output:
  format: json
```