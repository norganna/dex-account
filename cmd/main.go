package main

import (
	"fmt"
	"github.com/norganna/dex-account-api/cmd/serve"
	"os"

	"github.com/spf13/cobra"
)

func commandRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "dex",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(2)
		},
	}
	rootCmd.AddCommand(serve.CommandServe())
	return rootCmd
}

func main() {
	if err := commandRoot().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}

	serve.Wait()
}