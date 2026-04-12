# CLI Reference

## `pipectl run`

Run a pipeline against data from `stdin`.

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
| `--output <path>` | `-o` | Write pipeline output to a file instead of stdout. Step logs are always written to stdout regardless of this flag. |
| `--verbose` | `-v` | Enable verbose logging. Prints per-step debug output — record counts, field operations, sort results — to stdout. |
| `--dry-run` | | Validate the pipeline config and print the ordered step list without executing any steps or reading stdin. |

### Examples

Run a pipeline with stdin:

```bash
pipectl run pipeline.yaml < input.json
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
- Input is always read from `stdin`. If nothing is piped in, the runtime executes but most pipelines will fail when the input format is parsed.
- Step logs (`log`, `count`) are written to `stdout`. Only the final payload output is affected by `-o`.
