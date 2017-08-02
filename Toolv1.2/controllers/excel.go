package controllers

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/sciter-sdk/go-sciter/window"
	"github.com/tealeg/xlsx"
)

type FileMap struct {
	Key  string
	Name string
}

func ReadXLSX(w *window.Window, fields []string, file string, tdb map[string]string) map[string]string {
	xlFile, err := xlsx.OpenFile(file)
	if err != nil {
		MsgLog(w, 500, err.Error())
	}
	db := make(map[string]string)
	var rowstr string
	var key string
	for _, sheet := range xlFile.Sheets {
		for k, row := range sheet.Rows {
			rowstr = ""
			key = ""
			if k > 0 {
				for i, cell := range row.Cells {
					if i > 0 && i <= len(fields) {
						rowstr += fields[i] + "=" + cell.String() + "&"

					}
					if i == 0 {
						key = cell.String()
					}
				}
				if tdb[key] == key {
					db[key] = rowstr
					delete(tdb, key)
				} else {
					//errch <- key + "不存在"
					MsgLog(w, 500, key+"不存在")
				}
			}
		}
	}
	return db
}

func ReadCVS(w *window.Window, fields []string, cvsFile string, tdb map[string]string) map[string]string {
	var rowstr string
	var key string
	var k int
	db := make(map[string]string)
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
			if tdb[key] == key {
				db[key] = rowstr
				delete(tdb, key)
			} else {
				//errch <- key + "不存在"
				MsgLog(w, 500, key+"不存在")
			}

		}
		k++
	}
	return db
}

func SyncDataToHTTP(w *window.Window, ch chan FileMap, conferenceID string, fields []string, fmap map[string]string, susch, errch chan string) {
	for {
		fileMap, ok := <-ch
		if !ok {
			break
		} else {
			//fmt.Println(fileMap.Key)
			dbstr := fmap[fileMap.Key] + "conference_id=" + conferenceID + "&grade_id=1&" + fields[0] + "=" + fileMap.Name
			SyncHTTP(w, dbstr, susch, errch)
			delete(fmap, fileMap.Key)
		}
	}
	close(errch)
	close(susch)
}
