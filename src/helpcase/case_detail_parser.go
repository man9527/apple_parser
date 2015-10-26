package helpcase

import (
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
	"strconv"
	"unicode/utf8"
)

type CaseDetail struct {
	SerialNo		string
	DetailUrl 		string
	LikeCount		int
	WordCount		int
	ParagraphCount	int
	ImgCount		int
	Content			string
}

var reg, _ = regexp.Compile("([0-9]+)")
const fbLikeUrl = "http://www.facebook.com/plugins/like.php?app_id=&channel=http%3A%2F%2Fstatic.ak.facebook.com%2Fconnect%2Fxd_arbiter%2F6brUqVNoWO3.js%3Fversion%3D41%23cb%3Df15f91cb6c%26domain%3Dwww.appledaily.com.tw%26origin%3Dhttp%253A%252F%252Fwww.appledaily.com.tw%252Ff22e68b69%26relation%3Dparent.parent&container_width=21&href=http%3A%2F%2Fwww.appledaily.com.tw%2Fappledaily%2Farticle%2Fheadline%2F||id1||%2F||id2||%2F&layout=button_count&locale=zh_TW&sdk=joey&send=false&show_faces=false&width=80"

func (caseDetail *CaseDetail) Parse() bool {

	var ids = reg.FindAllString(caseDetail.DetailUrl, 2)

	if len(ids)<2 {
		return false
	}

	var url = strings.Replace(fbLikeUrl, "||id1||", ids[0], 1)
	url = strings.Replace(url, "||id2||", ids[1], 1)

	doc := InvokeHttpRequest(url, fallback(caseDetail.SerialNo))

	if doc == nil {
		return false
	}

	var fblike int
	doc.Find(".pluginCountTextDisconnected").Each(func(i int, s *goquery.Selection) {
		fblike, _ = strconv.Atoi(s.Text())
	})

	caseDetail.LikeCount=fblike

	maindoc := InvokeHttpRequest(caseDetail.DetailUrl, nil)

	if maindoc == nil {
		return false
	}

	var content string
	var paraCount	int
	maindoc.Find("#introid,#bcontent").Each(func(i int, s *goquery.Selection) {
		content += s.Text()
		paraCount++
	})

	var wordcount = utf8.RuneCountInString(content)

	caseDetail.Content=content
	caseDetail.ParagraphCount=paraCount
	caseDetail.WordCount=wordcount

	var figureCount	= 0
	maindoc.Find(".lbimg").Each(func(i int, s *goquery.Selection) {
		figureCount++
	})

	caseDetail.ImgCount=figureCount

	return true
}

