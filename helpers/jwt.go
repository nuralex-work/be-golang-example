package helpers

import (
	"be-golang/structs/general"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

var AppIDs = []interface{}{"WRdC75eHzvJT6FDw5fsfjX9TMuRCCXA3JbXhWcCXhB4VTKvU4ffCYCOH9P4x9nOaE34oeaIZ36c1bzFrMWS5CUWarmUHjUUeNzG4QXuPSAAIKHznN8fZOKKlNLsXuEwFODlyUYytGyz8umWsUCHGM6KR0JPSqPO7i2YQ2QeHAq9wa8RvZWzv4m2dh67T0bcU8e4zdfn9K4uCdWtozj9tb1sm3dTwfvvjN9e1ND678OWHRaCAxpV6zwj29Z6k2lu3g5JqPlUZ1aEC2p3Zir9j2PpSFpZeyGyPHzeNjR3MpzC32LySFXdQiYtTJdvRIPjtP9PKrlj3HsfYt42ucluuwhbsRMoMbKO6FI81cKToQj6YArKvOmBHnc0ItBwSMIplKOZ1iialRZCTb6WTm6oHFbK4YSdwHuiLGxkgHb4P9kz0kuU4wh6TqMY92tENIEGZqgg1hDA1tV6WdAhhwe044CcLMXpL6mRhAbL28dVI32BBYlf4oaiel56jTnZ1dSWF", "cXvkLsk1a7baKXziCeB4cJBBm6YKyKtr41uXGueQFQozx133XWM0iRiBkJqcSVfhWYxO0tDPzRU4lkFGh1m7GNfOPYDND6RGXWhC8C2R5Np7ROAXYtcYVYTRJRkH1uavXHnQomRZrfTRgySJUlEplrQczxer6s3vpj2VB0rV3X1PqRUAa0S7PKyBvXpLdjNsWBngY1CmeqhWgen5WhUwMGk8MhHUuh3V9dekVOwv2mCz9TaWj4wq4N369SPZPUkeTVLfjq7owFhYXi8tV27KpX0jRFZPdFHHHwGOY3qYZ9aRhKJfEThgDteHeHYkPAnISlmAm7h81dcpWtO6TbhmJbEX1CavljZ7sVVP05BEqov9N9PINU9qti3CsDRg1VzUpsNDEP8E0i3CyJmNQkohK0aRTV3Py4I0waAMXiVq3Vz3IcMB08aD4EiBpvoY89wGs3ay4HK2ZHvSQ2cSDjXvqVj0QuB9HnFBf7Z3vopBi4XoUuPbn3SXxhudfjDrJ4bW", "RV9i4QbAUe4qbMN+zYULbePP607WkNyLX8S3exX5QuhwB1CXSn95lBu8Jp981PQuD3efbOXUPY3e9U04lilGOHpKCROE/Ixi9bHz34sAOQG/wjfqA0lBbO3/IvacmPtDN2JuzRsOPUkYZOQMH0c+3azkr/Csn6h0qmXG9umdb4VTwEY9bwnrUIjWMrkeOd98cpEmCqNIWGGkeIR2bT6tCZbanuSUS+Ho7p0sS16p6dwrchUj7wjTBCszQmOH/NX5lohM/2LAUJtnEHruiqPjL47WXe6sWF6fl6a68TxENX7IDp1RLlN+Jd50zcAEvJGr26p9eJeRD+e92YjceywpxRkcCtZjmAZnN3F9dZdZImiTPTV6NslpJozA78WoZnsgIMxIBr+56ZRxxeR+WVAEEXuEqJL5m0GfPdoEZVTc4pX+gzUNKManUi19fs9/pbx0Ypq1tOaYgJ2p66Cpl9P6oOj0bprznhbPw2JtL6C3f6BNL62AY+6h4IPrLSKSi6VOhoNvTvCxMzF1S7ZeieiJSZeeWHuaeIQ1vcRzoTkryLRp4PLBew/rBK+91QOfJ1GNcGQlK1ibKIx8QQ+kb27RmK6BrCQpNTsKkaN2zxv00Lip2zMkl9ZGCAskNBYwDpmUQBilY+RFnBbLMwumrgALFEPO5C59gkZOlZfZElbxQTM"}
var JWTKeys = []interface{}{"iGSpuiPe6YubU7TRu0NZ8dHBUutMIXZvn194MHyfhFwQJ4VTNYNn1qErusfNdgTD", "r0EPsjAcrho8VHeSO2MYinhrnbFwgZIb", "a71b537ca697fd6493144dc73fd81231a2934023947777e362b6269364c7e9c6"}

func DecodeTokenJwt(jwtToken string) (jwt.MapClaims, bool) {

	// check if token contains Bearer
	if strings.Contains(jwtToken, "Bearer") {
		jwtToken = strings.Replace(jwtToken, "Bearer ", "", -1)
	}

	secret := []byte(os.Getenv("JWT_SECRET"))

	// Parse the token
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the token's signing method is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		log.Fatalf("Error parsing token: %v", err)
		return nil, false
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		return nil, false
	}
}
func GetUserId(tokenStr string) (string, bool) {
	claims, ok := DecodeTokenJwt(tokenStr)
	if ok {
		return claims["Id"].(string), true
	} else {
		return "", false
	}
}
func GetDataJwt(tokenStr string) (Id string, role string, ok bool) {
	claims, ok := DecodeTokenJwt(tokenStr)
	fmt.Println("claims: ", claims)
	if ok {
		return claims["Id"].(string), claims["Roles"].(string), true
	} else {
		return
	}
}
func GenerateToken(appId string, userId string, roles string, permission []string, overrideExpiration time.Time) (res interface{}, err error) {
	var tokenExpirationTime int
	_, _ = fmt.Sscanf("hours:"+os.Getenv("TOKEN_EXPIRATION_TIME"), "hours:%5d", &tokenExpirationTime)

	var expiredToken = time.Now().Local().Add(time.Hour * time.Duration(tokenExpirationTime))

	if !overrideExpiration.IsZero() {
		expiredToken = overrideExpiration
	}

	token := jwt.New(jwt.SigningMethodHS256)
	value := token.Claims.(jwt.MapClaims)

	value["Client"] = "Karya Group"
	value["Id"] = userId
	value["Roles"] = roles
	value["Permission"] = permission
	value["iat"] = time.Now().Local().Unix()
	value["Exp"] = expiredToken.Format("2006-01-02 15:04:05")

	var jwtKey string
	for i := 0; i < len(AppIDs); i++ {
		if AppIDs[i] == appId {
			jwtKey = JWTKeys[i].(string)
		}
	}
	if jwtKey == "" || jwtKey != os.Getenv("JWT_SECRET") {
		return nil, errors.New("AppId is Invalid")
	}

	tokenString, err := token.SignedString([]byte(jwtKey))
	dataToken := general.AccessToken{
		AccessToken: tokenString,
		CreatedAt:   time.Now().Local().Format("2006-01-02 15:04:05"),
		Exp:         value["Exp"].(string),
		TokenType:   "Bearer",
		Duration:    os.Getenv("TOKEN_EXPIRATION_TIME") + " Hour",
		Roles:       StringToSHA1(roles + "@2025"),
	}
	return dataToken, err
}
