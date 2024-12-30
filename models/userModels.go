package models

type User struct {
	ID           int     `json:"id"`
	FirstName    *string `json:"first_name" validate:"required,min=2,max=100"`
	LastName     *string `json:"last_name" validate:"required,min=2,max=100"`
	Password     *string `json:"password" validate:"required,min=6"`
	Email        *string `json:"email" validate:"required,email"`
	Phone        *string `json:"phone" validate:"required,min=10,max=10"`
	Token        *string `json:"token"`
	UserType     *string `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	RefreshToken *string `json:"refresh_token"`
	UserId       string  `json:"user_id"`
}
