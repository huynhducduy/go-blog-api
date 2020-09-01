package user

type User struct {
	Id         *int64  `json:"id,omitempty"`
	Username   *string `json:"username"`
	Password   *string `json:"password,omitempty"`
	OldPassword *string `json:"old_password,omitempty"`
	Name       *string `json:"name"`
	Email	   *string `json:"email"`
	Role      *string `json:"role"`
}