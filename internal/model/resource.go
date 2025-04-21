package model

type ResourceType string

const (
	ResourceTypeBrand        ResourceType = "BRAND"
	ResourceTypeComment      ResourceType = "COMMENT"
	ResourceTypeProductModel ResourceType = "PRODUCT_MODEL"
	ResourceTypeProduct      ResourceType = "PRODUCT"
	ResourceTypeRefund       ResourceType = "REFUND"
)
