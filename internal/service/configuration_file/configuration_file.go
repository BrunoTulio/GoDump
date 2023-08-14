package configuration_file

import (
	"os"
	"time"

	"github.com/BrunoTulio/GoDump/internal/constants"
	"github.com/BrunoTulio/GoDump/internal/domain"
	"github.com/BrunoTulio/GoDump/pkg/folder"
	"github.com/pkg/errors"

	"gopkg.in/yaml.v3"
)

type (
	ConfigurationFileService interface {
		Create() error
		Get() (domain.Config, error)
	}

	configurationFileService struct {
	}
)

func (s *configurationFileService) Create() error {

	err := folder.Create(constants.PathConfig)

	if err != nil {
		return err
	}

	config := domain.Config{
		Cron:                 "@every 1h30m",
		DurationFileInFolder: time.Duration(24 * time.Hour * 2),

		GoogleDrive: domain.GoogleDrive{
			Enable: false,
		},
		Alert: domain.Alert{
			SendMail: "sendtomail@mail.com",
			Enable:   true,
			Mail: domain.Mail{
				User:           "user@example.com",
				From:           "user@example.com",
				Password:       "password",
				Host:           "host.example.com",
				Port:           3500,
				InsureSecurity: false,
			},
		},
		Dumps: domain.Dumps{
			domain.Dump{
				Key:      "unique",
				Host:     "localhost",
				Port:     5432,
				Username: "username",
				Password: "password",
				Database: "database",
				Priority: domain.Normal,
			},
		},
	}

	yamlData, err := yaml.Marshal(&config)

	if err != nil {
		return errors.Wrap(err, "Error while Marshaling")
	}

	ff, err := os.Create(constants.PathConfigFile)
	if err != nil {
		return errors.Wrap(err, "Error  creating file config")
	}
	defer ff.Close()

	_, err = ff.WriteString(string(yamlData))

	if err != nil {
		return errors.Wrap(err, "Error  writing file config")
	}

	return nil
}

// Get implements ConfigurationFileService.
func (g *configurationFileService) Get() (domain.Config, error) {

	data, err := os.ReadFile(constants.PathConfigFile)

	if err != nil {
		return domain.Config{}, errors.Wrap(err, "Erro read configuration file")
	}

	config := domain.Config{}

	err = yaml.Unmarshal([]byte(data), &config)

	if err != nil {
		return domain.Config{}, errors.Wrap(err, "Erro yaml marshal")
	}

	return config, nil
}

func New() ConfigurationFileService {
	return &configurationFileService{}
}
