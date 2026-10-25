// Command gen generates synthetic customer-like records (CSV, JSON, or
// JSONL) for the pipectl performance benchmark. See benchmark/README.md.
package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/pipectl/pipectl/benchmark/internal/sizefmt"
)

var plans = []string{"free", "pro", "enterprise"}

type record struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Country    string `json:"country,omitempty"`
	Plan       string `json:"plan"`
	CreditCard string `json:"credit_card"`
	SignupDate string `json:"signup_date"`
	Amount     string `json:"amount"`
	Status     string `json:"status"`
}

func (r record) csvRow() []string {
	return []string{r.ID, r.FirstName, r.LastName, r.Email, r.Country, r.Plan, r.CreditCard, r.SignupDate, r.Amount, r.Status}
}

var csvHeader = []string{"id", "first_name", "last_name", "email", "country", "plan", "credit_card", "signup_date", "amount", "status"}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "gen:", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		rowsFlag int64
		sizeFlag string
		format   string
		outPath  string
		dupRate  float64
		seed     int64
	)

	flag.Int64Var(&rowsFlag, "rows", 0, "Number of records to generate (mutually exclusive with -size)")
	flag.StringVar(&sizeFlag, "size", "", "Target output file size, e.g. 100MB, 1GB (mutually exclusive with -rows)")
	flag.StringVar(&format, "format", "csv", "Output format: csv, json, or jsonl")
	flag.StringVar(&outPath, "out", "", "Output file path (required)")
	flag.Float64Var(&dupRate, "dup-rate", 0.02, "Fraction of records that reuse a prior id (feeds the dedupe step)")
	flag.Int64Var(&seed, "seed", 42, "Random seed, for reproducible generation")
	flag.Parse()

	if outPath == "" {
		return fmt.Errorf("-out is required")
	}
	if format != "csv" && format != "json" && format != "jsonl" {
		return fmt.Errorf("-format must be csv, json, or jsonl, got %q", format)
	}
	if rowsFlag > 0 && sizeFlag != "" {
		return fmt.Errorf("-rows and -size are mutually exclusive")
	}
	var targetBytes int64
	if sizeFlag != "" {
		b, err := sizefmt.Parse(sizeFlag)
		if err != nil {
			return err
		}
		targetBytes = b
	} else if rowsFlag <= 0 {
		rowsFlag = 1000
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer f.Close()

	bw := bufio.NewWriterSize(f, 1<<20)
	cw := &countingWriter{w: bw}

	rng := rand.New(rand.NewSource(seed))
	gofakeit.Seed(seed)

	gen := &generator{rng: rng, dupRate: dupRate}

	switch format {
	case "csv":
		err = writeCSV(cw, gen, rowsFlag, targetBytes)
	case "json":
		err = writeJSON(cw, gen, rowsFlag, targetBytes)
	case "jsonl":
		err = writeJSONL(cw, gen, rowsFlag, targetBytes)
	}
	if err != nil {
		return err
	}

	if err := bw.Flush(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "gen: wrote %s (%s), %d records\n", outPath, sizefmt.Format(cw.n), gen.count)
	return nil
}

// countingWriter tracks total bytes written, so generation can stop once a
// target file size is reached without pre-estimating average record size.
type countingWriter struct {
	w io.Writer
	n int64
}

func (c *countingWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.n += int64(n)
	return n, err
}

// generator produces records with occasional missing/mixed-case fields and
// duplicate ids, so every step in the benchmark pipeline has real work to do.
type generator struct {
	rng       *rand.Rand
	dupRate   float64
	count     int64
	recentIDs []string
}

const recentIDCap = 10000

var signupRangeStart = time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
var signupRangeEnd = time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

func (g *generator) next() record {
	g.count++

	id := strconv.FormatInt(g.count, 10)
	if len(g.recentIDs) > 0 && g.rng.Float64() < g.dupRate {
		id = g.recentIDs[g.rng.Intn(len(g.recentIDs))]
	}
	g.recentIDs = append(g.recentIDs, id)
	if len(g.recentIDs) > recentIDCap {
		g.recentIDs = g.recentIDs[1:]
	}

	country := gofakeit.Country()
	if g.rng.Float64() < 0.05 {
		country = ""
	}

	status := "inactive"
	if g.rng.Float64() < 0.75 {
		status = "active"
	}

	return record{
		ID:         id,
		FirstName:  mixCase(g.rng, gofakeit.FirstName()),
		LastName:   mixCase(g.rng, gofakeit.LastName()),
		Email:      mixCase(g.rng, gofakeit.Email()),
		Country:    mixCase(g.rng, country),
		Plan:       plans[g.rng.Intn(len(plans))],
		CreditCard: gofakeit.CreditCardNumber(nil),
		SignupDate: gofakeit.DateRange(signupRangeStart, signupRangeEnd).Format("2006-01-02"),
		Amount:     strconv.FormatFloat(gofakeit.Price(1, 5000), 'f', 2, 64),
		Status:     status,
	}
}

// mixCase randomly upper/lower/leaves a string as-is, so the normalize step
// (capitalize/lower/upper) has real casing inconsistencies to fix.
func mixCase(rng *rand.Rand, s string) string {
	if s == "" {
		return s
	}
	switch rng.Intn(3) {
	case 0:
		return strings.ToUpper(s)
	case 1:
		return strings.ToLower(s)
	default:
		return s
	}
}

func done(rows int64, target int64, cw *countingWriter, i int64) bool {
	if target > 0 {
		return cw.n >= target
	}
	return i >= rows
}

func writeCSV(cw *countingWriter, gen *generator, rows int64, target int64) error {
	w := csv.NewWriter(cw)
	if err := w.Write(csvHeader); err != nil {
		return err
	}
	for i := int64(0); !done(rows, target, cw, i); i++ {
		if err := w.Write(gen.next().csvRow()); err != nil {
			return err
		}
		w.Flush()
		if err := w.Error(); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(cw *countingWriter, gen *generator, rows int64, target int64) error {
	if _, err := cw.Write([]byte("[\n")); err != nil {
		return err
	}
	for i := int64(0); !done(rows, target, cw, i); i++ {
		if i > 0 {
			if _, err := cw.Write([]byte(",\n")); err != nil {
				return err
			}
		}
		b, err := json.Marshal(gen.next())
		if err != nil {
			return err
		}
		if _, err := cw.Write(b); err != nil {
			return err
		}
	}
	_, err := cw.Write([]byte("\n]\n"))
	return err
}

func writeJSONL(cw *countingWriter, gen *generator, rows int64, target int64) error {
	for i := int64(0); !done(rows, target, cw, i); i++ {
		b, err := json.Marshal(gen.next())
		if err != nil {
			return err
		}
		if _, err := cw.Write(b); err != nil {
			return err
		}
		if _, err := cw.Write([]byte("\n")); err != nil {
			return err
		}
	}
	return nil
}
