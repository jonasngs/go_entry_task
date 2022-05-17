package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"github.com/jonasngs/go_entry_task/httpserver/config"
	"github.com/gin-gonic/gin"
	pb "github.com/jonasngs/go_entry_task/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client pb.UserManagerClient

// Handler to load login page
func loginHandler(ctx *gin.Context) {
	session_token, err := ctx.Cookie("session-id")
	if err != nil {
		log.Printf("No user session found: %v\n", err)
		ctx.HTML(http.StatusOK, "login.html", gin.H{})
		return
	}
	_, err = client.LoadProfile(ctx, &pb.ProfileRequest{SessionToken: session_token})
	if err != nil {
		ctx.HTML(http.StatusOK, "login.html", gin.H{})
		return
	}
	ctx.Redirect(http.StatusFound, "/profile")
}

// Handler to load profile page
func profileHandler(ctx *gin.Context) {
	session_token, err := ctx.Cookie("session-id")
	if err != nil {
		log.Printf("No user session found: %v\n", err)
		ctx.Redirect(http.StatusFound, "/login")
		return
	}
	//fmt.Println(session_token)
	res, err := client.LoadProfile(ctx, &pb.ProfileRequest{SessionToken: session_token})
	if err != nil {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}
	ctx.HTML(http.StatusOK, "profile.html", gin.H{
		"Nickname":        res.User.Nickname,
		"Profile_Picture": res.User.ProfilePicture,
		"Username": res.User.Username,
	})
}

// Handler to load edit page
func updateHandler(ctx *gin.Context) {
	session_token, err := ctx.Cookie("session-id")
	if err != nil {
		log.Printf("No user session found: %v\n", err)
		ctx.Redirect(http.StatusFound, "/login")
		return
	}
	res, err := client.LoadProfile(ctx, &pb.ProfileRequest{SessionToken: session_token})
	if err != nil {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}
	ctx.HTML(http.StatusOK, "edit.html", gin.H{
		"Nickname":        res.User.Nickname,
		"Profile_Picture": res.User.ProfilePicture,
	})

}

// Handler for login post request
func loginRequestHandler(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	res, err := client.Login(ctx, &pb.LoginRequest{Username: username, Password: password})
	if err != nil {
		ctx.HTML(http.StatusForbidden, "login.html", gin.H{
			"res": "Incorrect Username/Password",
		})
		return
	}
	if res.SessionToken != "" {
		ctx.SetCookie("session-id", res.SessionToken, 360, "/", "localhost", false, true) //maxAge parameter in seconds
	}
	ctx.Redirect(http.StatusFound, "/profile")
	// if res.GetRc().String() == "Success" {
	// 	if res.SessionToken != "" {
	// 		ctx.SetCookie("session-id", res.SessionToken, 360, "/", "localhost", false, true) //maxAge parameter in seconds
	// 	}
	// 	ctx.Redirect(http.StatusFound, "/profile")
	// 	return
	// } else {
	// 	ctx.HTML(http.StatusForbidden, "login.html", gin.H{
	// 		"res": "Incorrect Username/Password",
	// 	})
	// 	return
	// }

}

// Handler for update post request
func updateRequestHandler(ctx *gin.Context) {
	session_token, err := ctx.Cookie("session-id")
	if err != nil {
		log.Printf("No user session found: %v\n", err)
		ctx.Redirect(http.StatusFound, "/login")
		return
	}
	image := &pb.Image{
		FileExtension: "",
		Data:          "",
	}
	nickname := ctx.PostForm("nickname")
	image_file, err := ctx.FormFile("image_profile")
	if err == nil {
		file_ext := filepath.Ext(image_file.Filename) //file extension
		file, err := image_file.Open()
		if err != nil {
			log.Printf("Could not open image file: %v\n", err)
			ctx.Redirect(http.StatusFound, "/edit")
			return
		}
		image_data, _ := ioutil.ReadAll(file)
		defer file.Close()
		encoded_image_bytes := base64.StdEncoding.EncodeToString(image_data) // this is without base64 encoding header
		image.FileExtension = file_ext
		image.Data = encoded_image_bytes
	} else {
		log.Printf("No image file found: %v\n", err)	
	}

	_, err = client.Update(ctx, &pb.UpdateRequest{Nickname: nickname, ProfilePicture: image, SessionToken: session_token})
	if err != nil {
		// Improve error handling
		ctx.HTML(http.StatusFound, "login.html", gin.H{})
		return
	}
	ctx.Redirect(http.StatusFound, "/profile")
}

func main() {

	config.OpenConfigFile()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//conn, err := grpc.Dial("localhost:4040", opts...)
	conn, err := grpc.Dial(config.GetTCPServer(), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client = pb.NewUserManagerClient(conn)

	router := gin.Default()
	router.LoadHTMLGlob("view/*.html")
	router.GET("/login", loginHandler)
	router.GET("/profile", profileHandler)
	router.POST("/login", loginRequestHandler)
	router.GET("/edit", updateHandler)
	router.POST("/edit", updateRequestHandler)
	router.Static("/resources", "../resources") // Static route for profile picture
	router.Run() // Using Default port 8080
}
