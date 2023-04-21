package categories

type FindCategoryInput struct {
	BusinessID int `form:"business_id" binding:"required"`
}
type AddCategoryInput struct {
	Category_name string `form:"category_name" binding:"required" json:"category_name"`
	ServiceID     string `form:"service_id" binding:"required" json:"service_id"`
	Tag           string `form:"tag" binding:"required" json:"tag"`
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
