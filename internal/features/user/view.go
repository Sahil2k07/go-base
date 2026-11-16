package user

type CurrentUserDataResponse struct {
	ID          uint     `json:"id"`
	Email       string   `json:"email"`
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	Phone       string   `json:"phone"`
	Gender      string   `json:"gender"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}
