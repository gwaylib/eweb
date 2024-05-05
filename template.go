package eweb

import (
	"html/template"
	"io"
	"io/fs"

	echo "github.com/labstack/echo/v4"
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
func FSTemplate(fsys fs.FS, patterns ...string) *Template {
	return &Template{template.Must(template.ParseFS(fsys, patterns...))}
}
func FilesTemplate(files ...string) *Template {
	return &Template{template.Must(template.ParseFiles(files...))}
}

// Implements echo.Renderer interface
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.ExecuteTemplate(w, name, data)
}
