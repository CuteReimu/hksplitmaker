package main

import (
	"encoding/json"
	"github.com/lxn/walk"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type jsonCategory struct {
	CategoryName           string                   `json:"categoryName"`
	SplitIds               []string                 `json:"splitIds"`
	Ordered                bool                     `json:"ordered"`
	EndTriggeringAutosplit bool                     `json:"endTriggeringAutosplit"`
	Names                  map[string]interface{}   `json:"names"`
	Icons                  map[string]interface{}   `json:"icons"`
	EndingSplit            *jsonCategoryEndingSplit `json:"endingSplit"`
	GameName               string                   `json:"gameName"`
}

type jsonCategoryEndingSplit struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type jsonCategoryInfo struct {
	FileName    string `json:"fileName"`
	DisplayName string `json:"displayName"`
}

var categoriesCache = make(map[string]*jsonCategory)

func initCategories() {
	buf, err := ioutil.ReadFile(filepath.Join("categories", "category-directory.json"))
	if err != nil {
		walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
		panic(err)
	}
	var categoryDirectoryCache map[string][]*jsonCategoryInfo
	err = json.Unmarshal(buf, &categoryDirectoryCache)
	if err != nil {
		walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
		panic(err)
	}
	for groupName, v := range categoryDirectoryCache {
		if groupName != "Individual Level" && groupName != "Main" && groupName != "Category Extensions" {
			continue
		}
		for _, info := range v {
			buf, err := ioutil.ReadFile(filepath.Join("categories", info.FileName+".json"))
			if err != nil {
				walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
				panic(err)
			}
			j := &jsonCategory{}
			err = json.Unmarshal(buf, j)
			if err != nil {
				walk.MsgBox(nil, "错误", err.Error(), walk.MsgBoxIconError)
				panic(err)
			}
			count := len(j.SplitIds)
			if !j.EndTriggeringAutosplit {
				if j.EndingSplit == nil {
					continue
				}
				if j.EndingSplit.Icon != "HollowKnightBoss" && j.EndingSplit.Icon != "RadianceBoss" {
					panic(info.FileName)
				}
				count++
			}
			foundPer := false
			for _, splitId := range j.SplitIds {
				if strings.Contains(splitId, "%") {
					foundPer = true
					break
				}
			}
			if foundPer {
				continue
			}
			if count >= 2 && j.Ordered /*&& len(j.SplitIds) <= 50*/ {
				categoriesCache[translate(info.DisplayName)] = j
			}
		}
	}
}

var categoryCurrent string

func onSelectCategory() {
	category := categoriesComboBox.Text()
	if len(category) == 0 || category == categoryCurrent {
		return
	}
	categoryCurrent = category
	j := categoriesCache[category]
	count := len(j.SplitIds)
	if !j.EndTriggeringAutosplit {
		count++
	}
	if count > 50 {
		if walk.MsgBox(mainWindow, "确认", "这个类别所含的片段较多，可能会加载很久，确定继续吗？", walk.MsgBoxYesNo) != walk.DlgCmdYes {
			return
		}
	}
	cleanAllLines()
	for i := len(lines); i < count-1; i++ {
		addLine(false)
	}
	reg, err := regexp.Compile(`{.*?}|\[[0-9DU, ]*]`)
	if err != nil {
		panic(err)
	}
	nameIndexCache := make(map[string]int)
	getNameFunc := func(splitId, description string) string {
		name := ""
		if names, ok := j.Names[splitId]; ok {
			if namestr, ok := names.(string); ok {
				name = strings.ReplaceAll(namestr, "%s", dropBrackets(description))
			} else if namearr, ok := names.([]interface{}); ok {
				if _, ok := nameIndexCache[splitId]; !ok {
					nameIndexCache[splitId] = 0
				}
				name = strings.ReplaceAll(namearr[nameIndexCache[splitId]].(string), "%s", dropBrackets(description))
				nameIndexCache[splitId]++
			}
		} else {
			name = dropBrackets(description)
		}
		return translate(strings.TrimSpace(reg.ReplaceAllString(name, "")))
	}
	if j.EndTriggeringAutosplit {
		for i, splitId := range j.SplitIds {
			splitId = strings.Trim(splitId, "-")
			if i < len(j.SplitIds)-1 {
				description := splitsDictIdToDescriptions[splitId]
				err := lines[i].splitId.SetText(description)
				if err != nil {
					walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
					return
				}
				err = lines[i].name.SetText(getNameFunc(splitId, description))
				if err != nil {
					walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
					return
				}
			} else {
				description := splitsDictIdToDescriptions[splitId]
				finalLine.endTrigger.SetChecked(false)
				err := finalLine.splitId.SetText(description)
				if err != nil {
					walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
					return
				}
				err = finalLine.name.SetText(getNameFunc(splitId, description))
				if err != nil {
					walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
					return
				}
			}
		}
	} else {
		for i, splitId := range j.SplitIds {
			splitId = strings.Trim(splitId, "-")
			description := splitsDictIdToDescriptions[splitId]
			err := lines[i].splitId.SetText(description)
			if err != nil {
				walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
				return
			}
			err = lines[i].name.SetText(getNameFunc(splitId, description))
			if err != nil {
				walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
				return
			}
		}
		finalLine.endTrigger.SetChecked(true)
		text := "空洞骑士"
		if j.EndingSplit.Icon == "RadianceBoss" {
			if j.EndingSplit.Name == "Absolute Radiance" {
				text = "无上辐光"
			} else {
				text = "辐光"
			}
		}
		err := finalLine.splitId2.SetText(text)
		if err != nil {
			walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
			return
		}
		err = finalLine.name.SetText(getNameFunc(j.EndingSplit.Icon, j.EndingSplit.Name))
		if err != nil {
			walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
			return
		}
	}
}

func dropBrackets(s string) string {
	rs := []rune(s)
	for i, r := range rs {
		if r == '（' {
			return string(rs[:i])
		}
	}
	return s
}
