/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/BrunoTulio/GoDump/internal/constants"
	"github.com/BrunoTulio/GoDump/internal/factory"
	"github.com/BrunoTulio/GoDump/pkg/logger"

	"github.com/spf13/cobra"
)

// autoCmd represents the dump command
var autoCmd = &cobra.Command{
	Use:   "auto",
	Short: "a",
	Long:  `Execute backup cron starting`,
	Run: func(cmd *cobra.Command, args []string) {
		factory.MakeLogger()
		c := factory.MakeCron()
		defer c.Stop()
		backupUseCase := factory.MakeDefaultBackupAutoUseCase(c)
		backupUseCase.Generate(context.Background())
		go c.Start()

		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt, os.Kill)
		<-sig
		logger.Info(fmt.Sprintf("%s closed", time.Now().Format(constants.LayoutDate)))
	},
}

func init() {
	rootCmd.AddCommand(autoCmd)
}
