package controllers

import (
	"fmt"
	"testing"
)

func TestReadCVS(t *testing.T) {
	fields := []string{"avatar", "guest_name", "mobile", "guest_email", "company_name", "position", "hint"}
	cvsFile := "/home/sensetime/Desktop/tst/test.csv"
	dataMap, err := ReadCVS(fields, cvsFile)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(dataMap)
}
