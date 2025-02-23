package rest

import (
	"crypto/subtle"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type AdminHandler struct {
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

// Create admin handler object
func NewAdminHandler(e *echo.Echo) {
	handler := &AdminHandler{}

	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte("admin")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("admin")) == 1 {
			return true, nil
		}
		return false, nil
	}))
	g.Use(middleware.Logger())
	// g.Use(middleware.Static("/public"))

	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/html/admin/*.html")),
	}

	g.GET("/", handler.AdminMainPage)
	g.GET("/update/:id", handler.AdminUpdatePage)
	g.GET("/create", handler.AdminCreatePage)
}

func (a *AdminHandler) AdminMainPage(c echo.Context) error {
	return c.Render(http.StatusOK, "main.html", nil)
}

func (a *AdminHandler) AdminUpdatePage(c echo.Context) error {
	return c.Render(http.StatusOK, "update.html", nil)
}

func (a *AdminHandler) AdminCreatePage(c echo.Context) error {
	return c.Render(http.StatusOK, "create.html", nil)
}
