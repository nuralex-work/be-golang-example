package routes

import (
	"be-golang/controllers/ctgeneral"
	"be-golang/routes/middleware"
	"github.com/labstack/echo/v4"
)

func BuildAccess(e *echo.Echo) {
	router := e.Group(prefix)
	auth := middleware.Auth
	RoleAdmin := ManageAccess("admin").RoleAccess

	//Master User
	router.GET("/m/users", ctgeneral.GetUser, auth, RoleAdmin)
	router.GET("/m/users/:id", ctgeneral.GetUserById, RoleAdmin)
	router.PUT("/m/users/:id", ctgeneral.UpdateUser, auth, RoleAdmin)
	router.POST("/m/users", ctgeneral.CreateUser, auth, RoleAdmin)
	//router.POST("/m/users-multi", ctgeneral.CreateUserMulti, auth, RoleAdmin)
	router.DELETE("/m/users/:id", ctgeneral.DeleteUser, auth, RoleAdmin)
	router.GET("/m/sales", ctgeneral.GetUserSales, auth)

	router.GET("/m/users/profile", ctgeneral.GetUserProfile, auth)

	//Master Roles
	router.GET("/m/roles", ctgeneral.GetRoles, auth, RoleAdmin)
}
