package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"fmt"
)

func loginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		
	})
}

func profileHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "profile.html", gin.H{
		
	})
}

func loginRequestHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	fmt.Printf("Username: %s LLLLL Password: %s", username, password)
	c.Redirect(http.StatusFound, "/profile")

}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("view/*.html")
	router.GET("/login", loginHandler)
	router.GET("/profile", profileHandler)
	router.POST("/login", loginRequestHandler)
	router.Run()
}