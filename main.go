package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats"
)

func main() {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	r.Methods("POST").
		Path("/{subject}").Handler(natsHandler(nc))
	log.Println("listening on 8080")
	http.ListenAndServe(":8080", r)
}

func natsHandler(nc *nats.Conn) http.Handler {
	result := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		subject := vars["subject"]
		msg, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = nc.Publish(subject, msg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return http.Handler(http.HandlerFunc(result))
}
