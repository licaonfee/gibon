# gibon

![Run test](https://github.com/licaonfee/gibon/workflows/Run%20test/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/licaonfee/gibon/badge.svg?branch=master)](https://coveralls.io/github/licaonfee/gibon?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/licaonfee/gibon)](https://goreportcard.com/report/github.com/licaonfee/gibon)


Create middleware pipelines for http, compatible with standard http library

## Example

```golang
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/licaonfee/gibon"
)

//log all requests
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.RequestURI)
        next.ServeHTTP(w, r)
    })
}

//authenticate valid users
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

//only if your browser can brew a coffee
func coffeeMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if strings.HasPrefix(r.UserAgent(), "CoffePot") {
            next.ServeHTTP(w, r)
        }
        http.Error(w, "I'm a teapot", http.StatusTeapot)
    })
}

// Some handlers
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

func main() {
    //First we want to log all request
    logPipeline := gibon.New().With(loggingMiddleware)
    //In some cases we want authenticantion
    //every call to With creates a new pipeline with all
    //middlewares from previous object
    authPipeline := logPipeline.With(authMiddleware)
    //Just serve requests from barista's browser
    coffeePipeline := logPipeline.With(coffeeMiddleware)
    //only serve baristas who are administrators
    secretCoffeePipeline := authPipeline.With(coffeeMiddleware)

    http.Handle("/", logPipeline.BuildFunc(index))
    http.Handle("/private", authPipeline.BuildFunc(privateA))
    http.Handle("/private/coffee", secretCoffeePipeline.BuildFunc(privateB))
    http.Handle("/coffee", coffeePipeline.BuildFunc(public))
    if err := http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
        log.Println(err)
    }
}

```
