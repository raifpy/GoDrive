package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/raifpy/Go/errHandler"
	"github.com/tcnksm/go-input"
	drive "google.golang.org/api/drive/v3"
)

func driveInit() {
	if _, err := os.Stat("credentials.json"); os.IsNotExist(err) {
		log.Fatalln("Clound't find credentials.json. Get it on Google Drive Console.")
	}
	a, err := getDrive("credentials.json")
	errHandler.HandlerExit(err)
	driveService = a

	if _, err := os.Stat(".drive"); err != nil {
		//driveMap := map[int]string{0: "your special drive"}
		driveArray := []string{"special"}
		log.Println("Listing Drives.")
		for _, value := range getDriveTeams(a) {
			//driveMap[index+1] = value.Name
			driveArray = append(driveArray, value.Id+":"+value.Name)
		}

		sonuc, err := inputUI.Select("Yayınlanacak bir drive seç;", driveArray, &input.Options{Loop: true, Required: true})
		errHandler.HandlerExit(err)

		id := strings.Split(sonuc, ":")[0]

		ioutil.WriteFile(".drive", []byte(id), os.ModePerm)
		teamDriverID = id
	} else {
		veri, err := ioutil.ReadFile(".drive")
		errHandler.HandlerExit(err)

		teamDriverID = string(veri)

	}
	go deleteDailyIPLimitMapSync()
	/*nLog("Crawling Drive")
	err = crawlDrive()
	if err != nil {
		nLog("Error on crawl: " + err.Error())
	}
	nLog("Crawling finished")

	go crawlLoop()*/ // Codeksion'a eklenecek.
}

/*
func crawlLoop() {
	for range time.NewTicker(time.Minute * time.Duration(crawlIntEveryMin)).C {
		nLog("Crawling Drive")
		err := crawlDrive()
		if err != nil {

			nLog("Error on crawl: " + err.Error())
		}
		nLog("Crawling finished")
	}
}
*/ // Codeksion'a eklenecek.
/*
func crawlDrive() error {
	files, err := getDriveFilesModular()
	if err != nil {

		return err
	}
	addDirDriveStoregeValue("/", files)
	for _, value := range files {

		crawlDriveSelf(value)
	}
	return nil
}

func crawlDriveSelf(files *drive.File) {
	nLog(files.Id + " : " + files.Name)
	if files.MimeType == googleDriveFolderMime {
		go crawlDriveFolder(files)
	} else {
		addFileDriveStoregeValue(files.Id, files)
	}
}

func crawlDriveFolder(file *drive.File) {

	oye, err := getDriveFolderValues(file.Id)
	if err != nil {

		nLog(err.Error() + " waiting 10 sec.")
		time.Sleep(time.Second * 10)
		crawlDriveFolder(file)
	}
	addDirDriveStoregeValue(file.Id, oye)
	for _, value := range oye {
		if value.MimeType == googleDriveFolderMime {
			go crawlDriveFolder(value)
		} else {
			go crawlDriveSelf(value)
		}
	}
}*/ // Codeksion'a eklenecek.

func listDriveFiles(drive *drive.Service) {
	e, err := drive.Files.List().Do()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Kindle: ", e.Kind)
	for _, file := range e.Files {
		fmt.Println(toStringJSON(file), "\n ")
	}
}

func getDriveFiles(drive *drive.Service) []*drive.File {
	e, err := drive.Files.List().Do()
	if err != nil {
		log.Println(err)
		return nil
	}
	return e.Files

}

func listDriversTeam(drive *drive.Service) {
	e, err := drive.Teamdrives.List().Do()
	if err != nil {
		log.Println(err)
		return
	}
	for _, tdrive := range e.TeamDrives {
		fmt.Println(tdrive.Name, tdrive.Id, tdrive.Kind, tdrive.ForceSendFields)
	}
}

func getDriveTeams(drive *drive.Service) []*drive.TeamDrive {
	e, err := drive.Teamdrives.List().Do()
	if err != nil {
		log.Println(err)
		return nil
	}
	/*for _, tdrive := range e.TeamDrives {
		fmt.Println(tdrive.Name, tdrive.Id, tdrive.Kind, tdrive.ForceSendFields)
	}*/
	return e.TeamDrives
}

func listDriveIDFolders(drive *drive.Service) {
	liste := drive.Files.List()
	liste.IncludeItemsFromAllDrives(true)
	liste.SupportsAllDrives(true)
	liste.Corpora("drive")
	liste.TeamDriveId("0AIQyAzGGMUwYUk9PVA")
	//liste.DriveId("0AIQyAzGGMUwYUk9PVA")
	liste.Q("mimeType = 'application/vnd.google-apps.folder'")

	ulan, err := liste.Do()
	if err != nil {
		log.Println(err)
		return
	}
	for _, value := range ulan.Files {
		fmt.Println(value)
	}
}

