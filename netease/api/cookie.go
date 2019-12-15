package api

import (
	"bytes"
	"encoding/gob"
	"github.com/spf13/afero"
)

func SaveCookies(path string, data map[string]string) (err error) {
	b := new(bytes.Buffer)
	if err = gob.NewEncoder(b).Encode(data); err == nil {
		fs := afero.NewOsFs()
		_ = fs.Remove(path)
		var f afero.File
		if f, err = fs.Create(path); err == nil {
			_, err = f.Write(b.Bytes())
			defer f.Close()
		}
	}
	return
}

func LoadCookies(path string) (ret map[string]string, err error) {
	var fs = afero.NewOsFs()
	var f afero.File
	if f, err = fs.Open(path); err == nil {
		if err = gob.NewDecoder(f).Decode(&ret); err == nil {
			defer f.Close()
		}
	}
	return
}
