package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

func UploadImageToBB(c echo.Context, paramName string, maxSize int64) (interface{}, error) {
	apiURL := "https://api.imgbb.com/1/upload"
	apiKey := "537a6a33fa7a7152e0191612fc3b99e5"

	maxFileSize := maxSize * 1024 * 1024
	file, err := c.FormFile(paramName)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to get image: " + err.Error()))
	}
	if file.Size > maxFileSize {
		return "", errors.New(fmt.Sprintf("File size exceeds %vMB limit", maxSize))
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		return "", errors.New(fmt.Sprintf("Invalid file type"))
	}

	src, err := file.Open()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to open image: " + err.Error()))
	}
	defer src.Close()
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	part, err := writer.CreateFormFile("image", file.Filename)
	if err != nil {
		return "", err
	}

	// Copy the file content to the multipart form
	if _, err := io.Copy(part, src); err != nil {
		return "", err
	}

	// Close the multipart writer to finalize the form
	if err := writer.Close(); err != nil {
		return "", err
	}
	writer.WriteField("key", apiKey)
	writer.Close()
	req, err := http.NewRequest("POST", apiURL, &b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if errs := json.Unmarshal(respBytes, &result); errs != nil {
		return "", errs
	}
	log.Println("res: ", result)
	if d, ok := result["data"]; ok {
		if newData, oke := d.(map[string]interface{}); oke {
			return strings.ReplaceAll(newData["url"].(string), "i.ibb.co", "i.ibb.co.com"), nil
		}
	}
	return "", nil
}
