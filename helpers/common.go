package helpers

import (
	"encoding/base64"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func Create_Dir(path string) error {
	var errnya error
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		errnya = err
	}
	return errnya
}

func B64Decode(b64string string, filepath string) error {
	var tr_err error
	dec, err := base64.StdEncoding.DecodeString(b64string)
	if err != nil {
		tr_err = err
	}

	f, err := os.Create(filepath)
	if err != nil {
		tr_err = err
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		tr_err = err
	}

	if err := f.Sync(); err != nil {
		tr_err = err
	}
	return tr_err
}

func BcryptEncode(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}
