package main

import (
	"encoding/csv"
	"log"
	"os"
	"time"

	"strings"

	"github.com/gocolly/colly"
)

func main() {
	currentDate := time.Now().Format("2006-01-02")
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
	writer.Write([]string{"Rank", "Name", "School"}) // "Position" removed at the moment.
	// TODO: scrape date from top right of page.  Looks like it's in an unnamed <p> element.

	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("www.drafttek.com"),
	)

	// Extract product details
	c.OnHTML(".BigBoardMainTable", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, t *colly.HTMLElement) {
			rank := t.ChildText("td.BigBoardTable2.boldplayerlabel")
			name := t.ChildText("span.boldplayerlabel")
			// I can only get the school by getting player name and school as concatenated text.  Annoying.
			nameSchool := t.ChildText("td.BigBoardTable.playername")
			nameIndex := strings.LastIndex(nameSchool, name)
			nameIndex += len(name)
			nameRunes := []rune(nameSchool)
			school := string(nameRunes[nameIndex:])
			// Position logic.  May be re-added at later time, as it doesn't work with multi-position prospects.
			//pos := t.ChildText("td.BigBoardTable2.stats20")
			//posIndex := strings.LastIndex(pos, "\n")
			//altIndex := strings.Index(pos, " ")
			//positionRunes := []rune(pos)
			//playerPosition := ""
			//if posIndex != -1 {
			//	playerPosition = string(positionRunes[:posIndex])
			//	playerPosition = strings.TrimRight(playerPosition, "\n\t")
			//}
			if len(name) > 0 {
				writer.Write([]string{
					rank,
					name,
					school,
					//playerPosition,
				})
			}
		})
	})

	c.Visit("https://www.drafttek.com/Top-100-NFL-Draft-Prospects-2019.asp")

	log.Printf("Scraping finished, check file %q for results\n", fName)

	// Display collector's statistics
	log.Println(c)
}
