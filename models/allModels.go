package models

type Category struct {
	ID           int64  `json:"id"`
	CategoryName string `json:"category_name"`
}

type Product struct {
	ID                 int64  `json:"id"`
	ProductName        string `json:"product_name"`
	GeneralDescription string `json:"general_description"`
	CategoryID         int64  `json:"category_id"`
}

type Tenant struct {
	ID                int64  `json:"id"`
	TenantName        string `json:"tenant_name"`
	RegisteredAddress string `json:"registered_address"`
	TenantContact     string `json:"tenant_contact"`
	TenantEmail       string `json:"tenant_email"`
}

type Hub struct {
	ID             int64  `json:"id"`
	TenantID       int64  `json:"tenant_id"`
	ManagerName    string `json:"manager_name"`
	ManagerContact string `json:"manager_contact"`
	ManagerEmail   string `json:"manager_email"`
}

type Seller struct {
	ID            int64  `json:"id"`
	HubID         int64  `json:"hub_id"`
	TenantID      int64  `json:"tenant_id"`
	SellerName    string `json:"seller_name"`
	SellerContact string `json:"seller_contact"`
	SellerEmail   string `json:"seller_email"`
}

type SKU struct {
	ID          int64  `json:"id"`
	HubID       int64  `json:"hub_id"`
	SellerID    int64  `json:"seller_id"`
	ProductID   int64  `json:"product_id"`
	Images      string `json:"images"`
	Description string `json:"description"`
	UnitPrice   int    `json:"unit_price"`
	Fragile     bool   `json:"fragile"`
	Dimensions  string `json:"dimensions"`
}

type Address struct {
	ID           int64  `json:"id"`
	EntityID     int64  `json:"entity_id"`
	EntityType   string `json:"entity_type"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	Pincode      string `json:"pincode"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
}

type ValidateOrderRequest struct {
	SKUID string `json:"sku_id"`
	HubID string `json:"hub_id"`
}

type ValidationResponse struct {
	IsValid bool
	Error   error
}
