package usecase_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tullo/invoice-mvp/database"
	"github.com/tullo/invoice-mvp/domain"
	"github.com/tullo/invoice-mvp/rest"
	"github.com/tullo/invoice-mvp/roles"
	"github.com/tullo/invoice-mvp/usecase"
)

func TestHttpCreateProjectUnauthorized(t *testing.T) {
	//=========================================================================
	// Setup
	r := database.NewFakeRepository()
	setupBaseData(r)
	createProject := usecase.NewCreateProject(r)
	userToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiR28gSW52b2ljZXIiLCJhZG1pbiI6ZmFsc2UsInN1YiI6ImY4YzM5YTMxLTljZWQtNDc2MS04YTMzLWI5YzYyOGE2NzUxMCJ9.pqpFYWFJdwZs7fXN8YAUriIyLuAfAB9uBjFfA4C6g9U"
	/* userToken payload:
	   {
	     "name": "Go Invoicer",
	     "admin": false,
	     "sub": "f8c39a31-9ced-4761-8a33-b9c628a67510"
	   }
	*/

	// Prepare HTTP-Request
	p := domain.Project{CustomerID: customer, Name: "Testing"}
	bs, _ := json.Marshal(&p)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/customers/%d/projects", customer), bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userToken))

	//=========================================================================
	// Add project using POST request
	res := httptest.NewRecorder()
	a := rest.NewAdapter()
	cp := a.CreateProjectHandler(createProject)
	cp = rest.JWTAuth(roles.AssertAdmin(cp))
	a.HandleFunc("/customers/{customerId:[0-9]+}/projects", cp).Methods("POST")
	a.R.ServeHTTP(res, req)

	//=========================================================================
	// Assert
	assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode)
}

func TestHttpCreateProjectAuthorized(t *testing.T) {
	//=========================================================================
	// Setup
	r := database.NewFakeRepository()
	setupBaseData(r)
	createProject := usecase.NewCreateProject(r)

	// Prepare HTTP-Request
	p := domain.Project{CustomerID: customer, Name: "Testing"}
	bs, _ := json.Marshal(&p)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/customers/%d/projects", customer), bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

	//=========================================================================
	// Add project using POST request
	res := httptest.NewRecorder()
	a := rest.NewAdapter()
	cp := a.CreateProjectHandler(createProject)
	cp = rest.JWTAuth(roles.AssertAdmin(cp))
	a.HandleFunc("/customers/{customerId:[0-9]+}/projects", cp).Methods("POST")
	a.R.ServeHTTP(res, req)

	//=========================================================================
	// Assert
	assert.Equal(t, http.StatusCreated, res.Result().StatusCode)
	assert.Equal(t, "/customers/1/projects/3", res.Result().Header["Location"][0])
	expected := domain.Project(domain.Project{ID: 3, CustomerID: p.CustomerID, Name: p.Name})
	actual := r.ProjectByID(expected.ID)
	assert.Equal(t, expected, actual)
}
