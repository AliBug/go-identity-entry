package restgin

import (
	"net/http"

	"github.com/alibug/go-identity/domain"
	"github.com/alibug/go-identity/user/repository/body"
	"github.com/gin-gonic/gin"
)

// UserHandler  represent the httphandler for user
type UserHandler struct {
	UserUsecase domain.UserUsecase
}

// NewUserHandler represent the httphandler for user
func NewUserHandler(route *gin.Engine, uc domain.UserUsecase) {
	handler := &UserHandler{
		UserUsecase: uc,
	}
	route.POST("/login", handler.Login)
	route.POST("/register", handler.RegisterUser)
	route.GET("/users/:id", handler.GetByID)
	// e.DELETE("/articles/:id", handler.Delete)
}

// Login 实现登录
func (u *UserHandler) Login(c *gin.Context) {
	var body body.LoginBody
	// 1、 校验 body 格式
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	// 2、校验用户名密码
	ctx := c.Request.Context()
	user, err := u.UserUsecase.CheckUsernameAndPass(ctx, body.Username, body.Password)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	// 3、⚠️ 此处是临时性的 设置 返回结果
	c.JSON(http.StatusOK, user)
}

// GetByID will get user by given id
func (u *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	ctx := c.Request.Context()

	user, err := u.UserUsecase.GetByID(ctx, id)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// RegisterUser will register user
func (u *UserHandler) RegisterUser(c *gin.Context) {
	var body body.RegisterBody
	// 1、 校验 body 格式
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	ctx := c.Request.Context()
	err := u.UserUsecase.RegisterUser(ctx, &body)
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
	default:
		return http.StatusInternalServerError
	}
}
