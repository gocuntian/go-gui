package controllers

import (
	"encoding/csv"
	"log"
	"os"
	"sync"

	"io"

	"errors"

	"github.com/sciter-sdk/go-sciter/window"
	"github.com/tealeg/xlsx"
)

var wg sync.WaitGroup
var wgp sync.WaitGroup

func ReadXlsx(fields []string, file string) (map[string]string, error) {
	data := make(map[string]string)
	xlFile, err := xlsx.OpenFile(file)
	if err != nil {
		return data, err
	}
	var rowstr string
	var key string
	for _, sheet := range xlFile.Sheets {
		for k, row := range sheet.Rows {
			if k > 0 {
				for i, cell := range row.Cells {
					if i > 0 && i <= len(fields) {
						rowstr += fields[i] + "=" + cell.String() + "&"

					}
					if i == 0 {
						key = cell.String()
					}
				}
				data[key] = rowstr
			}
		}
	}
	return data, nil
}

func ReadCVS(fields []string, cvsFile string) (map[string]string, error) {
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
				if i > 0 && i <= len(fields) {
					rowstr += fields[i] + "=" + cell + "&"

				}
				if i == 0 {
					key = cell
				}
			}
			data[key] = rowstr
		}
	}
	return data, nil
}

func ReadCVSAndSync(w *window.Window, fields []string, cvsFile, avatarFolder, conferenceID string) {

	go AvatarMap(w, avatarFolder)
	log.Println("start wait......")
	wg.Wait()
	go func(w *window.Window, fields []string, cvsFile, conferenceID string) {
		var rowstr string
		var key string
		var k int
		file, err := os.Open(cvsFile)
		if err != nil {
			MsgLog(w, 500, err.Error())
		}
		defer file.Close()
		csvr := csv.NewReader(file)
		for {
			rowstr = ""
			key = ""
			row, err := csvr.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				MsgLog(w, 500, err.Error())
				continue
			}
			if k > 0 {
				for i, cell := range row {
					if i > 0 && i <= len(fields) {
						rowstr += fields[i] + "=" + cell + "&"

					}
					if i == 0 {
						key = cell
					}
				}
				keyImg, ok := db[key]
				err = errors.New("头像不存在")
				if !ok {
					MsgLog(w, 500, err.Error())
					continue
				}
				rowstr += fields[0] + "=" + keyImg + "&conference_id=" + conferenceID + "&grade_id=1"
				SyncHTTP(w, rowstr)

			}
			k++
		}
	}(w, fields, cvsFile, conferenceID)

	wg.Wait()
	return
}

func CheckFileStat(file string) (string, error) {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("File does not exist.")
			return "", err
		}
		return "", err
	}
	return file, nil
}
