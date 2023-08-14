package factory

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/BrunoTulio/GoDump/internal/constants"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func MakeOAuth2Config() (*oauth2.Config, error) {
	b, err := os.ReadFile(constants.PathGoogleDriveCredentialAccount)
	if err != nil {
		log.Fatalf("Unable to read credentials.json file. Err: %v\n", err)
	}

	return google.ConfigFromJSON(b, constants.Scope)
}

func MakeGoogleDriveService() (*drive.Service, error) {
	token, err := tokenFromFile()

	if err != nil {
		return nil, err
	}

	config, err := MakeOAuth2Config()

	if err != nil {
		return nil, errors.Wrap(err, "Config Create oauth2.\n")
	}

	client := config.Client(context.Background(), token)

	service, err := drive.NewService(context.Background(), option.WithHTTPClient(client))

	if err != nil {
		return nil, errors.Wrap(err, "Create service drive.\n")
	}

	return service, nil
}

func tokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open(constants.PathGoogleDriveToken)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}
