package paypal

import (
	"fmt"
	"time"
)

type (

	// CreateProductResp struct response after creating a product
	CreateProductResp struct {
		ID          string    `json:"id,omitempty"`
		Name        string    `json:"name,omitempty"`
		Description string    `json:"description,omitempty"`
		Type        string    `json:"type,omitempty"`
		Category    string    `json:"category,omitempty"`
		ImageURL    string    `json:"image_url,omitempty"`
		HomeURL     string    `json:"home_url,omitempty"`
		CreateTime  time.Time `json:"create_time,omitempty"`
		UpdateTime  time.Time `json:"update_time,omitempty"`
		Links       []Link    `json:"links,omitempty"`
	}

	// ProductListInput struct
	ProductListInput struct {
		Page          string `json:"page,omitempty"`      //Default: 0.
		PageSize      string `json:"page_size,omitempty"` //Default: 10.
		TotalRequired string `json:"total_required,omitempty"`
	}

	// ListProductsResp struct
	ListProductsResp struct {
		Products   []ListProductsItem `json:"products,omitempty"`
		TotalItems string             `json:"total_items,omitempty"`
		TotalPages string             `json:"total_pages,omitempty"`
		Links      []Link             `json:"links,omitempty"`
	}
)

// GetProduct gets a single catalog product
func (c *client) GetProduct(productID string) (*ProductResp, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/catalogs/products/"+productID), nil)
	response := &ProductResp{}
	if err != nil {
		return response, err
	}
	err = c.sendWithAuth(req, response)
	return response, err
}

// ListProducts lists all catalog products
func (c *client) ListProducts(cplp *ProductListInput) (*ListProductsResp, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/catalogs/products"), nil)
	q := req.URL.Query()
	q.Add("page", cplp.Page)
	q.Add("page_size", cplp.PageSize)
	q.Add("total_required", cplp.TotalRequired)
	req.URL.RawQuery = q.Encode()
	response := &ListProductsResp{}
	if err != nil {
		return response, err
	}
	err = c.sendWithAuth(req, response)
	return response, err
}

// CreateProduct creates a new catalog product
func (c *client) CreateProduct(product *CreateProductInput) (*CreateProductResp, error) {
	req, err := c.newRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/catalogs/products"), product)
	response := &CreateProductResp{}
	if err != nil {
		return response, err
	}
	err = c.sendWithAuth(req, response)
	return response, err
}
