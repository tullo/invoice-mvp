package main

import (
	"github.com/tullo/invoice-mvp/database"
	"github.com/tullo/invoice-mvp/rest"
	"github.com/tullo/invoice-mvp/usecase"
)

func main() {
	repository := database.NewFakeRepository()
	a := rest.NewAdapter()

	// Activities
	activities := usecase.NewActivities(repository)
	a.HandleFunc("/activities", a.ActivitiesHandler(activities)).Methods("GET")

	createActivity := usecase.NewCreateActivity(repository)
	a.HandleFunc("/activities", a.CreateActivityHandler(createActivity)).Methods("POST")

	// Booking
	createBooking := usecase.NewCreateBooking(repository)
	cb := a.CreateBookingHandler(createBooking)
	a.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}/bookings", cb).Methods("POST")

	deleteBooking := usecase.NewDeleteBooking(repository)
	db := a.DeleteBookingHandler(deleteBooking)
	a.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}/bookings/{bookingId:[0-9]+}", db).Methods("DELETE")

	// Customer
	createCustomer := usecase.NewCreateCustomer(repository)
	cc := a.CreateCustomerHandler(createCustomer)
	a.HandleFunc("/customers", cc).Methods("POST")

	// Invoice
	createInvoice := usecase.NewCreateInvoice(repository)
	ci := a.CreateInvoiceHandler(createInvoice)
	a.HandleFunc("/customers/{customerId:[0-9]+}/invoices", ci).Methods("POST")

	updateInvoice := usecase.NewUpdateInvoice(repository)
	a.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", a.UpdateInvoiceHandler(updateInvoice)).Methods("PUT")

	invoice := usecase.NewGetInvoice(repository)
	a.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", a.GetInvoiceHandler(invoice)).Methods("GET")

	// Projekt
	createProject := usecase.NewCreateProject(repository)
	a.HandleFunc("/customers/{customerId:[0-9]+}/projects", a.CreateProjectHandler(createProject)).Methods("POST")

	// Stundensatz
	createRate := usecase.NewCreateRate(repository)
	a.HandleFunc("/customers/{customerId:[0-9]+}/projects/{projectId:[0-9]+}/rates", a.CreateRateHandler(createRate)).Methods("POST")

	// Webserver
	a.ListenAndServe()
}
