package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raifpy/Go/raiFile"
	input "github.com/tcnksm/go-input"
	"google.golang.org/api/drive/v3"
)

func ginLog(g *gin.Context) {

	if varJSON.ShowLogs {
		log.Println(g.Request.URL.Path, g.ClientIP())
	}
	if varJSON.SaveLogs {
		raiFile.WriteFileLines(saveLogsPath, fmt.Sprintf("Time: %s\nIp: %s\nPath: %s\n", time.Now().String(), g.ClientIP(), g.Request.URL.String()))
	}
}

func nLog(veri string) {
	if varJSON.ShowLogs {
		log.Println(veri)
	}
	if varJSON.SaveLogs {
		raiFile.WriteFileLines(saveLogsPath, fmt.Sprintf("Time: %s\nLog: %s\n\n", time.Now().String(), veri))
	}
}

func to32digit(passStr string) []byte {
	pass := []byte(passStr)
	if len(pass) == 32 {
		return pass
	}
	if len(pass) < 32 {
		var newPass = make([]byte, 32)
		copy(newPass, pass)
		for i := len(pass); i < 32; i++ {

			newPass[i] = 5
		}
		return newPass
	}
	if len(pass) > 32 {
		return pass[:32]
	}
	return pass

}

func getHashInput() []byte {

	if _, err := os.Stat("token.json"); os.IsNotExist(err) {
		hashTmp, err := inputUI.Ask("Create a password [Enter-use default]", &input.Options{HideDefault: true, Default: defaultHash})

		if err != nil {
			log.Fatalln(err)
		}

		if hashTmp != defaultHash {
			return to32digit(hashTmp)
		}

		f, err := os.Create(".default")
		if err != nil {
			nLog(err.Error() + "\nCloudn't create .default\nYOU NEED CREATE MANUEL\n\tType \"touch .default\"")
		} else {
			f.Close()
		}
		//firstOpen = true
		return to32digit(defaultHash)

	}

	// token.json mevcut ise;
	if _, err := os.Stat(".default"); err == nil {
		return to32digit(defaultHash)
	}
	hashTmp, err := inputUI.Ask("Give me pass?", &input.Options{HideDefault: true, Loop: true, Required: true, Hide: true})

	if err != nil {
		log.Fatalln(err)
	}

	return to32digit(hashTmp)

}

func toStringJSON(veri interface{}) string {
	ğ, err := json.Marshal(veri)
	if err != nil {
		log.Println(err)
		return ""
	}
	ğğ := string(ğ)
	ğ = nil
	return ğğ

}

func isInList(key string, list []string) bool {
	for _, value := range list {
		if value == key {
			return true
		}
	}
	return false
}

func filesToTemplateHTMLTdStructArray(files []*drive.File) []templateHTMLTdStruct {
	var veri = make([]templateHTMLTdStruct, len(files))
	for index, value := range files {
		mime := formatMimeType(value.MimeType)
		var path string
		if mime == "folder" {
			path = "/d"
		} else {
			path = "/f"
		}
		veri[index] = templateHTMLTdStruct{
			Href: path + "?id=" + value.Id,
			Name: value.Name,
			Mime: mime,
			Size: fmt.Sprint(value.Size),
		}
	}
	return veri
}

func formatMimeType(mimetype string) string {
	if mimetype == googleDriveFolderMime {
		return "folder"
	}

	spi := strings.Split(mimetype, "/")
	if len(spi) == 1 {
		return spi[0]
	}
	return spi[1]
}

func getDownloadIPMapItem(key string) int {
	value, ok := downloadIPMap.Load(key)
	if !ok {
		return -1
	}
	return value.(int)
}

func addDownloadIPMapItem(key string, digit int) {
	downloadIPMap.Store(key, digit)
}

func deleteDownloadIPMapItem(key string) {
	downloadIPMap.Delete(key)
}

func getFileDriveStorageValue(id string) (*drive.File, bool) {
	fileArrayInterface, ok := fileDriveStorageMap.Load(id)
	if !ok {
		return nil, false
	}
	return fileArrayInterface.(*drive.File), true
}

func getDirDriveStorageValue(id string) ([]*drive.File, bool) {
	fileArrayInterface, ok := dirDriveStorageMap.Load(id)
	if !ok {
		return nil, false
	}
	return fileArrayInterface.([]*drive.File), true
}

func addFileDriveStoregeValue(id string, files *drive.File) {
	fileDriveStorageMap.Store(id, files)
}

func addDirDriveStoregeValue(id string, files []*drive.File) {
	dirDriveStorageMap.Store(id, files)
}

func getDailyIPLimitMapValue(ip string) int {
	value, ok := dailyIPLimitMap.Load(ip)
	if !ok {
		return 0
	}
	return value.(int)
}
func addDailyIPLimitMapValue(ip string) {
	dailyIPLimitMap.Store(ip, getDailyIPLimitMapValue(ip)+1)
}
func deleteDailyIPLimitMapSync() {
	for range time.NewTicker(time.Hour * 24).C {
		dailyIPLimitMap = sync.Map{}
	}
}
func checkIPAddrOverDailyRequest(ip string) bool {
	if getDailyIPLimitMapValue(ip) >= varJSON.DailyIPLimit {
		nLog(ip + " blocked: daily limit")
		return true
	}
	addDailyIPLimitMapValue(ip)
	return false
}

func blockDownloadAfterDownload(ip string) bool {
	_, ok := blockDownloadAfterDownloadMap.Load(ip)
	if !ok {
		//nLog(ip + " download block: after download block config")
		blockDownloadAfterDownloadMap.Store(ip, nil)
		go deleteBlockDownloadAfterDownload(ip)
		//nLog(ip + " can download files now: unblocked")

	}
	return ok
}

func deleteBlockDownloadAfterDownload(ip string) {
	time.Sleep(time.Minute * time.Duration(varJSON.BlockDownloadAfterDownloadMin))
	blockDownloadAfterDownloadMap.Delete(ip)

}
