# cast

Converts field values to a specified type.

**Supported formats:** `json` `jsonl`

## Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `fields` | map | Yes | Map of field path to cast configuration |

Each field entry supports:

| Key | Type | Required | Description |
|-----|------|----------|-------------|
| `type` | string | Yes | Target type: `int`, `float`, `bool`, `string`, or `time` |
| `format` | string | No | Date/time parse format (Go layout string). Only valid for `type: time`. Defaults to RFC 3339. |
| `true_values` | list | No | Strings to treat as `true`. Only valid for `type: bool`. |
| `false_values` | list | No | Strings to treat as `false`. Only valid for `type: bool`. |

## Examples

Basic type conversion:

```yaml
- cast:
    fields:
      age:
        type: int
      price:
        type: float
      active:
        type: bool
```

Custom bool values:

```yaml
- cast:
    fields:
      active:
        type: bool
        true_values: [yes, "1", enabled]
        false_values: [no, "0", disabled]
```

Parse a date string:

```yaml
- cast:
    fields:
      created_at:
        type: time
        format: "2006-01-02"
```

Nested fields:

```yaml
- cast:
    fields:
      user.age:
        type: int
      tags[0]:
        type: string
```

## Notes

- Field paths support dot notation (`user.address`) and array indexing (`tags[0]`).
- Default `true` values for `bool`: `true`, `t`, `1`, `yes`, `y`, `on`.
- Default `false` values for `bool`: `false`, `f`, `0`, `no`, `n`, `off`.
- `true_values` and `false_values` must not overlap.
- For `type: time`, the `format` string uses [Go's time layout](https://pkg.go.dev/time#Layout). The reference time is `Mon Jan 2 15:04:05 MST 2006`.
- Array fields are cast element-by-element.
