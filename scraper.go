package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	ucd := time.Now()
	roundedCurrentDate := time.Date(ucd.Year(), ucd.Month(), ucd.Day(), 0, 0, 0, 0, ucd.Location())
	currentDate := roundedCurrentDate.Format("2006-01-02")
	fName := "draft-prospects-" + currentDate + ".csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// Write CSV header
	writer.Write([]string{"Rank", "Name", "School", "Pos1", "Pos2", "Height", "Weight", "Change", "Date"})
	// TODO: scrape date from page.

	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("www.drafttek.com"),
	)

	// Get input from table columns, garbage data is at the beginning. It will be filtered later.
	ranks := make([]string, 0)
	c.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
		rank := e.Text
		ranks = append(ranks, rank)
	})

	changes := make([]string, 0)
	c.OnHTML("tr td:nth-of-type(2)", func(e *colly.HTMLElement) {
		change := strings.TrimSpace(e.Text)
		changes = append(changes, change)
	})

	names := make([]string, 0)
	c.OnHTML("tr td:nth-of-type(3)", func(e *colly.HTMLElement) {
		name := e.Text
		names = append(names, name)
	})

	schools := make([]string, 0)
	c.OnHTML("tr td:nth-of-type(4)", func(e *colly.HTMLElement) {
		school := e.Text
		schools = append(schools, school)
	})

	pos1s := make([]string, 0)
	c.OnHTML("tr td:nth-of-type(5)", func(e *colly.HTMLElement) {
		pos1 := e.Text
		pos1s = append(pos1s, pos1)
	})

	heights := make([]string, 0)
	c.OnHTML("tr td:nth-of-type(6)", func(e *colly.HTMLElement) {
		height := strings.Replace(e.Text, "\"", "", -1)
		heights = append(heights, height)
	})

	weights := make([]string, 0)
	c.OnHTML("tr td:nth-of-type(7)", func(e *colly.HTMLElement) {
		weight := e.Text
		weights = append(weights, weight)
	})

	pos2s := make([]string, 0)
	c.OnHTML("tr td:nth-of-type(8)", func(e *colly.HTMLElement) {
		pos2 := e.Text
		pos2s = append(pos2s, pos2)
	})

	var updatedDate string
	// There's a div called "calloutwifnba" and the date.
	c.OnHTML("html body div#outer div#wrapper div#content div#calloutwifnba strong", func(e *colly.HTMLElement) {
		runes := []rune(e.Text)

		// Expect date between parenthesis
		startIndex := strings.Index(e.Text, "(") + 1
		endIndex := strings.Index(e.Text, ")")

		// Extract date
		date := string(runes[startIndex:endIndex])

		// Extracted date should be in this form
		const dateForm = "January 2, 2006"

		// Parse into time type
		t, err := time.Parse(dateForm, date)

		updatedDate = t.Format("2006-01-02")

		parseError := err != nil

		datesMatch := currentDate == updatedDate

		fmt.Println("Unformatted string to convert:", date)
		fmt.Println("Converted date:", updatedDate)
		fmt.Println("Current Date:", currentDate)
		fmt.Println("Dates match?:", datesMatch)
		fmt.Println("Parsing errors?:", parseError, err)
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		//fmt.Println("html here - no goquery")
		//fmt.Println(e.Text)
		// Categories: Rk,Chg,Player,College,P1,Ht,Wt,P2,Dif,BIO,SCT
		// ranks cleanup
		fmt.Println("ranks length:", len(ranks))
		dataStart := Find(ranks, "Rk")
		dataStart++
		cleanRanks := ranks[dataStart:]
		fmt.Println("ranks length after cleanup:", len(cleanRanks))

		//changes cleanup
		fmt.Println("changes length before cleanup:", len(changes))
		dataStart = Find(changes, "Chg")
		dataStart++
		cleanChanges := changes[dataStart:]
		fmt.Println("changes length after cleanup:", len(cleanChanges))

		// Names cleanup
		fmt.Println("names length:", len(names))
		dataStart = Find(names, "Player")
		dataStart++
		cleanNames := names[dataStart:]
		fmt.Println("names length after cleanup:", len(cleanNames))

		// Schools cleanup
		fmt.Println("schools length:", len(schools))
		dataStart = Find(schools, "College")
		dataStart++
		cleanSchools := schools[dataStart:]
		fmt.Println("Schools length after cleanup:", len(cleanSchools))

		// Positions cleanup (both primary and alternate positions)
		fmt.Println("pos1 length:", len(pos1s))
		dataStart = Find(pos1s, "P1")
		dataStart++
		cleanPos1s := pos1s[dataStart:]
		fmt.Println("Pos1 length after cleanup:", len(cleanPos1s))

		fmt.Println("pos2 length:", len(pos2s))
		dataStart = Find(pos2s, "P2")
		dataStart++
		cleanPos2s := pos2s[dataStart:]
		fmt.Println("Pos2 length after cleanup:", len(cleanPos2s))

		// Height/Weight cleanup
		fmt.Println("heights length:", len(heights))
		dataStart = Find(heights, "Ht")
		dataStart++
		cleanHeights := heights[dataStart:]
		fmt.Println("Heights length after cleanup:", len(cleanHeights))

		fmt.Println("weights length:", len(weights))
		dataStart = Find(weights, "Wt")
		dataStart++
		cleanWeights := weights[dataStart:]
		fmt.Println("Weights length after cleanup:", len(cleanWeights))

		fmt.Println("Sample athlete:", cleanRanks[0], cleanChanges[0], cleanNames[0], cleanSchools[0],
			cleanPos1s[0], cleanPos2s[0], cleanHeights[0], cleanWeights[0])
		fmt.Println("Ranks last updated:", updatedDate)
		// Actually write the data to the CSV.
		//"Rank", "Name", "School", "Pos1", "Pos2", "Height", "Weight", "Change", "Date"
		for i, rank := range cleanRanks {
			writer.Write([]string{
				rank,
				cleanNames[i],
				cleanSchools[i],
				cleanPos1s[i],
				cleanPos2s[i],
				cleanHeights[i],
				cleanWeights[i],
				cleanChanges[i],
				updatedDate,
			})
		}

		// These slices keep getting larger, so they'll probably need to be blanked out
		// before the next page is scraped.
		ranks = nil
		changes = nil
		names = nil
		schools = nil
		pos1s = nil
		pos2s = nil
		heights = nil
		weights = nil
	})

	c.Visit("https://www.drafttek.com/Top-100-NFL-Draft-Prospects-2019.asp")
	c.Visit("https://www.drafttek.com/Top-100-NFL-Draft-Prospects-2019-Page-2.asp")
	c.Visit("https://www.drafttek.com/Top-100-NFL-Draft-Prospects-2019-Page-3.asp")
	c.Visit("https://www.drafttek.com/Top-100-NFL-Draft-Prospects-2019-Page-4.asp")

	log.Printf("Scraping finished, check file %q for results\n", fName)

	// Display collector's statistics
	log.Println(c)
}

// Find returns the smallest index i at which x == a[i],
// or len(a) if there is no such index.
func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}
