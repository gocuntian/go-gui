package controllers

import (
	"encoding/csv"
	"os"

	"github.com/tealeg/xlsx"
)

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
						// if i == 2 || i == 3 || i == 4 || i == 5 {
						// 	val, _ := DesEncrypt([]byte(cell.String()), []byte(DES_KEY))
						// 	rowstr += fields[i] + "=" + val + "&"
						// } else {
						// 	rowstr += fields[i] + "=" + cell.String() + "&"
						// }

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
					// if i == 2 || i == 3 || i == 4 || i == 5 {
					// 	val, _ := DesEncrypt([]byte(cell), []byte(DES_KEY))
					// 	rowstr += fields[i] + "=" + val + "&"
					// } else {
					// 	rowstr += fields[i] + "=" + cell + "&"
					// }
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
