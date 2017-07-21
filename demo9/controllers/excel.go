package controllers

import (
	"encoding/csv"
	"os"
)

// func ReadExcel(fields []string, file string) (map[string]string, error) {
// 	//data := make(map[string]string)

// }

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
