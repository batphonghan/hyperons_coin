package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hyperon",
	Short: "hyperon coin deamon",
	Long:  "hyperon coin deamon",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
