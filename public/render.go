package public

import (
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"github.com/briskt/go-htmx-app/api"
)

// TemplRenderer renders TEMPL components for Echo.
type TemplRenderer struct{}

// Render renders a TEMPL component.
// The `data` must be of type templ.Component or a function returning one.
func (r *TemplRenderer) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	var comp templ.Component

	switch v := data.(type) {
	case templ.Component:
		comp = v
	case func() templ.Component:
		comp = v()
	default:
		return api.NewAppError(
			http.ErrNotSupported,
			api.ErrorRenderingTemplate,
			http.StatusInternalServerError,
		)
	}

	err := comp.Render(context.Background(), w)
	if err != nil {
		return api.NewAppError(err, api.ErrorRenderingTemplate, http.StatusInternalServerError)
	}
	return nil
}
