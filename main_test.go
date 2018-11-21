package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestPrepare(t *testing.T) {
	settings := map[string][]Settings{
		"valid": {
			{
				Filepath: "files/valid.csv",
				Encoding: "utf-8",
			},
			{
				Filepath: "files/valid.xlsx",
				Encoding: "utf-8",
			},
			{
				Filepath: "files/valid.xml",
				Encoding: "utf-8",
			},
			{
				Filepath: "files/valid.json",
				Encoding: "utf-8",
			},
			{
				Filepath: "files/valid_big.csv",
				Encoding: "utf-8",
			},
			{
				Filepath: "files/valid_verybig.csv",
				Encoding: "utf-8",
			},
			{
				Filepath: "files/valid_big.xml",
				Encoding: "utf-8",
			},
			{
				Filepath: "files/valid_1251.csv",
				Encoding: "windows-1251",
			},
		},
		"invalid": {
			{
				Filepath: "files/invalid.csg",
				Encoding: "utf-8",
			},
			{
				Filepath: "files/invalid.xml",
				Encoding: "utf-8",
			},
			{
				Filepath: "files/invalid.json",
				Encoding: "utf-8",
			},
			{
				Filepath: "files/invalid.xlsx",
				Encoding: "utf-8",
			},
		},
	}

	os.Mkdir("files/res", 0755)

	for k, v := range settings {
		for _, s := range v {
			start := time.Now()
			cfg := NewConfig(s)
			r, err := os.Open(s.Filepath)
			os.Mkdir("files/res/"+cfg.Extension, 0755)
			w, err := os.OpenFile("files/res/"+cfg.Extension+"/"+preparedFileName, os.O_RDWR, 0755)
			if err != nil {
				w, err = os.Create("files/res/" + cfg.Extension + "/" + preparedFileName)
			}
			if err != nil {
				t.Fatal(err)
			}
			err = prepare(r, w, cfg)
			if err == nil {
				err = w.Sync()
			}
			if k == "valid" && err != nil {
				t.Error(fmt.Sprintf("%s : %s : error : %s", k, cfg.Extension, err.Error()))
			} else if k == "invalid" && err == nil {
				t.Error(fmt.Sprintf("%s : %s : error : %s", k, cfg.Extension, err.Error()))
			} else if k == "valid" && err == nil {
				t.Log(fmt.Sprintf("%s : %s : success : %f", k, cfg.Extension, time.Since(start).Seconds()))
			} else {
				t.Log(fmt.Sprintf("%s : %s : success : %f : %s", k, cfg.Extension, time.Since(start).Seconds(), err.Error()))
			}
		}
	}

	os.RemoveAll("files/res")
}

func TestNewConfig(t *testing.T) {
	settings := map[string][]Settings{
		"valid": {
			{
				Filepath: "files/valid.csv",
				Encoding: "utf-8",
			},
		},
		"invalid": {
			{
				Filepath: "files/invalidcsg",
				Encoding: "utf-8",
			},
			{
				Filepath: "files/",
				Encoding: "utf-8",
			},
			{},
		},
	}

	for k, v := range settings {
		for _, s := range v {
			cfg := NewConfig(s)
			if k == "valid" && cfg.Extension != "" && cfg.Name != "" {
				t.Log(fmt.Sprintf("%s : %s : success", cfg.Name, cfg.Extension))
			} else if k == "invalid" && (cfg.Extension == "" || cfg.Name == "") {
				t.Log(fmt.Sprintf("%s : %s : success", cfg.Name, cfg.Extension))
			} else {
				t.Error(fmt.Sprintf("%v : %s : error", cfg, cfg.Extension))
			}
		}
	}

}
