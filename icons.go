package main

import (
	"bufio"
	"encoding/base64"
	"github.com/lxn/walk"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

var iconDict = make(map[string]string)

func init() {
	f, err := os.Open("icons/icons.ts")
	if err != nil {
		walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
		}
	}()
	re, err := regexp.Compile(`import\s+(\w+)\s+from\s+"(.*?)"\s*;`)
	if err != nil {
		walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
		panic(err)
	}
	rd := bufio.NewReader(f)
	line, isPrefix, err := rd.ReadLine()
	for ; err == nil; line, isPrefix, err = rd.ReadLine() {
		if isPrefix {
			walk.MsgBox(nil, "错误", "暂不支持这样的文件", walk.MsgBoxIconError)
			panic(err)
		}
		lineStr := strings.TrimSpace(string(line))
		if len(lineStr) < 2 || lineStr[:2] == "//" {
			continue
		}
		result := re.FindStringSubmatch(lineStr)
		if result != nil {
			iconDict[result[1]] = result[2]
		}
	}
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
		panic(err)
	}
}

func getIcon(splitId string) string {
	iconPath, ok := iconDict[splitId]
	if !ok {
		return ""
	}
	iconPath = path.Join("icons", iconPath)
	buf, err := ioutil.ReadFile(iconPath)
	if err != nil {
		return ""
	}
	s := base64.StdEncoding.EncodeToString(append([]byte{0, 2}, buf...))
	return livesplitFormatHeader + s
	//return "<![CDATA[" + livesplitFormatHeader + s + "]]>"
}

const livesplitFormatHeader = "AAEAAAD/////AQAAAAAAAAAMAgAAAFFTeXN0ZW0uRHJhd2luZywgVmVyc2lvbj00LjAuMC4wLCBDdWx0dXJlPW5ldXRyYWwsIFB1YmxpY0tleVRva2VuPWIwM2Y1ZjdmMTFkNTBhM2EFAQAAABVTeXN0ZW0uRHJhd2luZy5CaXRtYXABAAAABERhdGEHAgIAAAAJAwAAAA8DAAAAOw8A"
