package usecase

import "time"

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

// TokenDetailBody - implement domain.TokenDetail interface
type TokenDetailBody struct {
	tokenID string
	userID  string
}

// JwtParams -
type JwtParams struct {
	issueTime         time.Time
	expirationSeconds time.Duration
	secret            []byte
	jwtID             string
	audience          string
	issuer            string
}

// GetExpirationSeconds -
func (j *JwtParams) GetExpirationSeconds() time.Duration {
	return j.expirationSeconds
}

// GetIssuer -
func (j *JwtParams) GetIssuer() string {
	return j.issuer
}

// GetJwtID -
func (j *JwtParams) GetJwtID() string {
	return j.jwtID
}

// GetAudience -
func (j *JwtParams) GetAudience() string {
	return j.audience
}

// GetSecret -
func (j *JwtParams) GetSecret() []byte {
	return j.secret
}

// GetIssueTime -
func (j *JwtParams) GetIssueTime() time.Time {
	return j.issueTime
}

// NewJwtParams - Create New jwtParams
func NewJwtParams(issueTime time.Time, expirationSeconds time.Duration, secret []byte, issuer string, jwtID string, audience string) *JwtParams {
	return &JwtParams{
		issueTime:         issueTime,
		expirationSeconds: expirationSeconds,
		secret:            secret,
		jwtID:             jwtID,
		audience:          audience,
		issuer:            issuer,
	}
}

// TokensBody - implement domain.Token interface
type TokensBody struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// GetAccessToken - implement domain.Token interface
func (t *TokensBody) GetAccessToken() string {
	return t.AccessToken
}

// GetRefreshToken - implement domain.Token interface
func (t *TokensBody) GetRefreshToken() string {
	return t.RefreshToken
}
