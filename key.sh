#!/bin/bash
curl --no-progress-meter -H "Authorization: $API_KEY" \
-H'Content-Type: application/json' \
-d '{ "key": { 
        "algorithm": "RS256", 
        "name": "Invoice MVP Key - SHA-256 with RSA", 
        "length": 2048, 
        "issuer": "invoice.mvp",
        "subject": "CN=invoice.mvp"
}}' http://localhost:9011/api/key/generate/032f0a62-e111-4144-9f37-ca55112383ec | jq > key.json
