package services

type FindCategoryInput struct {
	BusinessID int `form:"business_id" binding:"required"`
}
type AddServiceInput struct {
	ServiceName string `form:"service_name" binding:"required" json:"service_name"`
	BusinessID  string `form:"business_id" binding:"required" json:"business_id"`
}

type BodyGetRequest struct {
	Search         string `form:"search"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	B2bToken       string `json:"b2b_token"`
	Draw           string `json:"draw"`
	OrderColumn    string `json:"order_column"`
	OrderDirection string `json:"order_direction"`
}
