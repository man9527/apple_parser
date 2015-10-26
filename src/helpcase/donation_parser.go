package helpcase
import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

const donation_url = "http://search.appledaily.com.tw/charity/projdetail/proj/"

type Donator struct {
	SerialNo	string
	Name		string
	Amount		int
	Date		string
	LongFour	int
}

type DonationDetail struct {
	SerialNo		string
	PublishDate		string
	Donators		[]Donator
	Url				string
}

func (donationDetail *DonationDetail) Parse() bool {
	var targetUrl = donation_url + donationDetail.SerialNo

	doc := InvokeHttpRequest(targetUrl, fallback(donationDetail.SerialNo))

	if doc == nil {
		return false
	}

	doc.Find("#inquiry3").Each(func(i int, s *goquery.Selection) {
		if (i==0) {
			s.Find("table tbody tr").Each(func(k int, s1 *goquery.Selection) {
				if (k>1) {
					var sn string
					var name string
					var amount	int
					var date string

					s1.Find("td").Each(func(m int, s3 *goquery.Selection) {
						switch m {
							case 0: sn=s3.Text()
							case 1: name=s3.Text()
							case 2: amount, _=strconv.Atoi(strings.TrimSpace(s3.Text()))
							case 3: date=s3.Text()
						}
					})

					donator:=&Donator{SerialNo:sn, Name:name, Amount:amount, Date:date}
					donationDetail.Donators = append(donationDetail.Donators, *donator)
				}
			})
		} else {
			s.Find(".bt").Each(
				func(j int, s2 *goquery.Selection) {
					if (j==0) {
						publishDate := s2.Text()
						donationDetail.PublishDate=publishDate
					}
				})
		}
	})

	return true
}