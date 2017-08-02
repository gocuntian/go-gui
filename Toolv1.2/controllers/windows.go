package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"time"

	"errors"

	"path/filepath"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

const URL = ""
const PathSeparator = `/`
const STATUS_CODE = 200

var CurCookies []*http.Cookie
var CurCookieJar *cookiejar.Jar

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
type UserData struct {
	Id           int32  `json:"id"`
	GuestName    string `json:"guest_name"`
	Mobile       string `json:"mobile"`
	Avatar       string `json:"avatar"`
	ConferenceId string `json:"conference_id"`
	GradeId      string `json:"grade_id"`
	UpdatedAt    string `json:"update_at"`
	CreatedAt    string `json:"create_at"`
}
type MsgLogin struct {
	Msg
	Data []LoginData `json:"data"`
}

type MsgUser struct {
	Msg
	Data []UserData `json:"data"`
}

type MsgRet struct {
	ErrCount int
	SucCount int
}

func init() {
	CurCookies = nil
	CurCookieJar, _ = cookiejar.New(nil)
}

func SetEventHandler(w *window.Window) {
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
		client := http.Client{
			CheckRedirect: nil,
			Jar:           CurCookieJar,
		}
		req, err := http.NewRequest("POST", URL+"/sm-login", strings.NewReader("email="+account.Account+"&password="+account.Password))
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
			CurCookies = CurCookieJar.Cookies(req.URL)
			w.LoadFile("views/index.html")
		}
		return sciter.NullValue()
	})

	w.DefineFunction("conference", func(args ...*sciter.Value) *sciter.Value {
		httpClient := &http.Client{
			Jar: CurCookieJar,
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
		conferenceID := args[2].String()
		excelFile := args[0].String()
		excelFile = strings.TrimLeft(excelFile, "file://")
		fields := []string{"avatar", "guest_name", "mobile", "guest_email", "company_name", "position", "hint"}
		ext := strings.ToLower(filepath.Ext(excelFile))
		if ext != ".csv" && ext != ".xlsx" {
			return sciter.NewValue("文件格式不允许")
		}
		ClearMsg(w, "result", "list")
		ClearMsg(w, "errmsg", "errlist")
		avatarFolder := args[1].String()
		ch := make(chan FileMap, 10)
		dir, err := ioutil.ReadDir(avatarFolder)
		if err != nil {
			MsgLog(w, 500, err.Error())
		}
		start := time.Now()
		var errcount int
		var suscount int
		var i int
		errch := make(chan string, 10)
		susch := make(chan string, 10)
		var fmap map[string]string
		tdb := make(map[string]string, len(dir))
		for _, fi := range dir {
			if fi.IsDir() {
				continue
			} else {
				tdb[fi.Name()] = fi.Name()
			}
		}
		if ext == ".csv" {
			fmap = ReadCVS(w, fields, excelFile, tdb)
		} else {
			fmap = ReadXLSX(w, fields, excelFile, tdb)
		}
		total := len(fmap)

		go AvatarMap(w, dir, avatarFolder, ch, fmap, errch)
		go SyncDataToHTTP(w, ch, conferenceID, fields, fmap, susch, errch)
		successclosed, errorclosed := false, false
		for {
			i++
			if i > 13 {
				ClearMsg(w, "result", "list")
				i = 0
			}
			if successclosed && errorclosed {
				break
			}
			select {
			case susone, ok := <-susch:
				if !ok {
					successclosed = true
				} else {
					suscount++
					MsgLog(w, STATUS_CODE, "["+susone+"]->已完成"+strconv.FormatFloat(float64(suscount)/float64(total)*100, 'f', 2, 64)+"%")
				}
			case errone, oks := <-errch:
				if !oks {
					errorclosed = true
				} else {
					MsgLog(w, 500, errone)
					errcount++
				}
			}
		}
		dis := time.Now().Sub(start).Seconds()
		return sciter.NewValue("同步成功:" + strconv.Itoa(suscount) + ", 失败:" + strconv.Itoa(errcount) + "总耗时:" + strconv.FormatFloat(dis, 'f', 5, 64) + "秒")
	})

}
