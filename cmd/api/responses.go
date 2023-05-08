package main

import data "ecom-api/inernal/data/models"

// A list of products
//
//swagger:response listProducts
type ProductsResponse struct {
	//in: body
	Body []data.Product
}

// A single product
//
//swagger:response
type ProductResponse struct {
	//in: body
	Body data.Product
}

// A list of categories
//
//swagger:response listCategories
type CategoriesResponse struct {
	//in: body
	Body []data.Category
}

// A single category
//
//swagger:response
type CategoryResponse struct {
	//in: body
	Body data.Category
}

// Returns token
//
//swagger:response
type AuthResponse struct {
	//in: body
	Body struct {
		Token string `json:"token"`
	}
}
