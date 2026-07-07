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
| `--output <path>` | `-o` | Write pipeline output to a file instead of stdout. Step logs are always written to stderr regardless of this flag. |
| `--verbose` | `-v` | Enable verbose logging. Prints per-step debug output — record counts, field operations, sort results — to stderr. |
| `--quiet` | `-q` | Suppress all diagnostic output. Only the final payload is written. Useful for scripting or when stderr noise is unwanted. |
| `--dry-run` | | Validate the pipeline config and print the ordered step list without executing any steps or reading input. |
| `--timing` | | Print a per-step table (duration, records in/out) to stderr after execution. Suppressed by `--quiet`. |
| `--var KEY=VALUE` | | Substitute `${VAR}` tokens in pipeline YAML before parsing. Repeatable. |
| `--max-input-size <size>` | | Maximum size of input read from stdin or `--input`, e.g. `64KB`, `256MB`, `1GB`. Input exceeding this size is rejected. Default `256MB`. |

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

Suppress diagnostic output:

```bash
pipectl run pipeline.yaml --quiet < input.json
```

Validate a pipeline without running it:

```bash
pipectl run pipeline.yaml --dry-run
```

Substitute variables in a pipeline:

```bash
pipectl run pipeline.yaml --var ENV=prod --var LIMIT=100 < input.json
```

### Notes

- `run` requires exactly one argument: the pipeline file path.
- Input can be provided via `--input <file>` or piped through stdin. If neither is provided, the runtime executes but most pipelines will fail when the input format is parsed.
- `--input` and piped stdin cannot be used together.
- Step logs (`log`, `count`) are written to `stderr`. Only the final payload output is affected by `-o`.
- `--var` substitutes `${KEY}` tokens in the pipeline YAML before parsing. Tokens left unresolved after all substitutions are applied cause an error at startup.

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

---

## `pipectl docs`

Show built-in documentation for pipeline steps.

```bash
pipectl docs [step]
```

With no argument, lists all available steps with a one-line description. With a step name, prints the full documentation for that step.

### Examples

List all steps:

```bash
pipectl docs
```

Show documentation for a specific step:

```bash
pipectl docs filter
pipectl docs cast
pipectl docs validate-json
```

### Notes

- Output is formatted for the terminal when stdout is a TTY. When piped or redirected, plain markdown is written instead.
- Step names match the keys used in pipeline YAML (e.g. `http-transform`, `validate-json`).
