package middleware

import (
	"context"
	"net/http"
	"net/url"
	"qingguo/utils"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	getMarkUrl = "https://jw.hebust.edu.cn/hbkjjw/student/xscj.stuckcj_data.jsp"
	logoutURL  = "https://jw.hebust.edu.cn/hbkjjw/DoLogoutServlet"
)

var (
	beginYear = 2024 // 2024-2025
	semester  = 0    // the frist semester
)

func GetMarkPage(client *http.Client, ctx *context.Context, redirectURL string) (*http.Response, error) {
	// Construct payload for getting marks
	payload := url.Values{
		"sjxz":             {"sjxz3"}, // Four-year program
		"ysyx":             {"yxcj"},  // Get vailed scores
		"zx":               {"1"},     // Enable major courses
		"fx":               {"1"},     // Enable minor courses
		"rxnj":             {"2023"},  // Admitted in 2023
		"ysyxS":            {"on"},
		"sjxzS":            {"on"},
		"zxC":              {"on"},
		"fxC":              {"on"},
		"xsjd":             {"1"},
		"menucode_current": {"S40303"}, // Menucode of score page
	}

	// 什么年代还在用 GBK，垃圾软件
	exportKey, err := utils.Utf8ToGBK("导出")
	if err != nil {
		log.Error("Failed to convert utf8 string: ", err)
		return nil, err
	}

	payload.Add("btnExport", exportKey) // Perform export action
	payload.Add("xn", strconv.Itoa(beginYear))
	payload.Add("xn1", strconv.Itoa(beginYear+1)) // 2024 to 2025 year
	payload.Add("xq", strconv.Itoa(semester))

	// New custom POST request
	req, err := http.NewRequestWithContext(*ctx, "POST", getMarkUrl, strings.NewReader(payload.Encode()))
	if err != nil {
		log.Error("Failed to create POST request: ", err)
		return nil, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "jw.hebust.edu.cn")
	req.Header.Set("Origin", "https://jw.hebust.edu.cn")
	req.Header.Set("Referer", "https://jw.hebust.edu.cn/hbkjjw/student/xscj.stuckcj.jsp?menucode=S40303")
	req.Header.Set("Sec-Fetch-Dest", "iframe")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Failed to post get mark request: ", err)
		return nil, err
	}

	// New custom GET request
	logoutRequest, err := http.NewRequestWithContext(*ctx, "GET", logoutURL, nil)
	if err != nil {
		log.Error("Failed to create POST request: ", err)
		return nil, err
	}

	logoutRequest.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	logoutRequest.Header.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	logoutRequest.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	logoutRequest.Header.Add("Cache-Control", "no-cache")
	logoutRequest.Header.Add("Connection", "keep-alive")
	logoutRequest.Header.Add("Host", "jw.hebust.edu.cn")
	logoutRequest.Header.Add("Pragma", "no-cache")
	logoutRequest.Header.Add("Referer", redirectURL)
	logoutRequest.Header.Add("Sec-Fetch-Dest", "document")
	logoutRequest.Header.Add("Sec-Fetch-Mode", "navigate")
	logoutRequest.Header.Add("Sec-Fetch-Site", "same-origin")
	logoutRequest.Header.Add("Sec-Fetch-User", "?1")
	logoutRequest.Header.Add("Upgrade-Insecure-Requests", "1")
	logoutRequest.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	logoutRequest.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"")
	logoutRequest.Header.Add("sec-ch-ua-mobile", "?0")
	logoutRequest.Header.Add("sec-ch-ua-platform", "\"Linux\"")

	// Logout when finished
	logoutResp, err := client.Do(logoutRequest)
	if err != nil {
		log.Error("Failed to logout: ", err)
	}
	defer logoutResp.Body.Close()
	return resp, nil
}
