package tools

import (
	"errors"
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type ExcelParserInterface interface {
	Parse() ([][]string, error)
}

type ExcelParser struct {
	fileName  string
	sheetName string
	ExcelParserInterface
}

func NewExcelParser(fileName, sheetName string) *ExcelParser {
	return &ExcelParser{
		fileName:  fileName,
		sheetName: sheetName,
	}
}

func (parser *ExcelParser) Parse() ([][]string, error) {
	var result [][]string
	fmt.Println("start parsing...")
	f, err := excelize.OpenFile(parser.fileName)
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	rows, _ := f.GetRows(parser.sheetName)
	dateI := 0
	priceI := 0
	currencyI := 0

	fmt.Println("start iterating by rows...")
	fmt.Println("Number of rows", len(rows))
	for i, row := range rows {
		fmt.Print(".")
		if i == 0 {
			for n, colCell := range row {
				switch colCell {
				case "DATE":
					dateI = n
				case "PRICE":
					priceI = n
				case "CURRENCY":
					currencyI = n
				}
			}

			fmt.Printf("Date: %v, Price: %v, Currency: %v \n", dateI, priceI, currencyI)
			if dateI == 0 && priceI == 0 && currencyI == 0 {
				fmt.Println("Spreadsheet does not contain required headers: 'DATE', 'CURRENCY' and 'PRICE'")
				return result, errors.New("spreadsheet does not contain required headers: 'DATE', 'CURRENCY' and 'PRICE'")
			}
			continue
		}

		if row[dateI] != "" && row[priceI] != "" {
			//fmt.Printf("%v \t %v \t %v \n", row[dateI], row[priceI], row[currencyI])

			if row[currencyI] == "RUB" {
				result = append(result, []string{row[dateI], row[priceI]})
			}
		}
	}
	fmt.Println(".")

	return result, nil
}
