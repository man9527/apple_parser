package main
import (
	"flag"
	"myproject"
	"fmt"
)

func main()  {
	fetch:=flag.Int("fetch",-1,"number of data fetch")
	isReset:=flag.Bool("reset",false, "reset database")

	flag.Parse()

	myproject.Run(func(s string) {fmt.Print(s)}, *fetch, *isReset)
}