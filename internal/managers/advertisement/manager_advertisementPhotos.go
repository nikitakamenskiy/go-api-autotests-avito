package advertisement

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	advertisementPhotos "api-tests-template/internal/client/http/advertisementPhotos"
	"api-tests-template/internal/managers/advertisement/models"
)

func GetAdvertisementPhotos(t *testing.T, token, advId string) []models.Photo {
	getPhotosRawResponse := advertisementPhotos.HttpGetAdvertisementPhotos(t, token, advId, http.StatusOK)
	defer getPhotosRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(getPhotosRawResponse.Body)
	require.NoError(t, errResponse)

	var photosResponse []models.Photo
	err := json.Unmarshal(respArray, &photosResponse)
	require.NoError(t, err)

	return photosResponse
}
