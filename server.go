package main

import (
	"be-golang/helpers"
	"be-golang/models/conn"
	"be-golang/structs/general"
	"fmt"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"time"

	//redisgo "be-golang/redis"
	"be-golang/routes"
	"os"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	errLoadingEnvFile := godotenv.Load()
	if errLoadingEnvFile != nil {
		helpers.HandleError("error loading the .env file", errLoadingEnvFile)
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		helpers.HandleError("error set timezone", err)
	}
	time.Local = loc

	conn.PostgresDB = conn.ConnectionPostgresDb(1)

	e := echo.New()

	e.Validator = &general.CustomValidator{Validator: validator.New()}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:4200", "https://sales-app-apergu.vercel.app", "https://salesapp.cvsm.co.id"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	e.Use(middleware.BodyLimit("100M"))
	e.Debug = true
	// Custom logger middleware for request tracing (logs method, path, status)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339}","method":"${method}","uri":"${uri}","status":${status},"error":"${error}}` + "\n",
		Output: e.Logger.Output(),
	}))

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok {
			if he.Code == http.StatusMethodNotAllowed {
				// Trace/log the 405 details
				log.Printf("405 Method Not Allowed: Method=%s, Path=%s, Allowed Methods=%s, User-Agent=%s",
					c.Request().Method,
					c.Request().URL.Path,
					c.Response().Header().Get("Allow"),
					c.Request().UserAgent())

				// Optional: Add custom response body for better client feedback
				c.JSON(he.Code, echo.Map{
					"error": "Method Not Allowed",
					"message": fmt.Sprintf("Method %s not allowed for %s. Allowed: %s",
						c.Request().Method, c.Request().URL.Path, c.Response().Header().Get("Allow")),
					"timestamp": time.Now().Format(time.RFC3339),
				})
				return
			}
		}
		// Fallback to default handler for other errors
		e.DefaultHTTPErrorHandler(err, c)
	}
	routes.BuildAccess(e)
	routes.BuildMaster(e)
	e.Logger.Fatal(e.Start(":" + os.Getenv("APP_PORT")))
}
