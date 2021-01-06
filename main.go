package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/lvl0nax/moex_portfolio/tools"
)

func main() {
	parser := tools.NewExcelParser("broker_report.xlsx", "report")

	res, err := parser.Parse()
	if err != nil {
		fmt.Println("Something went wrong!!!!")
	}
	//fmt.Printf("Result: %v", res)

	fetcher := tools.NewMoexFetcher("SBSPB")
	err = fetcher.Fetch("10.04.2019")
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(fetcher.TickerPrices)

	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "Date")
	f.SetCellValue("Sheet1", "B1", "Spent")
	f.SetCellValue("Sheet1", "C1", "SP500 price")
	f.SetCellValue("Sheet1", "D1", "SP500 value")

	for i, row := range res {
		n := i + 2
		date := row[0]
		spentS := strings.Replace(row[1], ",", ".", -1)
		spent, err := strconv.ParseFloat(spentS, 64)
		if err != nil {
			fmt.Println(err)
		}
		dateBroker, err := time.Parse(tools.BrokerDateLayout, date)
		if err != nil {
			fmt.Println(err)
		}
		moexDate := dateBroker.Format(tools.MoexDateLayout)

		price := fetcher.TickerPrices[moexDate]
		value := spent / price

		f.SetCellValue("Sheet1", fmt.Sprintf("A%v", n), date)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%v", n), spent)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%v", n), price)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%v", n), value)
	}

	if err := f.SaveAs("tmp_result.xlsx"); err != nil {
		fmt.Println(err)
	}
}
