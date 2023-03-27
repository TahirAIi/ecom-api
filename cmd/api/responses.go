package main

import data "ecom-api/inernal/data/models"

// A list of products
//swagger:response
type ProductsResponse struct {
	// All movies in system
	//
	//in: body
	Body	[]data.Product
}

