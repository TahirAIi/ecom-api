package main

// swagger:parameters listProducts
type ProductsParams struct {
	//
    // in: path
	// required: true
	Id int `json:"id"`
	
	// in: query
	Limit int `json:"limit"`
}