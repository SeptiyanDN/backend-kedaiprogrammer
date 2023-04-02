package businesses

type AddBusinessInput struct {
	Business_name string `form:"business_name" binding:"required"`
	Domain        string `form:"domain" binding:"required"`
}
