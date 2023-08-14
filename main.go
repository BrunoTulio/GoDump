/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import "github.com/BrunoTulio/GoDump/cmd"

func main() {
	cmd.Execute()
}

//	signal.Notify(globalOSSignalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
//<-globalOSSignalCh
/*
	// fileSh := fmt.Sprintf(".%s/backup.sh", dir())
	// cmdArgs := []string{"exec", "-it", "consiga_electoral_db", "bash", "-c", fileSh}




}

func dir() string {
	dir, _ := os.Getwd()
	return dir
}

*/
