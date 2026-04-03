package cmdutil

import "github.com/spf13/cobra"

func GetStringFlag(cmd *cobra.Command, name string) string {
	val, _ := cmd.Flags().GetString(name)
	return val
}

func GetBoolFlag(cmd *cobra.Command, name string) bool {
	val, _ := cmd.Flags().GetBool(name)
	return val
}
