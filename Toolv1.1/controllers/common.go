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
	"time"

	sciter "github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

var db map[string]string

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
	httpReq, _ = http.NewRequest("POST", URL+"/sm-guest/import", strings.NewReader(data))
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return "", err
	}
	var msguser MsgUser
	err = json.Unmarshal([]byte(body), &msguser)
	if err != nil {
		return "", err
	}
	if STATUS_CODE != msguser.Msg.StatusCode {
		return "", errors.New(msguser.Msg.Message)
	} else {
		return msguser.Data[0].GuestName, nil
	}
}

func SyncHTTP(w *window.Window, data string) {
	httpClient := &http.Client{
		Jar: CurCookieJar,
	}
	var httpReq *http.Request
	httpReq, _ = http.NewRequest("POST", URL+"/sm-guest/import", strings.NewReader(data))
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		MsgLog(w, 500, err.Error())
		return
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		MsgLog(w, 500, err.Error())
		return
	}
	var msguser MsgUser
	err = json.Unmarshal([]byte(body), &msguser)
	if err != nil {
		MsgLog(w, 500, err.Error())
		return
	}
	if STATUS_CODE != msguser.Msg.StatusCode {
		err = errors.New(msguser.Msg.Message)
		MsgLog(w, 500, err.Error())
		return
	} else {
		MsgLog(w, STATUS_CODE, " 嘉宾:("+msguser.Data[0].GuestName+")---同步成功")
		return
	}
}

func AvatarMap(w *window.Window, folderPath string) {
	db = make(map[string]string)
	dir, err := ioutil.ReadDir(folderPath)
	if err != nil {
		MsgLog(w, 500, err.Error())
	}
	wg.Add(len(dir))
	for _, fi := range dir {
		wg.Done()
		if fi.IsDir() {
			continue
		} else {
			keyImg, err := UploadFile(URL+"/upload", folderPath+PathSeparator+fi.Name())
			if err != nil {
				MsgLog(w, 500, err.Error())
			}
			db[fi.Name()] = keyImg
		}
	}
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
	return string(crypted), nil
}
func TripleDesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	fmt.Println(block)
	if err != nil {
		return nil, err
	}
	fmt.Println("=====================")
	//origData = PKCS5Padding(origData, block.BlockSize())
	origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
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

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

func NowTime() string {
	return time.Now().Format("2006/01/02 15:04:05.9999")
}

func AppendMsg(wd *window.Window, msg string) error {
	root, err := wd.GetRootElement()
	if err != nil {
		return err
	}

	resultElement, err := root.SelectById("result")
	if err != nil {
		return err
	}
	err = resultElement.SetHtml(msg, sciter.SIH_APPEND_AFTER_LAST)
	if err != nil {
		return err
	}
	root.Update(true)
	return nil
}

func ClearMsg(wd *window.Window) error {
	root, err := wd.GetRootElement()
	if err != nil {
		return err
	}

	resultElement, err := root.SelectById("result")
	if err != nil {
		return err
	}
	err = resultElement.SetHtml("<div id=\"result\" class=\"list\"></div>", sciter.SOH_REPLACE)
	if err != nil {
		return err
	}
	root.Update(true)
	return nil
}

func MsgLog(wd *window.Window, code int, msg string) {
	if code == STATUS_CODE {
		//fmt.Println("=====", msg)
		AppendMsg(wd, "<div style=\"color:#FF8C00\">"+NowTime()+"  [正在进行中...] 内容:["+msg+"]</div>")
	} else {
		//fmt.Println("=====", msg)
		AppendMsg(wd, "<div style=\"color:#FF0000\">"+NowTime()+"  [错误日志]  内容：["+msg+"]</div>")
	}
}
