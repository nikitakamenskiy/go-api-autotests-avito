package models

import authModels "api-tests-template/internal/managers/auth/models"

type Advertisement struct {
	CreatedAt   string          `json:"created_at,omitempty"`
	Description string          `json:"description,omitempty"`
	Id          string          `json:"id,omitempty"`
	Photos      []Photo         `json:"photos,omitempty"`
	Price       float64         `json:"price,omitempty"`
	Quantity    int             `json:"quantity,omitempty"`
	Title       string          `json:"title,omitempty"`
	UpdatedAt   string          `json:"updated_at,omitempty"`
	User        authModels.User `json:"user,omitempty"`
	UserId      string          `json:"user_id,omitempty"`
}

type Photo struct {
	CreatedAt string `json:"created_at,omitempty"`
	Id        string `json:"id,omitempty"`
	SortOrder int    `json:"sort_order,omitempty"`
	Url       string `json:"url,omitempty"`
}

type CreateErrorResponse struct {
	Details string `json:"details"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
