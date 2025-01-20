package middleware

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	webRootPath = "https://jw.hebust.edu.cn/hbkjjw/" // Education affair system root URL for Hebei University of Science and Technology
	EASRootPath = "https://jw.hebust.edu.cn"         // Education affair system root URL for Hebei University of Science and Technology
)

func showMessage(message string) {
	log.Info("Message:", message)
}

func doPostBarLogon(response *http.Response) string {
	// Response body include a gzip data stream
	gzReader, err := gzip.NewReader(response.Body)
	if err != nil {
		log.Error("Unable to decompress gzip data stream: ", err)
		return ""
	}
	defer gzReader.Close()

	uncompressedData, err := io.ReadAll(gzReader)
	if err != nil {
		log.Error("Cannot read from gzReader ", err)
		return ""
	}

	var data map[string]interface{}
	json.Unmarshal(uncompressedData, &data)

	status, ok := data["status"].(string)
	message, okMsg := data["message"].(string)
	if !ok || !okMsg {
		log.Error("Status or message not found in response")
		return ""
	}

	log.Info("Status:", status)
	log.Info("Message:", message)

	redirectURL := ""
	if status == "200" {
		result, ok := data["result"].(string)
		if ok {
			redirectURL = EASRootPath + result
			log.Info("Login succeed, redirecting to: ", redirectURL)
		}
	} else {
		if status == "407" {
			log.Warning("Alert:", message)
			showMessage("")
		} else {
			showMessage(message)
		}
	}
	return redirectURL
}

func doBarLogin(client *http.Client, ctx *context.Context, username string, password string) string {
	// Construct request URL
	requestUrl := webRootPath + "cas/logon.action"

	// Build request
	params := url.Values{}
	params.Add("username", username)
	params.Add("password", password)
	params.Add("loginmethod", "xiqueer")

	req, err := http.NewRequestWithContext(*ctx, "POST", requestUrl, strings.NewReader(params.Encode()))
	if err != nil {
		log.Error("Failed to build request for login: ", err)
	}

	req.Header.Add("Accept", "text/plain, */*; q=0.01")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Host", "jw.hebust.edu.cn")
	req.Header.Add("Origin", "https://jw.hebust.edu.cn")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Referer", "https://jw.hebust.edu.cn/hbkjjw/cas/login.action")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Linux\"")

	// Try to get login status from EAS
	response, err := client.Do(req)
	if err != nil {
		log.Error("Failed to send login request:", err)
		return ""
	}

	// Postprocess after login
	return doPostBarLogon(response)
}

// Generate QrCode content and track user
func Uuid(length int, radix int) string {
	chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" // UUID Seed from EAS
	charArray := strings.Split(chars, "")
	uuid := make([]string, 0, length)
	rand.New(rand.NewSource(time.Now().UnixNano())) // Random seed from website

	if length > 0 {
		// Compact form
		for i := 0; i < length; i++ {
			index := rand.Intn(radix)
			uuid = append(uuid, charArray[index])
		}
	} else {
		// RFC 4122, version 4 form
		uuid = make([]string, 36)
		uuid[8], uuid[13], uuid[18], uuid[23] = "-", "-", "-", "-"
		uuid[14] = "4"

		for i := 0; i < 36; i++ {
			if uuid[i] == "" {
				r := rand.Intn(16)
				if i == 19 {
					uuid[i] = charArray[(r&0x3)|0x8] // Set to RFC 4122 High Bit
				} else {
					uuid[i] = charArray[r]
				}
			}
		}
	}

	return "smdljwxt" + strings.Join(uuid, "") // Magic string from the web side
}

func LoginByQrCode(client *http.Client, ctx *context.Context, uuid16 string) (*http.Client, string, error) {
	scanflag := true
	defer func() {
		scanflag = false // Stop polling when login succeed
	}()

	// Construct request URL
	requestUrl := webRootPath + "frame/LoginBar.jsp"

	// Reversed from EAS web request
	payload := url.Values{
		"operate": {"query"},
		"qrCode":  {uuid16},
	}

	// Build login reqeust
	loginRequest, err := http.NewRequestWithContext(*ctx, "POST", requestUrl, strings.NewReader(payload.Encode()))
	if err != nil {
		log.Error("Failed to build login request: ", err)
		return client, "", err
	}

	loginRequest.Header.Add("Accept", "*/*")
	loginRequest.Header.Add("Accept-Encoding", "deflate, br, zstd")
	loginRequest.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	loginRequest.Header.Add("Cache-Control", "no-cache")
	loginRequest.Header.Add("Connection", "keep-alive")
	loginRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	loginRequest.Header.Add("Host", "jw.hebust.edu.cn")
	loginRequest.Header.Add("Origin", "https://jw.hebust.edu.cn")
	loginRequest.Header.Add("Pragma", "no-cache")
	loginRequest.Header.Add("Referer", "https://jw.hebust.edu.cn/hbkjjw/cas/login.action")
	loginRequest.Header.Add("Sec-Ch-Ua", "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"")
	loginRequest.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	loginRequest.Header.Add("Sec-Ch-Ua-Platform", "\"Linux\"")
	loginRequest.Header.Add("Sec-Fetch-Dest", "empty")
	loginRequest.Header.Add("Sec-Fetch-Mode", "cors")
	loginRequest.Header.Add("Sec-Fetch-Site", "same-origin")
	loginRequest.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	loginRequest.Header.Add("X-Requested-With", "XMLHttpRequest")

	redirectURL := ""

	for scanflag {
		resp, err := client.Do(loginRequest)
		if err != nil {
			log.Error("Error sending request:", err)
			return client, "", err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error("Error reading response:", err)
			return client, "", err
		}

		data := string(body)
		// log.Info("username: ", data)
		// log.Info("password: ", uuid16)
		if data != "" {
			scanflag = false
			username := data
			password := uuid16
			redirectURL = doBarLogin(client, ctx, username, password)
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	return client, redirectURL, nil
}
