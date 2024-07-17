package model

type Res struct {
	SkuId       int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description" swaggertype:"string"`
	Size        string `json:"size" swaggertype:"string"`
	Color       string `json:"color" swaggertype:"string"`
}

func Resource(product *ChildProductModel) *Res {
	if product == nil {
		return &Res{}
	}
	resource := &Res{
		SkuId:       product.ChildProduct.SkuId,
		Name:        product.ParentProduct.Name,
		Description: product.ParentProduct.Description,
		Size:        product.ChildProduct.Size,
		Color:       product.ChildProduct.Color,
	}

	return resource
}
