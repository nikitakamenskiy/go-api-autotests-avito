package advertisement

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/managers/advertisement"
	advModels "api-tests-template/internal/managers/advertisement/models"
	advList "api-tests-template/internal/managers/advertisements"
	advListModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData      models.LoginOkResponse
	currentDirName string
	createdAdvID   string
}

const (
	invalidRequestError         = "invalid_request"
	invalidinternalError        = "internal_error"
	unauthorizedError           = "unauthorized"
	invalidRequestErrorAssert   = "Неправильная ошибка"
	invalidRequestMessageAssert = "Неправильное сообщение об ошибке"
)

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Получаем текущую директорию")
	var err error
	s.currentDirName, err = os.Getwd()
	require.NoError(s.T(), err, "невозможно получить текущую директорию")

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownTest() {
	if s.createdAdvID != "" {
		base.Precondition("Удаляем созданное объявление с ID: " + s.createdAdvID)
		advertisement.DeleteAdvertisement(s.T(), s.loginData.Token, s.createdAdvID)
		s.createdAdvID = ""
	}
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// TestPostAdvertisementPositive проверяет создание объявления со всеми параметрами и его доступность
func (s *TestSuite) TestPostAdvertisementPositive() {
	params := newAdvertisementParams()
	photos := newAdvertisementPhotos(s.currentDirName)

	var advCreateResponse advModels.Advertisement
	s.Run("Создаём объявление со всеми параметрами", func() {
		advCreateResponse = advertisement.CreateAdvertisement(s.T(), s.loginData.Token, params, photos)
		s.createdAdvID = advCreateResponse.Id
	})

	s.Run("Проверяем, что все поля объявления соответствуют отправленным данным", func() {
		require.Equal(s.T(), params["title"], advCreateResponse.Title, "title не совпадает")
		require.Equal(s.T(), params["description"], advCreateResponse.Description, "description не совпадает")
		require.Equal(s.T(), params["price"], strconv.FormatFloat(advCreateResponse.Price, 'f', -1, 64), "price не совпадает")
		require.Equal(s.T(), params["quantity"], strconv.Itoa(advCreateResponse.Quantity), "quantity не совпадает")
		require.Len(s.T(), advCreateResponse.Photos, 3, "загружены не все фотографии")
	})

	var advertisementGet advModels.Advertisement
	s.Run("Получаем объявление по ID", func() {
		advertisementGet = advertisement.GetAdvertisement(s.T(), advCreateResponse.Id)
	})

	s.Run("Проверяем что вернулось тоже объявление что и при создании", func() {
		advertisementGetWithoutUser := advertisementGet
		advertisementGetWithoutUser.User = models.User{}
		require.Equal(s.T(), advCreateResponse, advertisementGetWithoutUser, "объекты объявлений не совпадают")
	})

	var advertisementPhotos []advModels.Photo
	s.Run("Получаем фотографии объявления", func() {
		advertisementPhotos = advertisement.GetAdvertisementPhotos(s.T(), s.loginData.Token, advCreateResponse.Id)
	})

	s.Run("Проверяем что фотографии совпадают с загруженными", func() {
		require.Len(s.T(), advertisementPhotos, 3, "количество фотографий не совпадает")
		for i, photo := range advCreateResponse.Photos {
			require.Equal(s.T(), photo, advertisementPhotos[i], "фотография не совпадает")
		}
	})

	var advertisementsList advListModels.AdvertisementsList
	s.Run("Получаем список объявлений  по полному названию", func() {
		advertisementsList = advList.GetAdvertisements(s.T(), advCreateResponse.Title)
	})

	s.Run("Проверяем что созданное объявление появилось в поиске", func() {
		require.Len(s.T(), advertisementsList.Items, 1, "количество найденных объявлений не равно 1")
		require.Equal(s.T(), advertisementGet, advertisementsList.Items[0], "вернулось другое объявление")
	})
}

// TestPostAdvertisementNoTitleNegative проверяет ошибку при создании объявления без title
func (s *TestSuite) TestPostAdvertisementNoTitleNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление без title", func() {
		params := newAdvertisementParams()
		delete(params, "title")
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, params, newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "title is required", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementEmptyTitleNegative проверяет ошибку при создании объявления с пустым title
func (s *TestSuite) TestPostAdvertisementEmptyTitleNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с пустым title", func() {
		params := newAdvertisementParams()
		params["title"] = ""
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, params, newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "title is required", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementLongTitleNegative проверяет ошибку при создании объявления с title длиннее 255 символов
func (s *TestSuite) TestPostAdvertisementLongTitleNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с title длиннее 255 символов", func() {
		params := newAdvertisementParams()
		params["title"] = utils.RandomString(256)
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, params, newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "title must not exceed 255 characters", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementWhitespaceTitleNegative проверяет ошибку при создании объявления с title из пробелов
// TODO: Known issue: bug001 - сервер принимает title из пробелов, ожидается 400.
func (s *TestSuite) TestPostAdvertisementWhitespaceTitleNegative() {
	s.T().Skip("Bug001: сервер принимает title состоящий из пробелов вместо возврата 400")

	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с title из пробелов", func() {
		params := newAdvertisementParams()
		params["title"] = "     "
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, params, newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "title is required", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementNoDescriptionNegative проверяет ошибку при создании объявления без description
func (s *TestSuite) TestPostAdvertisementNoDescriptionNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление без description", func() {
		params := newAdvertisementParams()
		delete(params, "description")
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, params, newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "description is required", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementEmptyDescriptionNegative проверяет ошибку при создании объявления с пустым description
func (s *TestSuite) TestPostAdvertisementEmptyDescriptionNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с пустым description", func() {
		params := newAdvertisementParams()
		params["description"] = ""
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, params, newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "description is required", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementLongDescriptionNegative проверяет ошибку при создании объявления с description длиннее 1000 символов
func (s *TestSuite) TestPostAdvertisementLongDescriptionNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с description длиннее 1000 символов", func() {
		params := newAdvertisementParams()
		params["description"] = utils.RandomString(1001)
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, params, newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "description must be between 1 and 1000 characters", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementWhitespaceDescriptionNegative проверяет ошибку при создании объявления с description из пробелов
// TODO: bug002 - сервер принимает description состоящий из пробелов, ожидается 400
func (s *TestSuite) TestPostAdvertisementWhitespaceDescriptionNegative() {
	s.T().Skip("Bug002: сервер принимает description состоящий из пробелов вместо возврата 400")

	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с description из пробелов", func() {
		params := newAdvertisementParams()
		params["description"] = "     "
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, params, newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "description is required", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementNoPhotosNegative проверяет ошибку при создании объявления без фото
func (s *TestSuite) TestPostAdvertisementNoPhotosNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление без фото", func() {
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, newAdvertisementParams(), []string{})
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "at least 1 photo is required", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementTooManyPhotosNegative проверяет ошибку при создании объявления с более чем 3 фото
func (s *TestSuite) TestPostAdvertisementTooManyPhotosNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с 4 фото", func() {
		photos := newAdvertisementPhotos(s.currentDirName)
		// добавляем 4 фотографию
		photos = append(photos, filepath.Join(s.currentDirName, "images", "advertisementImage1.jpg"))
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, newAdvertisementParams(), photos)
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "maximum 3 photos allowed", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementInvalidPhotoFormatNegative проверяет ошибку при создании объявления с файлом неверного формата
func (s *TestSuite) TestPostAdvertisementInvalidPhotoFormatNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с файлом .txt вместо фото", func() {
		photos := []string{filepath.Join(s.currentDirName, "images", "incorrectImage.txt")}
		errResponse = advertisement.CreateAdvertisementError500(s.T(), s.loginData.Token, newAdvertisementParams(), photos)
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidinternalError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "Failed to create advertisement with photos", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementLargePhotoNegative проверяет ошибку при создании объявления с фото > 5MB
func (s *TestSuite) TestPostAdvertisementLargePhotoNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с фото размером > 5MB", func() {
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, newAdvertisementParams(), newLargePhoto(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "photo 1 exceeds maximum size of 5MB", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementEmptyPhotoNegative проверяет ошибку при создании объявления с пустым файлом фото
// TODO: bug005 - сервер принимает пустой файл (0 байт) вместо возврата 400
func (s *TestSuite) TestPostAdvertisementEmptyPhotoNegative() {
	s.T().Skip("Bug005: сервер принимает пустой файл фото (0 байт) вместо возврата 400")

	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с пустым файлом фото", func() {
		photos := []string{filepath.Join(s.currentDirName, "images", "emptyImage.jpg")}
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, newAdvertisementParams(), photos)
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "photo is empty", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementNegativePriceNegative проверяет ошибку при создании объявления с отрицательной ценой
func (s *TestSuite) TestPostAdvertisementNegativePriceNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с отрицательной ценой", func() {
		params := newAdvertisementParams()
		params["price"] = "-1"
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, params, newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "price must be non-negative", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementStringPriceNegative проверяет ошибку при создании объявления со строкой в качестве цены
func (s *TestSuite) TestPostAdvertisementStringPriceNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление со строкой в качестве цены", func() {
		params := newAdvertisementParams()
		params["price"] = "abc"
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, params, newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "invalid price format", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementVeryLargePriceNegative проверяет ошибку при создании объявления с очень большой ценой
// TODO: bug003 - сервер возвращает 500 с сообщением про фото вместо валидационной ошибки цены
func (s *TestSuite) TestPostAdvertisementVeryLargePriceNegative() {
	s.T().Skip("Bug003: при цене 999999999999999 сервер возвращает 500 с сообщением про фото вместо ошибки валидации цены")

	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с очень большой ценой", func() {
		params := newAdvertisementParams()
		params["price"] = "999999999999999"
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, params, newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "price is too large", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementZeroQuantityNegative проверяет ошибку при создании объявления с quantity = 0
func (s *TestSuite) TestPostAdvertisementZeroQuantityNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с quantity = 0", func() {
		params := newAdvertisementParams()
		params["quantity"] = "0"
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, params, newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "quantity must be at least 1", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementNegativeQuantityNegative проверяет ошибку при создании объявления с отрицательным quantity
func (s *TestSuite) TestPostAdvertisementNegativeQuantityNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с отрицательным quantity", func() {
		params := newAdvertisementParams()
		params["quantity"] = "-1"
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, params, newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "quantity must be at least 1", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementFloatQuantityNegative проверяет ошибку при создании объявления с дробным quantity
func (s *TestSuite) TestPostAdvertisementFloatQuantityNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с дробным quantity", func() {
		params := newAdvertisementParams()
		params["quantity"] = "1.5"
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, params, newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "invalid quantity format", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementStringQuantityNegative проверяет ошибку при создании объявления со строкой в качестве quantity
func (s *TestSuite) TestPostAdvertisementStringQuantityNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление со строкой в качестве quantity", func() {
		params := newAdvertisementParams()
		params["quantity"] = "abc"
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, params, newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "invalid quantity format", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementOverflowQuantityNegative проверяет ошибку при создании объявления с переполнением quantity
// TODO: bug004 - сервер принимает очень большое quantity и возвращает переполненное значение вместо ошибки валидацииf
/*
ДЕТАЛИЗАЦИЯ БАГА ДЛЯ ИСПРАВЛЕНИЯ.
У Нас 2 разных поведения:
1. Слишком большое, но парсится
- при quantity = 2147483647(граница int32) объявление создаётся корректно
- при quantity = 2147483648(превышение границы) сохраняется некорректное значение (отрицательное) объявление создается -2147483648
- при ещё больших значениях (например 21474836480) может сохраняться 0 и создаваться объявление
2.Слишком большое, не парсится
- при quantity = 99999999999999999999 "message": "invalid quantity format"
Предположительно проблема возникает при превышении допустимого диапазона значений quantity типа int32
По ТЗ ограничение количество 100.

*/
func (s *TestSuite) TestPostAdvertisementOverflowQuantityNegative() {
	s.T().Skip("Bug004: при quantity > допустимого диапазона сервер сохраняет некорректные значения (отрицательные или 0) вместо ошибки валидации")

	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с quantity больше допустимого диапазона", func() {
		params := newAdvertisementParams()
		params["quantity"] = "999999999999"
		errResponse = advertisement.CreateAdvertisementErrorBadRequest(s.T(), s.loginData.Token, params, newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), invalidRequestError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "invalid quantity format", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementNoTokenNegative проверяет ошибку при создании объявления без токена авторизации
func (s *TestSuite) TestPostAdvertisementNoTokenNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление без токена авторизации", func() {
		errResponse = advertisement.CreateAdvertisementErrorNoToken(s.T(), newAdvertisementParams(), newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), unauthorizedError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "Authorization header required", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementEmptyBearerNegative проверяет ошибку при создании объявления с пустым Bearer токеном
func (s *TestSuite) TestPostAdvertisementEmptyBearerNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с пустым Bearer токеном", func() {
		errResponse = advertisement.CreateAdvertisementErrorEmptyBearer(s.T(), newAdvertisementParams(), newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), unauthorizedError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "Invalid authorization header format", errResponse.Message, invalidRequestMessageAssert)
	})
}

// TestPostAdvertisementInvalidTokenNegative проверяет ошибку при создании объявления с невалидным токеном
func (s *TestSuite) TestPostAdvertisementInvalidTokenNegative() {
	var errResponse advModels.CreateErrorResponse
	s.Run("Создаём объявление с невалидным токеном", func() {
		errResponse = advertisement.CreateAdvertisementErrorUnauthorized(s.T(), newAdvertisementParams(), newAdvertisementPhotos(s.currentDirName))
	})

	s.Run("Проверяем ошибку", func() {
		require.Equal(s.T(), unauthorizedError, errResponse.Error, invalidRequestErrorAssert)
		require.Equal(s.T(), "Invalid or expired token", errResponse.Message, invalidRequestMessageAssert)
	})
}
