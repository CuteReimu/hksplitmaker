package main

import (
	"fmt"
	"github.com/lxn/walk"
	"github.com/xuri/excelize/v2"
	"os"
	"strings"
)

var translateDict = make(map[string]string)

func init() {
	xlsx, err := excelize.OpenFile("translate.xlsx")
	if err != nil {
		walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}
	rows, err := xlsx.GetRows("Sheet1")
	if err != nil {
		walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
		os.Exit(-1)
	}
	for i, row := range rows {
		key, val := strings.TrimSpace(row[0]), strings.TrimSpace(row[1])
		if len(key) == 0 || len(val) == 0 {
			walk.MsgBox(nil, "警告", fmt.Sprintf("第%d行出现空数据", i+1), walk.MsgBoxIconWarning)
		} else {
			translateDict[row[0]] = row[1]
		}
	}
}

func translate(s string) string {
	s2 := doTranslate(s)
	if s2 != s {
		return strings.ReplaceAll(s2, " ", "")
	}
	return s
}

func doTranslate(s string) string {
	if val, ok := translateDict[s]; ok {
		return val
	}
	arr := strings.Split(s, " ")
	for n := len(arr) - 1; n >= 1; n-- {
		for i := 0; i+n <= len(arr); i++ {
			key := strings.Join(arr[i:i+n], " ")
			if val, ok := translateDict[key]; ok {
				return doTranslate(strings.ReplaceAll(s, key, val))
			}
		}
	}
	return s
}
