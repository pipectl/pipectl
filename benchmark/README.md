# pipectl performance benchmark

Measures pipectl's throughput and peak memory use on a realistic 10-step pipeline,
across CSV, JSON, and JSONL, at file sizes from 10MB to 1GB. Results are published
at [`website/docs/performance.md`](../website/docs/performance.md).

## Why 1GB, not bigger

pipectl reads the whole payload into memory — a `[]byte` buffer plus a fully
parsed in-memory record structure — rather than streaming. There's no way
around this for a pipeline that includes `sort`, `dedupe`, or a count-based
`assert`: those steps need to see every record before they can do anything,
so at least one stage forces full materialization regardless of engine
architecture. 1GB is a deliberately round ceiling that reflects pipectl's
actual target use case (small-to-medium batch/ETL jobs), not a hard technical
limit — see the "Reading the memory numbers" section on the results page for
what the peak-RSS multiplier means for sizing `--max-input-size` on your own
machine.

## Reproducing the benchmark

From the repo root:

```sh
go run ./benchmark/run
```

This does the whole thing end-to-end:

1. Builds `pipectl` fresh from the current working tree (not a stale binary)
2. Generates synthetic data at each size x format combination into
   `benchmark/data/` (gitignored — regenerated on demand, never committed).
   Existing files are reused on subsequent runs unless `-regenerate` is
   passed, so iterating is fast.
3. Runs the matching benchmark pipeline (`pipeline-csv.yaml` /
   `pipeline-json.yaml` / `pipeline-jsonl.yaml`) 5 times per size x format,
   capturing wall-clock time and peak RSS (via the OS's own process
   accounting — `syscall.Rusage`, not the external `time` binary) for each
   run, and reports the median
4. Rewrites the `<!-- BENCHMARK:START -->` / `<!-- BENCHMARK:END -->` block
   inside `website/docs/performance.md` in place — the surrounding
   methodology and positioning prose is hand-written and untouched by reruns

Flags (all optional):

| Flag | Default | Meaning |
|---|---|---|
| `-sizes` | `10MB,100MB,250MB,500MB,1GB` | Comma-separated sizes to benchmark |
| `-formats` | `csv,json,jsonl` | Comma-separated formats to benchmark |
| `-reps` | `5` | Repetitions per size x format (median is reported) |
| `-regenerate` | `false` | Regenerate data files even if they already exist |
| `-out` | `website/docs/performance.md` | File whose marker block gets rewritten — point this at a scratch file for a dry run |

Re-running this later (e.g. after a stretch of new features) regenerates
just the results block — `git diff website/docs/performance.md` shows
exactly what changed before you commit and redeploy the site.

A quick smoke test before committing to the full matrix:

```sh
go run ./benchmark/run -sizes 10MB -formats csv -reps 1 -out /tmp/perf.md
```

## The pipeline

`assert -> default -> normalize -> dedupe -> filter -> redact -> select ->
rename -> sort -> count`, identical across all three formats so the
comparison isolates format overhead rather than step selection. `cast`,
`convert`, and `validate-json` are JSON/JSONL-only steps (see
`website/docs/steps/*.md`) and are deliberately excluded so CSV doesn't
carry extra format-conversion cost the other two don't. `http-request` /
`http-transform` are excluded as I/O-bound — this benchmark measures
pipeline throughput, not network latency.

A real CSV-sourced pipeline that needs type casting would add a `convert`
step (csv -> json/jsonl) before `cast`; that cost isn't captured here since
it would make the CSV number incomparable to JSON/JSONL.

## The data generator

`go run ./benchmark/gen` produces synthetic customer-like records (id,
name, email, country, plan, credit card, signup date, amount, status) using
[gofakeit](https://github.com/brianvoe/gofakeit) for realistic values, with
deliberate imperfections so every pipeline step has real work to do:

- mixed-case names/emails/country (feeds `normalize`)
- ~5% missing country (feeds `default`)
- ~2% duplicate ids (feeds `dedupe`)
- ~25% inactive status (feeds `filter`)

```sh
go run ./benchmark/gen -size 100MB -format json -out benchmark/data/100mb.json
```

Use `-rows N` instead of `-size` for an exact record count, and `-seed N`
for reproducible generation (the runner always uses a fixed seed).
