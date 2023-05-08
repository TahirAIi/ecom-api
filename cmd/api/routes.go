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

	router.Group(func(r chi.Router) {
		r.Use(app.IsAdmin())
		r.Post("/categories", app.createCategoryHandler)
		r.Patch("/categories/{category_id}", app.updateCategoryHandler)
		r.Delete("/categories/{category_id}", app.deleteCategoryHandler)

		r.Post("/categories/{category_id}/products", app.createProductHandler)
		r.Patch("/categories/{category_id}/products/{product_id}", app.updateProductHandler)
		r.Delete("/categories/{category_id}/products/{product_id}", app.deleteProductHandler)
	})

	router.Get("/categories", app.listCategoryHandler)
	router.Get("/categories/{category_id}", app.getCategoryHandler)

	router.Get("/categories/{category_id}/products", app.listProductHandler)
	router.Get("/categories/{category_id}/products/{product_id}", app.getProductHandler)

	router.Post("/users", app.createUserHandler)
	router.Post("/authorize", app.authUserHandler)

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
