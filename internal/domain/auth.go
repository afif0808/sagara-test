package domain

type LoginCredentials struct {
	Identity string `json:"identity"` // e.g : ( email , username)
	Password string `json:"password"`
}
