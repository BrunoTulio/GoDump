/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/BrunoTulio/GoDump/internal/factory"
	"github.com/BrunoTulio/GoDump/pkg/logger"
	"github.com/spf13/cobra"
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "d",
	Long:  `Backup manual.`,
	Run: func(cmd *cobra.Command, args []string) {
		factory.MakeLogger()
		backupUseCase := factory.MakeDefaultBackupUseCase()

		key, _ := cmd.Flags().GetString("key")

		if key != "" {

			err := backupUseCase.GenerateByKey(key)

			if err != nil {
				logger.Fatal(err)
			}

			logger.Info("Finalizado backup")
			return
		}

		err := backupUseCase.Generate()
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("Finalizado backup")

	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)

	dumpCmd.Flags().StringP("key", "k", "", "Backup key bump config")
}
