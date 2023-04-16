package app

import (
	"github.com/julienschmidt/httprouter"
	"golang-restful-api/controller"
	"golang-restful-api/exception"
)

func NewRouter(categoryController controller.CategoryController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/categories", categoryController.GetAllCategory)
	router.GET("/api/categories/:categoryId", categoryController.GetCategoryById)
	router.POST("/api/categories", categoryController.CreateCategory)
	router.PUT("/api/categories/:categoryId", categoryController.UpdateCategory)
	router.DELETE("/api/categories/:categoryId", categoryController.DeleteCategory)

	router.PanicHandler = exception.ErrorHandler

	return router
}
