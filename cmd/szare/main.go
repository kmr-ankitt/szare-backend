package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mdp/qrterminal/v3"
)

func main() {
	router := gin.Default()
	port := "8080"
	// step 1: create a server
	// step 2: show link and qr code
	// step 3: get("/") here we will fetch all the files and dir of server computer
	// step 4: if clicked on file then it will download the file

	ShowQRCode(port)

	router.Use(cors.Default())
	router.GET("/", getHomepage)
	router.POST("/api/download", downloadFile)
	router.Run("localhost:" + port)
}

func downloadFile(ctx *gin.Context) {
	currWorkingDir, err := os.Getwd()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get current working directory"})
		return
	}

	fileName := ctx.Request.URL.Query().Get("filename")
	if fileName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Filename not provided"})
		return
	}

	ctx.FileAttachment(currWorkingDir+"/"+fileName, fileName)
}

func getHomepage(ctx *gin.Context) {
	sysFiles := getFiles()
	var fileNames []string
	for _, file := range sysFiles {
		fileNames = append(fileNames, file.Name())
	}
	ctx.JSON(http.StatusOK, fileNames)
}

func getFiles() []os.DirEntry {
	items, err := os.ReadDir("./")
	if err != nil {
		fmt.Println(err)
	}
	return items
}

/*
* ShowQRCode shows the server link and it's QR code
 */
func ShowQRCode(port string) {
	ip := GetLocalIP()
	hostIp := "http://" + ip + ":" + port
	fmt.Println("Enter " + hostIp + " in your browser" + "\nor\n" + "Scan this QR code on your phone")

	config := qrterminal.Config{
		Level:     qrterminal.M,
		Writer:    os.Stdout,
		BlackChar: qrterminal.WHITE,
		WhiteChar: qrterminal.BLACK,
		QuietZone: 1,
	}
	qrterminal.GenerateWithConfig(hostIp, config)
}

/*
* GetLocalIP returns the local IP address of the computer
 */
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}