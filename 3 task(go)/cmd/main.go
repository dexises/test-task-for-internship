package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Influencer struct {
	Rank      string
	Name      string
	Category  string
	Followers string
	Country   string
	EngAuth   string
	EngAvg    string
}

func main() {
	url := "https://hypeauditor.com/top-instagram-all-russia/"

	influencers, err := scrapeInfluencers(url)
	if err != nil {
		log.Fatal(err)
	}

	err = writeCSV(influencers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("CSV file created successfully.")
}

func scrapeInfluencers(url string) ([]Influencer, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to request the page: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var influencers []Influencer

	doc.Find(".row__top").Each(func(i int, s *goquery.Selection) {
		rank := s.Find("span").First().Text()
		influencerName := s.Find(".contributor__name-content").Text()
		category := s.Find(".tag.topic").First().Find(".tag__content").Text()
		followers := s.Find(".row-cell.subscribers").Text()
		country := s.Find(".row-cell.audience").Text()
		engAuth := s.Find(".row-cell.authentic").Text()
		engAvg := s.Find(".row-cell.engagement").Text()

		influencer := Influencer{
			Rank:      strings.TrimSpace(rank),
			Name:      strings.TrimSpace(influencerName),
			Category:  strings.TrimSpace(category),
			Followers: strings.TrimSpace(followers),
			Country:   strings.TrimSpace(country),
			EngAuth:   strings.TrimSpace(engAuth),
			EngAvg:    strings.TrimSpace(engAvg),
		}

		influencers = append(influencers, influencer)
	})

	return influencers, nil
}

func writeCSV(influencers []Influencer) error {
	file, err := os.Create("influencers.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Rank", "Influencer", "Category", "Followers", "Country", "Eng. (Auth.)", "Eng. (Avg.)"}
	err = writer.Write(header)
	if err != nil {
		return err
	}

	for _, influencer := range influencers {
		row := []string{
			influencer.Rank,
			influencer.Name,
			influencer.Category,
			influencer.Followers,
			influencer.Country,
			influencer.EngAuth,
			influencer.EngAvg,
		}

		err := writer.Write(row)
		if err != nil {
			return err
		}
	}

	return nil
}
