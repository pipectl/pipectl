# CLI Reference

## Global flags

| Flag | Short | Description |
|------|-------|-------------|
| `--version` | | Print the version and exit. |

```bash
pipectl --version
# pipectl version v1.2.0
```

---

## `pipectl run`

Run a pipeline against input data.

```bash
pipectl run <pipeline.yaml> [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `pipeline.yaml` | Required. Path to the pipeline YAML file. |

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--input <path>` | `-i` | Read pipeline input from a file. Alternative to piping via stdin. Cannot be combined with piped stdin. |
| `--output <path>` | `-o` | Write pipeline output to a file instead of stdout. Step logs are always written to stdout regardless of this flag. |
| `--verbose` | `-v` | Enable verbose logging. Prints per-step debug output — record counts, field operations, sort results — to stdout. |
| `--dry-run` | | Validate the pipeline config and print the ordered step list without executing any steps or reading input. |

### Examples

Run a pipeline with stdin:

```bash
pipectl run pipeline.yaml < input.json
```

Run a pipeline with an input file:

```bash
pipectl run pipeline.yaml --input input.json
```

Write output to a file:

```bash
pipectl run pipeline.yaml -o output.jsonl < input.csv
```

Enable verbose logging:

```bash
pipectl run pipeline.yaml --verbose < input.json
```

Validate a pipeline without running it:

```bash
pipectl run pipeline.yaml --dry-run
```

### Notes

- `run` requires exactly one argument: the pipeline file path.
- Input can be provided via `--input <file>` or piped through stdin. If neither is provided, the runtime executes but most pipelines will fail when the input format is parsed.
- `--input` and piped stdin cannot be used together.
- Step logs (`log`, `count`) are written to `stdout`. Only the final payload output is affected by `-o`.

---

## `pipectl validate`

Validate a pipeline file without executing it or reading any input data.

```bash
pipectl validate <pipeline.yaml>
```

Checks that the pipeline YAML is well-formed, all step configs are valid, and the pipeline can be compiled. Exits `0` silently on success. Exits non-zero and prints an error on failure.

### Arguments

| Argument | Description |
|----------|-------------|
| `pipeline.yaml` | Required. Path to the pipeline YAML file. |

### Examples

```bash
pipectl validate pipeline.yaml
```

```bash
pipectl validate bad-pipeline.yaml
# [line 4] input format must be one of: json, jsonl, csv
# exit status 1
```

### Notes

- No input data is required — useful for CI pre-flight checks.
- Validates both the pipeline spec (YAML structure, step configs) and the execution plan (step compilation).
