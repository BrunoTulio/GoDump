package drive

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/BrunoTulio/GoDump/internal/constants"
	"github.com/BrunoTulio/GoDump/internal/domain"
	"github.com/BrunoTulio/GoDump/pkg/logger"

	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"
)

type (
	FileInfo struct {
		ID   string
		Size int64
		Name string
	}

	DriveService interface {
		Upload(ctx context.Context, pathFile string, dump domain.Dump) error
		Get(ctx context.Context, dump domain.Dump) (*FileInfo, error)
		Download(ctx context.Context, dump domain.Dump) (string, error)
	}

	googleDriveService struct {
		service *drive.Service
	}
)

// Download implements DriveService.
func (g *googleDriveService) Download(ctx context.Context, dump domain.Dump) (string, error) {
	f, err := g.Get(ctx, dump)

	if err != nil {
		return "", err
	}

	res, err := g.service.Files.Get(f.ID).Download()

	if err != nil {
		return "", errors.Wrap(err, "failed to download google drive service")
	}

	defer res.Body.Close()
	fileName := path.Join(constants.PathGoogleDrive, nameFile(dump))

	file, err := os.Create(fileName)
	if err != nil {
		return "", errors.Wrap(err, "create file temp for google drive service")
	}
	defer file.Close()
	_, err = io.Copy(file, res.Body)

	if err != nil {
		return "", errors.Wrap(err, "failed to cop on google drive service")
	}

	return fileName, nil
}

// Get implements DriveService.
func (g *googleDriveService) Get(ctx context.Context, dump domain.Dump) (*FileInfo, error) {
	f, err := g.getFile(dump)

	if err != nil {
		return nil, err
	}

	return &FileInfo{
		ID:   f.Id,
		Name: f.Name,
		Size: f.Size,
	}, nil
}

// Upload implements UpdateDriveService.
func (u *googleDriveService) Upload(ctx context.Context, pathFile string, dump domain.Dump) error {

	logger.Infof("Initialize upload drive service path file: %s\n", pathFile)

	ff, err := os.Open(pathFile)

	if err != nil {
		return errors.Wrap(err, "Cannot open file")
	}

	defer ff.Close()

	logger.Info("Verificando arquivo no mydrive")

	fileDrive, err := u.getFile(dump)

	if err != nil {
		return err
	}

	_, err = u.createOrUpdateFile(ff, dump, fileDrive)

	if err != nil {
		return err
	}

	return nil

}

func (u *googleDriveService) getFile(dump domain.Dump) (*drive.File, error) {

	files, err := u.service.Files.List().Q(fmt.Sprintf(`name="%s"`, nameFile(dump))).Do()

	if err != nil {
		return nil, errors.Wrap(err, "Get file google drive failed")
	}

	if len(files.Files) > 0 {
		return files.Files[0], nil
	}

	return nil, nil
}

func nameFile(dump domain.Dump) string {
	return fmt.Sprintf("%s_%s.backup", dump.Key, dump.Database)
}

func (u *googleDriveService) createOrUpdateFile(ff io.Reader, dump domain.Dump, fileDrive *drive.File) (string, error) {

	if fileDrive != nil {
		logger.Info("Inicializando atualização do arquivo no mydrive")

		_, err := u.service.Files.Update(fileDrive.Id, &drive.File{
			MimeType: "application/octet-stream",
			Name:     nameFile(dump),
			DriveId:  dump.Key,
		}).Media(ff).ProgressUpdater(func(current, total int64) {
			logger.Infof("Progresso %d no total de %d", current, total)
		}).Do()

		if err != nil {
			return "", errors.Wrap(err, "Insert file failed google drive")
		}
		logger.Info("Atualizado arquivo no mydrive com sucesso")

		return fileDrive.Id, nil
	}
	logger.Info("Inicializando criação do arquivo no mydrive")

	file, err := u.service.Files.Create(&drive.File{
		MimeType: "application/octet-stream",
		Name:     nameFile(dump),
		DriveId:  dump.Key,
	}).Media(ff).ProgressUpdater(func(current, total int64) {
		logger.Infof("Progresso %d no total de %d", current, total)
	}).Do()

	if err != nil {
		return "", errors.Wrap(err, "Insert file failed google drive")
	}
	logger.Info("Arquivo criado no mydrive com sucesso")

	return file.Id, nil
}

func New(service *drive.Service) DriveService {
	return &googleDriveService{service}
}
