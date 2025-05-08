package middlewares

import (
	"fmt"

	"github.com/antoniomarfa/hexatools/api/utils"
	"github.com/gin-gonic/gin"
)

// Recover is a middleware function to defer and return an error response in case of panic during the handler execution
/*
func Recover(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			var err error
			if rec := recover(); rec != nil {
				switch t := rec.(type) {
				case error:
					err = t
				case string:
					err = fmt.Errorf(t)
				default:
					err = fmt.Errorf("unknown error ocurred")
				}
				utils.ResponseError(w, r, nil, err)
			}
		}()
		next(w, r)
	})
}
*/

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			var err error
			if rec := recover(); rec != nil {
				switch t := rec.(type) {
				case error:
					err = t
				case string:
					err = fmt.Errorf(t)
				default:
					err = fmt.Errorf("unknown error occurred")
				}
				utils.ResponseError(c.Writer, c.Request, nil, err)
				// En caso de error, abortamos la ejecución del handler.
				c.Abort()
			}
		}()
		// Continuamos con el siguiente middleware o handler.
		c.Next()
	}
}
