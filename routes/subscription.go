package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SubscriptionRequest struct {
	ID      string `json:"id" validate:"required"`
	TopicID string `json:"topicId" validate:"required"`
}

func (router *Router) registerSubsciberRoutes() {
	router.echo.GET("/api/subscriptions", func(c echo.Context) error {
		topicId := c.QueryParam("topicName")

		if topicId == "" {
			subscriptions, err := router.pubsub.GetSubscriptions(c.Request().Context())
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, subscriptions)
		}

		subscriptions, err := router.pubsub.GetSubscriptionsByTopic(c.Request().Context(), topicId)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, subscriptions)
	})

	router.echo.POST("/api/subscriptions", func(c echo.Context) error {
		var req SubscriptionRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		if err := c.Validate(&req); err != nil {
			return err
		}

		err := router.pubsub.CreateSubscription(c.Request().Context(), req.ID, req.TopicID)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, map[string]string{"id": req.ID, "topicId": req.TopicID})
	})

	router.echo.DELETE("/api/subscriptions/:subscriptionId", func(c echo.Context) error {
		id := c.Param("subscriptionId")
		err := router.pubsub.DeleteSubscription(c.Request().Context(), id)
		if err != nil {
			return err
		}
		return c.NoContent(http.StatusNoContent)
	})
}
