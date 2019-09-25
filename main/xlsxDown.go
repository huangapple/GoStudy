package main

import (
	"fmt"
	"github.com/huangapple/xlsx"
	"net"
	"net/http"
	"time"
)

func main() {

	startHttp(80)
}

func startHttp(port int) {
	//start server
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	//mux.HandleFunc("/", replay)

	mux.HandleFunc("/download3", download3)

	fmt.Printf("监听%d端口\n", port)
	err = http.Serve(l, mux)
	if err != nil {
		panic(err)
	}
}

//边生成xlsx， 边下载
func download3(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=测试慢慢边生成数据边导.xlsx")
	w.WriteHeader(200)

	build := xlsx.NewStreamFileBuilder(w)

	//err := build.AddSheet("sheet-test", []string{"header1", "header2", "header3", "header4"}, nil)

	sheet, err := build.AddSheetN("sheet-test")

	//build.AddSheetS("sheet-test",xlsx.StreamStyle{} )

	if err != nil {
		panic(err)
	}
	{
		row := sheet.AddRow()
		cell := row.AddCell()
		cell.Merge(3, 0)
		cell.Value = "标题"
	}

	stream, err := build.Build()
	if err != nil {
		panic(err)
	}
	//style := new(xlsx.StreamStyle)
	//cell := xlsx.NewStreamCell("xxxx", *style, xlsx.CellTypeString)
	//
	//err = stream.WriteS([]xlsx.StreamCell{cell})
	//if err != nil {
	//	panic(err)
	//}

	for i := 0; i < 1000; i++ {
		err = stream.Write([]string{fmt.Sprintf("%d-1", i), fmt.Sprintf("%d-2", i), fmt.Sprintf("%d-3", i), fmt.Sprintf("%d-4", i)})
		if err != nil {
			panic(err)
		}
		fmt.Println(i)
		time.Sleep(time.Millisecond * 10)
	}

	err = stream.Close()
	if err != nil {
		panic(err)
	}
}
