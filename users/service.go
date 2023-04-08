package users

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"kedaiprogrammer/helpers"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

type Services interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	IsUsernameAvailable(input CheckUsernameInput) (bool, error)
	GetUserByUUID(UUID string) (User, error)
	GetUserByToken(token string) (User, error)
	SaveToken(UUID string, token string) (User, error)
}

type services struct {
	repository Repository
}

func NewServices(repository Repository) *services {
	return &services{repository}
}

func (s *services) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	picture := ""
	if input.Picture != "" {
		decodedImage, err := base64.StdEncoding.DecodeString(input.Picture)
		// file := bytes.NewReader(decodedImage)
		fileType := http.DetectContentType(decodedImage)
		if fileType != "image/jpeg" && fileType != "image/jpg" && fileType != "image/png" {
			return user, err
		}
		_, waktu := helpers.TimeInLocal("Asia/Jakarta")
		tz := waktu.Format("20060102150405")
		waktuFile := tz
		picture = fmt.Sprintf("%s/%s.%s", input.Username, waktuFile, strings.Split(fileType, "/")[1])
	}
	user.Uuid = helpers.GenerateUUID()
	user_info := make(map[string]interface{})
	user_info["full_name"] = input.FullName
	user_info["telepon"] = input.Telepon
	user_info["Address"] = input.Address
	user_info["picture"] = picture
	user_info["business_inheritance"] = input.BusinessInheritance
	dataJson, _ := json.Marshal(user_info)
	JsonRaw := json.RawMessage(dataJson)
	user.UserInfo = postgres.Jsonb{JsonRaw}
	user.Username = input.Username
	user.Email = input.Email
	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	user.Password = string(password)
	// err = helpers.SaveObjectToS3(s3, "staging-smartpatrol/absensi/", imageIn, file)
	// if err != nil {
	// 	return false, err
	// }
	fmt.Println(user)
	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func (s *services) Login(input LoginInput) (User, error) {
	username := input.Username
	password := input.Password

	user, err := s.repository.FindByUsername(username)

	if err != nil {
		return user, err
	}
	if user.Uuid == "" {
		return user, errors.New("User Not Found on the Database")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *services) IsEmailAvailable(input CheckEmailInput) (bool, error) {
	email := input.Email
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}
	if user.Uuid == "" {
		return true, nil
	}
	return false, nil

}

func (s *services) IsUsernameAvailable(input CheckUsernameInput) (bool, error) {
	username := input.Username
	user, err := s.repository.FindByUsername(username)
	if err != nil {
		return false, err
	}
	if user.Uuid == "" {
		return true, nil
	}
	return false, nil
}

func (s *services) GetUsersByUUID(UUID string) (User, error) {
	user, err := s.repository.FindByUUID(UUID)
	if err != nil {
		return user, err
	}
	if user.Uuid == "" {
		return user, errors.New("User Not Found on the Database")
	}
	return user, nil

}

func (s *services) GetUserByToken(token string) (User, error) {
	user, err := s.repository.FindByToken(token)
	if err != nil {
		return user, err
	}
	if user.Uuid == "" {
		return user, errors.New("User Not Found on the Database")
	}
	return user, nil

}

func (s *services) SaveToken(UUID string, token string) (User, error) {
	user, err := s.repository.FindByIdAndUpdateToken(UUID, token)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *services) GetUserByUUID(UUID string) (User, error) {
	user, err := s.repository.FindByUUID(UUID)
	if err != nil {
		return user, err
	}
	return user, nil
}
