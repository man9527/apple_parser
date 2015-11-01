package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/skratchdot/open-golang/open"
	"os"
	"flag"
	"myproject"
	"time"
	"strconv"
)

var outTE *walk.TextEdit
var inTE *walk.TextEdit
var checkFilterDate *walk.CheckBox
var inputYear *walk.LineEdit
var inputMonth *walk.LineEdit

func main() {
	fetch := flag.Int("fetch", -1, "number of data fetch")
	isReset := flag.Bool("reset", false, "reset database")

	flag.Parse()

	var currentYear = time.Now().Year();
	var currentMonth = int(time.Now().Month())-1;

	if currentMonth<=0 {
		currentMonth=12
		currentYear=currentYear-1
	}

	window := &MainWindow{
		Title:   "蘋果基金會網站資料Parser",
		MinSize: Size{600, 400},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					TextEdit{AssignTo: &inTE, Text: "請點擊\"開始\"按鈕向蘋果網站抓取資料。\n點擊\"打開輸出目錄\"取得Excel檔案\n"},
					TextEdit{AssignTo: &outTE, ReadOnly: true, RowSpan:2},

					Composite{
						Layout: HBox{},
						Children:[]Widget{
							CheckBox{AssignTo: &checkFilterDate, Text:"只取得指定年月後的資料"},
							LineEdit{AssignTo: &inputYear, Text:strconv.Itoa(currentYear)},
							Label{Text:"年"},
							LineEdit{AssignTo: &inputMonth, Text:strconv.Itoa(currentMonth)},
							Label{Text:"月"},
						},
					},
					Label{},
					PushButton{
						Text: "開始",
						OnClicked: func() {
							outTE.SetText("開始時間:" + time.Now().Format("2006/01/02 03:04:05") + "\n \n")
							var isCheckDate = checkFilterDate.Checked()
							var startDate, _ = time.Parse("2006/1/2", "1900/1/1")
							var err error
							if isCheckDate {
								startDate, err = time.Parse("2006/1/2", inputYear.Text() + "/" + inputMonth.Text() + "/1")
								if err != nil {
									LogToWindow("日期格式有誤，請更正之後再按\"開始\"繼續")
									return
								} else {
									LogToWindow("\n只取得" + startDate.Format("2006/01/02") + "之後的資料\n \n")
								}
							}

							go myproject.Run(LogToWindow, *fetch, *isReset, isCheckDate, startDate)
						},
					},
					PushButton{
						Text: "打開輸出目錄取得Excel檔案",
						OnClicked: func() {
							dir, _ := os.Getwd()
							open.Run("file:///" + dir + "/output")
						},
					},
				},

			},
		},
	}
	window.Run()
}

func LogToWindow(output string) {
	outTE.AppendText(output)
}