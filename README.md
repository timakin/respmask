![test workflow](https://github.com/timakin/respmask/actions/workflows/test.yml/badge.svg)

# respmask

`respmask` is a Go middleware designed to dynamically mask specific fields in JSON responses, catering to the diverse needs of modern APIs dealing with sensitive data.

## Motivation

In the world of APIs, especially those that handle sensitive or private data, it's imperative to mask certain response fields without altering the actual data model. `respmask` offers a seamless solution, allowing developers to apply dynamic masking depending on the context of the request.

## Features

- **Dynamic Masking Rules**: Define and adapt masking functions based on request context.
- **Masking Modes**: Switch between `ExactMode` and `RecursiveMode` for granular control over response data masking.
- **Recursive Masking**: Efficiently traverse and mask nested JSON objects and arrays.
- **Extendability**: Craft custom masking functions tailored to specific use cases.
- **Predefined Masking Rules**: Benefit from built-in rules for typical masking scenarios, including emails, passwords, phone numbers, and credit card numbers.

## Installation

```bash
go get github.com/timakin/respmask
```

## Usage

### Basic Usage with Default Masking Rules

1. Leverage the in-built default masking rules.

```go
import "github.com/timakin/respmask"

func maskingConfig(r *http.Request) (map[string]respmask.MaskingFunc, respmask.MaskingMode) {
    return map[string]respmask.MaskingFunc{
        "email":    respmask.DefaultMaskingRules[respmask.EmailMasking],
        "password": respmask.DefaultMaskingRules[respmask.PasswordMasking],
        // ... use other default masking functions as needed ...
    }, respmask.RecursiveMode
}
```

2. Integrate the middleware with your HTTP server and choose a masking mode.

```go
func main() {
    http.Handle("/api/data", respmask.NewMaskingMiddleware(maskingConfig, http.HandlerFunc(handleData)))
    http.ListenAndServe(":8080", nil)
}
```

3. Observe specified fields in your JSON responses being masked according to predefined rules!

### Masking Modes

`respmask` offers two distinct masking modes:

- **ExactMode**: Masks keys strictly based on the hierarchy defined in configuration.
- **RecursiveMode**: Masks keys throughout the JSON structure, regardless of their depth.

For instance, given a masking function that targets the key "email" with `ExactMode` selected, only the top-level "email" key will be masked. On the other hand, with `RecursiveMode`, "email" keys across all nesting levels will be masked.

### Crafting Custom Masking Functions

Design your custom masking functions with ease.

```go
func customMaskFunc(input string) string {
    // Your unique masking logic here
    return "masked_value"
}
```

Subsequently, employ your custom function in the dynamic keys definition:

```go
func maskingConfig(r *http.Request) (map[string]respmask.MaskingFunc, respmask.MaskingMode) {
    return map[string]respmask.MaskingFunc{
        "custom_field": customMaskFunc,
        // ... other fields and functions ...
    }, respmask.RecursiveMode
}
```

## Testing

Refer to the in-package test cases for insights on how to effectively test the masking functionality.

## Contributing

1. Fork the repository.
2. Initiate a new branch for your feature or bugfix.
3. Implement your changes.
4. Push to your fork and open a pull request!
