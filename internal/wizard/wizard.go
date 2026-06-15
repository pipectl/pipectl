package wizard

import (
	"github.com/charmbracelet/huh"
)

// Run presents the interactive wizard and returns the collected answers.
func Run() (Result, error) {
	r := Result{
		ID:         "my-pipeline",
		OutputFile: "pipeline.yaml",
	}

	form := huh.NewForm(
		// Welcome screen.
		huh.NewGroup(
			huh.NewNote().
				Title("pipectl init").
				Description(
					"This wizard generates a starter pipeline YAML file.\n\n"+
						"You'll answer a few questions about your data and which processing\n"+
						"steps you need. A YAML file will then be written with sensible\n"+
						"placeholder values and inline comments — ready for you to customise.\n\n"+
						"Navigate with arrow keys. Press enter to confirm each selection.\n"+
						"Press ctrl+c at any time to exit without saving.",
				).
				Next(true).
				NextLabel("Get started →"),
		),

		// Pipeline identity and input.
		huh.NewGroup(
			huh.NewInput().
				Title("Pipeline ID").
				Description("A short name used in logs and error messages.").
				Placeholder("my-pipeline").
				Value(&r.ID),

			huh.NewSelect[string]().
				Title("Input format").
				Description("The format of the data this pipeline will read.").
				Options(
					huh.NewOption("JSON   — array of objects, or a single object", "json"),
					huh.NewOption("JSONL  — one JSON object per line (newline-delimited)", "jsonl"),
					huh.NewOption("CSV    — comma-separated values with a header row", "csv"),
				).
				Value(&r.InputFormat),
		),

		// Step selection.
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Steps").
				Description(
					"Choose the processing steps for your pipeline.\n"+
						"Steps run top-to-bottom in the order shown below.\n"+
						"You can reorder or remove them in the generated YAML.\n"+
						"Use space to select, ↑↓ to move, / to filter by name.",
				).
				Options(
					// Transformation
					huh.NewOption("select         — keep specific fields, drop everything else", "select"),
					huh.NewOption("rename         — rename fields", "rename"),
					huh.NewOption("cast           — convert field values to a different type (int, float, bool…)", "cast"),
					huh.NewOption("default        — fill missing or empty fields with a default value", "default"),
					huh.NewOption("normalize      — apply text transformations (trim, lower, upper, capitalize…)", "normalize"),
					huh.NewOption("redact         — mask or hash sensitive field values", "redact"),
					huh.NewOption("convert        — change the payload format mid-pipeline (e.g. JSON → CSV)", "convert"),
					// Filtering & sorting
					huh.NewOption("filter         — keep only records that match a condition", "filter"),
					huh.NewOption("dedupe         — remove duplicate records based on one or more fields", "dedupe"),
					huh.NewOption("sort           — sort records by a field (ascending or descending)", "sort"),
					huh.NewOption("limit          — truncate to the first N records", "limit"),
					// Validation
					huh.NewOption("validate-json  — validate records against a JSON Schema file", "validate-json"),
					huh.NewOption("assert         — fail the pipeline if record counts are out of range", "assert"),
					// Diagnostics
					huh.NewOption("count          — print record count to stderr (payload passes through)", "count"),
					huh.NewOption("log            — print a message and optional sample records to stderr", "log"),
					// HTTP
					huh.NewOption("http-request   — send payload to an HTTP endpoint and continue (fire-and-forget)", "http-request"),
					huh.NewOption("http-transform — send payload to an HTTP endpoint and replace it with the response", "http-transform"),
				).
				Filterable(true).
				Value(&r.Steps),
		),

		// Output settings.
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Output format").
				Description("The format this pipeline will produce.").
				Options(
					huh.NewOption("JSON   — array of objects, or a single object", "json"),
					huh.NewOption("JSONL  — one JSON object per line (newline-delimited)", "jsonl"),
					huh.NewOption("CSV    — comma-separated values with a header row", "csv"),
				).
				Value(&r.OutputFormat),

			huh.NewInput().
				Title("Output file").
				Description("Where to save the generated YAML. Leave empty to print to stdout.").
				Placeholder("pipeline.yaml").
				Value(&r.OutputFile),
		),
	)

	if err := form.Run(); err != nil {
		return Result{}, err
	}

	return r, nil
}
