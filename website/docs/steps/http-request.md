# http-request

Sends the payload to an HTTP endpoint and continues unchanged. Useful for webhooks and audit logging.

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
- `proxy`, `headers`, and `timeout` can be set once for all `http-request`/`http-transform` steps via the pipeline-level [`defaults.http`](../concepts#defaults) instead of repeating them on every step. A step's own value always overrides the default; `headers` are merged, with the step's value winning on a key conflict.
