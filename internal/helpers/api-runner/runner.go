package api_runner

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/steinfletcher/apitest"
)

const (
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36"
)

type ApiTest struct {
	token     string
	userAgent string
	host      string
}

func New() *ApiTest {
	return &ApiTest{
		userAgent: defaultUserAgent,
		host:      os.Getenv("API_URL"),
	}
}

func GetRunner() *ApiTest {
	return New()
}

// Auth Установка Bearer token-а для прохождения авторизации на приватные ручки
func (at *ApiTest) Auth(token string) *ApiTest {
	at.token = token
	return at
}

func (at *ApiTest) Create() *apitest.APITest {
	apitestNew := apitest.New().EnableNetworking()
	if os.Getenv("DEBUG") == "true" {
		apitestNew = apitestNew.Debug()
	}
	return apitestNew.
		Intercept(func(request *http.Request) {
			request.Header.Set("User-Agent", at.userAgent)

			if err := mergeServiceUrls(at.host, request.URL); err != nil {
				panic(fmt.Sprintf("Не удалось применить хост из API_URL: %s", err.Error()))
			}

			if len(at.token) != 0 {
				request.Header.Add("Authorization", "Bearer "+at.token)
			}
		})
}

func mergeServiceUrls(host string, r *url.URL) error {
	if host == "" {
		return fmt.Errorf("API_URL не задан: проверьте переменные окружения")
	}

	urlParsed, err := url.Parse(host)
	if err != nil {
		return fmt.Errorf("host cannot be parsed: %s", err.Error())
	}

	if urlParsed.Path != "" {
		r.Path = path.Join(urlParsed.Path, r.Path)
	}

	r.Scheme = urlParsed.Scheme
	r.Host = urlParsed.Host
	return nil
}
