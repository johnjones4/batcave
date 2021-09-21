package main

import (
	"hal9000/types"
	"net/http"
)

func jobHandler(runtime *types.Runtime) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

	}
}
