package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/raifpy/Go/errHandler"
	"github.com/raifpy/Go/saes"
	"github.com/tcnksm/go-input"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.

	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {

		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	authCode, err := inputUI.Ask("Go to the following link in your browser then type the "+
		"authorization code: \n"+authURL, &input.Options{Loop: true, Required: true})

	errHandler.HandlerExit(err)

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	fbyte, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	tok := &oauth2.Token{}
	var veri []byte
	if runOnBackgroundHash != nil {
		veri = runOnBackgroundHash
	} else {
		veri = getHashInput()
	}
	oye, err := saes.Decrypt(fbyte, veri)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.NewDecoder(bytes.NewReader(oye)).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
ulan:
	data, err := json.Marshal(token)
	errHandler.HandlerExit(err)
	ehe, err := saes.Encrypt(data, getHashInput())
	data = nil
	if errHandler.HandlerBool(err) {
		goto ulan
	}
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()

	f.Write(ehe)
	ehe = nil
}

func getDrive(credentialsPath string) (*drive.Service, error) {
	credentialsByte, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return nil, err
	}
	conf, err := google.ConfigFromJSON(credentialsByte, drive.DriveReadonlyScope)
	credentialsByte = nil
	if err != nil {
		return nil, err
	}
	cli := getClient(conf)
	//fmt.Println()

	return drive.New(cli)

}
