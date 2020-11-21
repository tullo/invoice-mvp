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
	ga := a.ActivitiesHandler(activities)
	a.HandleFunc("/activities", ga).Methods("GET")

	createActivity := usecase.NewCreateActivity(repository)
	ca := a.CreateActivityHandler(createActivity)
	a.HandleFunc("/activities", ca).Methods("POST")

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
	ui := a.UpdateInvoiceHandler(updateInvoice)
	a.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", ui).Methods("PUT")

	invoice := usecase.NewGetInvoice(repository)
	gi := a.GetInvoiceHandler(invoice)
	a.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", gi).Methods("GET")

	// Project
	createProject := usecase.NewCreateProject(repository)
	cp := a.CreateProjectHandler(createProject)
	a.HandleFunc("/customers/{customerId:[0-9]+}/projects", cp).Methods("POST")

	// Hourly rate
	createRate := usecase.NewCreateRate(repository)
	cr := a.CreateRateHandler(createRate)
	a.HandleFunc("/customers/{customerId:[0-9]+}/projects/{projectId:[0-9]+}/rates", cr).Methods("POST")

	// Webserver
	a.ListenAndServe()
}
