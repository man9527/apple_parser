package helpcase
import (
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
	"sync"
	"log"
)

var mainUrl = "http://search.appledaily.com.tw/charity/projlist/Page/"
var HelpcasetableChannel chan int
var HelpcaseDetailChannel chan Helpcase
var WaitGroupForMainTable sync.WaitGroup
var WaitGroupForDetail sync.WaitGroup


type HelpcaseTable struct {
	Url 			string
	TotalRecords 	int
	PerPage			int
	Records			[]Helpcase
}

type Helpcase struct {
	SerialNo		string
	Title			string
	Date 			string
	Status      	string
	Amount      	int
	DetailUrl      	string
	DonationUrl		string
	PublishDate		string
}

func New(url string, perpage int) (ht *HelpcaseTable) {
	return &HelpcaseTable{Url:url, PerPage:perpage}
}

func SubHelpcaseTableParserListener() {
	for pageNum := range HelpcasetableChannel {
		perPageParser(pageNum)
	}
}

func SubHelpcaseDetailParserListener(fetch int) {
	count := 0
	for helpcase := range HelpcaseDetailChannel {
		if (fetch==-1 || count<fetch) {
			BeginToProcessHelpcase(&helpcase)
			count++
		} else {
			WaitGroupForMainTable.Done()

			func() (ok bool) {
				defer func() {recover()}()
				WaitGroupForDetail.Done()
				WaitGroupForDetail.Done()
				return true
			}()

			close(HelpcasetableChannel)
			close(HelpcaseDetailChannel)
			break
		}
	}
}

func (ht *HelpcaseTable) Parse() bool {
	log.Println("Begin the program ...")

	doc := InvokeHttpRequest(ht.Url, func() {
		log.Println("無法連結至蘋果網站取得資料，程式中止，請打開瀏覽器嘗試看看，檢查網路連線，或稍後再試")
		log.Println("Failed Url:" + ht.Url)
	})

	if doc!=nil {
		doc.Find("#charity_day").Each(func(i int, s *goquery.Selection) {
			end := strings.Index(s.Text(), "筆")
			begin := strings.Index(s.Text(), "計")
			tt, _ := strconv.Atoi(strings.Trim(s.Text()[begin + 3:end], " "))
			ht.TotalRecords = tt
		})

		log.Printf("連結至蘋果網站，一共%d筆個案\n", ht.TotalRecords)

		for i := 1; i <= ht.TotalRecords / ht.PerPage + 1; i++ {
			//for i:=67; i<=67; i++ {
			WaitGroupForMainTable.Add(1)
			result := func()(ok bool) {
				defer func() {recover()}()
				HelpcasetableChannel <- i
				return true
			}()

			if !result {
				break
			}
		}

		return true
	} else {
		return false
	}
}

func perPageParser(pageNum int) {
	defer WaitGroupForMainTable.Done()
	log.Printf("Processing page number %d ...\n", pageNum)

	var url = mainUrl + strconv.Itoa(pageNum)

	doc := InvokeHttpRequest(url,  func() {
		log.Println("無法連結至蘋果網站取得Page資料，程式中止，請打開瀏覽器嘗試看看，檢查網路連線，或稍後再試")
		log.Println("Failed Url:" + url)
	})

	keepGoing:=true

	if doc!=nil {
		doc.Find("#inquiry3 table tbody tr").Each(func(i int, s *goquery.Selection) {
			if (i > 1) { // ignore header line
				var serial, title, date, status, detail string
				var amount int

				s.Find("td").Each(func(j int, is *goquery.Selection) {
					switch j {
					case 0:
						serial = is.Text()
						detail, _ = is.Find("a").Attr("href")
					case 1:
						title = is.Text()
					case 2:
						date = is.Text()
					case 3:
						status = is.Text()
					case 4:
						amount, _ = strconv.Atoi(strings.TrimSpace(is.Text()))
					}
				})

				hcase := Helpcase{
					SerialNo: serial,
					Title: title,
					Date: date,
					Status: status,
					Amount: amount,
					DetailUrl: detail}

				if keepGoing {
					keepGoing = func() (ok bool) {
						defer func() {recover()}()
						WaitGroupForDetail.Add(1)
						HelpcaseDetailChannel <- hcase
						return true
					}()
				}
			}
		})
	}
	log.Printf("Processing page number %d done.\n", pageNum)
}

func BeginToProcessHelpcase(helpcase *Helpcase) {
	defer WaitGroupForDetail.Done()
	log.Printf("Parsing case %s\n", helpcase.SerialNo)
	isUpdate := IsHelpcaseExist(helpcase)

	if isUpdate {
		casedetail := getCaseDetail(helpcase.SerialNo, helpcase.DetailUrl)
		if casedetail!=nil {
			CheckOrCreateCaseDetail(helpcase, casedetail)
			donationDetail := getDonationDetail(helpcase.SerialNo)

			if (donationDetail!=nil) {
				CheckOrCreateDonationDetail(helpcase, donationDetail)
			}
		}

	} else {
		log.Println("Case is closed. Skip updating detail.")
	}
}

func getCaseDetail(serialNo string, url string) *CaseDetail {
	log.Printf("Parsing case detail %s\n", serialNo)
	helpDetail := &CaseDetail{SerialNo:serialNo, DetailUrl:url}
	result := helpDetail.Parse()
	if result {
		return helpDetail
	} else {
		return nil
	}
}

func getDonationDetail(serialNo string) *DonationDetail{
	log.Printf("Parsing donation detail %s\n", serialNo)

	dd := &DonationDetail{SerialNo:serialNo}
	result := dd.Parse()

	if result {
		return dd
	} else {
		return nil
	}
}

