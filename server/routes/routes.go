// Package routes the routes package
package routes

import (
	"log"

	"github.com/phanirithvij/central_server/server/config"
)

var (
	// Endpoints ...
	endpoints = map[string]bool{}
)

// CheckEndpoints checks whether each route's endpoints are initialized
func CheckEndpoints() {
	done := true
	notDone := []string{}
	for _, v := range config.EndPointStrings {
		if c, ok := endpoints[v]; ok {
			done = done && c
		} else {
			done = false
			notDone = append(notDone, v)
		}
	}
	if !done {
		log.Println("[Warning] some endpoints are not registered")
		for _, v := range notDone {
			log.Printf("[Warning] endpoints of %s are not registered", v)
		}
	}
}

// RegisterSelf routes needs to use this method to register themselves
func RegisterSelf(name string) {
	endpoints[name] = true
}
