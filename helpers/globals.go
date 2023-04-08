package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
)

func GenerateUUID() string {
	u := uuid.New()
	return u.String()
}

func ParseTime(input string) time.Time {
	layout := "2006-01-02 15:04:05"
	loc, err := time.LoadLocation("Asia/Jakarta") // contoh lokasi
	if err != nil {
		fmt.Println(err)
	}
	parseResult, _ := time.ParseInLocation(layout, input, loc)
	return parseResult
}
func GetExcelHeaders(f string) ([]string, []map[int]string) {
	xlsx, err := excelize.OpenFile(f)
	if err != nil {
		log.Fatal("ERROR", err.Error())
	}

	map_header := []map[int]string{}

	sheetName := "Sheet1"

	data, err_ := xlsx.GetRows(sheetName)
	//dataHeader, errHeader := xlsx.GetH

	if err_ != nil {
		log.Panic(err_.Error())
	}

	//fmt.Println("len data excel ",data[0])

	for i := 0; i < len(data); i++ {
		if i > 0 {
			data_body := map[int]string{}

			len_data := len(data[i])

			for j := 0; j < len_data; j++ {

				if len(data[i]) < len(data[0]) {

					len_data := len(data[0]) - len(data[i])

					for k := 1; k <= len_data; k++ {
						data_body[j+k] = ""
					}
					data_body[j] = data[i][j]
				}

				data_body[j] = data[i][j]

			}

			//map_rows = append(map_rows, data_body)
			map_header = append(map_header, data_body)
		}

	}

	dataHeader := data[0]
	dataRows := map_header

	return dataHeader, dataRows
}

func DataInSlices(a string, list []string) (bool, int) {
	cnt := 0
	for _, b := range list {
		if cast.ToString(b) == a {
			cnt++
		}
	}
	return true, cnt
}

func EncryptUUID(text string) (string, error) {
	plaintext := []byte(text)
	key := []byte(viper.GetString("SECRET.SECRET_KEY_JWT"))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(cryptoText string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher([]byte(viper.GetString("SECRET.SECRET_KEY_JWT")))
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	fmt.Println("cipper" + string(ciphertext))
	return string(ciphertext), nil
}

func ValidateUUID(encrypt string, uuid string) bool {
	decripted, err := Decrypt(encrypt)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return decripted == uuid
}
func TimeInLocal(local string) (time.Time, time.Time) {
	t := time.Now()

	location, err := time.LoadLocation(local)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println("Location : ", location, " Time : ", t.In(location)) // America/New_York
	resp := t.In(location)

	return t, resp
}
