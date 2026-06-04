package helpers

import (
	"be-golang/structs/general"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateTokenJwt() (string, error) {
	// Create JWT claims
	claims := jwt.MapClaims{
		"iss":   os.Getenv("CLIENT_ID"),
		"scope": []string{"restlets", "rest_webservices"},
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour).Unix(),
		"aud":   "https://" + os.Getenv("ACCOUNT_ID") + ".suitetalk.api.netsuite.com/services/rest/auth/oauth2/v1/token",
	}

	// Create a new JWT token with RS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	token.Header["alg"] = "PS256"
	token.Header["kid"] = os.Getenv("CERT_ID")

	// Load RSA private key
	privateKeyBytes := []byte(Base64Decode(os.Getenv("PRIVATE_KEY")))
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		fmt.Println("Error parsing private key:", err)
		return "", nil
	}

	// Sign the token with the private key
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return "", nil
	}
	//
	// Parse and validate the token (optional)
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is RS256
		//if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		//	return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		//}
		return privateKey.Public(), nil
	})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		return "", nil
	}

	// Validate the token
	if _, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		//fmt.Println("Claims:")
		//for key, value := range claims {
		//	fmt.Printf("%s: %v\n", key, value)
		//}
	} else {
		fmt.Println("Invalid token")
	}

	return parsedToken.Raw, nil
}

func GetOauth2AccessToken() (res general.OauthResponse, err error) {
	token, _ := GenerateTokenJwt()
	fmt.Println("===get token==", token)
	httpClient := NewRestyClient("https://" + os.Getenv("ACCOUNT_ID") + ".suitetalk.api.netsuite.com")
	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Accept":       "application/json",
	}

	payload := map[string]string{
		"grant_type":            "client_credentials",
		"client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
		"client_assertion":      token,
	}

	if resp, err := httpClient.PostOauth2("/services/rest/auth/oauth2/v1/token", header, nil, payload); err != nil {
		return res, err
	} else {
		return resp, nil
	}
}

func GetAccessToken() (res string, error error) {
	//token := redisConn.Get("access_token")
	//fmt.Println(token, " ini token redis")
	//if token != nil {
	//	fmt.Println(token, " ini token redis nil")
	//	return token.(string), nil
	//} else {
	tokens, er := GetOauth2AccessToken()
	if er != nil {
		return "", error
	}
	//redisConn.Set("access_token", tokens.AccessToken, 3500)
	return tokens.AccessToken, er
	//}
}
