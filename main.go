package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tullo/invoice-mvp/database"
	"github.com/tullo/invoice-mvp/rest"
	"github.com/tullo/invoice-mvp/roles"
	"github.com/tullo/invoice-mvp/usecase"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		os.Exit(1)
	}

	repository := database.NewFakeRepository()
	a := rest.NewAdapter()

	// IDP redirects to this URI after user authentication
	a.HandleFunc("/auth/token", rest.OAuth2AccessCodeGrant(a.OAuth2AccessTokenHandler())).Methods("GET")

	// Activities
	activities := usecase.NewActivities(repository)
	ga := a.ActivitiesHandler(activities)
	ga = rest.JWTAuth(ga)
	a.HandleFunc("/activities", ga).Methods("GET")

	createActivity := usecase.NewCreateActivity(repository)
	ca := a.CreateActivityHandler(createActivity)
	ca = rest.JWTAuth(ca)
	a.HandleFunc("/activities", ca).Methods("POST")

	// Booking
	createBooking := usecase.NewCreateBooking(repository)
	cb := a.CreateBookingHandler(createBooking)
	cb = rest.JWTAuth(roles.AssertOwnsInvoice(cb, repository))
	a.HandleFunc("/book/{invoiceId:[0-9]+}", cb).Methods("POST")

	deleteBooking := usecase.NewDeleteBooking(repository)
	db := a.DeleteBookingHandler(deleteBooking)
	db = rest.JWTAuth(db)
	a.HandleFunc("/invoices/{invoiceId:[0-9]+}/bookings/{bookingId:[0-9]+}", db).Methods("DELETE")

	// Customer
	createCustomer := usecase.NewCreateCustomer(repository)
	cc := a.CreateCustomerHandler(createCustomer)
	cc = rest.JWTAuth(cc)
	a.HandleFunc("/customers", cc).Methods("POST")

	// Invoice
	createInvoice := usecase.NewCreateInvoice(repository)
	ci := a.CreateInvoiceHandler(createInvoice)
	ci = rest.JWTAuth(ci)
	a.HandleFunc("/customers/{customerId:[0-9]+}/invoices", ci).Methods("POST")

	updateInvoice := usecase.NewUpdateInvoice(repository)
	ui := a.UpdateInvoiceHandler(updateInvoice)
	ui = rest.JWTAuth(ui)
	a.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", ui).Methods("PUT")

	invoice := usecase.NewGetInvoice(repository)
	gi := a.GetInvoiceHandler(invoice)
	gi = rest.JWTAuth(gi)
	a.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", gi).Methods("GET")

	// Project
	createProject := usecase.NewCreateProject(repository)
	cp := a.CreateProjectHandler(createProject)
	cp = rest.JWTAuth(roles.AssertAdmin(cp))
	a.HandleFunc("/customers/{customerId:[0-9]+}/projects", cp).Methods("POST")

	// Hourly rate
	createRate := usecase.NewCreateRate(repository)
	cr := a.CreateRateHandler(createRate)
	cr = rest.JWTAuth(cr)
	a.HandleFunc("/customers/{customerId:[0-9]+}/projects/{projectId:[0-9]+}/rates", cr).Methods("POST")

	// Webserver
	a.ListenAndServeTLS()
}
