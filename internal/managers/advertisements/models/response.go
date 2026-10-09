package models

import "api-tests-template/internal/managers/advertisement/models"

type AdvertisementsList struct {
	HasNext    bool                   `json:"has_next"`
	Items      []models.Advertisement `json:"items"`
	NextCursor string                 `json:"next_cursor"`
	Total      int                    `json:"total"`
}
