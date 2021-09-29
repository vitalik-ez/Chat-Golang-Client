package entity

type signUp struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewSignUp(input map[string]string) *signUp {
	return &signUp{
		Name:     input["name"],
		Email:    input["email"],
		Password: input["password"],
	}
}

type signIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewSignIn(input map[string]string) *signIn {
	return &signIn{
		Email:    input["email"],
		Password: input["password"],
	}
}
