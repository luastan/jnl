/*
Copyright © 2022 Luastan <@luastan>
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/luastan/jnl/jnl"
	"github.com/spf13/cobra"
	"os"
)

// getHistoryCmd represents the getHistory command
var getHistoryCmd = &cobra.Command{
	Use:   "getHistory",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("getHistory called")
		executions, err := jnl.EveryExecution()
		if err != nil {
			jnl.ErrorLogger.Fatalln(err.Error())
		}

		switch *format {
		case "table":
			err := jnl.InfoInTable(executions, os.Stdout)
			if err != nil {
				jnl.ErrorLogger.Fatalln(err.Error())
			}
			break
		case "CSV":
		case "csv":
			err := jnl.InfoInCSV(executions, os.Stdout)
			if err != nil {
				jnl.ErrorLogger.Fatalln(err.Error())
			}
			break
		default:
			jnl.ErrorLogger.Fatalln(errors.New(fmt.Sprintf("unknown format: \"%s\"", *format)))
		}
	},
}

func init() {
	rootCmd.AddCommand(getHistoryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getHistoryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getHistoryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
