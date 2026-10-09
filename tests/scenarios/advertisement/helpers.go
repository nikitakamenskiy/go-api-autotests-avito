package advertisement

import (
	"path/filepath"
	"strconv"

	"api-tests-template/internal/utils"
)

// newAdvertisementParams генерирует уникальные параметры для создания объявления
func newAdvertisementParams() map[string]string {
	return map[string]string{
		"title":       "title_" + utils.RandomString(14),
		"description": "description_" + utils.RandomString(14),
		"price":       strconv.Itoa(utils.RandomInt(1, 10000)),
		"quantity":    strconv.Itoa(utils.RandomInt(1, 100)),
	}
}

// newAdvertisementPhotos возвращает пути к тестовым фотографиям
func newAdvertisementPhotos(dir string) []string {
	base := filepath.Join(dir, "images")

	return []string{
		filepath.Join(base, "advertisementImage1.jpg"),
		filepath.Join(base, "advertisementImage2.jpg"),
		filepath.Join(base, "advertisementImage3.jpg"),
	}
}

// newLargePhoto возвращает путь к файлу фото размером > 5MB
func newLargePhoto(dir string) []string {
	return []string{filepath.Join(dir, "images", "largeImage.png")}
}
