package main

import data "ecom-api/inernal/data/models"

// A list of products
//swagger:response ProductsResponse
type ProductsResponse struct {
	// All movies in system
	//
	//in: body
	Body	[]data.Product
}

// A single product
//swagger:response
type ProductResponse struct {
	//in: body
	Body data.Product
}

