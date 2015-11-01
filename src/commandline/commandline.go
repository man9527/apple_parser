package main
import (
	"flag"
	"myproject"
	"fmt"
	"time"
)

func main()  {
	fetch:=flag.Int("fetch",-1,"number of data fetch")
	isReset:=flag.Bool("reset",false, "reset database")
	isCheckDate:=flag.Bool("isFilterDate", false, "is filter date")
	startDateStr:=flag.String("startDate", "1900/1/1", "date")

	flag.Parse()

	var startDate time.Time
	var err error
	if *isCheckDate {
		startDate, err = time.Parse("2006/1/2", *startDateStr)

		if err != nil {
			fmt.Println("日期格式有誤，請更正之後再按\"開始\"繼續")
			return
		} else {
			fmt.Println("\n只取得" + startDate.Format("2006/01/02") + "之後的資料\n \n")
		}
	}

	myproject.Run(func(s string) {fmt.Print(s)}, *fetch, *isReset, *isCheckDate, startDate)
}