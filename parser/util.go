package parser

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
	"github.com/lexkong/log"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func IsVaildUrl(url string) bool {
	return govalidator.IsURL(url)
}

func HasHttpPrefix(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// 返回编码
func DetermineEncode(read io.Reader) encoding.Encoding {
	bs, err := bufio.NewReader(read).Peek(1024)
	if err != nil {
		log.Errorf(err, "", "")
		return unicode.UTF8
	}

	encode, _, _ := charset.DetermineEncoding(bs, "")
	return encode
}

func getDoc(url string) (*goquery.Document, error) {
	vaildurl := url
	if !IsVaildUrl(url) {
		vaildurl = "https://" + strings.TrimPrefix(url, "//")
	}
	log.Infof("url %s", vaildurl)
	reps, err := http.Get(vaildurl)
	defer func() {
		if reps != nil {
			reps.Body.Close()
		}
	}()
	//trycount := 0
	//// 重试
	//for err == http.ErrHandlerTimeout && trycount < MAXTRYCOUNT {
	//	reps, err = http.Get(vaildurl)
	//	trycount++
	//}
	if err != nil {
		log.Errorf(err, "url %s", url)
		return nil, err
	}
	if reps.StatusCode != http.StatusOK {
		log.Errorf(errors.New("staus code not ok"), "url %s get statu code %v\n", url, reps.StatusCode)
		return nil, errors.New("statusCode not ok")
	}

	encode := determineEncode(reps.Body)
	tsbytes := transform.NewReader(reps.Body, encode.NewEncoder())
	body, _ := ioutil.ReadAll(tsbytes)
	return goquery.NewDocumentFromReader(bytes.NewReader(body))
}

func determineEncode(read io.Reader) encoding.Encoding {
	bs, err := bufio.NewReader(read).Peek(1024)
	if err != nil {
		log.Errorf(err, "%v")
		return unicode.UTF8
	}
	encode, _, _ := charset.DetermineEncoding(bs, "")
	return encode
}
