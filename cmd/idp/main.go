package main

import (
	"html/template"
	"io"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/ory/hydra-client-go/client"
	"github.com/sirupsen/logrus"

	"github.com/ihippik/idp-example/provider"
)

const (
	hydraAdmin = "http://localhost:4445"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("views/login.html")),
	}

	adminURL, err := url.Parse(hydraAdmin)
	if err != nil {
		logrus.WithError(err).Fatalln("url parse error")
	}
	hydraClient := client.NewHTTPClientWithConfig(
		nil,
		&client.TransportConfig{
			Schemes:  []string{adminURL.Scheme},
			Host:     adminURL.Host,
			BasePath: adminURL.Path,
		},
	)
	srv := provider.NewService(hydraClient)
	e.GET("/login", srv.LoginHandler)
	e.GET("/consent", srv.ConsentHandler)
	e.POST("/signIn", srv.SignInHandler)
	e.Logger.Fatal(e.Start(":3000"))
}
