package serve

import "github.com/spf13/cobra"

var (
	runningServer *serveOptions
)
// CommandServe returns the command handler.
func CommandServe() *cobra.Command {
	options := &serveOptions{}

	cmd := &cobra.Command{
		Use:     "serve [flags] [config]",
		Short:   "Launch Dex",
		Example: "dex serve",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				options.configFile = args[0]
			}
			runningServer = options
			return runServe(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.WebHTTPAddr, "web-http-addr", options.WebHTTPAddr, "Web HTTP address")
	flags.StringVar(&options.WebHTTPSAddr, "web-https-addr", options.WebHTTPSAddr, "Web HTTPS address")
	flags.StringVar(&options.GrpcAddr, "grpc-addr", options.GrpcAddr, "gRPC API address")

	return cmd
}
