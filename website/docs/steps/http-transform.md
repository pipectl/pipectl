# http-transform

Sends the payload to an HTTP endpoint and replaces it with the response. Useful for enrichment and validation.

**Supported formats:** `json` `jsonl` `csv`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `url` | string | Yes | Target URL |
| `method` | string | Yes | HTTP method: `POST`, `PUT`, `PATCH`, or `DELETE` |
| `headers` | map | No | Request headers as key/value strings |
| `timeout` | integer | No | Request timeout in seconds. Must be between 1 and 300. Defaults to 60. |
| `proxy` | string | No | Proxy URL |
| `expect-format` | string | No | Expected response format: `json`, `jsonl`, or `csv`. Defaults to `json`. |

## Example

```yaml
- http-transform:
    url: https://api.example.com/enrich
    method: POST
    timeout: 30
    headers:
      Authorization: Bearer ${API_TOKEN}
      Content-Type: application/json
    expect-format: json
```

## Notes

- Only HTTP `200 OK` responses are accepted. Any other status code fails the pipeline.
- The response `Content-Type` must match `expect-format`. Set it explicitly if the endpoint requires it.
- For `POST`, `PUT`, `PATCH`, and `DELETE`, the current payload is sent as the request body.
- For JSONL payloads, the step sends `application/x-ndjson` as the `Content-Type` unless you override it in `headers`.
- For CSV payloads, the step sends `text/csv` as the `Content-Type` unless you override it in `headers`.
- For JSON payloads, the step sends `application/json` as the `Content-Type` unless you override it in `headers`.
- Environment variables in header values (e.g. `${API_TOKEN}`) are not automatically expanded. Use your shell or a secrets manager to inject values before running the pipeline.
- `proxy`, `headers`, and `timeout` can be set once for all `http-request`/`http-transform` steps via the pipeline-level [`defaults.http`](../concepts#defaults) instead of repeating them on every step. A step's own value always overrides the default; `headers` are merged, with the step's value winning on a key conflict.
