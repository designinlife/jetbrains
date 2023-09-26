package cmd

import (
	"fmt"
	"github.com/designinlife/jetbrains/common"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

// versionCmd represents the ls command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("%s %s built on %s (%s)", common.Name, common.Version, common.BuiltOn, runtime.Version()))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
