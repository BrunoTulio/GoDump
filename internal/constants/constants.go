package constants

import "google.golang.org/api/drive/v2"

const (
	PathGoogleDrive                  = "./.google_drive"
	PathGoogleDriveCredentialAccount = PathGoogleDrive + "/credentials.json"
	PathToken                        = "./token"
	PathGoogleDriveToken             = PathToken + "/token.json"
	PathConfig                       = "./config"
	PathConfigFile                   = PathConfig + "/config.yaml"
	PathBackup                       = "./backup"

	Scope = drive.DriveScope

	LayoutDate = "2006-01-02T15:04:05-0700"
)
