// Command run is the pipectl performance benchmark orchestrator: it builds
// pipectl fresh, generates synthetic data across a size x format matrix,
// runs the benchmark pipelines, and rewrites the results block inside
// website/docs/performance.md in place. See benchmark/README.md.
//
// Usage (from the repo root): go run ./benchmark/run
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pipectl/pipectl/benchmark/internal/sizefmt"
)

type sizeSpec struct {
	label string
	bytes int64
}

const (
	startMarker  = "<!-- BENCHMARK:START -->"
	endMarker    = "<!-- BENCHMARK:END -->"
	defaultSizes = "10MB,100MB,250MB,500MB,1GB"
	genSeed      = "42"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		sizesFlag   string
		formatsFlag string
		reps        int
		regenerate  bool
		outPath     string
	)
	flag.StringVar(&sizesFlag, "sizes", defaultSizes, "Comma-separated list of sizes to benchmark, e.g. 10MB,1GB")
	flag.StringVar(&formatsFlag, "formats", "csv,json,jsonl", "Comma-separated list of formats to benchmark")
	flag.IntVar(&reps, "reps", 5, "Repetitions per size x format; the median is reported")
	flag.BoolVar(&regenerate, "regenerate", false, "Regenerate data files even if they already exist")
	flag.StringVar(&outPath, "out", filepath.Join("website", "docs", "performance.md"), "File whose BENCHMARK marker block is rewritten with results")
	flag.Parse()

	if _, err := os.Stat("go.mod"); err != nil {
		return fmt.Errorf("must be run from the repo root (go.mod not found in current directory)")
	}

	sizes, err := parseSizes(sizesFlag)
	if err != nil {
		return err
	}
	formats := strings.Split(formatsFlag, ",")

	fmt.Fprintln(os.Stderr, "run: building pipectl...")
	pipectlBin, cleanup, err := buildPipectl()
	if err != nil {
		return err
	}
	defer cleanup()

	dataDir := filepath.Join("benchmark", "data")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return err
	}

	results := map[string][]sizeResult{}
	for _, format := range formats {
		pipelineFile := filepath.Join("benchmark", fmt.Sprintf("pipeline-%s.yaml", format))
		if _, err := os.Stat(pipelineFile); err != nil {
			return fmt.Errorf("no pipeline definition for format %q (expected %s)", format, pipelineFile)
		}

		var formatResults []sizeResult
		for _, size := range sizes {
			dataFile := filepath.Join(dataDir, fmt.Sprintf("%s-%s.%s", size.label, format, format))
			if regenerate || !fileExists(dataFile) {
				fmt.Fprintf(os.Stderr, "run: generating %s...\n", dataFile)
				if err := generate(size.label, format, dataFile); err != nil {
					return err
				}
			}

			fmt.Fprintf(os.Stderr, "run: benchmarking %s / %s (%d reps)...\n", format, size.label, reps)
			res, err := benchmark(pipectlBin, pipelineFile, dataFile, size, reps)
			if err != nil {
				return fmt.Errorf("%s/%s: %w", format, size.label, err)
			}
			formatResults = append(formatResults, res)
		}
		results[format] = formatResults
	}

	report := renderReport(formats, results)
	if err := rewriteMarkedBlock(outPath, report); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "run: updated %s\n", outPath)
	return nil
}

func parseSizes(s string) ([]sizeSpec, error) {
	var out []sizeSpec
	for _, label := range strings.Split(s, ",") {
		label = strings.TrimSpace(label)
		b, err := sizefmt.Parse(label)
		if err != nil {
			return nil, fmt.Errorf("--sizes: %w", err)
		}
		out = append(out, sizeSpec{label: label, bytes: b})
	}
	return out, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func buildPipectl() (bin string, cleanup func(), err error) {
	tmpDir, err := os.MkdirTemp("", "pipectl-bench")
	if err != nil {
		return "", nil, err
	}
	cleanup = func() { os.RemoveAll(tmpDir) }

	bin = filepath.Join(tmpDir, "pipectl")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/pipectl")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("build pipectl: %w", err)
	}
	return bin, cleanup, nil
}

