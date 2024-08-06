# reqpretty


<img src="https://i.imgur.com/4PxgK9B.png" width="350" height="300">

`reqpretty` is a Go middleware package designed to beautify and log HTTP requests and responses in a structured and readable format. It provides detailed logging of request and response headers, bodies, query parameters, and more, with customization options to suit your needs.

## Features

- **Request Logging**: Log HTTP method, URL, headers, query parameters, and body.
- **Response Logging**: Log HTTP status code, headers, and body.
- **Customization**: Configure what details to log, including request and response headers, bodies, and query parameters.
- **Context Attributes**: Extract and log specific context attributes.
- **Colorized Output**: Optionally colorize log output for better readability.

## Installation

To install the package, run:

```sh
go get github.com/1saif/reqpretty
```

## Usage

Here's a simple example of how to use `reqpretty` in your Go application.

### Setup Logger

First, configure the logger:

```go
package main

import (
"log/slog"
"github.com/1saifj/reqpretty"
)

func main() {
logger := &reqpretty.Logger{}
reqpretty.Configure(logger)
// Your application code
}
```

### Middleware Example

Next, use the `reqpretty` middleware in your HTTP server:

```go
package main

import (
"net/http"
"github.com/1saifj/reqpretty"
)

func main() {
opts := reqpretty.Options{
IncludeRequest:            true,
IncludeRequestHeaders:     true,
IncludeRequestQueryParams: true,
IncludeRequestBody:        true,
IncludeResponse:           true,
IncludeResponseHeaders:    true,
IncludeResponseBody:       true,
ContextAttributes:         []string{"request_id", "user_id"},
}

    mux := http.NewServeMux()
    mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, world!"))
    })

    loggedMux := reqpretty.DebugHandler(opts, mux)

    http.ListenAndServe(":8080", loggedMux)
}
```

## Configuration

### Options

The `Options` struct allows you to customize what details are logged:

- **IncludeRequest**: Log the request details (default: `false`).
- **IncludeRequestHeaders**: Log request headers (default: `false`).
- **IncludeRequestQueryParams**: Log request query parameters (default: `false`).
- **IncludeRequestBody**: Log request body (default: `false`).
- **IncludeResponse**: Log the response details (default: `false`).
- **IncludeResponseHeaders**: Log response headers (default: `false`).
- **IncludeResponseBody**: Log response body (default: `false`).
- **ContextAttributes**: List of context attributes to log (default: `nil`).

### Logger

The `Logger` struct is used to configure the logger:

- **clone()**: Create a copy of the logger.

## Example: Custom Logger

You can customize the logger further by implementing your own `slog.Handler`:

```go
package main

import (
"context"
"log/slog"
"github.com/1saifj/reqpretty"
)

type CustomHandler struct {
handler slog.Handler
}

func (h CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
return h.handler.Enabled(ctx, level)
}

func (h CustomHandler) Handle(ctx context.Context, record slog.Record) error {
// Custom log handling
return h.handler.Handle(ctx, record)
}

func (h CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
return CustomHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h CustomHandler) WithGroup(name string) slog.Handler {
return CustomHandler{handler: h.handler.WithGroup(name)}
}

func main() {
logger := &reqpretty.Logger{}
reqpretty.Configure(logger)
customHandler := CustomHandler{handler: slog.Default().Handler()}
slog.SetDefault(slog.New(customHandler))

    // Your application code
}
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License.
