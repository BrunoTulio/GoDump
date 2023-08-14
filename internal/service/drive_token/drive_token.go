package drive_token

import (
	"context"
	"encoding/json"

	"fmt"

	"os"

	"github.com/BrunoTulio/GoDump/internal/constants"
	"github.com/BrunoTulio/GoDump/pkg/folder"
	"github.com/BrunoTulio/GoDump/pkg/logger"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Constants

type (
	DriveTokenService interface {
		Save(ctx context.Context) error
	}

	driveTokenService struct {
		config *oauth2.Config
	}
)

func (s *driveTokenService) Save(ctx context.Context) error {
	authURL := s.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	fmt.Println("Paste Authorization code here :")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return fmt.Errorf("Unable to read authorization code %v", err)
	}

	token, err := s.config.Exchange(ctx, authCode)
	if err != nil {
		return errors.Wrap(err, "Unable to retrieve token from web")
	}

	err = storeToken(token)

	if err != nil {
		return errors.Wrap(err, "Error save token.json")
	}

	return nil
}

func storeToken(token *oauth2.Token) error {

	err := folder.Create(constants.PathToken)

	if err != nil {
		return err
	}

	logger.Infof("Saving credential file to: %s\n", constants.PathGoogleDriveToken)
	f, err := os.OpenFile(constants.PathGoogleDriveToken, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

func New(config *oauth2.Config) DriveTokenService {
	return &driveTokenService{
		config,
	}
}
