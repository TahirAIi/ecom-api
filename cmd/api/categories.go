package main

import (
	data "ecom-api/inernal/data/models"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Println(err)
	}

	var parentId *int
	var description *string
	title := r.PostForm.Get("title")
	tempDescription := r.PostForm.Get("description")
    tempParentId := r.PostForm.Get("parent_id") 
    
	if len(tempDescription) != 0 {
		description = &tempDescription
	}

	if len(tempParentId) > 0 {
		*parentId = app.convertToInt(tempParentId, 0)
	}

	if len(title) == 0 {
		app.sendResponse(w, response{"message": "Title can not be empty"}, http.StatusUnprocessableEntity)
	}

	category := &data.Category{
		Uuid:        uuid.New().String(),
		Title:       title,
		ParentId:    parentId,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = app.models.Category.Insert(category)

	if err != nil {
		app.logger.Println(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	err = app.sendResponse(w, response{"category": category}, http.StatusOK)
	if err != nil {
		app.logger.Println(err)
		app.sendInternalServerErrorResponse(w)
	}
}

func (app *application) listCategoryHandler(w http.ResponseWriter, r *http.Request) {
	limit := 20
	offset := 0
	err := r.ParseForm()
	if err != nil {
		app.logger.Fatal(err)
	}

	page := app.convertToInt(r.PostForm.Get("page"), 0)

	offset = page * limit
	returnTotalCount, _ := strconv.ParseBool(r.Form.Get("includeTotalCount"))
	response := make(map[string]interface{})
	if returnTotalCount != false {
		totalCategories, err := app.models.Category.GetTotalCount()
		if err != nil {
			app.logger.Println(err)
			app.sendInternalServerErrorResponse(w)
			return
		}

		response["total"] = totalCategories
	}

	if err != nil {
		app.logger.Println(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	categories, err := app.models.Category.GetAll(limit, offset)
	if err != nil {
		app.logger.Println(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	response["categories"] = categories

	app.sendResponse(w, response, http.StatusOK)

}

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(2048)

	if err != nil {
		app.logger.Println("Unable to parse form")
	}

	uuid := chi.URLParam(r, "uuid")

	category, err := app.models.Category.Get(uuid)
	if err != nil {
		app.logger.Printf("No category found for %s", uuid)
		return
	}

	title := r.MultipartForm.Value["title"][0]
	description := r.PostForm.Get("description")

	if len(title) > 0 {
		category.Title = title
	}

	if len(description) > 0 {
		category.Description = &description
	}

	err = app.models.Category.Update(category)
	if err != nil {
		app.logger.Printf("Unable to update category %s", err)
		return
	}

	app.sendResponse(w, response{"category": category}, http.StatusOK)
	return
}

func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	err := app.models.Category.Delete(uuid)

	if err != nil {
		app.logger.Printf("Unable to delete %s", err)
		return
	}
	app.sendResponse(w, response{}, http.StatusNoContent)

	return
}
