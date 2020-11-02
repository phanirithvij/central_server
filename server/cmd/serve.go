package cmd

import (
	"github.com/phanirithvij/central_server/server/serve"
	"github.com/spf13/cobra"
)

var (
	port  int
	debug bool
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server",
	Run: func(cmd *cobra.Command, args []string) {
		serve.Serve(port, debug)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().IntVarP(&port, "port", "p", 9090, "The port to serve the server")
	serveCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "If debug or release")
}
