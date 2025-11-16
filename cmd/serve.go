package cmd

import (
	"github.com/spf13/cobra"

	"codeberg.org/aryak/mozhi/serve"
)

var host string = ""
var port string = "3000"

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server.",
	Long:  `Start the web server.`,
	Run: func(cmd *cobra.Command, args []string) {
		serve.Serve(host, port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&host, "host", "H", "", "The host Mozhi will listen on. Defaults to listening on all interfaces, and overrides the MOZHI_HOST environment variable.")
	serveCmd.Flags().StringVarP(&port, "port", "p", "", "The port Mozhi will listen on. Defaults to 3000, and overrides the MOZHI_PORT environment variable.")

	// set variables to the value of the flags
	host = serveCmd.Flag("host").Value.String()
	port = serveCmd.Flag("port").Value.String()
}
