package configreader

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

const mongoURLTemplate = "mongodb://%s:%s@%s/%s"
const redisURLTemplate = "redis://:%s@%s/%s"

// 设置 要从环境变量中要读取的 变量前缀（自动转为大写)
func init() {
	viper.SetEnvPrefix("config")

	// 绑定 配置文件名 环境变量
	viper.BindEnv("name")
	// 绑定 配置文件路径 环境变量
	viper.BindEnv("path")
	// 绑定 配置文件类型 环境变量
	viper.BindEnv("type")

	// 设置 配置文件名 缺省值
	viper.SetDefault("name", "config")
	// 设置 配置文件路径 缺省值

	// ⚠️ 下一次 把 缺省路径 放到 绝对路径 /conf
	viper.SetDefault("path", "/conf")
	// 设置 配置文件
	viper.SetDefault("type", "yaml")
	configFileName := viper.GetString("name")
	configFilePath := viper.GetString("path")
	configFileType := viper.GetString("type")

	// 确定要读取的配置文件
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(configFilePath)

	// 处理读取错误
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Read config file error: %s", err)
	}
}

// ReadMongoConfig - return mongodb conn url
func ReadMongoConfig() string {
	// 实际读取相应参数
	mongoUser := viper.GetString("mongo.user")
	mongoPass := viper.GetString("mongo.pass")
	mongoHost := viper.GetString("mongo.host")
	mongoDbName := viper.GetString("mongo.database")
	return fmt.Sprintf(mongoURLTemplate, mongoUser, mongoPass, mongoHost, mongoDbName)
}

// ReadRedisConfig - return redis conn url
func ReadRedisConfig() string {
	// 读取 redis 参数
	redisHost := viper.GetString("redis.host")
	redisPass := viper.GetString("redis.pass")
	redisDbName := viper.GetString("redis.database")
	return fmt.Sprintf(redisURLTemplate, redisPass, redisHost, redisDbName)
}

/*
type TokenConfig interface {
	GetAccessTokenSecret() []byte
	GetRefreshTokenSecret() []byte
	GetIssuer() string
	GetAccessExpirationSeconds() time.Duration
	GetRefreshExpirationSeconds() time.Duration
}
*/

// TokenConfig -
type TokenConfig struct {
	accessExpirationSeconds  time.Duration
	refreshExpirationSeconds time.Duration
	accessTokenSecret        []byte
	refreshTokenSecret       []byte
	issuer                   string
}

// GetIssuer -
func (t *TokenConfig) GetIssuer() string {
	return t.issuer
}

// GetAccessTokenSecret -
func (t *TokenConfig) GetAccessTokenSecret() []byte {
	return t.accessTokenSecret
}

// GetRefreshTokenSecret -
func (t *TokenConfig) GetRefreshTokenSecret() []byte {
	return t.refreshTokenSecret
}

// GetAccessExpirationSeconds -
func (t *TokenConfig) GetAccessExpirationSeconds() time.Duration {
	return t.accessExpirationSeconds
}

// GetRefreshExpirationSeconds -
func (t *TokenConfig) GetRefreshExpirationSeconds() time.Duration {
	return t.refreshExpirationSeconds
}

// ReadTokenConfig - read Token Config
func ReadTokenConfig() *TokenConfig {
	accessTokenSecret := []byte(viper.GetString("token.accessSecret"))
	refreshTokenSecret := []byte(viper.GetString("token.refreshSecret"))
	accessExpirationSeconds := time.Duration(viper.GetInt("token.expiresSeconds"))
	refreshExpirationSeconds := time.Duration(viper.GetInt("token.refreshSeconds"))
	issuer := viper.GetString("token.issuer")
	return &TokenConfig{
		accessTokenSecret:        accessTokenSecret,
		refreshTokenSecret:       refreshTokenSecret,
		accessExpirationSeconds:  accessExpirationSeconds,
		refreshExpirationSeconds: refreshExpirationSeconds,
		issuer:                   issuer,
	}
}

// CookieConfig - contaiin cookie setting
type CookieConfig struct {
	secure             bool
	httpOnly           bool
	accessTokenMaxAge  int
	refreshTokenMaxAge int
	domain             string
}

// GetAccessTokenMaxAge -
func (c *CookieConfig) GetAccessTokenMaxAge() int {
	return c.accessTokenMaxAge
}

// GetRefreshTokenMaxAge -
func (c *CookieConfig) GetRefreshTokenMaxAge() int {
	return c.refreshTokenMaxAge
}

// GetDomain -
func (c *CookieConfig) GetDomain() string {
	return c.domain
}

// GetSecure -
func (c *CookieConfig) GetSecure() bool {
	return c.secure
}

// GetHTTPOnly -
func (c *CookieConfig) GetHTTPOnly() bool {
	return c.httpOnly
}

/*
secure             bool
	httpOnly           bool
	accessTokenMaxAge  int
	refreshTokenMaxAge int
	domain             string
*/

// ReadCookieConfig -
func ReadCookieConfig() *CookieConfig {
	secure := viper.GetBool("cookie.secure")
	httpOnly := viper.GetBool("cookie.httpOnly")
	accessTokenMaxAge := viper.GetInt("cookie.accessTokenMaxAge")
	refreshTokenMaxAge := viper.GetInt("cookie.refreshTokenMaxAge")
	domain := viper.GetString("cookie.domain")
	return &CookieConfig{
		secure:             secure,
		httpOnly:           httpOnly,
		accessTokenMaxAge:  accessTokenMaxAge,
		refreshTokenMaxAge: refreshTokenMaxAge,
		domain:             domain,
	}
}
