package randomuser

type UserResponse struct {
	Results []User `json:"results"`
}

func (u *UserResponse) Error() string {
	return "api error message"
}

type User struct {
	Name UserName `json:"name"`
}

type UserName struct {
	First string `json:"first"`
	Last  string `json:"last"`
}
