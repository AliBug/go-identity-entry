package body

// TokenBody - implement domain.Token interface
type TokenBody struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// GetAccessToken - implement domain.Token interface
func (t *TokenBody) GetAccessToken() string {
	return t.AccessToken
}

// GetRefreshToken - implement domain.Token interface
func (t *TokenBody) GetRefreshToken() string {
	return t.RefreshToken
}

// TokenDetailBody - implement domain.TokenDetail interface
type TokenDetailBody struct {
	tokenID string
	userID  string
}

// NewTokenDetailBody - new a TokenDetailBody
func NewTokenDetailBody(tokenID string, userID string) *TokenDetailBody {
	return &TokenDetailBody{tokenID, userID}
}

// GetTokenID - implement domain.TokenDetail interface
func (t *TokenDetailBody) GetTokenID() string {
	return t.tokenID
}

// GetUserID - implement domain.TokenDetail interface
func (t *TokenDetailBody) GetUserID() string {
	return t.userID
}
