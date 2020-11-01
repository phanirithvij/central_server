package cmd

import (
	"github.com/phanirithvij/btp/central/server/serve"
	"github.com/spf13/cobra"
)

var (
	port int
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server",
	Run: func(cmd *cobra.Command, args []string) {
		serve.Serve(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().IntVarP(&port, "port", "p", 9090, "The port to serve the server")
}
