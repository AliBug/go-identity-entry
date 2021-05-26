module github.com/alibug/go-identity-entry

go 1.16

require (
	github.com/alibug/go-identity-utils v0.1.10
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.7.1
	github.com/go-redis/redis/v8 v8.8.2
	github.com/google/uuid v1.2.0
	go.mongodb.org/mongo-driver v1.5.2
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
)

// replace github.com/alibug/go-identity-utils => ../go-identity-utils
