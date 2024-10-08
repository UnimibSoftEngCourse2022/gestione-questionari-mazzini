package internal

import (
	common "form_management/common/logger"
	"form_management/internal/auth"
	"form_management/internal/question"
	"form_management/internal/quiz"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Server() {
	e := echo.New()

	common.NewLogger()
	e.Use(common.LoggingMiddleware)
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/views/**/*.html")),
	}
	e.Renderer = renderer
	e.Static("/static", "public/static")

	e.GET("/*", func(c echo.Context) error { return c.Render(http.StatusNotFound, "error.html", "404 not found") })

	e.GET("/login", LoginPageHandler)
	e.GET("/register", RegisternPageHandler)

	e.GET("/", HomePageHandler, auth.AuthMiddleware)
	e.GET("/questions", QuestionPageHandler, auth.AuthMiddleware)
	e.GET("/quizs", QuizPageHandler, auth.AuthMiddleware)
	e.GET("/quiz-edit", QuizEditPageHandler, auth.AuthMiddleware)

	authRoute := e.Group("/auth")
	formRoute := e.Group("/form", auth.AuthMiddleware)
	quizRoute := e.Group("/quiz", auth.AuthMiddleware)
	auth.Route(authRoute)
	question.Route(formRoute)
	quiz.Route(quizRoute)

	common.Logger.LogInfo().Msg(e.Start(":" + os.Getenv("APP_PORT")).Error())
}
