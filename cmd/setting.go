/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/BrunoTulio/GoDump/internal/factory"
	"github.com/BrunoTulio/GoDump/pkg/logger"
	"github.com/spf13/cobra"
)

// settingCmd
var (
	isInitDriveConfiguration bool
)

var settingCmd = &cobra.Command{
	Use:   "setting",
	Short: "s",
	Long:  `Create configuration app system backup`,

	Run: func(cmd *cobra.Command, args []string) {

		if isInitDriveConfiguration {
			settingDriverService := factory.MakeDriveTokenService()

			err := settingDriverService.Save(context.Background())

			if err != nil {
				logger.Errorf("Token drive google, Err: %v\n", err)
			}

			return
		}
		createConfigurationService := factory.MakeConfigurationFileService()
		err := createConfigurationService.Create()

		if err != nil {
			logger.Errorf("Generate file configuration, Err: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(settingCmd)
	settingCmd.Flags().BoolVarP(&isInitDriveConfiguration, "drive", "d", false, "Initialize the google drive configuration token")
}
