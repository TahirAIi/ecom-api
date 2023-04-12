package main

// swagger:parameters listProducts
type ProductsParams struct {
	// in: path
	// required: true
	Id int `json:"id"`

	// in: query
	Limit int `json:"limit"`
}

//
//swagger:parameters createProduct updateProduct
type ProductBody struct {
	//in: formData
	//required: true
	Title string `json:"title"`

	// in:formData
	//required: false
	Description string `json:"description"`

	// in:formData
	//required: true
	Price int `json:"price"`

	//swagger:file
	//in:formData
	//required: false
	MainPicture []byte `json:"main_picture"`
}
