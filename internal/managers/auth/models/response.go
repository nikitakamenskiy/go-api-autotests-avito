package models

type LoginOkResponse struct {
	Token string `json:"token,omitempty"`
	User  User   `json:"user,omitempty"`
}

type User struct {
	BirthDate string `json:"birth_date,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	Id        string `json:"id,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}
