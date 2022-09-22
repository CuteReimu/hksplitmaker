package main

import (
	_ "embed"
	"encoding/json"
	"github.com/lxn/walk"
	"os"
	"path"
	"regexp"
	"strings"
)

type jsonCategory struct {
	CategoryName             string                   `json:"categoryName"`
	SplitIds                 []string                 `json:"splitIds"`
	Ordered                  bool                     `json:"ordered"`
	StartTriggeringAutosplit string                   `json:"startTriggeringAutosplit"`
	EndTriggeringAutosplit   bool                     `json:"endTriggeringAutosplit"`
	Names                    map[string]interface{}   `json:"names"`
	Icons                    map[string]interface{}   `json:"icons"`
	EndingSplit              *jsonCategoryEndingSplit `json:"endingSplit"`
	GameName                 string                   `json:"gameName"`
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
	buf, err := assets.ReadFile(path.Join(hkSplitMakerDir, "categories", "category-directory.json"))
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
	for _, v := range categoryDirectoryCache {
		for _, info := range v {
			buf, err := assets.ReadFile(path.Join(hkSplitMakerDir, "categories", info.FileName+".json"))
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
					continue
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
			if count >= 2 && j.Ordered {
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
	resetLines(count - 1)
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
			name = strings.TrimSpace(reg.ReplaceAllString(name, ""))
		}
		if len(name) == 0 {
			name = dropBrackets(description)
		}
		return translate(name)
	}
	if startTrigger, ok := splitsDictIdToDescriptions[j.StartTriggeringAutosplit]; ok {
		startTriggerCheckBox.SetChecked(true)
		err := startTriggerComboBox.SetText(startTrigger)
		if err != nil {
			walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
			return
		}
		startTriggerComboBox.SetEnabled(true)
	} else {
		startTriggerCheckBox.SetChecked(false)
		startTriggerComboBox.SetEnabled(false)
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
	categoryName = j.CategoryName
	saveTimeCheckBox.SetEnabled(false)
	saveTimeCheckBox.SetChecked(false)
	coloNotice(j)
}

func dropBrackets(s string) string {
	idx := strings.LastIndex(s, "(")
	if idx > 0 {
		return s[:idx]
	}
	return s
}

//go:embed blank-colo-save.dat
var blankColoSave []byte

func coloNotice(j *jsonCategory) {
	if func() bool {
		for i := range j.SplitIds {
			splitId := strings.Trim(j.SplitIds[i], "-")
			if splitId == "BronzeEnd" || splitId == "SilverEnd" || splitId == "GoldEnd" {
				return true
			}
		}
		return false
	}() {
		if walk.DlgCmdOK == walk.MsgBox(mainWindow, "提示", "愚人斗兽场相关的Splits需要专用的存档才能正常使用。是否要生成专用存档？", walk.MsgBoxOKCancel|walk.MsgBoxIconInformation) {
			if err := os.WriteFile("user4.dat", blankColoSave, 0644); err != nil {
				walk.MsgBox(mainWindow, "错误", err.Error(), walk.MsgBoxIconError)
			} else {
				walk.MsgBox(mainWindow, "提示", `已成功生成 user4.dat ，请自行放入空洞骑士的存档目录`, walk.MsgBoxIconInformation)
			}
		}
	}
}
