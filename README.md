# 🎨 reqpretty

<div align="center">

<img src="logo.png" width="350" height="300">

**A beautiful Go middleware for HTTP request/response logging** ✨

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.19-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/1saif/reqpretty)](https://goreportcard.com/report/github.com/1saif/reqpretty)

</div>

---

`reqpretty` is a Go middleware package designed to **beautify and log HTTP requests and responses** in a structured and readable format. It provides detailed logging of request and response headers, bodies, query parameters, and more, with customization options to suit your needs.

## ✨ Features

| Feature | Description |
|---------|-------------|
| 📥 **Request Logging** | Log HTTP method, URL, headers, query parameters, and body |
| 📤 **Response Logging** | Log HTTP status code, headers, and body |
| ⚙️ **Customization** | Configure what details to log, including request and response headers, bodies, and query parameters |
| 🔍 **Context Attributes** | Extract and log specific context attributes |
| 🌈 **Colorized Output** | Optionally colorize log output for better readability |

## 🚀 Installation

```bash
go get github.com/1saif/reqpretty
```

## 📖 Usage

Here's a simple example of how to use `reqpretty` in your Go application.

### 🔧 Setup Logger

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

### 🛠️ Middleware Example

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

## ⚙️ Configuration

### 📋 Options

The `Options` struct allows you to customize what details are logged:

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `IncludeRequest` | `bool` | `false` | Log the request details |
| `IncludeRequestHeaders` | `bool` | `false` | Log request headers |
| `IncludeRequestQueryParams` | `bool` | `false` | Log request query parameters |
| `IncludeRequestBody` | `bool` | `false` | Log request body |
| `IncludeResponse` | `bool` | `false` | Log the response details |
| `IncludeResponseHeaders` | `bool` | `false` | Log response headers |
| `IncludeResponseBody` | `bool` | `false` | Log response body |
| `ContextAttributes` | `[]string` | `nil` | List of context attributes to log |

### 🔧 Logger

The `Logger` struct is used to configure the logger:

- **`clone()`**: Create a copy of the logger

## 🎯 Example: Custom Logger

You can customize the logger further by implementing your own `slog.Handler`:

<details>
<summary>Click to expand custom logger example</summary>

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

</details>

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Made with ❤️ by the 1saifj**

⭐ **Star this repo if you find it helpful!** ⭐

</div>
