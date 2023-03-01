package jwt

import (
	logger "go_rest_api_with_mysql/pkg/log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger = logger.GetLogger().Sugar()

func CreateSignedJwt(signingMethod string, claims map[string]interface{}, signKey string) (string, error) {
	signMethod := jwt.GetSigningMethod(signingMethod)

	claimsMap := jwt.MapClaims{}
	// copy passed claims values
	for k, v := range claims {
		claimsMap[k] = v
	}

	// set iat and jti values
	claimsMap["iat"] = time.Now()
	claimsMap["jat"] = time.Now()

	token := jwt.NewWithClaims(signMethod, claimsMap)
	jwtSignedToken, err := token.SignedString([]byte(signKey))
	if err != nil {
		return "", nil
	}
	log.Debugf("jwtSignedToken: %s", jwtSignedToken)
	return jwtSignedToken, nil
}
