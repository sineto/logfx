package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	fromDir string
	toDir   string
	fctype  string
)

var jsonHandler = slog.NewJSONHandler(os.Stderr, nil)
var logger = slog.New(jsonHandler)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logfx",
		Short:   "logfx - CLI to process file compression",
		PreRunE: cmdRunErr,
		RunE:    run,
	}

	cmd.PersistentFlags().StringVar(&fromDir, "from", "", "path directory to compress logs (required)")
	cmd.MarkPersistentFlagRequired("from")

	cmd.PersistentFlags().StringVar(&toDir, "to", ".", "path directory to move the compressed logs")
	cmd.PersistentFlags().StringVar(&fctype, "type", "targz", "compression file type - use (targz|zip)")

	cmd.AddCommand(CronCmd())
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
		return targz(&fromDir, &toDir)
	}

	return nil
}

func cmdRunErr(cmd *cobra.Command, args []string) error {
	switch fctype {
	case "targz", "zip":
	default:
		return fmt.Errorf("unknown compression type: %s", fctype)
	}

	return nil
}
