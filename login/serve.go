package login

import "github.com/gin-gonic/gin"

func Serve() {
	router := gin.Default()
	
	router.GET("/embed", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/auth", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/login", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/www", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/static", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

/*
		cs := []*http.Cookie{
			&http.Cookie{Name: "sessions_gog_com", Value: (*s).session},
			&http.Cookie{Name: "gog-al", Value: (*s).al},
		}
*/
	router.GET("/", func(c *gin.Context) {
		cli := getClient([]*http.Cookie{})
		res, err := cli.Get("https://auth.gog.com/auth?client_id=46899977096215655&amp;redirect_uri=https://embed.gog.com/on_login_success?origin=client&amp;response_type=code&amp;layout=default&amp;brand=gog")
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		} else {
			processGogResponse(res, c)
		}
	})

	router.Run("localhost:8080")
}
