package users

import "time"

type UserFormatter struct {
	Uuid       string    `json:"uuid"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Token      string    `json:"token"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

func FormatUser(user User, token string, generateUUID string) UserFormatter {
	formatter := UserFormatter{
		Uuid:       generateUUID,
		Username:   user.Username,
		Email:      user.Email,
		Token:      token,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	return formatter
}

type loginFormatter struct {
	Token string `json:"token"`
}

func FormatUserLogin(token string) loginFormatter {
	formatter := loginFormatter{
		Token: token,
	}
	return formatter
}
