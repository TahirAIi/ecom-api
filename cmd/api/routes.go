package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Get("/categories", app.listCategoryHandler)
	router.Get("/categories/{id}", app.listCategoryHandler)
	router.Post("/categories", app.createCategoryHandler)
	router.Patch("/categories/{id}", app.updateCategoryHandler)
	router.Delete("/categories/{id}", app.deleteCategoryHandler)

	router.Get("/categories/{id}/products", app.listProductHandler)
	router.Get("/categories/{id}/products/{product_id}", app.getProductHandler)
	router.Post("/categories/{id}/products", app.createProductHandler)
	router.Patch("/categories/{id}/products/{product_id}", app.updateProductHandler)
	router.Delete("/categories/{id}/products/{product_id}", app.deleteProductHandler)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "uploads"))
	FileServer(router, "/files", filesDir)

	return router
}

func FileServer(router chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		router.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	router.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
