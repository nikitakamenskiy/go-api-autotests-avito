package advertisements

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

func HttpGetAdvertisements(t *testing.T, search string, statusCode int) *http.Response {
	return apiRunner.GetRunner().Create().Get(path.AdvertisementsPath).
		Query("search", search).
		Expect(t).
		Status(statusCode).
		End().Response
}
