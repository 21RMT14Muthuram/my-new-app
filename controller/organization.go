package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	Config "github.com/21RMT14Muthuram/my-new-app/database"
	models "github.com/21RMT14Muthuram/my-new-app/model"
	"github.com/gin-gonic/gin"
)


func CreateOrganization(c *gin.Context){

	adminID := c.Param("id") 

	var orgdetails  models.OrganizationType

	// Bind request body to struct
    if err := c.ShouldBindJSON(&orgdetails); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Invalid JSON",
            "error":   err.Error(),
        })
        return
    }

	fmt.Println("new Organization", orgdetails)
	_, err := Config.DB.Exec(context.Background(),
		`INSERT INTO organizationmgmt.organization (admin_id, org_id, org_name)VALUES ($1, $2, $3) `,
		 adminID, orgdetails.OrgCode, orgdetails.OrgName)

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create organization"})
		return
	}

	_, err = Config.DB.Exec(
        context.Background(),
        `UPDATE usermgmt.users 
         SET is_admin = $1 
         WHERE id = $2`,
        true, adminID,
    )

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Admin not found!"})
		return
	}

	c.JSON(200, gin.H{"data":orgdetails})

	fmt.Print(err)
}

func UserJoin(c *gin.Context){

	userIDStr := c.Param("id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Invalid user ID",
            "error":   err.Error(),
        })
        return
    }

	var org models.OrganizationType
	if err := c.ShouldBindJSON(&org); err != nil{
		 c.JSON(http.StatusBadRequest, gin.H{
            "message": "Invalid JSON",
            "error":   err.Error(),
        })

		return
	}

	err = Config.DB.QueryRow(
		context.Background(),
		`SELECT id FROM organizationmgmt.organization 
		WHERE org_name = $1 and org_id = $2`,
		org.OrgName, org.OrgCode).Scan(&org.ID)
	
	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{"message": "Organization not Found!"})
		return
	}

	_, err = Config.DB.Exec(
        context.Background(),
        `INSERT INTO usermgmt.roles (org_id, user_id) 
         VALUES ($1, $2)`,
        org.ID, userID,
    )

	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "message": "Failed to join user",
            "error":   err.Error(),
        })
        return
	}

	c.JSON(http.StatusOK, gin.H{
        "message": "User joined successfully",
        "org_id":  org.ID,
        "user_id": userID,
    })

}