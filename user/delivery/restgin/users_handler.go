package restgin

import (
	"net/http"

	"github.com/alibug/go-identity-entry/domain"
	tokenBody "github.com/alibug/go-identity-entry/token/repository/body"
	userBody "github.com/alibug/go-identity-entry/user/repository/body"
	"github.com/alibug/go-identity-utils/config"
	"github.com/alibug/go-identity-utils/status"
	"github.com/gin-gonic/gin"
)

// UsersHandler  represent the httphandler for user
type UsersHandler struct {
	userUsecase   domain.UserUsecase
	tokensUsecase domain.TokensUseCase
	cookieConfig  config.CookieConfig
}

// NewUsersHandler represent the httphandler for user
func NewUsersHandler(route *gin.Engine, uuc domain.UserUsecase, tuc domain.TokensUseCase, cc config.CookieConfig) {
	handler := &UsersHandler{
		userUsecase:   uuc,
		tokensUsecase: tuc,
		cookieConfig:  cc,
	}

	// ⚠️ login
	route.POST("/login", handler.mustNotLoginInterceptor(), handler.Login)
	route.POST("/register", handler.mustNotLoginInterceptor(), handler.RegisterUser)
	route.POST("/logout", handler.Logout)
}

// Logout -
func (u *UsersHandler) Logout(c *gin.Context) {
	// 1、从 cookie 中 获取 token
	tokens := u.getTokenFromCookie(c)
	if tokens == nil {
		c.JSON(status.GetStatusCode(status.ErrUnauthorized), status.ResponseError{Message: "You are not logged in"})
		return
	}

	// 3、Delete access token
	ctx := c.Request.Context()
	err := u.tokensUsecase.CheckTokensAndLogout(ctx, tokens)
	if err != nil {
		c.JSON(status.GetStatusCode(err), status.ResponseError{Message: err.Error()})
		return
	}

	// 4、正确删除 token 后， 清理 cookie
	u.clearAccessTokenInCookie(c)
	u.clearUserInfoInCookie(c)
	c.JSON(http.StatusOK, gin.H{"logout": true})
}

// Login 实现登录 ⚠️ 多次 重复登录的问题 是否要检查 ？？？？
func (u *UsersHandler) Login(c *gin.Context) {
	var body userBody.LoginBody
	// 1、 校验 body 格式
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(status.GetStatusCode(err), status.ResponseError{Message: err.Error()})
		return
	}

	// 3、校验用户名密码
	ctx := c.Request.Context()
	user, err := u.userUsecase.CheckAccountAndPassUC(ctx, body.Account, body.Password)
	if err != nil {
		c.JSON(status.GetStatusCode(err), status.ResponseError{Message: err.Error()})
		return
	}

	// 4.1 、创建 Tokens
	tokens, err := u.tokensUsecase.CreateTokens(ctx, user.GetUserID())
	if err != nil {
		c.JSON(status.GetStatusCode(err), status.ResponseError{Message: err.Error()})
		return
	}

	// 5、写入 cookie
	u.setTokenToCookie(c, tokens)
	u.setUserInfoToCookie(c, user)

	// 6、⚠️ 此处是临时性的 设置 返回结果
	c.JSON(http.StatusOK, gin.H{"displayname": user.GetDisplayName()})
}

// GetByID will get user by given id
func (u *UsersHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	ctx := c.Request.Context()

	user, err := u.userUsecase.GetByIDUC(ctx, id)
	if err != nil {
		c.JSON(status.GetStatusCode(err), status.ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// RegisterUser will register user
func (u *UsersHandler) RegisterUser(c *gin.Context) {
	var body userBody.RegisterBody
	// 1、 校验 body 格式
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(status.GetStatusCode(err), status.ResponseError{Message: err.Error()})
		return
	}

	ctx := c.Request.Context()
	err := u.userUsecase.RegisterUserUC(ctx, &body)
	if err != nil {
		c.JSON(status.GetStatusCode(err), status.ResponseError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"ok": true})
}

func (u *UsersHandler) setTokenToCookie(c *gin.Context, tokens domain.Tokens) {
	c.SetCookie(u.cookieConfig.GetAccessTokenField(), tokens.GetAccessToken(), u.cookieConfig.GetAccessTokenMaxAge(), "/", u.cookieConfig.GetDomain(), u.cookieConfig.GetSecure(), u.cookieConfig.GetHTTPOnly())
	c.SetCookie(u.cookieConfig.GetRefreshTokenField(), tokens.GetRefreshToken(), u.cookieConfig.GetRefreshTokenMaxAge(), "/", u.cookieConfig.GetDomain(), u.cookieConfig.GetSecure(), u.cookieConfig.GetHTTPOnly())
}

func (u *UsersHandler) setUserInfoToCookie(c *gin.Context, user domain.User) {
	c.SetCookie(u.cookieConfig.GetDisplayNameField(), user.GetDisplayName(), u.cookieConfig.GetRefreshTokenMaxAge(), "/", u.cookieConfig.GetDomain(), u.cookieConfig.GetSecure(), false)
	// c.SetCookie("userID", user.GetUserID(), u.cookieConfig.GetRefreshTokenMaxAge(), "/", u.cookieConfig.GetDomain(), u.cookieConfig.GetSecure(), u.cookieConfig.GetHTTPOnly())
	c.SetCookie(u.cookieConfig.GetUserIDField(), user.GetUserID(), 86400, "/", u.cookieConfig.GetDomain(), u.cookieConfig.GetSecure(), u.cookieConfig.GetHTTPOnly())
}

func (u *UsersHandler) clearUserInfoInCookie(c *gin.Context) {
	c.SetCookie(u.cookieConfig.GetDisplayNameField(), "", -1, "/", u.cookieConfig.GetDomain(), u.cookieConfig.GetSecure(), false)
	c.SetCookie(u.cookieConfig.GetUserIDField(), "", -1, "/", u.cookieConfig.GetDomain(), u.cookieConfig.GetSecure(), u.cookieConfig.GetHTTPOnly())
}

func (u *UsersHandler) clearAccessTokenInCookie(c *gin.Context) {
	c.SetCookie(u.cookieConfig.GetAccessTokenField(), "", -1, "/", u.cookieConfig.GetDomain(), u.cookieConfig.GetSecure(), u.cookieConfig.GetHTTPOnly())
	c.SetCookie(u.cookieConfig.GetRefreshTokenField(), "", -1, "/", u.cookieConfig.GetDomain(), u.cookieConfig.GetSecure(), u.cookieConfig.GetHTTPOnly())
}

func (u *UsersHandler) getTokenFromCookie(c *gin.Context) domain.Tokens {
	accessToken, err := c.Cookie(u.cookieConfig.GetAccessTokenField())
	if err != nil {
		accessToken = ""
	}
	refreshToken, err := c.Cookie(u.cookieConfig.GetRefreshTokenField())
	if err != nil {
		refreshToken = ""
	}
	if accessToken == "" && refreshToken == "" {
		return nil
	}
	return &tokenBody.TokenBody{AccessToken: accessToken, RefreshToken: refreshToken}
}

func (u *UsersHandler) mustNotLoginInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := u.getTokenFromCookie(c)
		if token != nil {
			c.JSON(status.GetStatusCode(status.ErrForbidden), status.ResponseError{Message: "You have logged in"})
			c.Abort()
		}
		c.Next()
	}
}
