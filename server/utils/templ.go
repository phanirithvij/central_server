// Package utils a utility package
package utils

import (
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/markbates/pkger"
)

// LoadTemplates loads the templates used by this package
// @t the template instance
// @dirname the directory where templates are residing
func LoadTemplates(t *template.Template, dirname string) (*template.Template, error) {

	// https://gin-gonic.com/docs/examples/bind-single-binary-with-template/
	// https://github.com/gin-gonic/examples/commit/c5a87f03d39fdb9e0f6312344c21ccdd55140293

	err := pkger.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}
		file, err := pkger.Open(path)
		if err != nil {
			return err
		}
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		// log.Println("Regster templ", path)
		// templates can be chained together
		t, err = t.New(path).Parse(string(h))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return t, nil
}
