package server

import "net/http"

type Auth struct {
	handler http.HandlerFunc
	all     bool
}
