package main

import (
    "github.com/lxn/walk"
    . "github.com/lxn/walk/declarative"
    "github.com/skratchdot/open-golang/open"
    "os"
    "flag"
    "myproject"
)

var outTE *walk.TextEdit
var inTE *walk.TextEdit

func main() {
    fetch := flag.Int("fetch", -1, "number of data fetch")
    isReset := flag.Bool("reset", false, "reset database")

    flag.Parse()

    window := &MainWindow{
        Title:   "蘋果基金會網站資料Parser",
        MinSize: Size{600, 400},
        Layout:  VBox{},
        Children: []Widget{
            VSplitter{
                Children: []Widget{
                    HSplitter{
                        Children: []Widget{
                            TextEdit{AssignTo: &inTE, Text: "請點擊\"開始\"按鈕向蘋果網站抓取資料。\n點擊\"打開輸出目錄\"取得Excel檔案\n"},
                            TextEdit{AssignTo: &outTE, ReadOnly: true},
                        },
                    },
                    HSplitter{
                        Children: []Widget{
                            PushButton{
                                MaxSize: Size{300, 100},
                                Text: "開始",
                                OnClicked: func() {
                                    outTE.SetText("")
                                    go myproject.Run(LogToWindow, *fetch, *isReset)
                                },
                            },
                            PushButton{
                                MaxSize: Size{300, 100},
                                Text: "打開輸出目錄取得Excel檔案",
                                OnClicked: func() {
                                    dir, _ := os.Getwd()
                                    open.Run("file:///" + dir + "/output")
                                },
                            },
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