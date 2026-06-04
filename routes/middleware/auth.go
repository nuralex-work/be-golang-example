// create jwt authentification
// check valid, check expired, and check signature

package middleware

import (
	"be-golang/helpers"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type Roles struct {
	Name string `json:"name" form:"name"`
}

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		if authorization := e.Request().Header.Get("Authorization"); authorization != "" {
			log.Println("========== Request log ==========")
			log.Println()
			log.Printf("URL: %s%s", e.Request().Host, e.Request().URL.Path)
			log.Println("=========== Request Log =============")
			log.Println()

			authorizationToken := strings.Split(authorization, " ")

			if len(authorizationToken) != 2 {
				return helpers.HandleResponse(e, http.StatusUnauthorized, "Unauthorized", nil)
			}

			if authorizationToken[0] != "Bearer" {
				return helpers.HandleResponse(e, http.StatusUnauthorized, "Unauthorized", nil)
			}

			token, err := jwt.Parse(authorizationToken[1], func(token *jwt.Token) (interface{}, error) {
				_, _ = token.Method.(*jwt.SigningMethodHMAC)

				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if err != nil {
				return helpers.HandleResponse(e, http.StatusUnauthorized, "authorization token credentials do not match", nil)
			}

			claims, ok := token.Claims.(jwt.MapClaims)

			if !ok || !token.Valid {
				return helpers.HandleResponse(e, http.StatusUnauthorized, "invalid authorization token credentials", nil)
			}

			tokenExpired, _ := time.Parse("2006-01-02 15:04:05", claims["Exp"].(string))
			currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))

			if currentTime.After(tokenExpired) {
				return helpers.HandleResponse(e, http.StatusUnauthorized, "Token Expired", nil)
			}

			e.Set("Authorization", claims)
			return next(e)
		}
		return helpers.HandleResponse(e, http.StatusUnauthorized, "Unauthorized", nil)
	}
}
func BasicAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		if authorization := e.Request().Header.Get("Authorization"); authorization != "" {

			uname, pass, ok := e.Request().BasicAuth()
			if ok && uname == os.Getenv("BASIC_USERNAME") && pass == os.Getenv("BASIC_PASS") {
				return next(e)
			}
			return helpers.HandleResponse(e, http.StatusUnauthorized, "unauthorized", nil)

		}
		return helpers.HandleResponse(e, http.StatusUnauthorized, "Unauthorized", nil)
	}
}
func (r Roles) RoleAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		if authorization := e.Request().Header.Get("Authorization"); authorization != "" {
			_, roles, ok := helpers.GetDataJwt(authorization)
			if ok && strings.Contains(r.Name, roles) {
				return next(e)
			}
			return helpers.HandleResponse(e, http.StatusForbidden, "Forbidden Access For This Role", nil)
		}
		return helpers.HandleResponse(e, http.StatusUnauthorized, "Unauthorized", nil)
	}
}
func LocalAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		if os.Getenv("ENV") != "local" {
			return helpers.HandleResponse(e, http.StatusForbidden, "Forbidden Access", nil)
		}
		return next(e)
	}
}
