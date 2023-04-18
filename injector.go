//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/julienschmidt/httprouter"
	"golang-restful-api/app"
	"golang-restful-api/controller"
	"golang-restful-api/middleware"
	"golang-restful-api/repository"
	"golang-restful-api/server"
	"golang-restful-api/service"
	"net/http"
)

var categorySet = wire.NewSet(
	repository.NewCategoryRepositoryImplementation,
	wire.Bind(new(repository.CategoryRepository), new(*repository.CategoryRepositoryImplementation)),
	service.NewCategoryServiceImplementation,
	wire.Bind(new(service.CategoryService), new(*service.CategoryServiceImplementation)),
	controller.NewCategoryControllerImplementation,
	wire.Bind(new(controller.CategoryController), new(*controller.CategoryControllerImplementation)),
)

func InitializedServer() *http.Server {
	wire.Build(
		app.NewDB,
		validator.New,
		categorySet,
		app.NewRouter,
		wire.Bind(new(http.Handler), new(*httprouter.Router)),
		middleware.NewAuthMiddleware,
		server.NewServer,
	)

	return nil
}
