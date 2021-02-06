package main

import (
	"golang.org/x/sync/syncmap"
)

var dailyIPLimitMap syncmap.Map
var blockDownloadAfterDownloadMap syncmap.Map

var downloadIPMap syncmap.Map
var dirDriveStorageMap syncmap.Map
var fileDriveStorageMap syncmap.Map

type templateHTMLTdStruct struct {
	Href string
	Name string
	Mime string
	Size string
}

type templateHTMLStruct struct {
	PublicName string                 `json:"PublicName"`
	Tbody      []templateHTMLTdStruct `json:"Tbody"`
}

type errorHTMLStruct struct {
	ErrorCode string
	ErrorDesc string
	ErrorCaus string
}
type varStruct struct {
	PublicName                       string `json:"PublicName"`
	CredentialsPath                  string `json:"CredentialsPath"`
	DailyIPLimit                     int    `json:"DailyIPLimit"`
	SameIPCantInstallMultiFiles      bool   `json:"SameIPCantInstallMultiFiles"`
	EnableBlockDownloadAfterDownload bool   `json:"EnableBlockDownloadAfterDownload"`
	BlockDownloadAfterDownloadMin    int    `json:"BlockDownloadAfterDownloadMin"`
	UploadDownoadLimitByte           int    `json:"UploadDownoadLimit_Byte"`
	ShowLogs                         bool   `json:"ShowLogs"`
	SaveLogs                         bool   `json:"SaveLogs"`
	RoutePort                        string `json:"RoutePort"`
}
