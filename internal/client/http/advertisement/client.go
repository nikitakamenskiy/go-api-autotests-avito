package advertisement

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

func HttpPostAdvertisement(t *testing.T, token string, statusCode int, formDataParams map[string]string, photosPath []string) *http.Response {
	request := apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath)

	for key, value := range formDataParams {
		request = request.MultipartFormData(key, value)
	}
	return request.
		MultipartFile("photos", photosPath...).
		Expect(t).
		Status(statusCode).
		End().Response
}

func HttpGetAdvertisement(t *testing.T, statusCode int, advertisementID string) *http.Response {
	return apiRunner.GetRunner().Create().Get(path.AdvertisementPath).
		Query("id", advertisementID).
		Expect(t).
		Status(statusCode).
		End().Response
}

func HttpDeleteAdvertisement(t *testing.T, token, advId string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", advId).
		Expect(t).
		Status(http.StatusNoContent).
		End().Response
}
