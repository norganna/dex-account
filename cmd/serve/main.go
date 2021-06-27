package serve

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	runningServer *serveOptions
)
// CommandServe returns the command handler.
func CommandServe() *cobra.Command {
	options := &serveOptions{}

	cmd := &cobra.Command{
		Use:     "serve [flags] [config]",
		Short:   "Launch Dex account API",
		Example: "dex-account serve",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				options.configFile = args[0]
			}
			runningServer = options
			return runServe(options)
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringVar(&options.WebHTTPAddr, "web-http-addr", options.WebHTTPAddr, "Web HTTP address")
	flags.StringVar(&options.WebHTTPSAddr, "web-https-addr", options.WebHTTPSAddr, "Web HTTPS address")
	flags.StringVar(&options.GrpcAddr, "grpc-addr", options.GrpcAddr, "gRPC API address")

	viper.BindPFlag("web-http-addr", flags.Lookup("web-http-addr"))
	viper.BindPFlag("web-https-addr", flags.Lookup("web-https-addr"))
	viper.BindPFlag("grpc-addr", flags.Lookup("grpc-addr"))

	return cmd
}