func listDriveIDValues(drive *drive.Service) {
	liste := drive.Files.List()
	liste.IncludeItemsFromAllDrives(true)
	liste.SupportsAllDrives(true)
	liste.Corpora("drive")
	ulan, err := liste.DriveId("0AIQyAzGGMUwYUk9PVA").Do()
	if err != nil {
		log.Println(err)
		return
	}
	for index, value := range ulan.Files {
		fmt.Println(index, value.Name,
			"|", value.IconLink, "|",
			value.Id, value.Size,
			value.ExportLinks,
			value.CreatedTime,
			value.FileExtension,
			value.WebContentLink,
			value.WebViewLink,
			value.FullFileExtension, "ğ", value.Id,
			value.Spaces, value.MimeType)
	}
}

func getDriveIDValues(drive *drive.Service) []*drive.File {
	liste := drive.Files.List()
	liste.IncludeItemsFromAllDrives(true)
	liste.SupportsAllDrives(true)
	liste.Corpora("drive")
	ulan, err := liste.DriveId("0AIQyAzGGMUwYUk9PVA").Do()
	if err != nil {
		log.Println(err)
		return nil
	}
	return ulan.Files
}

func listFolderValues(driver drive.Service) {
	//dirid := "19_5hPymCTcLykpUToTtEG-lDgywzkLxW"
	dirid := "1C1PdWFo2SgPYN1VKxqsuWdnzUr5cRhNp"
	klasor := driver.Files.List()
	klasor.IncludeItemsFromAllDrives(true)
	klasor.IncludeTeamDriveItems(true)
	klasor.Corpora("drive")
	klasor.TeamDriveId("0AIQyAzGGMUwYUk9PVA")
	klasor.Q(fmt.Sprintf("'%s' in parents", dirid))
	/*klasor.Pages(context.Background(), func(d *drive.FileList) error {
		fmt.Println(d.Files)
		return nil
	})*/
	klasor.SupportsAllDrives(true)
	file, err := klasor.Do()
	if err != nil {
		log.Println(err)
		return
	}
	for index, value := range file.Files {
		fmt.Println(index, value)
	}

}

func getFolderValues(driver drive.Service, dirid string) []*drive.File {
	//dirid := "19_5hPymCTcLykpUToTtEG-lDgywzkLxW"
	//dirid := "1C1PdWFo2SgPYN1VKxqsuWdnzUr5cRhNp"
	klasor := driver.Files.List()
	klasor.IncludeItemsFromAllDrives(true)
	klasor.IncludeTeamDriveItems(true)
	klasor.Corpora("drive")
	klasor.TeamDriveId("0AIQyAzGGMUwYUk9PVA")
	klasor.Q(fmt.Sprintf("'%s' in parents", dirid))
	/*klasor.Pages(context.Background(), func(d *drive.FileList) error {
		fmt.Println(d.Files)
		return nil
	})*/
	klasor.SupportsAllDrives(true)
	file, err := klasor.Do()
	if err != nil {
		log.Println(err)
		return nil
	}
	return file.Files

}

func downloadDriveValues(drive drive.Service) {

	file := drive.Files.Get("1ZuLtjiONWEmp9VLjpgnbbX8P8ZUjqlDk")

	fmt.Println(file.Header())
	/*f, err := file.Do()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(f.Size)*/
	response, err := file.Download()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(response.Header)

	aboFile, _ := os.Create("abo.7z")
	io.Copy(aboFile, response.Body)
	aboFile.Close()
	response.Body.Close()
	fmt.Println("bitti ")

}
func getDownloadMediaResponse(fileid string) (*http.Response, *drive.File, error) {

	file := driveService.Files.Get(fileid)
	file.SupportsTeamDrives(true)

	dfile, err := file.Do()
	if err != nil {
		return nil, nil, err
	}
	resp, err := file.Download()
	if err != nil {
		return nil, nil, err
	}

	return resp, dfile, nil

}

func getDriveFilesModular() ([]*drive.File, error) {
	list := driveService.Files.List()
	ilgiliMakamaGöreListeNiceliklerEkleyenFonksiyon(list)
	files, err := list.Do()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return files.Files, nil
}

func getDriveFolderValues(dirid string) ([]*drive.File, error) {
	//dirid := "19_5hPymCTcLykpUToTtEG-lDgywzkLxW"
	//dirid := "1C1PdWFo2SgPYN1VKxqsuWdnzUr5cRhNp"
	list := driveService.Files.List()
	ilgiliMakamaGöreListeNiceliklerEkleyenFonksiyon(list)
	list.Q(fmt.Sprintf("'%s' in parents", dirid))

	file, err := list.Do()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return file.Files, nil

}

func ilgiliMakamaGöreListeNiceliklerEkleyenFonksiyon(list *drive.FilesListCall) {
	if teamDriverID != "special" {
		list.IncludeItemsFromAllDrives(true)
		list.SupportsAllDrives(true)
		list.Corpora("drive")
		list.DriveId(teamDriverID)
	}
}
