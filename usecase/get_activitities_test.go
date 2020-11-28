package usecase_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/tullo/invoice-mvp/domain"
	"github.com/tullo/invoice-mvp/identityprovider/fusionauth"

	"github.com/stretchr/testify/assert"
	"github.com/tullo/invoice-mvp/database"
	"github.com/tullo/invoice-mvp/rest"
	"github.com/tullo/invoice-mvp/usecase"
)

func TestHttpGetActivities(t *testing.T) {
	//=========================================================================
	// Setup
	r := database.NewFakeRepository()
	setupBaseData(r)
	activities := usecase.NewActivities(r)
	userID := r.CustomerByID(customer).UserID

	// Login to IDM
	data := make(url.Values)
	data.Set("loginId", os.Getenv("MVP_USERNAME"))
	data.Set("password", os.Getenv("MVP_PASSWORD"))
	auth, err := fusionauth.Login(data)
	if err != nil {
		t.Error(err)
	}

	// Prepare HTTP-Request
	req, _ := http.NewRequest("GET", "/activities", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth.AccessToken))
	req.Header.Set("Content-Type", "application/json")

	//=========================================================================
	// Get request
	a := rest.NewAdapter()
	ga := a.ActivitiesHandler(activities)
	ga = rest.JWTAuth(ga)
	a.Handle("/activities", ga).Methods("GET")

	res := httptest.NewRecorder()
	a.R.ServeHTTP(res, req)

	//================== UNCACHED ACTIVITIES ==================================
	// Assert
	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	assert.Equal(t, []string{"public, max-age=0"}, res.Result().Header["Cache-Control"])
	assert.Equal(t, []string{"application/json"}, res.Result().Header["Content-Type"])
	var as []domain.Activity
	as = append(as, domain.Activity{ID: 1, Name: "Programming", UserID: userID})
	as = append(as, domain.Activity{ID: 2, Name: "Quality control", UserID: userID})
	as = append(as, domain.Activity{ID: 3, Name: "Project management", UserID: userID})
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var alist []domain.Activity
	if err := json.Unmarshal(body, &alist); err != nil {
		t.Fatal(err)
	}
	assert.ElementsMatch(t, as, alist)

	//================== CACHED ACTIVITIES ====================================

	// Prepare HTTP-Request
	req, _ = http.NewRequest("GET", "/activities", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth.AccessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Last-Modified-Since", res.Result().Header["Last-Modified"][0])

	res = httptest.NewRecorder()
	a.R.ServeHTTP(res, req)

	//=========================================================================
	// Assert
	assert.Equal(t, http.StatusNotModified, res.Result().StatusCode)

	//================== RELOADED ACTIVITIES ==================================

	// Prepare HTTP-Request
	req, _ = http.NewRequest("GET", "/activities", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth.AccessToken))
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")

	res = httptest.NewRecorder()
	a.R.ServeHTTP(res, req)

	//=========================================================================
	// Assert
	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	var nocache []domain.Activity
	if err := json.Unmarshal(body, &nocache); err != nil {
		t.Fatal(err)
	}
	assert.ElementsMatch(t, as, nocache)
}
