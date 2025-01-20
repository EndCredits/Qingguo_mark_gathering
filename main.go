package main

import (
	"compress/gzip"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"qingguo/middleware"
	"qingguo/utils"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	directory = ""
)

func getMarkbyUUID(uuid16 string, ctx *context.Context) ([]byte, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("Failed to create Cookie Jar", err)
		return nil, err
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout: 60 * time.Second,
		},
	}

	retClient, redirectURL, err := middleware.LoginByQrCode(client, ctx, uuid16)
	if err != nil {
		log.Error("Error in getting http client: ", err)
		return nil, err
	}

	markResp, err := middleware.GetMarkPage(retClient, ctx, redirectURL)
	if err != nil {
		log.Error("Failed get mark: ", err)
	}

	defer markResp.Body.Close()
	retRespBody := markResp.Body
	gzReader, err := gzip.NewReader(markResp.Body)
	if err != nil {
		respBody, _ := io.ReadAll(retRespBody)
		log.Errorf("Respbody: %v", respBody)
		log.Error("Failed to create gzip reader: ", err)
		return nil, err
	}

	uncompressedData, err := io.ReadAll(gzReader)
	if err != nil {
		log.Error("Failed to read decompressed data from gzip stream: ", err)
	}

	markPage, err := utils.GBKToUtf8(utils.AsciiToString(uncompressedData))
	if err != nil {
		log.Error("Failed to get mark page: ", err)
	}

	studentID, parsedJson, err := utils.ParsePage(markPage)
	if err != nil {
		log.Error("Failed to parse page: ", err)
	}

	markResp.Body.Close()
	retClient.CloseIdleConnections()

	fileHandler, err := os.Create(directory + "/" + studentID + ".json")
	if err != nil {
		log.Fatal("Unable to save student marks json")
		panic(err)
	}
	defer fileHandler.Close()

	fileHandler.Write(parsedJson)

	return parsedJson, nil
}

func getMarkHandler(w http.ResponseWriter, r *http.Request) {
	uuid16 := r.URL.Query().Get("uuid16")
	ctx := r.Context()
	if len(uuid16) == 0 {
		log.Error("Request didn't include a uuid, rejected")
		http.Error(w, "UUID is required", http.StatusBadRequest)
		return
	}

	retJson, err := getMarkbyUUID(uuid16, &ctx)
	if err != nil {
		http.Error(w, "Unable to get mark by uuid: ", http.StatusTeapot)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	written, err := w.Write(retJson)
	if err != nil {
		log.Fatalf("Error writing response: %v", err)
	}
	log.Infof("Successfully responsed: %v", written)
}

func main() {
	flag.StringVar(&directory, "directory", "", "Define where your file saved")
	flag.Parse()

	if len(directory) == 0 {
		log.Fatal("Please define where you want your files to go (by --directory)")
		os.Exit(-1)
	}

	_, err := utils.HasWritePermission(directory)
	if err != nil {
		log.Fatal("No permission to write: ", err)
		os.Exit(-2)
	}

	http.HandleFunc("/getmark", getMarkHandler)

	log.Info("Server is listening on port 32958...")
	if err := http.ListenAndServe(":32958", nil); err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}
