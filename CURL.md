# CURL

## Interactions

GET /contacts

```sh
curl -i http://localhost:8080/contacts \
  -H 'cache-control: no-cache'

HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 46

[{"Firstname":"Andreas","Lastname":"Amstutz"}]
```

GET /contacts/{id}

```sh
curl -i http://localhost:8080/contacts/1

HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 44

{"Firstname":"Andreas","Lastname":"Amstutz"}
```

POST /contacts

```sh
curl -i -X POST http://localhost:8080/contacts \
  -H 'Content-Type: application/json' \
  -d '{
    "Firstname": "Meister",
    "Lastname": "Propper"
}'

HTTP/1.1 201 Created
Location: /contacts/2
...
```

DELETE /contacts/{id}

```sh
curl -i -X DELETE http://localhost:8080/contacts/1

HTTP/1.1 204 No Content
```

PUT /contacts/{id}

```sh
curl -i -X PUT http://localhost:8080/contacts/1 \
  -H 'Content-Type: application/json' \
  -d '{
    "Firstname": "Andreas",
    "Lastname": "Amstutz"
}'

HTTP/1.1 204 No Content
```
