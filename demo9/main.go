package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	"errors"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

const URL = "https://apipre.bi.sensetime.com"
const PathSeparator = `/`
const STATUS_CODE = 200

var gCurCookies []*http.Cookie
var gCurCookieJar *cookiejar.Jar
var tdb map[string]string
var admin_id int32

type Account struct {
	Account  string
	Password string
}

type Msg struct {
	StatusCode int32  `json:"status_code"`
	Message    string `json:"message"`
}

type UploadData struct {
	ImageKey string `json:"image_key"`
}

type MsgUpload struct {
	Msg
	Data []UploadData `json:"data"`
}

type LoginData struct {
	Id    int32  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
type MsgLogin struct {
	Msg
	Data []LoginData `json:"data"`
}

func init() {
	gCurCookies = nil
	gCurCookieJar, _ = cookiejar.New(nil)
	tdb = make(map[string]string)
	admin_id = 0
}

func userLogin(w *window.Window) {
	//注册dump函数方便在tis脚本中打印数据
	w.DefineFunction("dump", func(args ...*sciter.Value) *sciter.Value {
		for _, v := range args {
			fmt.Print(v.String() + " ")
		}
		fmt.Println()
		return sciter.NullValue()
	})

	//login函数，用于用户登录逻辑，这里只是简单的把数据打印出来
	w.DefineFunction("login", func(args ...*sciter.Value) *sciter.Value {
		reqJson := args[0].String()
		var account Account
		err := json.Unmarshal([]byte(reqJson), &account)
		if err != nil {
			return sciter.NewValue(err.Error())
		}

		reqUrl := URL + "/sm-login"
		client := http.Client{
			CheckRedirect: nil,
			Jar:           gCurCookieJar,
		}
		req, err := http.NewRequest("POST", reqUrl, strings.NewReader("email="+account.Account+"&password="+account.Password))
		if err != nil {
			return sciter.NewValue(err.Error())
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		if err != nil {
			return sciter.NewValue(err.Error())
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return sciter.NewValue(err.Error())
		}
		var msglogin MsgLogin
		err = json.Unmarshal([]byte(body), &msglogin)
		if err != nil {
			return sciter.NewValue(err.Error())
		}
		if msglogin.Msg.StatusCode != STATUS_CODE {
			err = errors.New(msglogin.Msg.Message)
			return sciter.NewValue(err.Error())
		} else {
			gCurCookies = gCurCookieJar.Cookies(req.URL)
			w.LoadFile("view/index.html")
		}
		return sciter.NullValue()
	})

	w.DefineFunction("conference", func(args ...*sciter.Value) *sciter.Value {
		httpClient := &http.Client{
			Jar: gCurCookieJar,
		}
		var httpReq *http.Request
		httpReq, _ = http.NewRequest("GET", URL+"/sm-conference?per_page=100", nil)
		httpResp, err := httpClient.Do(httpReq)
		if err != nil {
			fmt.Println(err)
		}
		defer httpResp.Body.Close()
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(string(body))
		return sciter.NewValue(string(body))
	})

	//上传信息
	w.DefineFunction("upload", func(args ...*sciter.Value) *sciter.Value {
		conferenceId := args[2].String()
		fmt.Println(conferenceId)
		excelFile := args[0].String()
		excelFile = strings.TrimLeft(excelFile, "file://")
		fields := []string{"avatar", "guest_name", "mobile"}
		cellsMap, err := readCVS(fields, excelFile)
		if err != nil {
			return sciter.NewValue(err.Error())
		}
		//fmt.Println(cellsMap)
		folderPath := args[1].String()
		dir, err := ioutil.ReadDir(folderPath)
		if err != nil {
			return sciter.NewValue(err.Error())
		}
		var mapString string
		//fmt.Println(len(dir))
		for _, fi := range dir {
			if fi.IsDir() {
				continue
			} else {
				keyImg, err := uploadFile(URL+"/upload", folderPath+PathSeparator+fi.Name())
				if err != nil {
					return sciter.NewValue(err.Error())
				}
				mapString = cellsMap[fi.Name()] + fields[0] + "=" + keyImg + "&conference_id=" + conferenceId + "&grade_id=1"
				fmt.Println(mapString)
				retstr, _ := insert(mapString)
				fmt.Println(retstr)
			}

		}
		//fmt.Println(cellsMap)
		return sciter.NewValue("ok")
	})

}

func insert(data string) (string, error) {
	httpClient := &http.Client{
		Jar: gCurCookieJar,
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

func readCVS(fields []string, cvsFile string) (map[string]string, error) {
	data := make(map[string]string)
	file, err := os.Open(cvsFile)
	if err != nil {
		return data, err
	}
	defer file.Close()
	csvr := csv.NewReader(file)
	d, err := csvr.ReadAll()
	if err != nil {
		return data, err
	}
	var rowstr string
	var key string
	for k, row := range d {
		rowstr = ""
		key = ""
		if k > 0 {
			for i, cell := range row {
				if i > 0 {
					rowstr += fields[i] + "=" + cell + "&"
				} else {
					key = cell
				}
			}
			data[key] = rowstr
		}
	}
	return data, nil
}

func uploadFile(url, file string) (string, error) {
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

func main() {
	w, err := window.New(sciter.DefaultWindowCreateFlag, &sciter.Rect{448, 184, 1000, 540})
	if err != nil {
		log.Fatal(err)
	}
	w.LoadFile("view/login.html")
	w.SetTitle("用户登录")
	userLogin(w)
	w.Show()
	fmt.Println("start run")
	w.Run()
	fmt.Println("end run !!!")
}
