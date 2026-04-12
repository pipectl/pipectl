# Core Concepts

## Pipeline

A pipeline is a YAML file that describes how data should be transformed. It has four top-level fields:

```yaml
id: my-pipeline        # identifier used in console output

input:
  format: csv          # how to parse stdin

steps:                 # ordered list of transformations
  - normalize:
      fields:
        email: lower

output:
  format: jsonl        # how to write the result
```

The runtime flow is:

1. Parse `stdin` using `input.format`
2. Run each step in order, passing the payload from one to the next
3. Write the final payload to `stdout` using `output.format`

## Step

A step is a single transformation applied to the payload. Each step in the `steps` list is a single-key object where the key is the step type and the value is its configuration:

```yaml
steps:
  - normalize:          # step type
      fields:           # step configuration
        email: lower
  - redact:
      fields: [password]
      strategy: mask
```

Steps run in order. If any step fails, the pipeline stops and an error is printed.

See the [Step Reference](./steps/) for all available steps.

## Payload

The payload is the data flowing through the pipeline. It starts as the parsed contents of `stdin` and is passed from step to step. Three payload formats are supported:

| Format | Description |
|--------|-------------|
| `json` | A single JSON object or a JSON array of objects |
| `jsonl` | One JSON object per line (newline-delimited JSON) |
| `csv` | Comma-separated values with a header row |

Every step declares which payload types it supports. If a step receives an unsupported payload type, the pipeline fails with a clear error.

Use the [`convert`](./steps/convert) step to change the payload format mid-pipeline.

## Input and output format

The `input.format` tells pipectl how to parse `stdin`. The `output.format` tells it how to serialise the final payload to `stdout`.

They do not need to match — you can read CSV and write JSONL, for example.

```yaml
input:
  format: csv
output:
  format: jsonl
```

### CSV delimiter

If your CSV uses a non-standard separator, set `delimiter` under `input`:

```yaml
input:
  format: csv
  delimiter: ";"
```

## Records

Steps that filter, sort, or count operate on *records* — the individual items within the payload:

| Payload | Records |
|---------|---------|
| JSON array | Each element of the array |
| JSON object | A single record |
| JSONL | Each line |
| CSV | Each data row (excluding the header) |

## Nested fields

For JSON and JSONL payloads, some steps support dot notation and array indexing to reference nested fields:

```yaml
- cast:
    fields:
      user.age:
        type: int
      tags[0]:
        type: string
```

CSV headers with dots (e.g. `user.name`) are automatically treated as nested objects when converting to JSON/JSONL.
