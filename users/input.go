package users

type RegisterUserInput struct {
	Username            string `form:"username" binding:"required" json:"username"`
	Email               string `form:"email" binding:"required,email" json:"email"`
	Password            string `form:"password" binding:"required" json:"password"`
	BusinessID          string `form:"business_id" json:"business_id"`
	Role                string `form:"role" json:"role"`
	FullName            string `form:"full_name" binding:"required" json:"full_name"`
	Telepon             string `form:"telepon" binding:"required" json:"telepon"`
	Address             string `form:"address" json:"address"`
	Picture             string `form:"picture" json:"picture"`
	BusinessInheritance bool   `form:"business_in_heritance" json:"business_inheritance"`
}

type LoginInput struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type CheckEmailInput struct {
	Email string `form:"email" binding:"required,email"`
}

type CheckUsernameInput struct {
	Username string `form:"username" binding:"required"`
}
