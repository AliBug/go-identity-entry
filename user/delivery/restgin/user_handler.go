package restgin

import (
	"log"
	"net/http"

	"github.com/alibug/go-identity/domain"
	tokenBody "github.com/alibug/go-identity/token/repository/body"
	userBody "github.com/alibug/go-identity/user/repository/body"
	"github.com/gin-gonic/gin"
)

// UserHandler  represent the httphandler for user
type UserHandler struct {
	UserUsecase  domain.UserUsecase
	TokenUsecase domain.TokenUsecase
	cookieConfig domain.CookieConfig
}

// NewUserHandler represent the httphandler for user
func NewUserHandler(route *gin.Engine, uc domain.UserUsecase, tc domain.TokenUsecase, cookie domain.CookieConfig) {
	handler := &UserHandler{
		UserUsecase:  uc,
		TokenUsecase: tc,
		cookieConfig: cookie,
	}

	// ⚠️ login
	route.POST("/login", mustNotLoginInterceptor(), handler.Login)
	route.POST("/register", mustNotLoginInterceptor(), handler.RegisterUser)
	route.POST("/logout", handler.Logout)
	// route.GET("/users/:id", handler.GetByID)
	// e.DELETE("/articles/:id", handler.Delete)
}

// Logout -
func (u *UserHandler) Logout(c *gin.Context) {
	// 1、从 cookie 中 获取 token
	token := getTokenFromCookie(c)
	if token == nil {
		c.JSON(getStatusCode(domain.ErrUnauthorized), ResponseError{Message: "You are not logged in"})
		return
	}

	// 2、Delete token
	err := u.TokenUsecase.DeleteTokenUc(c, token)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	// 3、正确删除 token 后， 清理 cookie
	u.clearAccessTokenInCookie(c)
	u.clearUserInfoInCookie(c)
	c.JSON(http.StatusOK, gin.H{"logout": true})
}

// Login 实现登录 ⚠️ 多次 重复登录的问题 是否要检查 ？？？？
func (u *UserHandler) Login(c *gin.Context) {
	var body userBody.LoginBody
	// 1、 校验 body 格式
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	// 3、校验用户名密码
	ctx := c.Request.Context()
	user, err := u.UserUsecase.CheckUsernameAndPassUc(ctx, body.Username, body.Password)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	// 4、创建 Token
	token, err := u.TokenUsecase.CreateTokenUc(ctx, user.GetUserID())
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	// 5、写入 cookie
	u.setTokenToCookie(c, token)
	u.setUserInfoToCookie(c, user)

	// 6、⚠️ 此处是临时性的 设置 返回结果
	c.JSON(http.StatusOK, gin.H{"displayname": user.GetDisplayName()})
}

// GetByID will get user by given id
func (u *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	ctx := c.Request.Context()

	user, err := u.UserUsecase.GetByIDUc(ctx, id)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// RegisterUser will register user
func (u *UserHandler) RegisterUser(c *gin.Context) {
	var body userBody.RegisterBody
	// 1、 校验 body 格式
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	ctx := c.Request.Context()
	err := u.UserUsecase.RegisterUserUc(ctx, &body)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"ok": true})
}

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	// logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrBadParamInput:
		return http.StatusBadRequest
	case domain.ErrUnauthorized:
		return http.StatusUnauthorized
	case domain.ErrForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func (u *UserHandler) setTokenToCookie(c *gin.Context, token domain.Token) {
	c.SetCookie(domain.AccessTokenField, token.GetAccessToken(), u.cookieConfig.GetAccessTokenMaxAge(), "/", u.cookieConfig.GetDomain(), u.cookieConfig.GetSecure(), u.cookieConfig.GetHTTPOnly())
	c.SetCookie(domain.RefreshTokenField, token.GetRefreshToken(), u.cookieConfig.GetRefreshTokenMaxAge(), "/", u.cookieConfig.GetDomain(), u.cookieConfig.GetSecure(), u.cookieConfig.GetHTTPOnly())
}

func (u *UserHandler) setUserInfoToCookie(c *gin.Context, user domain.User) {
	c.SetCookie("displayname", user.GetDisplayName(), u.cookieConfig.GetRefreshTokenMaxAge(), "/", u.cookieConfig.GetDomain(), u.cookieConfig.GetSecure(), false)
}

func (u *UserHandler) clearUserInfoInCookie(c *gin.Context) {
	c.SetCookie("displayname", "", 0, "/", u.cookieConfig.GetDomain(), u.cookieConfig.GetSecure(), false)
	log.Println("clear user info")
}

func (u *UserHandler) clearAccessTokenInCookie(c *gin.Context) {
	c.SetCookie(domain.AccessTokenField, "", 0, "/", u.cookieConfig.GetDomain(), u.cookieConfig.GetSecure(), u.cookieConfig.GetHTTPOnly())
	c.SetCookie(domain.RefreshTokenField, "", 0, "/", u.cookieConfig.GetDomain(), u.cookieConfig.GetSecure(), u.cookieConfig.GetHTTPOnly())
}

func getTokenFromCookie(c *gin.Context) domain.Token {
	accessToken, err := c.Cookie(domain.AccessTokenField)
	if err != nil {
		accessToken = ""
	}
	refreshToken, err := c.Cookie(domain.RefreshTokenField)
	if err != nil {
		refreshToken = ""
	}
	if accessToken == "" && refreshToken == "" {
		return nil
	}
	return &tokenBody.TokenBody{AccessToken: accessToken, RefreshToken: refreshToken}
}

func mustNotLoginInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getTokenFromCookie(c)
		if token != nil {
			c.JSON(getStatusCode(domain.ErrForbidden), ResponseError{Message: "You have logged in"})
			c.Abort()
		}
		c.Next()
	}
}
