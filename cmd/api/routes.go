package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Get("/categories", app.listCategoryHandler)
	router.Get("/categories/{uuid}", app.listCategoryHandler)
	router.Post("/categories", app.createCategoryHandler)
	router.Patch("/categories/{uuid}", app.updateCategoryHandler)
	router.Delete("/categories/{uuid}", app.deleteCategoryHandler)

	router.Get("categories/{uuid}/prdoducts", app.getProductHandler)
	router.Post("categories/{uuid}/prdoducts", app.createProductHandler)
	router.Patch("categories/{uuid}/prdoducts", app.updateProductHandler)
	router.Delete("categories/{uuid}/prdoducts", app.deleteProductHandler)

	return router
}
