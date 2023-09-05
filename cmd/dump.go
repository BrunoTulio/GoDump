/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/BrunoTulio/GoDump/internal/domain"
	"github.com/BrunoTulio/GoDump/internal/factory"
	"github.com/spf13/cobra"
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "d",
	Long:  `Backup manual.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := factory.MakeLogger()
		backupUseCase := factory.MakeDefaultBackupUseCase()

		key, _ := cmd.Flags().GetString("key")
		typeFlag, _ := cmd.Flags().GetString("type")

		var fileType *domain.Type
		if typeFlag != "" {
			t, _ := domain.ParseFileType(typeFlag)
			fileType = &t
		}

		if key != "" {

			err := backupUseCase.GenerateByKey(fileType, key)

			if err != nil {
				logger.Fatal(err)
			}

			logger.Info("Finalizado backup")
			return
		}

		err := backupUseCase.Generate(fileType)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("Finalizado backup")

	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)

	dumpCmd.Flags().StringP("key", "k", "", "Backup key bump config")
	dumpCmd.Flags().StringP("type", "t", "", "Backup type file")
}
