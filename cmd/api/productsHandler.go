package main

import (
	"database/sql"
	data "ecom-api/inernal/data/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

//swagger:route POST /categories/{category_id}/products Products createProduct
// Creates a product.
//
//Consumes:
//	- multipart/form-data
//
//Produces:
//	- application/json
//Parameters:
//	ProductBody
//responses:
//	200: ProductResponse
//swagger:response

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	product := &data.Product{}

	err := r.ParseMultipartForm(int64(app.config.multipartFormSize))
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	title := r.FormValue("title")
	tempDescription := r.FormValue("description")
	price, err := strconv.ParseFloat(r.FormValue("price"), 32)

	if err != nil {
		app.sendInternalServerErrorResponse(w)
		return
	}

	fileName, err := app.uploadFile("main_picture", r)

	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	if fileName != nil {
		product.MainPicture = *fileName
	}

	val, err := app.convertToInt(chi.URLParam(r, "category_id"))

	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	if val < 1 {
		app.sendResponse(w, response{"message": "Category id is required"}, http.StatusUnprocessableEntity)
		return
	}
	product.CategoryId = val

	if len(title) == 0 {
		app.sendResponse(w, response{"message": "Title can't be empty"}, http.StatusUnprocessableEntity)
		return
	}
	product.Title = title

	if len(tempDescription) > 0 {
		product.Description = &tempDescription
	}

	if price <= 0 {
		app.sendResponse(w, response{"message": "Price can't be less than 0"}, http.StatusUnprocessableEntity)
		return
	}

	product.Price = int32(price * 100)
	err = app.models.Product.Insert(product)
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	app.sendResponse(w, response{"product": product}, http.StatusCreated)
}

// Updates a product.
//
// Consumes:
//   - multipart/form-data
//
// Parameters:
//
//	ProductBody
//
// responses:
//
//	200: ProductResponse
//
//swagger:route PATCH /categories/{category_id}/products Products updateProduct
//swagger:response
func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	var product data.Product

	productId, err := app.convertToInt(chi.URLParam(r, "product_id"))
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	categoryId, err := app.convertToInt(chi.URLParam(r, "category_id"))
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	err = app.models.Product.Get(categoryId, productId, &product)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.sendResponse(w, response{}, http.StatusNotFound)
			return
		default:
			app.log(err)
			app.sendInternalServerErrorResponse(w)
			return
		}
	}

	err = r.ParseMultipartForm(int64(app.config.multipartFormSize))
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	if (len(r.FormValue("category_id"))) > 0 {
		val, err := app.convertToInt(r.FormValue("category_id"))
		if err != nil {
			app.log(err)
			app.sendInternalServerErrorResponse(w)
			return
		}
		product.CategoryId = val
	}

	if len(r.FormValue("title")) > 0 {
		product.Title = r.FormValue("title")
	}

	if r.PostForm.Has("description") {
		tempDescription := r.FormValue("description")
		if len(tempDescription) > 0 {
			product.Description = &tempDescription
		}
	}

	if len(r.FormValue("price")) > 0 {
		price, err := strconv.ParseFloat(r.FormValue("price"), 32)
		if err != nil {
			app.sendInternalServerErrorResponse(w)
			return
		}

		product.Price = int32(price * 100)
	}
	if _, _, err = r.FormFile("main_picture"); err != http.ErrMissingFile {
		fileName, err := app.uploadFile("main_picture", r)

		if err != nil && !errors.Is(err, http.ErrMissingFile) {
			app.log(err)
			app.sendInternalServerErrorResponse(w)
			return
		}

		if fileName != nil {
			product.MainPicture = *fileName
		}
	}

	err = app.models.Product.Update(categoryId, &product)
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	app.sendResponse(w, response{"product": product}, http.StatusCreated)
}

// List single product.
//
// Produces:
//   - application/json
//
// Parameters:
//
//   - name: category_id
//     in: path
//     description: Id of category
//     required: true
//     type: integer
//     format: int32
//
//   - name: product_id
//     in: path
//     description: Id of product
//     required: true
//     type: integer
//     format: int32
//
// responses:
//
//	200: ProductResponse
//
//swagger:route GET /categories/{category_id}/products/{product_id} Products GetProductItem
//swagger:response
func (app *application) getProductHandler(w http.ResponseWriter, r *http.Request) {
	productId, err := app.convertToInt(chi.URLParam(r, "product_id"))
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	categoryId, err := app.convertToInt(chi.URLParam(r, "category_id"))
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	var product data.Product
	err = app.models.Product.Get(categoryId, productId, &product)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.sendResponse(w, response{"message": "Resource not found"}, http.StatusNotFound)
			return
		default:
			app.log(err)
			app.sendInternalServerErrorResponse(w)
			return
		}
	}

	app.sendResponse(w, response{"product": product}, http.StatusOK)
}

// Lists products of a category.
//
//	Consumes:
//	- application/json
//
// responses:
//
//	200: listProducts
//
//swagger:route GET /categories/{id}/products Products listProducts
//swagger:response
func (app *application) listProductHandler(w http.ResponseWriter, r *http.Request) {
	categoryId, err := app.convertToInt(chi.URLParam(r, "category_id"))
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	if categoryId < 1 {
		app.sendResponse(w, response{"message": "Category id should not be less than 1"}, http.StatusUnprocessableEntity)
		return
	}

	searchTerm := ""
	if len(r.URL.Query().Get("search_term")) > 0 {
		searchTerm = r.URL.Query().Get("search_term")
	}
	

	var page int32
	page = 0
	if len(r.URL.Query().Get("page")) > 0 {
		page, err = app.convertToInt(r.URL.Query().Get("page"))
		if err != nil {
			app.log(err)
			app.sendInternalServerErrorResponse(w)
			return
		}
	}

	limit := int32(20)
	offset := (page - 1) * limit
	products, err := app.models.Product.GetAll(categoryId, limit, offset, searchTerm)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.sendResponse(w, response{}, http.StatusNotFound)
		default:
			app.log(err)
			app.sendInternalServerErrorResponse(w)
		}
		return
	}
	for _, product := range products {
		product.MainPicture = app.GenerateFileUrl(product.MainPicture)
	}

	totalCount, err := app.models.Product.Count(int(categoryId))
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	app.sendResponse(w, response{"products": products, "total": totalCount}, http.StatusOK)
}

//	swagger:route DELETE /categories/{category_id}/products/{product_id} Products deleteProductItem
//	Delete single product.
//
//	Produces:
//		- application/json
// 	Parameters:
//		+ name: category_id
//		in: path
//		description: Id of category
//		required: true
//		type: integer
//		format: int32
//
//		+ name: product_id
//		in: path
//		description: Id of product
//		required: true
//		type: integer
//		format: int32

// responses:
//
//	204: []
//
// swagger:response
func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	productId, err := app.convertToInt(chi.URLParam(r, "product_id"))
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	categoryId, err := app.convertToInt(chi.URLParam(r, "category_id"))
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	if categoryId < 1 {
		app.sendResponse(w, response{"message": "Category id can not be less than 1"}, http.StatusUnprocessableEntity)
		return
	}
 
	if productId < 1 {
		app.sendResponse(w, response{"message": "Product id can not be less than 1"}, http.StatusUnprocessableEntity)
		return
	}

	err = app.models.Product.Delete(categoryId, productId)
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	app.sendResponse(w, response{}, http.StatusNoContent)
}
