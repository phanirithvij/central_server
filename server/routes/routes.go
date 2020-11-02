// Package routes the routes package
package routes

import (
	"fmt"

	"github.com/phanirithvij/central_server/server/routes/api"
	"github.com/phanirithvij/central_server/server/routes/home"
	"github.com/phanirithvij/central_server/server/routes/register"
)

// CheckEndpoints checks whether each route's endpoints are initialized
func CheckEndpoints() {
	done := false
	done = done || home.EndpointsRegistered
	done = done || register.EndpointsRegistered
	done = done || api.EndpointsRegistered
	if !done {
		fmt.Println("[Warning] some endpoints are not registered")
	}
}
