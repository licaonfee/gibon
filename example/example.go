package main

import (
	"fmt"
	"log"
	"net/http"

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
	fmt.Fprint(w, "Hello World")
}

func privateA(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Confidential")
}

func privateB(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Top Secret")
}

func public(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Public Domain")
}

/* test with
curl -vvv http://localhost:3000/index
curl -vvv http://localhost:3000/public
curl -vvv http://admin:pass123@localhost:3000/priv
curl -vvv http://admin:XXXX@localhost:3000/secret
*/

func main() {

	logPipeline := gibon.New().Add(loggingMiddleware)
	authPipeline := logPipeline.Add(authMiddleware)

	http.Handle("/", logPipeline.BuildFunc(index))
	http.Handle("/priv", authPipeline.BuildFunc(privateA))
	http.Handle("/secret", authPipeline.BuildFunc(privateB))
	http.Handle("/public", logPipeline.BuildFunc(public))
	if err := http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		log.Println(err)
	}
}
