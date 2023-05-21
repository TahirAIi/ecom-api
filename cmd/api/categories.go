package main

import (
	"database/sql"
	data "ecom-api/inernal/data/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Creates a category.
//
// Consumes:
//   - multipart/form-data
//
// Produces:
//   - application/json
//
// Parameters:
//
//	CategoryBody
//
// responses:
//
//	200: CategoryResponse
//
//swagger:route POST /categories Categories createCategory
//swagger:response
func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(int64(app.config.multipartFormSize))
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	var parentId *int32
	var description *string
	title := r.FormValue("title")
	tempDescription := r.FormValue("description")
	tempParentId := r.FormValue("parent_id")

	if len(tempDescription) != 0 {
		description = &tempDescription
	}

	if len(tempParentId) > 0 {
		*parentId, err = app.convertToInt(tempParentId)
		if err != nil {
			app.log(err)
			app.sendInternalServerErrorResponse(w)
			return
		}
	}

	if len(title) == 0 {
		app.sendResponse(w, response{"message": "Title can not be empty"}, http.StatusUnprocessableEntity)
	}

	category := &data.Category{
		Title:       title,
		ParentId:    parentId,
		Description: description,
	}
	err = app.models.Category.Insert(category)

	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	err = app.sendResponse(w, response{"category": category}, http.StatusOK)
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
	}
}

// swagger:route GET /categories Categories listCategories
// Lists all the categories.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// responses:
//
//	200: listCategories
//
//swagger:response

func (app *application) listCategoryHandler(w http.ResponseWriter, r *http.Request) {
	limit := int32(20)
	offset := int32(0)
	page := int32(1)

	var err error
	if len(chi.URLParam(r, "page")) > 0 {
		page, err = app.convertToInt(chi.URLParam(r, "page"))
	}

	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	offset = (page - 1) * limit
	returnTotalCount, _ := strconv.ParseBool(chi.URLParam(r,"includeTotalCount"))
	response := make(map[string]interface{})

	if returnTotalCount != false {
		totalCategories, err := app.models.Category.GetTotalCount()
		if err != nil {
			app.log(err)
			app.sendInternalServerErrorResponse(w)
			return
		}

		response["total"] = totalCategories
	}

	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	categories, err := app.models.Category.GetAll(limit, offset)
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	response["categories"] = categories

	app.sendResponse(w, response, http.StatusOK)

}

//swagger:route GET /categories/{category_id} Categories GetCategoryItem
//Get single category.
//
//Produces:
//	- application/json
// Parameters:
//	+ name: category_id
//	in: path
//	description: Id of category
//	required: true
//	type: integer
//	format: int32
//
// responses:
//
//	200: CategoryResponse
//
//swagger:response

func (app *application) getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryId, err := app.convertToInt(chi.URLParam(r, "category_id"))
	if err != nil {
		app.sendInternalServerErrorResponse(w)
	}

	category, err := app.models.Category.Get(categoryId)
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

	app.sendResponse(w, response{"category": category}, http.StatusOK)
}

//swagger:route PATCH /categories/{category_id} Categories updateCategory
// Updates a category.
//
// Consumes:
//   - multipart/form-data
//
// Produces:
//   - application/json
//
// Parameters:
//
//	CategoryBody
//
// responses:
//
//	200: CategoryResponse
//
//swagger:response

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(int64(app.config.multipartFormSize))

	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	id, err := app.convertToInt(chi.URLParam(r, "category_id"))

	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	category, err := app.models.Category.Get(id)
	if err != nil {
		app.sendResponse(w, response{}, http.StatusNotFound)
		return
	}

	title := r.FormValue("title")
	if r.Form.Has("description") {
		description := r.FormValue("description")
		category.Description = &description
	}

	if len(title) > 0 {
		category.Title = title
	}

	err = app.models.Category.Update(category)
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	app.sendResponse(w, response{"category": category}, http.StatusOK)
	return
}

//swagger:route DELETE /categories/{category_id} Categories deleteCategoryItem
//Delete single category.
//
//Produces:
//	- application/json
// Parameters:
//	+ name: category_id
//	in: path
//	description: Id of category
//	required: true
//	type: integer
//	format: int32

// responses:
//
//	204: []
//
//swagger:response
func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.convertToInt(chi.URLParam(r, "category_id"))

	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	err = app.models.Category.Delete(id)

	if err != nil {
		app.log(err)
		return
	}
	app.sendResponse(w, response{}, http.StatusNoContent)

	return
}
