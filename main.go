package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/tullo/invoice-mvp/database"
	"github.com/tullo/invoice-mvp/domain"
)

var repository = database.NewRepository()

func init() {
	rand.Seed(time.Now().UnixNano())
	// customer
	cid := repository.AddCustomer("3skills")
	// project
	pid := repository.AddProject("Instantfoo.com", cid)
	// invoice
	repository.CreateInvoice(domain.Invoice{Month: 6, Year: 2018, CustomerID: cid})
	// activity
	aid := repository.AddActivity("Programming")
	// rate
	var r domain.Rate
	r.ProjectID = pid
	r.ActivityID = aid
	r.Price = 60.55
	repository.AddRate(r)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/customers", customers).Methods("GET")
	r.HandleFunc("/customers/{customerId:[0-9]+}/projects", projects).Methods("GET")
	r.HandleFunc("/activities", activities).Methods("GET")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices", createInvoice).Methods("POST")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}/bookings", createBooking).Methods("POST")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}/bookings/{bookingId:[0-9]+}", deleteBooking).Methods("DELETE")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", updateInvoice).Methods("PUT")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", readInvoiceHandler).Methods("GET")

	fmt.Println("Restvoice started on http://localhost:8080...")
	_ = http.ListenAndServe(":8080", r)
}

func customers(w http.ResponseWriter, r *http.Request) {
	cs := repository.Customers()
	b, _ := json.Marshal(cs)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(b)
}

func projects(w http.ResponseWriter, r *http.Request) {
	cid, _ := strconv.Atoi(mux.Vars(r)["customerId"])
	ps := repository.Projects(cid)
	b, _ := json.Marshal(ps)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(b)
}

func activities(w http.ResponseWriter, r *http.Request) {
	as := repository.Activities()
	b, _ := json.Marshal(as)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(b)
}

func createInvoice(w http.ResponseWriter, r *http.Request) {
	// Read invoice data from request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Unmarshal json payload
	var i domain.Invoice
	_ = json.Unmarshal(body, &i)

	// extract customerId from URI
	i.CustomerID, _ = strconv.Atoi(mux.Vars(r)["customerId"])

	// Create invoice
	created, _ := repository.CreateInvoice(i)

	// Marshal invoice to json
	b, err := json.Marshal(created)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Write response
	location := fmt.Sprintf("%s/%d", r.URL.String(), created.ID)
	w.Header().Set("Location", location)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(b)
}

func createBooking(w http.ResponseWriter, r *http.Request) {
	// Read booking data from request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create booking booking and marshal it to JSON
	var b domain.Booking
	if err := json.Unmarshal(body, &b); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	b.InvoiceID, _ = strconv.Atoi(mux.Vars(r)["invoiceId"])
	created, _ := repository.CreateBooking(b)
	bs, err := json.Marshal(created)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response
	location := fmt.Sprintf("%s/%d", r.URL.String(), created.ID)
	w.Header().Set("Location", location)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(bs)
}

func deleteBooking(w http.ResponseWriter, r *http.Request) {
	bid, _ := strconv.Atoi(mux.Vars(r)["bookingId"])
	repository.DeleteBooking(bid)
	w.WriteHeader(http.StatusNoContent)
}

func updateInvoice(w http.ResponseWriter, r *http.Request) {
	// Read invoice data from request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Unmarshal and update invoice
	var i domain.Invoice
	if err := json.Unmarshal(body, &i); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	i.ID, _ = strconv.Atoi(mux.Vars(r)["invoiceId"])
	i.CustomerID, _ = strconv.Atoi(mux.Vars(r)["customerId"])

	// Aggregate positions
	if i.Status == "ready for aggregation" {
		bs := repository.BookingsByInvoiceID(i.ID)
		for _, b := range bs {
			a := repository.ActivityByID(b.ActivityID)
			rate := repository.RateByProjectIDAndActivityID(b.ProjectID, b.ActivityID)
			i.AddPosition(b.ProjectID, a.Name, b.Hours, rate.Price)
		}
		i.Status = "payment expected"
		i.Updated = time.Now().UTC()
	}

	repository.Update(i)

	// Write response
	w.WriteHeader(http.StatusNoContent)
}

func readInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["invoiceId"])
	i, _ := repository.FindByID(id)
	accept := r.Header.Get("Accept")
	switch accept {
	case "application/pdf":
		content := bytes.NewReader(i.ToPDF())
		http.ServeContent(w, r, "invoice.pdf", i.Updated, content)
		// ServeContent sets the content-type header and makes sure that the
		// following headers get the correct values as well:
		// - If-Match
		// - If-Unmodified-Since
		// - If-None-Match
		// - If-Modified-Since
		// - If-Range
	case "application/json":
		bs, _ := json.Marshal(i)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(bs)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}
