# respmask

`respmask` is a Go middleware for dynamically masking specific fields in JSON responses.

## Motivation

While working with APIs that deal with sensitive data, it becomes crucial to mask specific fields in the response without altering the actual data model. The aim of `respmask` is to provide a simple and efficient way to dynamically mask any field in the response based on the context of the request.

## Features

- **Dynamic Masking Rules**: Define masking functions on the fly based on request context.
- **Recursive Masking**: Traverse and mask nested JSON objects and arrays.
- **Extendable**: Easily extend with custom masking functions for different use cases.
- **Predefined Masking Rules**: Out-of-the-box support for common masking scenarios like emails, passwords, phone numbers, and credit card numbers.

## Installation

```bash
go get github.com/timakin/respmask
```

## Usage

### Basic Usage with Default Masking Rules

1. Use the provided default masking rules.

```go
import "github.com/timakin/respmask"

func dynamicKeysAndMaskingFuncs(r *http.Request) map[string]respmask.MaskingFunc {
    return map[string]respmask.MaskingFunc{
        "email":    respmask.DefaultMaskingRules[respmask.EmailMasking],
        "password": respmask.DefaultMaskingRules[respmask.PasswordMasking],
        // ... use other default masking functions as needed ...
    }
}
```

2. Set up the middleware with your HTTP server.

```go
func main() {
	http.Handle("/api/data", respmask.NewMaskingMiddleware(dynamicKeysAndMaskingFuncs, http.HandlerFunc(handleData)))
	http.ListenAndServe(":8080", nil)
}
```

3. Watch specific fields in your JSON responses get masked using the predefined rules!

### Custom Masking Functions

You can easily extend the functionality by defining your own masking functions.

```go
func customMaskFunc(input string) string {
    // Your custom masking logic
    return "masked_value"
}
```

Then, use your custom function in the dynamic keys function:

```go
func dynamicKeysAndMaskingFuncs(r *http.Request) map[string]respmask.MaskingFunc {
    return map[string]respmask.MaskingFunc{
        "custom_field": customMaskFunc,
        // ... other fields and functions ...
    }
}
```

## Testing

Refer to the provided test cases in the package to see how to effectively test the masking functionality.

## Contributing

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Make your changes.
4. Push to your fork and submit a pull request!
