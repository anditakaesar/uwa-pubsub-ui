package routes

import (
	"github.com/anditakaesar/uwa-pubsub-ui/xpubsub"
	"github.com/labstack/echo/v4"
)

type Router struct {
	echo   *echo.Echo
	pubsub xpubsub.Interface
}

func (router *Router) Initialize() {
	router.registerTopicRoutes()
	router.registerSubsciberRoutes()
}

func NewRouter(e *echo.Echo, psClient xpubsub.Interface) *Router {
	return &Router{
		echo:   e,
		pubsub: psClient,
	}
}
