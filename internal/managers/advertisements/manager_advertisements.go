package advertisements

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
)

// GetAdvertisements получает список объявлений по поисковому запросу и возвращает десериализованный ответ
func GetAdvertisements(t *testing.T, search string) models.AdvertisementsList {
	getAdvertisementsRawResponse := advertisements.HttpGetAdvertisements(t, search, http.StatusOK)
	defer getAdvertisementsRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(getAdvertisementsRawResponse.Body)
	require.NoError(t, errResponse)

	advListResponse := models.AdvertisementsList{}
	err := json.Unmarshal(respArray, &advListResponse)
	require.NoError(t, err)

	return advListResponse
}
