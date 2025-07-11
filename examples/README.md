# reqpretty Examples

This directory contains examples and tests to verify that **reqpretty** is working correctly.

## 🧪 Testing Methods

### 1. **Unit Tests** (Basic verification)
Run the built-in tests:
```bash
cd ..
go test -v
```

### 2. **Live Server Testing** (Recommended)
Start the test server and make real HTTP requests:

**Terminal 1** - Start the server:
```bash
cd examples
go run test_server.go
```

**Terminal 2** - Run test commands:
```bash
# Manual curl commands
curl -X GET "http://localhost:8080/hello?test=true"

# Or run the automated test script
./test_commands.sh
```

### 3. **Configuration Testing**
Test different logging configurations:
```bash
cd examples
go run config_test.go
```

## 📋 What to Look For

When **reqpretty** is working correctly, you should see:

### ✅ **Request Logs** (with ⤴ REQUEST ⤴)
```
INFO ⤴ REQUEST ⤴ method=GET url=http://localhost:8080/hello?test=true headers=map[...] query_params=map[test:[true]]
```

### ✅ **Response Logs** (with emoji + RESPONSE ⤵)
```
INFO ✅ RESPONSE ⤵ status="200 OK" duration=1.2ms body={"message":"Hello, World!"}
```

### ✅ **Pretty JSON Formatting**
Request/response bodies should be nicely indented:
```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

### ✅ **Context Attributes**
If you add context values, they should appear in logs:
```
INFO ⤴ REQUEST ⤴ trace_id=trace-789 user_id=user-456 method=POST ...
```

### ✅ **Error Handling**
Error responses should show with error emoji:
```
INFO ❌ RESPONSE ⤵ status="500 Internal Server Error" duration=0.5ms
```

## 🔧 Available Test Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/hello` | GET | Simple JSON response with query params |
| `/users` | POST | Accepts JSON body, returns created user |
| `/error` | GET | Returns 500 error for testing error logs |

## 🎛️ Configuration Examples

The examples demonstrate different configuration options:

- **Full logging**: Everything enabled
- **Minimal logging**: Only basic request/response info
- **Headers only**: Just headers, no bodies
- **Context attributes**: Custom context values in logs

## 🚨 Troubleshooting

**No logs appearing?**
- Check that `reqpretty.Configure(logger)` is called
- Verify the handler is wrapped with `reqpretty.DebugHandler()`

**Logs look weird?**
- Make sure you're using the correct `Options` configuration
- Check that emojis are set (SuccessEmoji, ErrorEmoji)

**Server won't start?**
- Make sure port 8080 is available
- Run `go mod tidy` in the examples directory 