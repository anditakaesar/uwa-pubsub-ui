package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/anditakaesar/uwa-pubsub-ui/internal/env"
	"github.com/anditakaesar/uwa-pubsub-ui/internal/features"
	"github.com/anditakaesar/uwa-pubsub-ui/routes"
	"github.com/anditakaesar/uwa-pubsub-ui/xpubsub"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	publicFoler string = "public"
	portAddr    string = ":8080"
)

func main() {
	e := echo.New()

	// Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.Static("/", publicFoler)
	features.InitializeValidator(e)

	env.InitEnv()

	psClient := xpubsub.NewPubSubClient(&env.PubSubConfig{})

	if env.InitPubSubFromFile() {
		topicSpec, err := xpubsub.LoadTopic(env.PubSubInitFile())
		if err != nil {
			slog.Error("failed to load topic", "error", err)
		}

		err = psClient.CreateTopicAndSubscriptions(context.Background(), topicSpec)
		if err != nil {
			slog.Error("failed to create topic and subscriptions", "error", err)
		}
	}

	router := routes.NewRouter(e, psClient)
	router.Initialize()

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/topics.html")
	})

	err := e.Start(portAddr)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