func generate(sizeLabel, format, outPath string) error {
	cmd := exec.Command("go", "run", "./benchmark/gen",
		"-size", sizeLabel, "-format", format, "-out", outPath, "-seed", genSeed)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type sizeResult struct {
	size          sizeSpec
	medianElapsed time.Duration
	medianRSS     int64 // 0 if unavailable (non-Unix)
	inputRecords  int64
	timingTable   string
}

func benchmark(pipectlBin, pipelineFile, dataFile string, size sizeSpec, reps int) (sizeResult, error) {
	var durations []time.Duration
	var rss []int64
	var lastTiming string
	var inputRecords int64

	for i := 0; i < reps; i++ {
		cmd := exec.Command(pipectlBin, "run", pipelineFile,
			"-i", dataFile, "-o", os.DevNull, "--timing", "--max-input-size", "4GB")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		start := time.Now()
		err := cmd.Run()
		elapsed := time.Since(start)
		if err != nil {
			return sizeResult{}, fmt.Errorf("pipectl run failed: %w\n%s", err, stderr.String())
		}
		durations = append(durations, elapsed)

		if r, ok := peakRSSBytes(cmd.ProcessState); ok {
			rss = append(rss, r)
		}

		table, firstIn := parseTimingTable(stderr.String())
		lastTiming = table
		inputRecords = firstIn
	}

	return sizeResult{
		size:          size,
		medianElapsed: medianDuration(durations),
		medianRSS:     medianInt64(rss),
		inputRecords:  inputRecords,
		timingTable:   lastTiming,
	}, nil
}

// parseTimingTable extracts the "--timing" table from pipectl's stderr
// output and the IN record count of its first row (the total input record
// count for that run).
func parseTimingTable(stderr string) (table string, firstIn int64) {
	lines := strings.Split(stderr, "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "STEP") && strings.Contains(line, "DURATION") {
			var tableLines []string
			tableLines = append(tableLines, line)
			for _, l := range lines[i+1:] {
				if strings.TrimSpace(l) == "" {
					break
				}
				tableLines = append(tableLines, l)
			}
			table = strings.Join(tableLines, "\n")
			if len(tableLines) > 1 {
				fields := strings.Fields(tableLines[1])
				if len(fields) >= 2 {
					if n, err := strconv.ParseInt(fields[len(fields)-2], 10, 64); err == nil {
						firstIn = n
					}
				}
			}
			return table, firstIn
		}
	}
	return "", 0
}

func medianDuration(d []time.Duration) time.Duration {
	if len(d) == 0 {
		return 0
	}
	sorted := append([]time.Duration(nil), d...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	return sorted[len(sorted)/2]
}

func medianInt64(v []int64) int64 {
	if len(v) == 0 {
		return 0
	}
	sorted := append([]int64(nil), v...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	return sorted[len(sorted)/2]
}

func renderReport(formats []string, results map[string][]sizeResult) string {
	var b strings.Builder

	fmt.Fprintf(&b, "_Last generated by `go run ./benchmark/run` on %s/%s, %d cores, %s._\n\n",
		runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), time.Now().UTC().Format("2006-01-02"))

	for _, format := range formats {
		fmt.Fprintf(&b, "### %s\n\n", strings.ToUpper(format))
		fmt.Fprintln(&b, "| Size | Records | Median duration | Throughput | Peak RSS | Peak RSS / input size |")
		fmt.Fprintln(&b, "|---|---|---|---|---|---|")

		var largest sizeResult
		for _, r := range results[format] {
			mbPerSec := float64(r.size.bytes) / r.medianElapsed.Seconds() / (1 << 20)
			ratio := "n/a"
			rssStr := "n/a"
			if r.medianRSS > 0 {
				rssStr = sizefmt.Format(r.medianRSS)
				ratio = fmt.Sprintf("%.1fx", float64(r.medianRSS)/float64(r.size.bytes))
			}
			fmt.Fprintf(&b, "| %s | %d | %s | %.1f MB/s | %s | %s |\n",
				r.size.label, r.inputRecords, fmtDuration(r.medianElapsed), mbPerSec, rssStr, ratio)

			if r.size.bytes >= largest.size.bytes {
				largest = r
			}
		}
		b.WriteString("\n")

		if largest.timingTable != "" {
			fmt.Fprintf(&b, "<details>\n<summary>Per-step timing at %s (%s)</summary>\n\n```\n%s\n```\n\n</details>\n\n",
				largest.size.label, format, largest.timingTable)
		}
	}

	return strings.TrimRight(b.String(), "\n") + "\n"
}

func fmtDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func rewriteMarkedBlock(path, content string) error {
	existing, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	start := strings.Index(string(existing), startMarker)
	end := strings.Index(string(existing), endMarker)
	if start == -1 || end == -1 || end < start {
		return fmt.Errorf("%s: missing %s / %s markers — add them before running the benchmark", path, startMarker, endMarker)
	}

	before := string(existing)[:start+len(startMarker)]
	after := string(existing)[end:]
	updated := before + "\n\n" + content + "\n" + after

	return os.WriteFile(path, []byte(updated), 0o644)
}
