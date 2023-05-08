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

// swagger:parameters listCategories
type CategoriesParams struct {
	// in: path
	// required: true
	Id int `json:"id"`

	// in: query
	Limit int `json:"limit"`
}
//swagger:parameters createCategory updateCategory
type CategoryBody struct {
	//in: formData
	//required: true
	Title string `json:"title"`

	// in:formData
	//required: false
	Description string `json:"description"`

	// in:formData
	//required: true
	ParentId int `json:"parent_id"`
}

//swagger:parameters authUser
type AuthBody struct {
	//in: formData
	//required: true
	Email string `json:"email"`

	// in:formData
	//required: true
	Password string `json:"password"`
}
