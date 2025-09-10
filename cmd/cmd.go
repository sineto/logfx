package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var filepath string
var fctype string

func RootCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "logfx",
		Short:   "logfx - CLI to process file compression",
		PreRunE: prerunerr,
		RunE:    run,
	}

	cmd.Flags().StringVar(&filepath, "log-path", "", "path to log file or directory")
	cmd.Flags().StringVar(&fctype, "type", "targz", "compression file type - use (targz|zip)")

	return cmd
}

func Execute() {
	if err := RootCmd().Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error while running command %v\n", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	switch fctype {
	case "targz":
		return targz(&filepath)
	}

	return nil
}

func prerunerr(cmd *cobra.Command, args []string) error {
	if filepath == "" {
		return fmt.Errorf("filepath is required")
	}

	switch fctype {
	case "targz", "zip":
	default:
		return fmt.Errorf("unknown compression type: %s", fctype)
	}

	return nil
}
