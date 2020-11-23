# CURL

## Interactions

---

## POST /activities

Activities:

```sh
curl -i https://127.0.0.1:8443/activities \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiR28gSW52b2ljZXIiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiZjhjMzlhMzEtOWNlZC00NzYxLThhMzMtYjljNjI4YTY3NTEwIn0.WI6cRXYnYqUAV6qqNtf4B8PdGMgKuHqENQP5N_iCZL8' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Programming" 
}'

# response             
HTTP/1.1 201 Created
Content-Type: application/json
Location: /activities/1
```

## GET /activities

```sh
curl -i https://127.0.0.1:8443/activities \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiR28gSW52b2ljZXIiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiZjhjMzlhMzEtOWNlZC00NzYxLThhMzMtYjljNjI4YTY3NTEwIn0.WI6cRXYnYqUAV6qqNtf4B8PdGMgKuHqENQP5N_iCZL8' \
  -H 'Content-Type: application/json'

# response             
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 47

[{"id":1,"name":"Programming","userId":"f8c39a31-9ced-4761-8a33-b9c628a67510"}]
```

---

## POST /customers

```sh
curl -i https://127.0.0.1:8443/customers \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiR28gSW52b2ljZXIiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiZjhjMzlhMzEtOWNlZC00NzYxLThhMzMtYjljNjI4YTY3NTEwIn0.WI6cRXYnYqUAV6qqNtf4B8PdGMgKuHqENQP5N_iCZL8' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "3skills" 
}'

# response
HTTP/1.1 201 Created
Location: /customers/1
```

---

## POST /customers/{customerId}/invoices

Invoices:

> The API currently only knows of the customer with ID 1, that is used in the URI for creating an invoice.

```sh
curl -i https://127.0.0.1:8443/customers/1/invoices \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiR28gSW52b2ljZXIiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiZjhjMzlhMzEtOWNlZC00NzYxLThhMzMtYjljNjI4YTY3NTEwIn0.WI6cRXYnYqUAV6qqNtf4B8PdGMgKuHqENQP5N_iCZL8' \
  -H 'Content-Type: application/json' \
  -d '{
    "month": 9,
    "year": 2020
}'
```

```sh
# response
HTTP/1.1 201 Created
Location: /customers/1/invoices/1
```

---

## POST /customers/{customerId}/projects

Projects:

```sh
curl -i https://127.0.0.1:8443/customers/1/projects \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiR28gSW52b2ljZXIiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiZjhjMzlhMzEtOWNlZC00NzYxLThhMzMtYjljNjI4YTY3NTEwIn0.WI6cRXYnYqUAV6qqNtf4B8PdGMgKuHqENQP5N_iCZL8' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Instantfoo.com"
}'

# response
HTTP/1.1 201 Created
Location: /customers/1/projects/1
```

---

## POST /customers/{customerId}/projects/{projectId}/rates

Rates:

```sh
curl -i https://127.0.0.1:8443/customers/1/projects/1/rates \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiR28gSW52b2ljZXIiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiZjhjMzlhMzEtOWNlZC00NzYxLThhMzMtYjljNjI4YTY3NTEwIn0.WI6cRXYnYqUAV6qqNtf4B8PdGMgKuHqENQP5N_iCZL8' \
  -H 'Content-Type: application/json' \
  -d '{
    "activityId": 1,
    "price": 67.80
}'

# response
HTTP/1.1 201 Created
Location: /customers/1/projects/1/rates/activity/1
```

---

## POST /customers/{customerId}/invoices/{invoiceId}/bookings

### Add bookings

> A new `Booking` references an `Invoice` with `invoiceId` variable in the URL.

The customer currently only has one project `Instantfoo.com` with ID 1.
The list of activities contains the activity `Programming` with ID 1.

The activity is now getting booked on this project:

```sh
# customer   : 1 (3skills)
# invoice    : 2 (2020-09)
# projectId  : 1 (Instantfoo.com)
# activityId : 1 (Programming)

curl -i https://127.0.0.1:8443/book/1 \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiR28gSW52b2ljZXIiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiZjhjMzlhMzEtOWNlZC00NzYxLThhMzMtYjljNjI4YTY3NTEwIn0.WI6cRXYnYqUAV6qqNtf4B8PdGMgKuHqENQP5N_iCZL8' \
  -H 'Content-Type: application/json' \
  -d '{
    "day": 31,
    "hours": 2.5,
    "description": "Front: bugfix #6789",
    "projectId": 1,
    "activityId": 1
}'

# response
HTTP/1.1 201 Created
Location: /customers/1/invoices/1/bookings/1
```

> The `Booking` of type `Programming` was succesfully created for the `Project` *instantfoo.com* on `Invoice` 1.

---

## DELETE /customers/{customerId}/invoices/{invoiceId}/bookings/{bookingId}

### Delete bookings

Currently it is not possible to correct bookings but required in the scope of this MVP.

The simple solution here is to just delete the booking and create a new one.

The error handling while deleting a non-existing booking will be ignored using a nil-operation, a function that does nothing.

