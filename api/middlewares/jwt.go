package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWT is a middleware function to check the authorization JWT Bearer token header of the request
func JWT(secret string, claims jwt.MapClaims) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("el método de inicio de sesión no es válido")
					}
					return []byte(secret), nil
				})
				if _, ok := token.Claims.(jwt.MapClaims); err != nil || !ok || !token.Valid {
					c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("token invalido : %s", err.Error())})
					c.Abort()
					return
				}
				for name, value := range claims {
					if claim, ok := token.Claims.(jwt.MapClaims)[name]; !(ok && claim == value) {
						c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Reclamo requerido% s no encontrado o incorrecto", name)})
						c.Abort()
						return
					}
				}
				// Pasar las claims decodificadas al contexto de Gin
				c.Set("decoded", token.Claims)
				c.Next()
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "El encabezado de autorización no tiene el formato correcto, debería ser Bearer + {token}"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Se requiere un encabezado de autorización"})
			c.Abort()
			return
		}
	}
}
