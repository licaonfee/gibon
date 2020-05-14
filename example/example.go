package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/licaonfee/gibon"
)

//just copied from gorilla mux
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u, p, ok := r.BasicAuth(); ok {
			if u == "admin" && p == "pass123" {
				log.Printf("Access granted to %s", u)
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Access Denied", http.StatusForbidden)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	return
}

func privateA(w http.ResponseWriter, r *http.Request) {
	return
}
func privateB(w http.ResponseWriter, r *http.Request) {
	return
}

func public(w http.ResponseWriter, r *http.Request) {
	return
}

/* test with
curl -vvv http://admin:pass1x23@localhost:3000/public
curl -vvv http://admin:NOPASS1122@localhost:3000/public
curl -vvv http://admin:pass123@localhost:3000/priv
curl -vvv http://admin:pass1x23@localhost:3000/public
curl -vvv http://admin:pass1x23@localhost:3000/priv
*/

func main() {
	logRequest := gibon.Add(loggingMiddleware)
	auth := logRequest.Add(authMiddleware)

	r := mux.NewRouter()
	r.Handle("/", logRequest.Build(http.HandlerFunc(index)))
	r.Handle("/priv", auth.Build(http.HandlerFunc(privateA)))
	r.Handle("/secret", auth.Build(http.HandlerFunc(privateB)))
	r.Handle("/public", logRequest.Build(http.HandlerFunc(public)))
	if err := http.ListenAndServe("0.0.0.0:3000", r); err != nil {
		log.Println(err)
	}
}
