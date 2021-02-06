package main

import (
	"fmt"
	"io"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"github.com/raifpy/Go/errHandler"
)

func route() error {
	if !varJSON.ShowLogs {
		gin.SetMode(gin.ReleaseMode)
	}
	root := gin.Default()
	root.LoadHTMLGlob(filepath.Join("template", "*"))
	root.GET("/", routeHome)
	root.GET("/style.css", routeCSS)

	root.GET("/d", routeDir)
	root.GET("/f", routeFile)

	root.NoRoute(route404)
	return root.Run(varJSON.RoutePort)

}
func route404(g *gin.Context) {
	ginLog(g)
	g.HTML(404, "error.html", errorHTMLStruct{
		ErrorCode: "404",
		ErrorDesc: "Not found",
		ErrorCaus: g.Request.URL.Path + " not found or moved",
	})
	return

}

func routeHome(g *gin.Context) {
	ginLog(g)
	if checkIPAddrOverDailyRequest(g.ClientIP()) {

		g.HTML(403, "error.html", errorHTMLStruct{
			"403", "Access denied", "You banned from server. Reason: To much requests",
		})
		return
	}
	/*files, ok := getDirDriveStorageValue("/")
	if !ok {
		g.HTML(404, "error.html", errorHTMLStruct{
			"404", "404 Not Found", "the file (you are looking) for may have been moved or removed",
		})
		return
	}*/ // Codeksion'a eklenecek
	files, err := getDriveFilesModular()
	if errHandler.HandlerBool(err) {
		g.HTML(500, "error.html", errorHTMLStruct{
			"500", err.Error(), "Server-in error. You can report this.",
		})
		return
	}

	g.HTML(200, "template.html", templateHTMLStruct{
		varJSON.PublicName,
		filesToTemplateHTMLTdStructArray(files),
	})
}

func routeError(g *gin.Context) {
	ginLog(g)
	g.HTML(404, "error.html", errorHTMLStruct{
		ErrorCode: "404",
		ErrorDesc: "404 Not Found",
		ErrorCaus: "the file (you are looking) for may have been moved or removed",
	})
}

func routeCSS(g *gin.Context) {
	g.HTML(200, "style.css", nil)
}

func routeDir(g *gin.Context) {
	ginLog(g)
	if checkIPAddrOverDailyRequest(g.ClientIP()) {

		g.HTML(403, "error.html", errorHTMLStruct{
			"403", "Access denied", "You banned from server. Reason: To much requests",
		})
		return
	}
	id, ok := g.GetQuery("id")
	if !ok || id == "" {
		nLog("Wrong dir-id ;" + g.Request.URL.Path + " => " + g.ClientIP())
		g.HTML(404, "error.html", errorHTMLStruct{
			ErrorCode: "404",
			ErrorDesc: "404 Not Found",
			ErrorCaus: "the file (you are looking) for may have been moved or removed",
		})
		return
	}
	/*files, ok := getDirDriveStorageValue(id)
	if !ok {
		g.HTML(404, "error.html", errorHTMLStruct{
			"404", "404 Not Found", "the file (you are looking) for may have been moved or removed",
		})
		return
	}*/ // Codeksion'a ekleyeceğim bunları.

	files, err := getDriveFolderValues(id)
	if errHandler.HandlerBool(err) {
		g.HTML(500, "error.html", errorHTMLStruct{
			"500", err.Error(), "Server-in error. You can report this.",
		})
		return
	}

	g.HTML(200, "template.html", templateHTMLStruct{
		PublicName: varJSON.PublicName, Tbody: filesToTemplateHTMLTdStructArray(files),
	})
}

func routeFile(g *gin.Context) {
	ginLog(g)
	ip := g.ClientIP()
	if checkIPAddrOverDailyRequest(ip) {

		g.HTML(403, "error.html", errorHTMLStruct{
			"403", "Access denied", "You banned from server. Reason: To much requests",
		})
		return
	}
	block := blockDownloadAfterDownload(ip)
	if block {
		g.HTML(403, "error.html", errorHTMLStruct{
			"403",
			"Access denied",
			"Your download request blocked. You have to wait " + strconv.Itoa(varJSON.BlockDownloadAfterDownloadMin) + "min.",
		})
		return
	}

	id, ok := g.GetQuery("id")
	if !ok || id == "" {
		nLog("Wrong file-id ;" + g.Request.URL.Path + " => " + g.ClientIP())
		g.HTML(404, "error.html", errorHTMLStruct{
			ErrorCode: "404",
			ErrorDesc: "404 Not Found",
			ErrorCaus: "the file (you are looking) for may have been moved or removed",
		})
		return
	}
	/*if _, ok := getFileDriveStorageValue(id); !ok {
		g.HTML(404, "error.html", errorHTMLStruct{
			ErrorCode: "404",
			ErrorDesc: "404 Not Found",
			ErrorCaus: "the file (you are looking) for may have been moved or removed",
		})
		return
	}*/ // Codeksion'a eklenecek.

	if varJSON.SameIPCantInstallMultiFiles {
		ipInt := getDownloadIPMapItem(ip)
		if ipInt != -1 {
			g.HTML(403, "error.html", errorHTMLStruct{
				"403", "Access denied for multiple file downloads", "Looks like you have an ongoing download from this server. This server allows each user to download to 1 file at a time.",
			})

			return
		}
		addDownloadIPMapItem(ip, 1)

	}

	response, file, err := getDownloadMediaResponse(id)
	if errHandler.HandlerBool(err) {
		g.HTML(500, "error.html", errorHTMLStruct{
			"500", err.Error(), "Server-in error. You can report this.",
		})
		return
	}

	bucket := ratelimit.NewBucketWithRate(float64(varJSON.UploadDownoadLimitByte), int64(varJSON.UploadDownoadLimitByte))
	//log.Println("Dosya content-type: ", response.Header.Get("Content-Type"))
	nLog(fmt.Sprintf("%s downloading %s", ip, file.Name))
	g.Header("Content-Type", response.Header.Get("Content-Type"))
	g.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s", file.Name, file.Name))
	g.Header("Content-Length", fmt.Sprint(response.ContentLength))

	g.Stream(func(w io.Writer) bool {
		_, err := io.Copy(w, ratelimit.Reader(response.Body, bucket))
		if errHandler.HandlerBool(err) {
			return false
		}
		return true
	})
	if varJSON.SameIPCantInstallMultiFiles {
		deleteDownloadIPMapItem(ip)
	}
}
