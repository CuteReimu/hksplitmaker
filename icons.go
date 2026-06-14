package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"embed"
	"encoding/base64"
	"errors"
	"io"
	fs2 "io/fs"
	"path"
	"regexp"
	"strings"
)

const hkSplitMakerDir = "hk-split-maker/src/asset/hollowknight"

//go:embed hk-split-maker/src/asset/hollowknight/icons
var iconFs embed.FS

var iconDict = make(map[string]string)

func init() {
	f, err := iconFs.Open(hkSplitMakerDir + "/icons/icons.ts")
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()
	re, err := regexp.Compile(`import\s+(\w+)\s+from\s+"(.*?)"\s*;`)
	if err != nil {
		panic(err)
	}
	rd := bufio.NewReader(f)
	line, isPrefix, err := rd.ReadLine()
	for ; err == nil; line, isPrefix, err = rd.ReadLine() {
		if isPrefix {
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
	if errors.Is(err, io.EOF) {
		err = nil
	}
	if err != nil {
		panic(err)
	}
}

func zipIcons() ([]byte, error) {
	b := &bytes.Buffer{}
	zipWriter := zip.NewWriter(b)
	err := fs2.WalkDir(iconFs, hkSplitMakerDir+"/icons", func(filePath string, d fs2.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() == "icons.ts" {
			return nil
		}
		fw, err := zipWriter.Create(strings.Replace(filePath, "hk-split-maker/src/asset/hollowknight/", "", 1))
		if err != nil {
			return err
		}
		buf, err := iconFs.ReadFile(filePath)
		if err != nil {
			return err
		}
		_, err = fw.Write(buf)
		return err
	})
	if err != nil {
		return nil, err
	}
	if err = zipWriter.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func getIconHtmlFormat(splitId string) string {
	iconPath, ok := iconDict[splitId]
	if !ok {
		return ""
	}
	iconPath = path.Join(hkSplitMakerDir, "icons", iconPath)
	buf, err := iconFs.ReadFile(iconPath)
	if err != nil {
		return ""
	}
	s := base64.StdEncoding.EncodeToString(buf)
	return htmlFormatHeader + s
}

func convertIconToHtmlFormat(icon string) string {
	if len(icon) == 0 {
		return ""
	}
	buf, err := base64.StdEncoding.DecodeString(icon)
	if err != nil {
		return ""
	}
	headerIndex := bytes.Index(buf, []byte("\x89PNG\x0D\x0A\x1A\x0A"))
	if headerIndex < 0 {
		return ""
	}
	buf = buf[headerIndex:]
	return htmlFormatHeader + base64.StdEncoding.EncodeToString(buf)
}

func convertIconToLiveSplitFormat(icon string) string {
	if !strings.HasPrefix(icon, htmlFormatHeader) {
		return ""
	}
	icon = icon[len(htmlFormatHeader):]
	buf, err := base64.StdEncoding.DecodeString(icon)
	if err != nil || len(buf) == 0 {
		return ""
	}
	return livesplitFormatHeader + base64.StdEncoding.EncodeToString(append([]byte{0, 2}, buf...))
}

const (
	htmlFormatHeader      = "data:image/png;base64,"
	livesplitFormatHeader = "AAEAAAD/////AQAAAAAAAAAMAgAAAFFTeXN0ZW0uRHJhd2luZywgVmVyc2lvbj00LjAuMC4wLCBDdWx0dXJlPW5ldXRyYWwsIFB1YmxpY0tleVRva2VuPWIwM2Y1ZjdmMTFkNTBhM2EFAQAAABVTeXN0ZW0uRHJhd2luZy5CaXRtYXABAAAABERhdGEHAgIAAAAJAwAAAA8DAAAAOw8A"
)
