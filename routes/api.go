package routes

import (
	"be-golang/routes/middleware"
)

var prefix = "/api/v1"

func ManageAccess(param string) *middleware.Roles {
	Roles := new(middleware.Roles)
	Roles.Name = param
	return Roles
}
