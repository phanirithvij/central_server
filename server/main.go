package main

import "github.com/phanirithvij/central_server/server/cmd"

//go:generate rm -f pkged.go
//go:generate pkger -o server
//go:generate mv pkged.go pkged_g.go

func main() {
	cmd.Execute()
}
