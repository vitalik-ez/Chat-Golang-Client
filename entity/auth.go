package entity

type signUp struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewSignUp(input map[string]string) *signUp {
	return &signUp{
		Name:     input["name"],
		Email:    input["email"],
		Password: input["password"],
	}
}

type signIn struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewSignIn(input map[string]string) *signIn {
	return &signIn{
		Email:    input["email"],
		Password: input["password"],
	}
}
