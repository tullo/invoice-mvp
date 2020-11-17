package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Contact holds first- and lastname of a contact.
type Contact struct {
	Firstname string
	Lastname  string
}

var cm = make(map[int]Contact)

func main() {
	r := mux.NewRouter()
	cm[1] = Contact{Firstname: "Andreas", Lastname: "Amstutz"}

	r.HandleFunc("/contacts", contacts).Methods("GET")
	r.HandleFunc("/contacts/{id:[0-9]+}", contact).Methods("GET")
	r.HandleFunc("/contacts", addContact).Methods("POST")
	r.HandleFunc("/contacts/{id:[0-9]+}", updateContact).Methods("PUT")
	r.HandleFunc("/contacts/{id:[0-9]+}", deleteContact).Methods("DELETE")
	log.Printf("Service listening on http://localhost:8080...")
	_ = http.ListenAndServe(":8080", r)
}

func contacts(w http.ResponseWriter, r *http.Request) {
	var s []Contact
	for _, c := range cm {
		s = append(s, c)
	}
	b, err := json.Marshal(s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(b)
}

func contact(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	id, _ := strconv.Atoi(v["id"])
	if _, ok := cm[id]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := json.Marshal(cm[id])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(b)
}

func addContact(w http.ResponseWriter, r *http.Request) {
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		if len(b) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id := nextID()
		var c Contact
		_ = json.Unmarshal(b, &c)
		cm[id] = c
		url := r.URL.String()
		w.Header().Set("Location", fmt.Sprintf("%s/%d", url, id))
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func updateContact(w http.ResponseWriter, r *http.Request) {
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		if len(b) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		v := mux.Vars(r)
		id, _ := strconv.Atoi(v["id"])
		var c Contact
		_ = json.Unmarshal(b, &c)
		cm[id] = c
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	id, err := strconv.Atoi(v["id"])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, ok := cm[id]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	delete(cm, id)
	w.WriteHeader(http.StatusNoContent)
}

func nextID() int {
	id := 1
	for k := range cm {
		if k >= id {
			id = k + 1
		}
	}
	return id
}
