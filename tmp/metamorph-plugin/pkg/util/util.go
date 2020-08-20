package util

import (
	"net/http"
	"go.uber.org/zap"
	"os"
	"io"
	"github.com/manojkva/metamorph-plugin/pkg/logger"
)

func DownloadUrl(filepath string, url string) error {
	logger.Log.Info("DownloadUrl()", zap.String("filepath", filepath), zap.String("url", url))
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		logger.Log.Error("Failed to download..",
			zap.String("url", url),
			zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		logger.Log.Error("Failed to create file", zap.String("filepath", filepath))
		return err
	}
	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
