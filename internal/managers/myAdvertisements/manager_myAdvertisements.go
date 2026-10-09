package myAdvertisements

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/myAdvertisements"
)

// GetMyAdvertisements возвращает список объявлений авторизованного пользователя и проверяет ожидаемый статус-код
func GetMyAdvertisements(t *testing.T, token string, expectedStatusCode int, limit int) string {
	myAdvertisementsRawResponse := myAdvertisements.HttpGetMyAdvertisements(t, token, limit)
	defer myAdvertisementsRawResponse.Body.Close()

	require.Equalf(t, expectedStatusCode, myAdvertisementsRawResponse.StatusCode,
		"HTTP status code должен быть %d", expectedStatusCode)

	myAdvertisementsResponse, errResponse := io.ReadAll(myAdvertisementsRawResponse.Body)
	require.NoError(t, errResponse)

	return string(myAdvertisementsResponse)
}
