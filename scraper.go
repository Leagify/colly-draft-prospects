package main

import (
	"encoding/csv"
	"fmt"
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

	// Categories: Rk,Chg,Player,College,P1,Ht,Wt,P2,Dif,BIO,SCT
	categories := [11]string{"Rk", "Chg", "Player", "College", "P1", "Ht", "Wt", "P2", "Dif", "BIO", "SCT"}

	// I'm getting some garbage data at the beginning of the table rows, as it appears that there are multiple tables.
	// I need to either figure out a way to ignore the first table or ignore input until the text matches one of the category headers above.

	//verify if the category has been reached.
	categoryReached := false
	ranks := make([]string, 0)
	c.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		rank := e.Text
		ranks = append(ranks, rank)
		writer.Write([]string{rank})
	})
	fmt.Println("ranks contains:", ranks)

	categoryReached = false
	c.OnHTML("tr td:nth-of-type(2)", func(e *colly.HTMLElement) {
		for _, y := range categories {
			if e.Text == y {
				categoryReached = true
			}
		}
		if categoryReached {
			fmt.Println("Change:", e.Text)
		}
	})

	categoryReached = false
	c.OnHTML("tr td:nth-of-type(3)", func(e *colly.HTMLElement) {
		for _, y := range categories {
			if e.Text == y {
				categoryReached = true
			}
		}
		if categoryReached {
			fmt.Println("Name:", e.Text)
		}
	})

	categoryReached = false
	c.OnHTML("tr td:nth-of-type(4)", func(e *colly.HTMLElement) {
		for _, y := range categories {
			if e.Text == y {
				categoryReached = true
			}
		}
		if categoryReached {
			fmt.Println("School:", e.Text)
		}
	})

	categoryReached = false
	c.OnHTML("tr td:nth-of-type(5)", func(e *colly.HTMLElement) {
		for _, y := range categories {
			if e.Text == y {
				categoryReached = true
			}
		}
		if categoryReached {
			fmt.Println("Position:", e.Text)
		}
	})

	categoryReached = false
	c.OnHTML("tr td:nth-of-type(6)", func(e *colly.HTMLElement) {
		for _, y := range categories {
			if e.Text == y {
				categoryReached = true
			}
		}
		if categoryReached {
			fmt.Println("Height:", e.Text)
		}
	})

	categoryReached = false
	c.OnHTML("tr td:nth-of-type(7)", func(e *colly.HTMLElement) {
		for _, y := range categories {
			if e.Text == y {
				categoryReached = true
			}
		}
		if categoryReached {
			fmt.Println("Weight:", e.Text)
		}
	})

	categoryReached = false
	c.OnHTML("tr td:nth-of-type(8)", func(e *colly.HTMLElement) {
		for _, y := range categories {
			if e.Text == y {
				categoryReached = true
			}
		}
		if categoryReached {
			fmt.Println("Alternate position:", e.Text)
		}
	})

	// There's a div called "calloutwifnba" and the date. I can't seem to access it yet.
	c.OnHTML(".calloutwifnba", func(e *colly.HTMLElement) {
		fmt.Println("Section containing date:", e.Text)
	})

	// This portion no longer works.  It is from the old page and produces no output.
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
	c.Visit("https://www.drafttek.com/Top-100-NFL-Draft-Prospects-2019-Page-2.asp")
	c.Visit("https://www.drafttek.com/Top-100-NFL-Draft-Prospects-2019-Page-3.asp")
	//c.Visit("https://www.drafttek.com/Top-100-NFL-Draft-Prospects-2019-Page-4.asp")

	log.Printf("Scraping finished, check file %q for results\n", fName)

	// Display collector's statistics
	log.Println(c)
}
