package body

// LoginBody - User login structure
type LoginBody struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}
