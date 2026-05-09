# http-request

Sends the current payload to an HTTP endpoint and continues the pipeline with the **same payload unchanged**. Useful for webhooks, notifications, audit logging, or any fire-and-continue side effect mid-pipeline.

**Supported formats:** `json` `jsonl` `csv`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `url` | string | Yes | Target URL |
| `method` | string | Yes | HTTP method: `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD`, or `OPTIONS` |
| `headers` | map | No | Request headers as key/value strings |
| `timeout` | integer | No | Request timeout in seconds. Must be between 1 and 300. Defaults to 60. |
| `proxy` | string | No | Proxy URL |

## Example

```yaml
- http-request:
    url: https://hooks.example.com/pipeline-event
    method: POST
    timeout: 10
    headers:
      Authorization: Bearer ${WEBHOOK_TOKEN}
      Content-Type: application/json
```

## Notes

- Any `2xx` response is accepted. Non-2xx status codes fail the pipeline.
- The response body is discarded. The pipeline payload is not modified.
- For `POST`, `PUT`, `PATCH`, and `DELETE`, the current payload is sent as the request body.
- For JSONL payloads, the step sends `application/x-ndjson` as the `Content-Type` unless you override it in `headers`.
- For CSV payloads, the step sends `text/csv` as the `Content-Type` unless you override it in `headers`.
- For JSON payloads, the step sends `application/json` as the `Content-Type` unless you override it in `headers`.
- Environment variables in header values (e.g. `${WEBHOOK_TOKEN}`) are not automatically expanded. Use your shell or a secrets manager to inject values before running the pipeline.
