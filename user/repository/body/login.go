package body

// LoginBody - User login structure
type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
