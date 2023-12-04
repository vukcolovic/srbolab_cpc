package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"srbolab_cpc/logoped"
	"strings"
)

var (
	FileService fileServiceInterface = &fileService{}
	RootPath    string
)

type fileService struct {
}

type fileServiceInterface interface {
	WriteFile(path string, content string) error
	GetFile(folder, clientID, filename string) (string, error)
	DeleteFile(filepath string) error
	GetPath(folder, clientID, filename string) string
}

func (f *fileService) WriteFile(path string, content string) error {
	idx := strings.Index(content, ";base64,")
	if idx < 0 {
		logoped.ErrorLog.Println("Error writting file, unexpected format")
		return errors.New("unexpected format")
	}
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(content[idx+8:]))
	buff := bytes.Buffer{}
	_, err := buff.ReadFrom(reader)
	if err != nil {
		logoped.ErrorLog.Println("Error writting file, error: ", err.Error())
		return err
	}
	err = os.WriteFile(RootPath+path, buff.Bytes(), 0644)
	if err != nil {
		logoped.ErrorLog.Println("Error writting file, error: ", err.Error())
		return err
	}

	return nil
}

func (f *fileService) GetFile(folder, clientID, filename string) (string, error) {
	b, err := os.ReadFile(RootPath + f.GetPath(folder, clientID, filename))
	if err != nil {
		logoped.ErrorLog.Println("Error getting file, error: ", err.Error())
		return "", err
	}

	encodedContent := base64.StdEncoding.EncodeToString(b)

	return encodedContent, err
}

func (f *fileService) DeleteFile(filepath string) error {
	err := os.Remove(RootPath + fmt.Sprintf("%s", filepath))
	if err != nil {
		logoped.ErrorLog.Println("Error deleting file, error: ", err.Error())
		return err
	}

	return nil
}

func (f *fileService) GetPath(folder, clientID, filename string) string {
	return fmt.Sprintf("/%s/%s_%s", folder, clientID, filename)
}
