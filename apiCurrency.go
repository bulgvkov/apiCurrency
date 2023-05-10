package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date"`
	Valute  []struct {
		Name  string `xml:"Name"`
		Value string `xml:"Value"`
	} `xml:"Valute"`
}

func getDates() []string {
	var result []string
	for i := 89; i >= 0; i-- {
		currentTime := time.Now()
		currentTime = currentTime.AddDate(0, 0, -i)
		url := "http://www.cbr.ru/scripts/XML_daily_eng.asp?date_req=" + currentTime.Format("02-01-2006")
		result = append(result, url)
	}
	return result
}

func get(url string) (ValCurs, error) {
	var result ValCurs
	resp, err := http.Get(url)
	if err != nil {
		return ValCurs{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ValCurs{}, fmt.Errorf("status error: %v", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ValCurs{}, fmt.Errorf("read body: %v", err)
	}
	err = xml.Unmarshal(data, &result)
	if err != nil {
		return ValCurs{}, err
	}
	return result, nil
}

func main() {
	var avgCurrency, counter float64
	var min, max float64
	var minDate, maxDate, valuteNameMin, valuteNameMax string
	var err error
	var result ValCurs
	var urls []string
	urls = getDates()
	for i := 0; i < len(urls); i++ {
		if result, err = get(urls[i]); err != nil {
			log.Printf("Failed to get XML: %v", err)
		}
		min, err = strconv.ParseFloat(result.Valute[0].Value, 64)
		if err != nil {
			panic(err)
			return
		}
		max, err = strconv.ParseFloat(result.Valute[0].Value, 64)
		if err != nil {
			panic(err)
			return
		}
		for j := 0; j < len(result.Valute); j++ {
			currentElement, err := strconv.ParseFloat(result.Valute[j].Value, 64)
			if err != nil {
				panic(err)
				return
			}
			if min > currentElement {
				min = currentElement
				valuteNameMin = result.Valute[j].Name
				minDate = result.Date
			}
			if max < currentElement {
				max = currentElement
				valuteNameMax = result.Valute[j].Name
				maxDate = result.Date
			}
			avgCurrency += currentElement
			counter++
		}

	}
	avgCurrency /= counter

	fmt.Printf("Максимальное значение %f, название %s, дата %s", max, valuteNameMax, maxDate)
	fmt.Printf("Минимальное значение %f, название %s, дата %s", min, valuteNameMin, minDate)
	fmt.Printf("Среднее значение курса рубля за весь период по всем валютам: %f", avgCurrency)
}
