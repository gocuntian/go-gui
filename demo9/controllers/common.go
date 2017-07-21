package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func UploadFile(url, file string) (string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("image", file)
	if err != nil {
		fmt.Println("error writing to buffer")
		return "", err
	}
	fh, err := os.Open(file)
	if err != nil {
		fmt.Println("error opening file")
		return "", err
	}
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return "", err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var msgupload MsgUpload
	err = json.Unmarshal([]byte(data), &msgupload)
	if err != nil {
		return "", err
	}
	if STATUS_CODE != msgupload.Msg.StatusCode {
		return "", errors.New(msgupload.Msg.Message)
	} else {
		return msgupload.Data[0].ImageKey, nil
	}

}

func SyncData(data string) (string, error) {
	httpClient := &http.Client{
		Jar: CurCookieJar,
	}
	var httpReq *http.Request
	httpReq, _ = http.NewRequest("POST", URL+"/sm-guest/create", strings.NewReader(data))
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(string(body))
	return string(body), nil
}
