package public

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/briskt/go-htmx-app/api"
)

type Renderer struct {
	templates *template.Template
}

// Render is a custom Echo renderer for templates. It buffers the result in case there's an error. Otherwise, a portion
// of the page will be sent with the error.
func (r *Renderer) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	if err := r.templates.ExecuteTemplate(w, name, data); err != nil {
		return api.NewAppError(err, api.ErrorRenderingTemplate, http.StatusInternalServerError)
	}
	return nil
}

func NewRenderer() *Renderer {
	return &Renderer{
		templates: template.Must(template.ParseFS(EFS(), "view/*.gohtml", "view/*/*.gohtml")),
	}
}
