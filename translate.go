package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	APPID 			= "你的百度翻译APPID"
	KEY				= "你的百度翻译密钥"
	LAN_FROM 		= "en"
	LAN_TO			= "zh"
	BDTRANSLATE_URI = "http://api.fanyi.baidu.com/api/trans/vip/translate"
)

type Result struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type Response struct {
	From string `json:"from"`
	To string `json:"to"`
	Result []Result `json:"trans_result"`
}

func genSalt() string {
	rand.Seed(time.Now().UnixNano())
	return strings.Trim(strings.Replace(fmt.Sprint(rand.Perm(10)), " ", "", -1), "[]")
}

func genSign(q string) string {
	str := fmt.Sprintf("%s%s%s%s", APPID, q, genSalt(), KEY)
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func checkTranslateErr(err error)  {
	if err != nil {
		log.Fatal("-TranslateErr: ", err)
	}
}

func Translate(q string) *Response {
	data := make(url.Values)
	data["q"] = []string{q}
	data["from"] = []string{LAN_FROM}
	data["to"] = []string{LAN_TO}
	data["appid"] = []string{APPID}
	data["salt"] = []string{genSalt()}
	data["sign"] = []string{genSign(q)}
	res, err := http.PostForm(BDTRANSLATE_URI, data)
	checkTranslateErr(err)
	defer res.Body.Close()
	response := new(Response)
	err = json.NewDecoder(res.Body).Decode(response)
	checkTranslateErr(err)
	return response
}