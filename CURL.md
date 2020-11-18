# CURL

## Interactions

## GET /customers

> To create an invoice a `customer ID` is requested.

```sh
curl -i http://localhost:8080/customers

# response
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 27

[{"id":1,"name":"3skills"}]
```

---

## POST /customers/{customerId}/invoices

> The API currently only knows of the customer with ID 1, that is used in the URI for creating an invoice.

```sh
curl -i http://localhost:8080/customers/1/invoices \
  -H 'Content-Type: application/json' \
  -d '{
    "month": 6,
    "year": 2018
}'
```

```sh
# response
HTTP/1.1 201 Created
Content-Type: application/json
Location: /customers/1/invoices/1
Content-Length: 104

{
  "id": 2,
  "month": 6,
  "year": 2018,
  "status": "open",
  "customerId": 1
}
```

---

## Projects & Activities

> The Invoice-Use-Case needs an `invoice ID` and an `activity ID` in order to create bookings.

Projects:

```sh
curl -s http://localhost:8080/customers/1/projects | jq

# response
[
  {
    "id": 1,
    "customerId": 1,
    "name": "Instantfoo.com"
  }
]
```

Activities:

```sh
curl -s http://localhost:8080/activities | jq

# response
[
  {
    "id": 1,
    "name": "Programming",
    "userId": ""
  }
]
```

---

## POST /customers/{customerId}/invoices/{invoiceId}/bookings

### Add bookings

> A new `Booking` references an `Invoice` with `invoiceId` variable in the URL.

The customer only has one project `Instantfoo.com` with ID 1.
The list of activities contains the activity `Programming` with ID 1.

The activity is now booked on this project:

```sh
# customer   : 1 (3skills)
# invoice    : 2 (2018-06)
# projectId  : 1 (Instantfoo.com)
# activityId : 1 (Programming)

curl -i http://localhost:8080/customers/1/invoices/2/bookings \
  -H 'Content-Type: application/json' \
  -d '{
    "day": 31,
    "hours": 2.5,
    "projectId": 1,
    "activityId": 1
}'

# response
HTTP/1.1 201 Created
Content-Type: application/json
Location: /customers/1/invoices/2/bookings/1
Content-Length: 82

{"day":31,"hours":2.5,"description":"","invoiceId":2,"projectId":1,"activityId":1}
```

> The `Booking` of type `Programming` was succesfully created for the `Project` *instantfoo.com* on `Invoice` 2.

---

## DELETE /customers/{customerId}/invoices/{invoiceId}/bookings/{bookingId}

### Delete bookings

Currently it is not possible to correct bookings but required in the scope of this MVP.

The simple solution here is to just delete the booking and create a new one.

The error handling while deleting a non-existing booking will be ignored using a nil-operation, a function that does nothing.

```sh
curl -i -X DELETE http://localhost:8080/customers/1/invoices/2/bookings/1

# response
HTTP/1.1 204 No Content
```

---

## PUT /customers/{customerId}/invoices/{invoiceId}

### Finalize an invoice

The handling of an invoice finalization is implemented using a PUT-Request, as an existing resource is getting updated.

```sh
curl -i -X PUT http://localhost:8080/customers/1/invoices/2 \
  -H 'Content-Type: application/json' \
  -d '{
    "month": 6,
    "year": 2018,
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
curl -s http://localhost:8080/customers/1/invoices/2 -H 'Accept: application/json' | jq

# response
{
  "id": 2,
  "month": 6,
  "year": 2018,
  "status": "payment expected",
  "customerId": 1,
  "positions": {
    "1": {
      "Programming": {
        "Hours": 2.5,
        "Price": 151.375
      }
    }
  }
}
```

---

