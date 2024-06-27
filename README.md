
# reqpretty 💅

Tired of plain, boring logs? Give your HTTP requests and responses a stylish makeover with \`reqpretty\`!

[Image/GIF of a plain log transforming into a colorful, structured log with reqpretty]

## Why reqpretty? ✨

* **Eye Candy:** Beautiful ASCII borders and emojis (✅/❌) for instant status checks.
* **Customizable:** Choose your colors, control the level of detail, and format the way you like it.
* **Framework-Agnostic:** Works with all your favorite Go web frameworks. Just wrap and go!
* **Structured Logging:** Powered by \`slog\` for efficient and insightful logging.

## Get Started 🚀

1. **Install:**

   ```bash
   go get github.com/1saifj/reqpretty
   ```

2. **Wrap & Log:**

   ```go
   import (
       "github.com/1saifj/reqpretty"
       "net/http"
   )
    
   func main() {
       // Create a reqpretty handler 
       http.Handle("/", reqpretty.DebugHandler(reqpretty.Options{}, yourHandler))

       // ... (your other server setup)
    }
   ```

## Configuration 🎨

```go
reqpretty.Config(reqpretty.Options{
    IncludeRequestHeaders:   true,
    IncludeRequestBody:      true,
    IncludeResponseHeaders:  true,
    IncludeResponseBody:     true,
    Colorer:                 &reqpretty.DefaultColorer{},
    EnableColor:             true, // Or false to disable colors
    SuccessEmoji: "✅", // Customize success emoji
    ErrorEmoji:   "❌", // Customize error emoji
})
```

## Example Output 📸

```
┌───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│  ⤴ REQUEST ⤴                                                                                                         │
│                                                                                                                       │
│ POST /api/users                                                                                                       │
│ Content-Type: application/json                                                                                        │
│                                                                                                                       │
│   {                                                                                                                   │
│     "name": "Alice"                                                                                                   │
│   }                                                                                                                   │
│                                                                                                                       │
│ ✅ RESPONSE [200/OK] [Time elapsed: 123 ms]⤵                                                                         │  
│                                                                                                                       │
│ Content-Type: [application/json]                                                                                      │
│                                                                                                                       │
│   {                                                                                                                   │
│     "message": "User created successfully"                                                                            │
│   }                                                                                                                   │
└───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

## Make It Your Own 🛠️

* **Custom Colorer:** Create your own \`Colorer\` implementation for unique color schemes.
* **Extend:** Contribute new features or formats – we welcome PRs!

## License 📄

MIT License - Go wild and be creative!

