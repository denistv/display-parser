package pipeline

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"net/http"

	"displayCrawler/internal/domain"
)

func NewDevicesCollector(logger *zap.Logger) *DevicesCollector {
	return &DevicesCollector{logger: logger}
}

type DevicesCollector struct {
	logger *zap.Logger
}

func (d *DevicesCollector) ItemsIndex(in <-chan domain.Brand) <-chan domain.Device {
	out := make(chan domain.Device)

	go func() {
		for brand := range in {
			d.collectBrandDevices(brand, out)
		}
	}()

	return out
}
func (d *DevicesCollector) collectBrandDevices(brand domain.Brand, out chan domain.Device) {
	res, err := http.Get(brand.Href)
	if err != nil {
		d.logger.Error(fmt.Errorf("getting brand: %w", err).Error())

		return
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		d.logger.Error(fmt.Errorf("creating document from reader: %w", err).Error())

		return
	}

	doc.
		Find(".model-listing-container-80").
		Each(func(i int, s *goquery.Selection) {
			s.Find("a").
				Each(func(i int, s *goquery.Selection) {
					modelName := s.Text()
					device := domain.Device{
						Brand:        brand,
						Name:         modelName,
					}

					out <- device
				})
		})

	close(out)

	return
}

func (d *DevicesCollector) Item(in chan domain.Device, out domain.Device) error {
	return nil
}
