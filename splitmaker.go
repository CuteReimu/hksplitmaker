package main

import (
	"embed"
	"encoding/json"
	"errors"
)

//go:embed hk-split-maker/src/asset/hollowknight/categories/*
var fs embed.FS

type CategoryDirectoryData struct {
	FileName    string `json:"fileName"`
	DisplayName string `json:"displayName"`
}

func GetAllFiles() (allFiles []Option) {
	file, _ := fs.ReadFile(hkSplitMakerDir + "/categories/category-directory.json")
	var v map[string][]*CategoryDirectoryData
	if err := json.Unmarshal(file, &v); err != nil {
		panic(err)
	}
	for _, categoryName := range []string{"Main", "Individual Level", "Category Extensions"} {
		for _, f := range v[categoryName] {
			allFiles = append(allFiles, Option{
				Value: f.FileName + ".json",
				Label: translate(f.DisplayName),
			})
		}
	}
	return
}

type SplitMakerConfig struct {
	CategoryName    string                   `json:"categoryName"`
	StartTriggering string                   `json:"startTriggeringAutosplit"`
	Ids             []string                 `json:"splitIds"`
	Names           map[string]StringOrSlice `json:"names"`
	Icons           map[string]StringOrSlice `json:"icons"`
	EndTriggering   bool                     `json:"endTriggeringAutosplit"`
	EndingSplit     *struct {
		Name string `json:"name"`
		Icon string `json:"icon"`
	} `json:"endingSplit"`
}

func GetSplitMakerConfig(fileName string) (*SplitMakerConfig, error) {
	buf, err := fs.ReadFile(hkSplitMakerDir + "/categories/" + fileName)
	if err != nil {
		return nil, err
	}
	var result SplitMakerConfig
	err = json.Unmarshal(buf, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type StringOrSlice []string

func (s *StringOrSlice) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = []string{str}
		return nil
	}

	var slice []string
	if err := json.Unmarshal(data, &slice); err == nil {
		*s = slice
		return nil
	}

	return errors.New("field must be string or []string")
}
