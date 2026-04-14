# Getting Started

## Installation

### macOS (Homebrew)

```bash
brew install pipectl/pipectl/pipectl
```

### Windows

**Direct download:** grab `pipectl_<version>_windows_amd64.zip` from the [Releases page](https://github.com/pipectl/pipectl/releases/latest), extract, and add the folder to your `PATH`.

**Scoop** (if you have Scoop installed):

```powershell
scoop bucket add pipectl https://github.com/pipectl/scoop-pipectl
scoop install pipectl
```

### Linux

Download a `.deb` or `.rpm` from the [Releases page](https://github.com/pipectl/pipectl/releases/latest), or use the `.tar.gz` archive on any distribution.

**Debian/Ubuntu:**
```bash
sudo dpkg -i pipectl_<version>_linux_amd64.deb
```

**Fedora/RHEL:**
```bash
sudo rpm -i pipectl_<version>_linux_amd64.rpm
```

### Go install (from source)

```bash
go install github.com/pipectl/pipectl/cmd/pipectl@latest
```

Verify the installation:

```bash
pipectl --help
```

## Your first pipeline

Create a pipeline file:

```yaml
# greet.yaml
id: greet
input:
  format: json
steps:
  - normalize:
      fields:
        name: capitalize
  - select:
      fields: [name, email]
output:
  format: json
```

Create some input:

```bash
echo '[{"name":"alice smith","email":"ALICE@EXAMPLE.COM"},{"name":"bob jones","email":"BOB@EXAMPLE.COM"}]' > people.json
```

Run it:

```bash
pipectl run greet.yaml < people.json
```

Or use `--input` instead of stdin redirection:

```bash
pipectl run greet.yaml --input people.json
```

Output:

```json
[{"email":"ALICE@EXAMPLE.COM","name":"Alice Smith"},{"email":"BOB@EXAMPLE.COM","name":"Bob Jones"}]
```

## Validate without running

Use `--dry-run` to check a pipeline is valid and see the planned steps without executing anything:

```bash
pipectl run greet.yaml --dry-run
```

```
Pipeline: greet
Steps:
  1. normalize
  2. select
```

## Write output to a file

By default pipectl writes to `stdout`. Use `-o` to write to a file instead:

```bash
pipectl run greet.yaml -o output.json < people.json
```

## Enable verbose logging

Use `--verbose` to see per-step detail — record counts, field operations, sort results — written to stderr:

```bash
pipectl run greet.yaml --verbose < people.json
```

## Next steps

- Read [Core Concepts](./concepts) to understand how pipelines, steps, and payloads fit together
- Browse the [Step Reference](./steps/) to see all available steps
- Explore [Example Pipelines](./examples/) for real-world patterns
