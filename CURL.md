# CURL

## Interactions

---

## POST /activities

Activities:

```sh
curl -i http://localhost:8080/activities   -H 'Content-Type: application/json'   -d '{
    "name": "Programming" 
}'

# response             
HTTP/1.1 201 Created
Content-Type: application/json
Location: /activities/1
```

## GET /activities

```sh
curl -i http://localhost:8080/activities   -H 'Content-Type: application/json'

# response             
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 47

[{"id":1,"name":"Programming","userId":"1234"}]
```

---

## POST /customers

```sh
curl -i http://localhost:8080/customers   -H 'Content-Type: application/json'   -d '{
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
curl -i http://localhost:8080/customers/1/invoices \
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
curl -i http://localhost:8080/customers/1/projects \
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
curl -i http://localhost:8080/customers/1/projects/1/rates \
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

curl -i http://localhost:8080/customers/1/invoices/1/bookings \
  -H 'Content-Type: application/json' \
  -d '{
    "day": 31,
    "hours": 2.5,
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
curl -i -X DELETE http://localhost:8080/customers/1/invoices/1/bookings/1

# response
HTTP/1.1 204 No Content
```

---

## PUT /customers/{customerId}/invoices/{invoiceId}

### Finalize an invoice

The handling of invoice finalization is implemented using a PUT-Request as an existing resource is getting updated.

```sh
curl -i -X PUT http://localhost:8080/customers/1/invoices/1 \
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
curl -s http://localhost:8080/customers/1/invoices/1 -H 'Accept: application/json' | jq

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
}
```

### Retrieve an Invoice with booking details

```sh
curl -s http://localhost:8080/customers/1/invoices/1 -H 'Accept: application/json' | jq

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
      "description": "",
      "invoiceId": 1,
      "projectId": 1,
      "activityId": 1
    }
  ],
}
```
