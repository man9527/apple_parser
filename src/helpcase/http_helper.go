package helpcase
import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"log"
)

var FailureQueue = make([]string, 0, 500)

func InvokeHttpRequest(url string, fallback func()) *goquery.Document {
	transport := http.Transport{ DisableKeepAlives : false }

	client := http.Client{
		Transport:&transport}

	resp, err := client.Get(url)

	if err != nil {
		log.Println("無法連結此網址取得資料，請打開瀏覽器嘗試看看，或檢查網路連線")
		log.Println("網址：" + url)
		log.Println("錯誤：" + err.Error())

		if fallback!=nil {
			fallback()
		}

		return nil
	}

	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromResponse(resp)

	return doc
}

func fallback(serailNo string) func() {
	return func() {
		FailureQueue = append(FailureQueue, serailNo)
	}
}
