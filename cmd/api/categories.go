package main

import (
	data "ecom-api/inernal/data/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(int64(app.config.multipartFormSize))
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	var parentId *int
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

func (app *application) listCategoryHandler(w http.ResponseWriter, r *http.Request) {
	limit := 20
	offset := 0
	page := 1

	err := r.ParseMultipartForm(int64(app.config.multipartFormSize))
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	if len(r.PostForm.Get("page")) > 0 {
		page, err = app.convertToInt(r.PostForm.Get("page"))
	}

	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	offset = (page - 1) * limit
	returnTotalCount, _ := strconv.ParseBool(r.Form.Get("includeTotalCount"))
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
