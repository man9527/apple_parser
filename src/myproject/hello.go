package myproject

import (
	_ "github.com/mattn/go-sqlite3"
	"helpcase"
	"fmt"
	"log"
	"github.com/tealeg/xlsx"
	"time"
	"strconv"
	"os"
	"path/filepath"
	"sync"
)

var f *os.File
var outfolder	string
var outputWindow func(output string)
var isRunning bool

func init2() {
	now :=time.Now()
	outfolder="output/" + now.Format("2006_01_02_15_04_05")
	os.MkdirAll("." + string(filepath.Separator) + outfolder ,0777);

	var err error
	f, err = os.OpenFile(outfolder + "/running.log", os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	// assign it to the standard logger
	log.SetOutput(f)

	helpcase.HelpcasetableChannel = make(chan int)
	helpcase.HelpcaseDetailChannel = make(chan helpcase.Helpcase)
	helpcase.WaitGroupForMainTable = sync.WaitGroup{}
	helpcase.WaitGroupForDetail = sync.WaitGroup{}
}

func printProgress() {
	var i=0

	for isRunning {
		time.Sleep(1 * time.Second)
		outputWindow(".")
		if i>80 {
			outputWindow("")
			i=0
		}
		i++
	}
}

func stopProcess() {
	isRunning=false
}
func Run(ow func(output string), fetch int, isReset bool) {
	if (isRunning) {
		log.Println("isRunning is true ... Skip the exeuction")
		return;
	} else {
		isRunning = true
		defer stopProcess()
	}

	init2()

	outputWindow = ow
	outputWindow("處理中\n")

	go printProgress()

	defer f.Close()

	if isReset {
		log.Println("Cleaning database ...")
		helpcase.CleanDB()
		log.Println("Cleaning database done")
	}

	go helpcase.SubHelpcaseTableParserListener()
	go helpcase.SubHelpcaseDetailParserListener(fetch)

	a := helpcase.New("http://search.appledaily.com.tw/charity/projlist/", 20)
	var result = a.Parse()

	if (!result) {
		return
	}

	helpcase.WaitGroupForMainTable.Wait()
	helpcase.WaitGroupForDetail.Wait()

	log.Println("Begin to retry failed request ...")

	for _, sno := range helpcase.FailureQueue {
		hp:=helpcase.GetHelpcase(sno)
		log.Println("retry " + hp.SerialNo)
		if hp!= nil {
			helpcase.WaitGroupForDetail.Add(1)
			helpcase.BeginToProcessHelpcase(hp)
		}
	}

	var file *xlsx.File
	var sheet1 *xlsx.Sheet

	log.Println("Generating Excel ...")

	file = xlsx.NewFile()
	sheet1 = file.AddSheet("Sheet1")
	createSheet1(sheet1)

	sheet2 := file.AddSheet("Sheet2")
	createSheet2(sheet2)

	err := file.Save(outfolder + "/main.xlsx")
	if err != nil {
		outputWindow(err.Error())
	}

	log.Println("Generating Excel Done")

	createDonation()

	stopProcess()

	outputWindow("完成!\n")
	dir, _ := os.Getwd()
	outputWindow("請在下列目錄取得檔案：" + dir + outfolder + "\n")
}

func createSheet1(sheet *xlsx.Sheet) {
	r := helpcase.GetAllHelpcase()

	var row *xlsx.Row
	var cell *xlsx.Cell

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "編號"
	cell = row.AddCell()
	cell.Value = "報導標題"
	cell = row.AddCell()
	cell.Value = "刊登日期"
	cell = row.AddCell()
	cell.Value = "狀態"
	cell = row.AddCell()
	cell.Value = "累計(元)"
	cell = row.AddCell()
	cell.Value = "捐款明細"

	font := &xlsx.Font{Color:"blue", Underline:true}
	style := xlsx.NewStyle()
	style.Font=*font

	for _, helpcase := range r {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = helpcase.SerialNo
		cell = row.AddCell()
		cell.Value = helpcase.Title
		cell = row.AddCell()
		cell.Value = helpcase.Date
		cell = row.AddCell()
		cell.Value = helpcase.Status
		cell = row.AddCell()
		cell.SetInt(helpcase.Amount)
		cell.NumFmt = "#,##0 ;(#,##0)"
		cell = row.AddCell()
		cell.SetStyle(style)
		cell.SetFormula("HYPERLINK(\"http://search.appledaily.com.tw/charity/projdetail/proj/" + helpcase.SerialNo +"\",\"明細\")")
	}
}

func createSheet2(sheet *xlsx.Sheet) {
	r := helpcase.GetAllHelpcaseDetail()

	var row *xlsx.Row
	var cell *xlsx.Cell

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "編號"
	cell = row.AddCell()
	cell.Value = "報導標題"
	cell = row.AddCell()
	cell.Value = "刊登日期"
	cell = row.AddCell()
	cell.Value = "按讚數"
	cell = row.AddCell()
	cell.Value = "段落數"
	cell = row.AddCell()
	cell.Value = "報導字數"
	cell = row.AddCell()
	cell.Value = "報導內圖片數"
	cell = row.AddCell()
	cell.Value = "報導URL"
	cell = row.AddCell()
	cell.Value = "報導內容全部"

	for _, helpcase := range r {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = helpcase.SerialNo
		cell = row.AddCell()
		cell.Value = helpcase.Title
		cell = row.AddCell()
		cell.Value = helpcase.Date
		cell = row.AddCell()
		cell.SetInt(helpcase.LikeCount)
		cell = row.AddCell()
		cell.SetInt(helpcase.ParagraphCount)
		cell = row.AddCell()
		cell.SetInt(helpcase.WordCount)
		cell = row.AddCell()
		cell.SetInt(helpcase.ImgCount)
		cell = row.AddCell()
		cell.SetFormula("HYPERLINK(\"" + helpcase.DetailUrl +"\",\"" + helpcase.DetailUrl +"\")")
		cell = row.AddCell()
		cell.Value = helpcase.Content
	}
}

func createDonation() {
	var currentYear = 0
	var currentMonth = 0

	r := helpcase.GetAllHelpcase()

	var file *xlsx.File
	var sheet1 *xlsx.Sheet
	log.Println("donation export begin")

	for _, hp := range r {

		dt := helpcase.GetAllDonationDetail(hp.SerialNo)

		if len(dt)==0 {
			continue
		}

		test, _ := time.Parse("2006/1/2", hp.Date)

		var isCurrent = isCurrentYearMonthMatch(currentYear, currentMonth, test.Year(), int(test.Month()))

		if !isCurrent {

			if file!=nil {
				var err error
				if currentMonth<10 {
					err = file.Save( outfolder + "/donation_" + strconv.Itoa(currentYear) + "0" + strconv.Itoa(currentMonth) + ".xlsx")
				} else {
					err = file.Save( outfolder + "/donation_" + strconv.Itoa(currentYear) + strconv.Itoa(currentMonth) + ".xlsx")
				}
				if err != nil {
					fmt.Printf(err.Error())
				}
			}

			file = xlsx.NewFile()

			currentYear = test.Year()
			currentMonth = int(test.Month())
		}

		sheet1 = file.AddSheet(hp.SerialNo)

		var publishDate, _ = time.Parse("2006/1/2", hp.Date)

		var cell *xlsx.Cell
		var row *xlsx.Row
		addHeader(sheet1)
		for _, donator := range dt {
			row = sheet1.AddRow()
			cell = row.AddCell()
			cell.SetString(donator.SerialNo)
			cell = row.AddCell()
			cell.Value = donator.Name
			cell = row.AddCell()
			cell.SetInt(donator.Amount)
			cell.NumFmt = "#,##0 ;(#,##0)"
			cell = row.AddCell()
			cell.Value=donator.Date

			var dDate, _ = time.Parse("2006/1/2", donator.Date)
			duration := dDate.Sub(publishDate)
			cell = row.AddCell()
			cell.Value=strconv.Itoa(int(duration.Hours()/24))
			cell = row.AddCell()
			if donator.LongFour==1 {
				cell.Value = "YES"
			} else {
				cell.Value = "NO"
			}
		}

		row = sheet1.AddRow()
		row.AddCell()
		row = sheet1.AddRow()
		cell = row.AddCell()
		cell.Value="捐款頁面的URL"
		cell = row.AddCell()
		cell.SetFormula("HYPERLINK(\"http://search.appledaily.com.tw/charity/projdetail/proj/" + hp.SerialNo +"\",\"http://search.appledaily.com.tw/charity/projdetail/proj/"+hp.SerialNo+"\")")

		row = sheet1.AddRow()
		cell = row.AddCell()
		cell.Value="出刊日期"
		cell = row.AddCell()
		cell.Value=hp.Date
		cell = row.AddCell()
		cell.Value="專案狀況"
		cell = row.AddCell()
		cell.Value=hp.Status
		row = sheet1.AddRow()
		cell = row.AddCell()
		cell.Value="捐款總計"
		cell = row.AddCell()
		cell.SetInt(hp.Amount)
		cell.NumFmt = "#,##0 ;(#,##0)"
		cell = row.AddCell()
		cell.Value="捐款筆數"
		cell = row.AddCell()
		cell.Value=strconv.Itoa(len(dt)) + "筆"
	}

	if file!=nil {
		var err error
		if currentMonth<10 {
			err = file.Save(outfolder + "/donation_" + strconv.Itoa(currentYear) + "0" + strconv.Itoa(currentMonth) + ".xlsx")
		} else {
			err = file.Save(outfolder + "/donation_" + strconv.Itoa(currentYear) + strconv.Itoa(currentMonth) + ".xlsx")
		}
		if err != nil {
			fmt.Printf(err.Error())
		}
	}

	log.Println("donation export done")
}

func isCurrentYearMonthMatch(currentYear int, currentMonth int, newYear int, newMonth int) bool {
	if currentMonth==0 || currentYear==0 || currentYear!=newYear || currentMonth!=newMonth {
		return false
	}

	return true
}

func addHeader(sheet *xlsx.Sheet) {
	var cell *xlsx.Cell
	var row *xlsx.Row
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetString("筆數")
		cell = row.AddCell()
		cell.Value = "捐款人姓名"
		cell = row.AddCell()
		cell.Value="累計(元)	"
		cell = row.AddCell()
		cell.Value="捐款明細 / 捐款日期"
		cell = row.AddCell()
		cell.Value="捐款日期與報導日期間隔差時間"
		cell = row.AddCell()
		cell.Value="捐款人姓名> 4個字"
}