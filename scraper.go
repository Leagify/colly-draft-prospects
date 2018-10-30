package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

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
	writer.Write([]string{"Rank", "Change", "Name", "School", "Pos1", "Pos2", "Height", "Weight"}) // "Position" removed at the moment.
	// TODO: scrape date from top right of page.  Looks like it's in an unnamed <p> element.

	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("www.drafttek.com"),
	)

	// Categories: Rk,Chg,Player,College,P1,Ht,Wt,P2,Dif,BIO,SCT
	categories := [11]string{"Rk", "Chg", "Player", "College", "P1", "Ht", "Wt", "P2", "Dif", "BIO", "SCT"}
	fmt.Println(categories)

	// I'm getting some garbage data at the beginning of the table rows, as it appears that there are multiple tables.
	// I need to either figure out a way to ignore the first table or ignore input until the text matches one of the category headers above.

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		//fmt.Println("this may be the second table")
	})

	ranks := make([]string, 0)
	c.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		rank := e.Text
		ranks = append(ranks, rank)
	})

	changes := make([]string, 0)
	c.OnHTML("tr td:nth-of-type(2)", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		change := e.Text
		changes = append(changes, change)
	})

	names := make([]string, 0)
	c.OnHTML("tr td:nth-of-type(3)", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		name := e.Text
		names = append(names, name)
	})

	schools := make([]string, 0)
	c.OnHTML("tr td:nth-of-type(4)", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		school := e.Text
		schools = append(schools, school)
	})

	pos1s := make([]string, 0)
	c.OnHTML("tr td:nth-of-type(5)", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		pos1 := e.Text
		pos1s = append(pos1s, pos1)
	})

	heights := make([]string, 0)
	c.OnHTML("tr td:nth-of-type(6)", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		height := e.Text
		heights = append(heights, height)
	})

	weights := make([]string, 0)
	c.OnHTML("tr td:nth-of-type(7)", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		weight := e.Text
		weights = append(weights, weight)
	})

	pos2s := make([]string, 0)
	c.OnHTML("tr td:nth-of-type(8)", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		pos2 := e.Text
		pos2s = append(pos2s, pos2)
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		fmt.Println("html here - no goquery")
		//fmt.Println(e.Text)
		fmt.Println("ranks length:", len(ranks))
		ranksDataStart := Find(ranks, "Rk")
		fmt.Println("actual data starts at", ranksDataStart)
		ranksDataStart += 1
		fmt.Println("first rank value is", ranks[ranksDataStart])
		fmt.Println("changes length:", len(changes))
		fmt.Println("names length:", len(names))
		fmt.Println("schools length:", len(schools))
		fmt.Println("pos1 length:", len(pos1s))
		fmt.Println("heights length:", len(heights))
		fmt.Println("weights length:", len(weights))
		fmt.Println("pos2 length:", len(pos2s))
		// These slices keep getting larger, so they'll probably need to be blanked out
		// before the next page is scraped.
	})

	// There's a div called "calloutwifnba" and the date. I can't seem to access it yet.
	c.OnHTML("calloutwifnba", func(e *colly.HTMLElement) {
		fmt.Println("Section containing date:", e.Text)
	})

	// This portion no longer works.  It is from the old page and produces no output.
	// Extract product details
	// c.OnHTML(".BigBoardMainTable", func(e *colly.HTMLElement) {
	// 	e.ForEach("tr", func(i int, t *colly.HTMLElement) {
	// 		rank := t.ChildText("td.BigBoardTable2.boldplayerlabel")
	// 		name := t.ChildText("span.boldplayerlabel")
	// 		// I can only get the school by getting player name and school as concatenated text.  Annoying.
	// 		nameSchool := t.ChildText("td.BigBoardTable.playername")
	// 		nameIndex := strings.LastIndex(nameSchool, name)
	// 		nameIndex += len(name)
	// 		nameRunes := []rune(nameSchool)
	// 		school := string(nameRunes[nameIndex:])
	// 		// Position logic.  May be re-added at later time, as it doesn't work with multi-position prospects.
	// 		//pos := t.ChildText("td.BigBoardTable2.stats20")
	// 		//posIndex := strings.LastIndex(pos, "\n")
	// 		//altIndex := strings.Index(pos, " ")
	// 		//positionRunes := []rune(pos)
	// 		//playerPosition := ""
	// 		//if posIndex != -1 {
	// 		//	playerPosition = string(positionRunes[:posIndex])
	// 		//	playerPosition = strings.TrimRight(playerPosition, "\n\t")
	// 		//}
	// 		if len(name) > 0 {
	// 			writer.Write([]string{
	// 				rank,
	// 				name,
	// 				school,
	// 				//playerPosition,
	// 			})
	// 		}
	// 	})
	// })

	c.Visit("https://www.drafttek.com/Top-100-NFL-Draft-Prospects-2019.asp")
	c.Visit("https://www.drafttek.com/Top-100-NFL-Draft-Prospects-2019-Page-2.asp")
	c.Visit("https://www.drafttek.com/Top-100-NFL-Draft-Prospects-2019-Page-3.asp")
	//c.Visit("https://www.drafttek.com/Top-100-NFL-Draft-Prospects-2019-Page-4.asp")

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
