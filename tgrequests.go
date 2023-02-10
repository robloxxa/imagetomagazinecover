package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TelegramGetFileRequest struct {
	Ok     bool         `json:"ok"`
	Result TelegramFile `json:"result"`
}

type TelegramFile struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	FilePath     string `json:"file_path"`
}

var client = &http.Client{}

func getFileURL(fileID string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", botToken, fileID)
}

func getFileLinkURL(filePath string) string {
	return fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", botToken, filePath)
}

func requestFilePath(fileID string) (string, error) {
	r, err := client.Get(getFileURL(fileID))
	if err != nil {
		return "", err
	}
	var file TelegramGetFileRequest
	err = json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	return file.Result.FilePath, nil
}

func getLinkToPhoto(fileID string) (string, error) {
	path, err := requestFilePath(fileID)
	if err != nil {
		return "", err
	}
	return getFileLinkURL(path), nil
}
