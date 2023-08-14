/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/BrunoTulio/GoDump/internal/factory"
	"github.com/spf13/cobra"
)

var (
	restoreDriverEnable bool
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "r",
	Long:  `Restore databases`,
	Run: func(cmd *cobra.Command, args []string) {
		log := factory.MakeLogger()

		restoreAutoUseCase := factory.MakeDefaultRestoreAutoUseCase()

		err := restoreAutoUseCase.Restore(restoreDriverEnable)

		if err != nil {
			log.Fatal(err)
		}

		log.Info("Restored databases successfully")

	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().BoolVarP(&restoreDriverEnable, "driver", "d", false, "Enable restore from backup driver")
}
