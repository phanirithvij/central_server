package main

import "github.com/phanirithvij/central_server/server/cmd"

//go:generate pkger -o server

func main() {
	cmd.Execute()
}
