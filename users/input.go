package users

type RegisterUserInput struct {
	Username string `form:"username" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
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
