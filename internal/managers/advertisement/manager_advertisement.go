package advertisement

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisement"
	"api-tests-template/internal/managers/advertisement/models"
)

// CreateAdvertisement создаёт объявление и возвращает десериализованный ответ
func CreateAdvertisement(t *testing.T, token string, formDataParams map[string]string, photos []string) models.Advertisement {
	createAdvertisementRawResponse := advertisement.HttpPostAdvertisement(t, token, http.StatusCreated, formDataParams, photos)
	defer createAdvertisementRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(createAdvertisementRawResponse.Body)
	require.NoError(t, errResponse)

	advResponse := models.Advertisement{}
	err := json.Unmarshal(respArray, &advResponse)
	require.NoError(t, err)

	return advResponse
}

// GetAdvertisement получает объявление по ID и возвращает десериализованный ответ
func GetAdvertisement(t *testing.T, advertisementID string) models.Advertisement {
	getAdvertisementRawResponse := advertisement.HttpGetAdvertisement(t, http.StatusOK, advertisementID)
	defer getAdvertisementRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(getAdvertisementRawResponse.Body)
	require.NoError(t, errResponse)

	advResponse := models.Advertisement{}
	err := json.Unmarshal(respArray, &advResponse)
	require.NoError(t, err)

	return advResponse
}

// DeleteAdvertisement удаляет объявление по ID
func DeleteAdvertisement(t *testing.T, token, advertisementID string) {
	advertisement.HttpDeleteAdvertisement(t, token, advertisementID)
}

// CreateAdvertisementError отправляет запрос на создание объявления и возвращает ошибку с указанным статусом
func CreateAdvertisementError(t *testing.T, token string, statusCode int, formDataParams map[string]string, photos []string) models.CreateErrorResponse {
	createAdvertisementRawResponse := advertisement.HttpPostAdvertisement(t, token, statusCode, formDataParams, photos)
	defer createAdvertisementRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(createAdvertisementRawResponse.Body)
	require.NoError(t, errResponse)

	errResp := models.CreateErrorResponse{}
	err := json.Unmarshal(respArray, &errResp)
	require.NoError(t, err)

	return errResp
}

// CreateAdvertisementErrorBadRequest отправляет запрос на создание объявления и возвращает ошибку 400
func CreateAdvertisementErrorBadRequest(t *testing.T, token string, formDataParams map[string]string, photos []string) models.CreateErrorResponse {
	return CreateAdvertisementError(t, token, http.StatusBadRequest, formDataParams, photos)
}

// CreateAdvertisementErrorUnauthorized отправляет запрос с невалидным токеном и возвращает ошибку 401
func CreateAdvertisementErrorUnauthorized(t *testing.T, formDataParams map[string]string, photos []string) models.CreateErrorResponse {
	return CreateAdvertisementError(t, "invalid_token", http.StatusUnauthorized, formDataParams, photos)
}

// CreateAdvertisementError500 отправляет запрос на создание объявления и возвращает ошибку 500
func CreateAdvertisementError500(t *testing.T, token string, formDataParams map[string]string, photos []string) models.CreateErrorResponse {
	return CreateAdvertisementError(t, token, http.StatusInternalServerError, formDataParams, photos)
}

// CreateAdvertisementErrorNoToken отправляет запрос без токена авторизации
func CreateAdvertisementErrorNoToken(t *testing.T, formDataParams map[string]string, photos []string) models.CreateErrorResponse {
	return CreateAdvertisementError(t, "", http.StatusUnauthorized, formDataParams, photos)
}

// CreateAdvertisementErrorEmptyBearer отправляет запрос с пустым Bearer токеном
func CreateAdvertisementErrorEmptyBearer(t *testing.T, formDataParams map[string]string, photos []string) models.CreateErrorResponse {
	return CreateAdvertisementError(t, " ", http.StatusUnauthorized, formDataParams, photos)
}
