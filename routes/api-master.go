package routes

import (
	"be-golang/controllers"
	"be-golang/routes/middleware"
	"github.com/labstack/echo/v4"
)

func BuildMaster(e *echo.Echo) {
	router := e.Group(prefix)
	auth := middleware.Auth
	RoleAdmin := ManageAccess("admin").RoleAccess
	RoleAdminAllIndoor := ManageAccess("admin, sales-indoor, accounting ,finance").RoleAccess

	//Master Customer
	router.GET("/master/customers", controllers.GetMasterCustomer, auth)
	router.GET("/master/customers/:id", controllers.GetMasterCustomerById, auth)
	router.PUT("/master/customers/:id", controllers.UpdateMasterCustomer, auth, RoleAdminAllIndoor)
	router.POST("/master/customers", controllers.CreateMasterCustomer, auth)
	router.POST("/master/customers-multi", controllers.CreateMasterCustomerMulti, auth)
	router.DELETE("/master/customers/:id", controllers.DeleteMasterCustomer, auth, RoleAdminAllIndoor)
	router.PUT("/master/customers/approval/:id", controllers.ApprovalMasterCustomer, auth, RoleAdmin)
}
