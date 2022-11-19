package helper

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/go-playground/validator/v10"
	"github.com/itchyny/base58-go"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message, status string, code int, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

func FormatValidationError(err error) []string {
	var errors []string

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}

func sha256Of(originUrl string) ([]byte, error) {
	algorithm := sha256.New()
	_, err := algorithm.Write([]byte(originUrl))
	if err != nil {
		return nil, err
	}
	return algorithm.Sum(nil), nil
}

func base58Encoded(bytes []byte) (string, error) {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func GenerateShortUrl(originUrl string) (string, error) {
	urlHashBytes, err := sha256Of(originUrl)
	if err != nil {
		return "", err
	}

	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString, err := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	if err != nil {
		return "", err
	}

	return finalString[:8], nil
}
