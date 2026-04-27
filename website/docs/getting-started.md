# Getting Started

## Installation

### macOS (Homebrew)

```bash
brew tap pipectl/pipectl
brew install --cask pipectl
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

### Docker

Images are published to GitHub Container Registry for both `amd64` and `arm64`:

```bash
docker pull ghcr.io/pipectl/pipectl:latest
```

Run a pipeline with stdin:

```bash
echo '[...]' | docker run --rm -i ghcr.io/pipectl/pipectl:latest run pipeline.yaml
```

Run with local files mounted:

```bash
docker run --rm -i \
  -v $(pwd):/data \
  ghcr.io/pipectl/pipectl:latest run /data/pipeline.yaml --input /data/input.json
```

### Go install (from source)

```bash
go install github.com/pipectl/pipectl/cmd/pipectl@latest
```

Verify the installation:

```bash
pipectl --help
```

## Upgrading

### macOS (Homebrew)

```bash
brew upgrade --cask pipectl
```

### Windows

**Direct download:** download the new `pipectl_<version>_windows_amd64.zip` from the [Releases page](https://github.com/pipectl/pipectl/releases/latest), extract, and replace the existing binary in your `PATH`.

**Scoop:**

```powershell
scoop update pipectl
```

### Linux

Download the new `.deb` or `.rpm` from the [Releases page](https://github.com/pipectl/pipectl/releases/latest).

**Debian/Ubuntu:**
```bash
sudo dpkg -i pipectl_<version>_linux_amd64.deb
```

**Fedora/RHEL:**
```bash
sudo rpm -U pipectl_<version>_linux_amd64.rpm
```

### Docker

```bash
docker pull ghcr.io/pipectl/pipectl:latest
```

### Go install (from source)

```bash
go install github.com/pipectl/pipectl/cmd/pipectl@latest
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

## Substitute variables

Use `--var KEY=VALUE` to substitute `${VAR}` tokens in pipeline YAML before it is parsed. This lets you write reusable pipelines and supply environment-specific values at runtime:

```yaml
# pipeline.yaml
id: example
input:
  format: ${INPUT_FORMAT}
steps:
  - limit:
      count: ${LIMIT}
output:
  format: json
```

```bash
pipectl run pipeline.yaml --var INPUT_FORMAT=jsonl --var LIMIT=50 < data.jsonl
```

`--var` can be repeated as many times as needed. Any `${VAR}` token left unresolved after all substitutions are applied causes an error at startup.

## Next steps

- Read [Core Concepts](./concepts) to understand how pipelines, steps, and payloads fit together
- Browse the [Step Reference](./steps/) to see all available steps
- Explore [Example Pipelines](./examples/) for real-world patterns
