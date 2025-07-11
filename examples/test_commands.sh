#!/bin/bash

echo "🧪 Testing reqpretty middleware..."
echo "Make sure the test server is running: go run test_server.go"
echo ""

echo "1️⃣ Testing GET request..."
curl -X GET "http://localhost:8080/hello?name=test&debug=true" \
  -H "User-Agent: reqpretty-test" \
  -H "X-Request-ID: req-123"

echo -e "\n\n2️⃣ Testing POST request with JSON..."
curl -X POST "http://localhost:8080/users" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer fake-token" \
  -H "X-User-ID: user-456" \
  -d '{"name": "John Doe", "email": "john@example.com"}'

echo -e "\n\n3️⃣ Testing error response..."
curl -X GET "http://localhost:8080/error" \
  -H "User-Agent: reqpretty-test"

echo -e "\n\n4️⃣ Testing invalid JSON POST..."
curl -X POST "http://localhost:8080/users" \
  -H "Content-Type: application/json" \
  -d '{"invalid": json}'

echo -e "\n\n✅ Test completed! Check the server logs for beautiful output." 