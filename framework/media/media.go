package media

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/DmitryVesenniy/go-rest-framework/common"
)

type MediaServise struct {
	MediaPath string
}

type MediaEntity struct {
	FileName string
	Path     string
}

func (mediaEntity *MediaEntity) Url() string {
	return filepath.Join(mediaEntity.Path, mediaEntity.FileName)
}

func (media *MediaServise) SaveMediaFile(fileName string, content []byte, userPath string) (MediaEntity, error) {
	result := MediaEntity{
		FileName: fileName,
		Path:     media.MediaPath,
	}

	if userPath != "" {
		filePath := filepath.Join(media.MediaPath, userPath)
		result.Path = filePath
	}

	if isExists, _ := common.Exists(result.Url()); isExists {
		uid := uuid.New()
		fileNameSplitted := strings.Split(fileName, ".")
		fileNameSplitted[0] = fmt.Sprintf("%s_%s", fileNameSplitted[0], uid)
		result.FileName = strings.Join(fileNameSplitted, ".")
	}

	f, err := os.Create(result.Url())
	if err != nil {
		return result, err
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		return result, err
	}

	err = f.Sync()
	if err != nil {
		return result, err
	}

	return result, nil
}

func (media *MediaServise) DeleteMediaFile(url string) error {
	return os.Remove(url)
}
