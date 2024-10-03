package quiz

import (
	common "form_management/common/logger"
	"form_management/db"
	"form_management/internal/quiz/entities"
	service "form_management/internal/quiz/services"

	"github.com/labstack/echo/v4"
)

func Route(e *echo.Group) {
	logger := common.Logger

	db := db.Init()
	db.AutoMigrate(
		&entities.Quiz{},
		&entities.QuizClosedQuestion{},
		&entities.QuizOpenQuestion{},
	)

	service := service.NewQuizService(&logger, db)
	handler := NewQuizHanlder(service)

	e.GET("/findAll", handler.ListQuiz)
	e.GET("/find", handler.FindQuiz)
	e.POST("/create", handler.CreateQuiz)
	e.POST("/update", handler.UpdateQuiz)
	e.DELETE("/delete", handler.DeleteQuiz)

}
