package service

import (
	"archive/zip"
	"bytes"
	"context"
	"github.com/apulis/bmod/aistudio-aom/internal/utils"
	"github.com/apulis/sdk/go-utils/logging"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	certConfigmap = "cert-cm"
)

func DownloadCerts() ([]byte, *utils.CodeMessage) {
	ctx := context.Background()
	cm, err := clientset.CoreV1().ConfigMaps("default").Get(ctx, certConfigmap, metav1.GetOptions{})
	if err != nil {
		utils.ErrorConfigmapOp.Message = err.Error()
		return nil, utils.ErrorConfigmapOp
	}

	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)

	for fn, content := range cm.Data {
		entry, err := zipWriter.Create(fn)
		if err != nil {
			utils.ErrZipEntryError.Message = err.Error()
			return nil, utils.ErrZipEntryError
		}
		if _, err := entry.Write([]byte(content)); err != nil {
			utils.ErrZipWriteError.Message = err.Error()
			return nil, utils.ErrZipWriteError
		}
	}

	if err := zipWriter.Close(); err != nil {
		logging.Error(err).Msg("zip writer close error")
		utils.ErrZipWriteError.Message = err.Error()
		return nil, utils.ErrZipWriteError
	}

	return zipBuffer.Bytes(), nil
}
