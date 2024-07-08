package model

type ChildSku struct {
	SkuId           int
	ParentProductID int
	Size            string
	Color           string
	Price           string
	Quantity        int
}
