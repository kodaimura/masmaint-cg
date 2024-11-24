package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"masmaint/config"
	"masmaint/internal/core/jwt"
	"masmaint/internal/core/errs"
)


func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cf := config.GetConfig()

		user, pass, ok := c.Request.BasicAuth()
		if !ok || user != cf.BasicAuthUser || pass != cf.BasicAuthPass {
			c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}


func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := jwt.Auth(c); err != nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}


func ApiJwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := jwt.Auth(c); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}


func ApiResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			switch e := err.(type) {
			case errs.BadRequestError:
				c.JSON(http.StatusBadRequest, gin.H{
					"error": e.Error(), 
					"details": gin.H{ "field": e.Field },
				})
			case errs.UnauthorizedError:
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": e.Error(),
					"details": gin.H{},
				})
			case errs.ForbiddenError:
				c.JSON(http.StatusForbidden, gin.H{
					"error": e.Error(),
					"details": gin.H{},
				})
			case errs.NotFoundError:
				c.JSON(http.StatusNotFound, gin.H{
					"error": e.Error(),
					"details": gin.H{},
				})
			case errs.ConflictError:
				c.JSON(http.StatusConflict, gin.H{
					"error": e.Error(),
					"details": gin.H{ "column": e.Column },
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": e.Error(),
					"details": gin.H{},
				})
			}
		}
	}
}