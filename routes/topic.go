package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type TopicRequest struct {
	ID string `json:"id" validate:"required"`
}

func (router *Router) registerTopicRoutes() {
	router.echo.GET("/api/topics", func(c echo.Context) error {
		topics, err := router.pubsub.GetTopics(c.Request().Context())
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, topics)
	})

	router.echo.DELETE("/api/topics/:topicId", func(c echo.Context) error {
		id := c.Param("topicId")
		err := router.pubsub.DeleteTopic(c.Request().Context(), id)
		if err != nil {
			return err
		}
		return c.NoContent(http.StatusNoContent)
	})

	router.echo.POST("/api/topics", func(c echo.Context) error {
		var req TopicRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		if err := c.Validate(&req); err != nil {
			return err
		}

		err := router.pubsub.CreateTopic(c.Request().Context(), req.ID)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, map[string]string{"id": req.ID})
	})
}
