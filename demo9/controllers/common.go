package controllers

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
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

//Des 加密
func DesEncrypt(data, key []byte) (string, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	data = PKCS5Padding(data, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(data))
	blockMode.CryptBlocks(crypted, data)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

//解密
func DesDecrypt(crypted string, key []byte) (string, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(crypted)
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(data))
	blockMode.CryptBlocks(origData, data)
	origData = PKCS5UnPadding(origData)
	return string(origData), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
