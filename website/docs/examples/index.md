# Example Pipelines

Real-world pipeline patterns showing how steps compose to solve common data problems.

| Example | Input | Output | What it shows |
|---------|-------|--------|---------------|
| [Customer Signup](./customer-signup) | JSON | JSON | Validate, normalize, redact, enrich via HTTP |
| [CSV Intake](./csv-intake) | CSV | JSONL | Clean and filter a partner CSV export |
| [Audit Export](./audit-export) | JSONL | JSONL | SHA-256 hash PII fields for safe export |
| [Service-to-Service](./service-to-service) | JSON | CSV | POST to an HTTP service, process CSV response |

All examples include runnable YAML and sample input data. The `examples/` directory in the repository contains these pipelines with real input files.
