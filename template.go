package eweb

import (
	"io"
	"io/fs"
	"text/template"

	"github.com/labstack/echo"
)

type Template struct {
	*template.Template
}

func NewTemplate(tpl *template.Template) *Template {
	return &Template{tpl}
}
func GlobTemplate(filePath string) *Template {
	return &Template{template.Must(template.ParseGlob(filePath))}
}
func FSTemplate(fsys fs.FS, patterns ...string) (*Template, error) {
	tpl, err := template.ParseFS(fsys, patterns...)
	if err != nil {
		return nil, err
	}
	return &Template{tpl}, nil
}
func FilesTemplate(files ...string) (*Template, error) {
	tpl, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}
	return &Template{tpl}, nil
}

// Implements Renderer interface
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.ExecuteTemplate(w, name, data)
}
