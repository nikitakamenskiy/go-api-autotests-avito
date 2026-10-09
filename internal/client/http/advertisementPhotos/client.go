package advertisementPhotos

import (
	"net/http"
	"strings"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

func HttpGetAdvertisementPhotos(t *testing.T, token, advId string, statusCode int) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(strings.ReplaceAll(path.AdvertisementsPhotosPath, "{id}", advId)).
		Expect(t).
		Status(statusCode).
		End().Response
}
