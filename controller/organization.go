package controller

import (
	"context"
	"fmt"
	"net/http"

	Config "github.com/21RMT14Muthuram/my-new-app/database"
	models "github.com/21RMT14Muthuram/my-new-app/model"
	"github.com/gin-gonic/gin"
)


func CreateOrganization(c *gin.Context){

	adminID := c.Param("id") 

	var orgdetails  models.OrganizationType

	fmt.Println("new Organization", orgdetails)
	_, err := Config.DB.Exec(context.Background(),
		`INSERT INTO organizationmgmt.organization (admin_id, org_id, org_name)VALUES ($1, $2, $3) `,
		 adminID, orgdetails.OrgCode, orgdetails.OrgName)

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create organization"})
		return
	}

	c.JSON(200, gin.H{"message": "Organization Created Successfully!"})

	fmt.Print(err)
}