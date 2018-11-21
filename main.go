package main

import (
	"encoding/csv"
	"encoding/json"
	"github.com/clbanning/x2j"
	"github.com/tealeg/xlsx"
	"go/types"
	"gopkg.in/iconv.v1"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultDir       = ""
	preparedFileName = "prepared.json"
	defaultEncoding  = "utf-8"
)

type Config struct {
	s         Settings
	Name      string
	Extension string
	Dir       string
	csv       struct {
		Comma       rune
		FieldsCount int
		LazyQuotes  bool
	}
}

type Settings struct {
	Id        string `json:"id"`
	Filepath  string `json:"filepath"`
	Encoding  string `json:"encoding"`
	Extension string `json:"extension"`
}

func main() {
	var (
		err      error
		settings Settings
	)

	settings, err = getConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(settings)

	err = prepareWrapper(&settings)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("done")
	}
}

func getConfig() (Settings, error) {
	var (
		err      error
		settings Settings
	)

	buffer, err := ioutil.ReadFile("config.json")
	if err != nil {
		return settings, err
	}

	err = json.Unmarshal(buffer, &settings)
	if err != nil {
		return settings, err
	}

	return settings, nil
}

func prepareWrapper(s *Settings) (err error) {
	var (
		f *os.File
		w *os.File
		//res *http.Response
	)

	cfg := NewConfig(*s)
	err = os.Mkdir(cfg.Dir, 0755)
	if err != nil {
		return err
	}

	if f, err = os.Create(cfg.Dir + "/" + cfg.Name); err == nil {
		defer f.Close()
		w, err = os.Create(cfg.Dir + "/" + preparedFileName)
		if err != nil {
			return err
		}
		defer w.Close()
		//if res, err = http.Get(s.Filepath); err == nil {
		//	if _, err = io.Copy(f, res.Body); err == nil {
		//if err = f.Sync(); err == nil {
		if f, err = os.Open(cfg.Dir + "/" + cfg.Name); err == nil {
			err = prepare(f, w, cfg)
			if err == nil {
				err = w.Sync()
			}
		}
		//}
		//}
		//}
	}

	return err
}

func prepare(r io.Reader, w io.Writer, cfg *Config) (err error) {
	var (
		res io.Reader
		cd  iconv.Iconv
	)
	if cfg.s.Encoding == defaultEncoding {
		res = r
	} else if cd, err = iconv.Open(defaultEncoding, cfg.s.Encoding); err == nil {
		defer cd.Close()
		res = iconv.NewReader(cd, r, 32*1024)
	} else {
		return err
	}

	return toJson(res, w, cfg)
}

func toJson(r io.Reader, w io.Writer, cfg *Config) (err error) {
	var (
		res interface{}
	)
	//fmt.Printf("cfg: %v\n", cfg)
	switch cfg.Extension {
	case "csv":
		var b *csv.Reader
		b = csv.NewReader(r)
		b.Comma = cfg.csv.Comma
		b.LazyQuotes = cfg.csv.LazyQuotes
		b.FieldsPerRecord = cfg.csv.FieldsCount

		res, err = b.ReadAll()
		if err != nil {
			return err
		}
		break
	case "xlsx":
		var (
			b []byte
			f *xlsx.File
		)
		if b, err = ioutil.ReadAll(r); err == nil {
			if f, err = xlsx.OpenBinary(b); err == nil {
				res, err = f.ToSlice()
			}
		}
		if err != nil {
			return err
		}
		break
	case "xml":
		var b []byte
		b, err = ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		_, err = x2j.ByteDocToTree(b)
		if err != nil {
			return err
		}
		res, err = x2j.ByteDocToMap(b)
		if err != nil {
			return err
		}
		break
	case "json":
		var b []byte
		b, _ = ioutil.ReadAll(r)
		valid := json.Valid(b)
		if !valid {
			return types.Error{Msg: "Invalid json format"}
		}
		_, err = io.Copy(w, r)
		return err
	default:
		return types.Error{Msg: "Invalid file extension"}
	}

	err = json.NewEncoder(w).Encode(res)
	return err
}

func NewConfig(s Settings) *Config {
	fname := filepath.Base(s.Filepath)
	ext := ""
	arFname := strings.Split(fname, ".")
	if len(arFname) > 1 {
		ext = arFname[len(arFname)-1]
	}

	return &Config{
		s,
		fname,
		ext,
		defaultDir + s.Id,
		struct {
			Comma       rune
			FieldsCount int
			LazyQuotes  bool
		}{Comma: ';', FieldsCount: -1, LazyQuotes: true}}
}