```sh
curl -i -X DELETE \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiR28gSW52b2ljZXIiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiZjhjMzlhMzEtOWNlZC00NzYxLThhMzMtYjljNjI4YTY3NTEwIn0.WI6cRXYnYqUAV6qqNtf4B8PdGMgKuHqENQP5N_iCZL8' \
https://127.0.0.1:8443/customers/1/invoices/1/bookings/1

# response
HTTP/1.1 204 No Content
```

---

## PUT /customers/{customerId}/invoices/{invoiceId}

### Finalize an invoice

The handling of invoice finalization is implemented using a PUT-Request as an existing resource is getting updated.

```sh
curl -i -X PUT https://127.0.0.1:8443/customers/1/invoices/1 \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiR28gSW52b2ljZXIiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiZjhjMzlhMzEtOWNlZC00NzYxLThhMzMtYjljNjI4YTY3NTEwIn0.WI6cRXYnYqUAV6qqNtf4B8PdGMgKuHqENQP5N_iCZL8' \
  -H 'Content-Type: application/json' \
  -d '{
    "month": 9,
    "year": 2020,
    "status": "ready for aggregation"
}'

# response
HTTP/1.1 204 No Content
```

The `UpdateInvoice` call on the repository implementation saves the now aggregated invoice.

---

## GET /customers/{customerId}/invoices/{invoiceId}

### Retrieve an Invoice

```sh
curl -s https://127.0.0.1:8443/customers/1/invoices/1 \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiR28gSW52b2ljZXIiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiZjhjMzlhMzEtOWNlZC00NzYxLThhMzMtYjljNjI4YTY3NTEwIn0.WI6cRXYnYqUAV6qqNtf4B8PdGMgKuHqENQP5N_iCZL8' \
  -H 'Accept: application/json' | jq

# response
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 164


{
  "id": 1,
  "month": 9,
  "year": 2020,
  "status": "payment expected",
  "customerId": 1,
  "positions": {
    "1": {
      "Programming": {
        "Hours": 2.5,
        "Price": 169.5
      }
    }
  }
  "_links": {
    "bookings": {
      "href": "/invoice/1/bookings"
    },
    "payment": {
      "href": "/payment/1"
    },
    "self": {
      "href": "/invoice/1"
    }
  }
}
```

### Retrieve an Invoice with booking details

```sh
curl -s https://127.0.0.1:8443/customers/1/invoices/1?expand=bookings \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiR28gSW52b2ljZXIiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiZjhjMzlhMzEtOWNlZC00NzYxLThhMzMtYjljNjI4YTY3NTEwIn0.WI6cRXYnYqUAV6qqNtf4B8PdGMgKuHqENQP5N_iCZL8' \
  -H 'Accept: application/json' | jq

# response
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 271

{
  "id": 1,
  "month": 9,
  "year": 2020,
  "status": "payment expected",
  "customerId": 1,
  "positions": {
    "1": {
      "Programming": {
        "Hours": 2.5,
        "Price": 169.5
      }
    }
  },
  "bookings": [
    {
      "id": 1,
      "day": 31,
      "hours": 2.5,
      "description": "Front: bugfix #6789",
      "invoiceId": 1,
      "projectId": 1,
      "activityId": 1
    }
  ],
}
```

----

## Basic Auth

Initial request:

```sh
curl -i https://127.0.0.1:8443/activities

HTTP/1.1 401 Unauthorized
Www-Authenticate: Basic realm="invoice.mvp"
```

Follow up request using Basic Auth credentials:

```sh
curl -i --user go:time https://127.0.0.1:8443/activities
HTTP/1.1 200 OK
```

## Digest Auth

Initial request:

```sh
curl -i https://127.0.0.1:8443/customers \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "3skills" 
}'

HTTP/1.1 401 Unauthorized
Www-Authenticate: Digest realm="invoice.mvp", nonce="UAZs1dp3wX5BtXEpoCXKO2lHhap564rX", opaque="XF3tAJ3483jUUAUJJQJJAHDQP01MJHD", qop="auth", algorithm="SHA-256"
```

- `qop` (Quality of Protection)
- `nonce` : random server generated sequence of chars (used by client to calculate the response hash)
- `opaque` : random server generated sequence of chars (sent back unchanged in header)

Follow up request using Digest Auth credentials:

```sh
curl -i --digest --user go:time https://127.0.0.1:8443/customers \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "3skills" 
}'

HTTP/1.1 201 Created
```

## JWT Auth

Initial request:

```sh
curl -i https://127.0.0.1:8443/customers/1/invoices \
  -H 'Content-Type: application/json' \
  -d '{
    "month": 9,
    "year": 2020
}'

HTTP/1.1 401 Unauthorized
Www-Authenticate: Bearer realm="invoice.mvp"
```

- `qop` (Quality of Protection)
- `nonce` : random server generated sequence of chars (used by client to calculate the response hash)
- `opaque` : random server generated sequence of chars (sent back unchanged in header)

Follow up request using Digest Auth credentials:

```sh
curl -i https://127.0.0.1:8443/customers/1/invoices \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiR28gSW52b2ljZXIiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiZjhjMzlhMzEtOWNlZC00NzYxLThhMzMtYjljNjI4YTY3NTEwIn0.WI6cRXYnYqUAV6qqNtf4B8PdGMgKuHqENQP5N_iCZL8' \
  -d '{
    "month": 9,
    "year": 2020
}'

HTTP/1.1 201 Created
```
