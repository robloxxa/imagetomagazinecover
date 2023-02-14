package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type TelegramGetFileRequest struct {
	Ok          bool         `json:"ok"`
	Description string       `json:"description,omitempty"`
	Result      TelegramFile `json:"result,omitempty"`
}

type TelegramFile struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	FilePath     string `json:"file_path"`
}

var client = &http.Client{}

func formatGetFileUrl(r MagazineRequest) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", r.Token, r.FileId)
}

func formatTelegramDownloadLink(r MagazineRequest, filePath string) string {
	return fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", r.Token, filePath)
}

func requestFilePath(mreq MagazineRequest) (string, error) {
	r, err := client.Get(formatGetFileUrl(mreq))
	if err != nil {
		return "", err
	}
	var file TelegramGetFileRequest
	err = json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		return "", err
	}
	if !file.Ok {
		return "", errors.New(file.Description)
	}
	defer r.Body.Close()

	return file.Result.FilePath, nil
}

func getLinkToPhoto(r MagazineRequest) (string, error) {
	path, err := requestFilePath(r)
	if err != nil {
		return "", err
	}
	return formatTelegramDownloadLink(r, path), nil
}
