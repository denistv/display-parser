package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)
import "go.uber.org/zap"

type Brand struct {
	ID   int64
	Name string
}

type Device struct {
	ID           int64
	Brand        Brand
	Name         string
	Kind         DeviceKind
	ModelYear    int
	Size         float32
	PanelType    string
	Resolution   string
	PixelDensity int
}

type DeviceKind string

const (
	DeviceKindMonitor DeviceKind = "monitor"
	DeviceKindTV      DeviceKind = "tv"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(255)
	}

	const siteURL = "https://www.displayspecifications.com"

	res, err := http.Get(siteURL)
	if err != nil {
		logger.Fatal(err.Error())
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logger.Fatal(err.Error())
	}


	// 1. получить список брендов
	// 2. получить список моделей бренда
	// 3. обойти каждую модель, собрать данные и создать сущность устройства
	// 4. полученные элемиенты записать в базу


	doc.
		Find(".brand-listing-container-frontpage").
		Each(func(i int, s *goquery.Selection) {
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				brand := s.Text()
				href, _ := s.Attr("href")
				logger.Info(fmt.Sprintf("brand: %s, link: %s", brand, href))
			})
		})
}
