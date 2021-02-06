package main

//Last edited Sun Feb  7 01:38:33 +03 2021

// Author @raifpy
// Author News t.me/raifBlog

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/raifpy/Go/errHandler"
	input "github.com/tcnksm/go-input"
	"google.golang.org/api/drive/v3"
)

// Know Bugs; Size always 0byte
const defaultHash = "That's GoDrive's default hash."
const googleDriveFolderMime = "application/vnd.google-apps.folder"

var (
	//publicName      string `json: "PublicName"`
	driveService *drive.Service
	//credentialsPath string = "credentials.json"
	inputUI = input.UI{Writer: os.Stdin, Reader: os.Stdout}
	//
	//dailyIPLimit                     int  = 1000
	//sameIPCantInstallMultiFiles      bool = true // if true IDM can't work but this blocks google drive api limit.
	//enableBlockDownloadAfterDownload bool = true
	//blockDownloadAfterDownloadMin    int  = 10

	// crawlIntEveryMin int = 60 // Codeksion ile gelecek
	//
	//upDownLimit  int64 = 5000 * 1024
	teamDriverID string
	//
	//showLogs     bool   = true
	//saveLogs     bool   = true
	saveLogsPath string = time.Now().Format("Mon Jan 2 15:04:05 2006")
	//routePort    string = ":6070"
	//
	runOnBackground     bool
	runOnBackgroundHash []byte = nil

	runningOnBackground bool

	varJSON varStruct
)

func init() {
	file, err := os.Open("config.json")
	errHandler.HandlerExit(err)

	var tmpVarStruct varStruct

	err = json.NewDecoder(file).Decode(&tmpVarStruct)
	errHandler.HandlerExit(err)
	file.Close()
	varJSON = tmpVarStruct

	/*//hash := flag.String("hash", "This is goDrive's default hash!", "Sadece bir hash belirleyin!") // unsafe

	upLimit := flag.Int64("limit", 5_000, "1000 ve katları olması önerilir. İndirme ve yükleme hızı buna göre limitlenecektir.")
	port := flag.String("port", ":6070", "Yayın yapılacak port")
	logs := flag.Bool("unlog", false, "logları yazmaz ehe")
	logPath := flag.String("logPath", time.Now().String(), "Log adı dostum")
	iplimit := flag.Int("iplimit", 500, "X ip'sinin bir günde ziyaret edebileceği miktar")
	umulti := flag.Bool("unmulti", false, "Aynı IP adresinin (aynı anda) birden fazla dosya indirmesine engel olur. Bir dosya tamamen inmeden ya da iptal etmeden öbürü indirilemez. idm'in çalışmasını engelleyecektir!")
	oback := flag.Bool("background", false, "Programı arka planda başlatmak için cart curt")
	help := flag.Bool("help", false, "")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	upDownLimit = *upLimit * 1024
	if !strings.Contains(*port, ":") {
		*port = ":" + *port
	}
	routePort = *port
	saveLogs = !*logs
	saveLogsPath = *logPath
	dailyIPLimit = *iplimit
	sameIPCantInstallMultiFiles = *umulti
	runOnBackground = *oback*/ //Codeksion'a ui olarak eklenecek

	background := flag.Bool("background", false, "on background")
	flag.Parse()

	runOnBackground = *background
	if tmpVarStruct.SaveLogs {
		ioutil.WriteFile(saveLogsPath, []byte("Log file created"), os.ModePerm)
	}

}

func umain() {

	if runOnBackground {
		sifre := string(getHashInput())
		for index, value := range os.Args {
			if value == "-background" || value == "--background" {
				os.Args[index] = ""
				break
			}
		}

		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Env = append(cmd.Env, "hash="+string(sifre))
		//cmd.Stdout = os.Stdout
		cmd.Start()
		os.Exit(0)

	}

	if sifre := os.Getenv("hash"); sifre != "" {

		log.Println("goDrive PID:  ", os.Getpid())
		runOnBackgroundHash = []byte(sifre)
		runningOnBackground = true
	}
	driveInit()
	err := route()
	if err != nil {
		log.Println(err)

	}

}

func main() {
	umain()
}

// ÇORBAAA
// Tekrardan uğraşırmıyım bilmiyorum ama uğraşmaya karar verirsem projeyi doğru düzgün kodlayıp UI ekleyeceğim. @Codeksion'da yayınlanacak.

// SOUP :)
// I don't know am i rewrite this program but if i decide; SOME BLA BLA. I don't have perfect English. Use Turkish translate :D
