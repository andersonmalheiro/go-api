package user

import (
	app "go-api/application/entities/user"
	"go-api/oops"

	"github.com/gin-gonic/gin"
)

// add is the handler function to POST requests on /users endpoint
func add(c *gin.Context) {
	var in app.INUser

	if err := c.ShouldBindJSON(&in); err != nil {
		oops.Handling(err, c)
		return
	}

	id, err := app.Add(&in)
	if err != nil {
		oops.Handling(err, c)
		return
	}

	c.JSON(201, id)
}
