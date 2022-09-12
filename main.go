package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/hooligram/kifu"
)

func main() {
	file, err := os.Create("clearquran.csv")
	if err != nil {
		kifu.Fatal(err.Error())
	}

	file.Write([]byte(""))
	writer := csv.NewWriter(file)

	for surah_id := 1; surah_id <= 114; surah_id++ {
		kifu.Info("Scraping: surah_id=%v", surah_id)
		lines := scrape(surah_id)

		for i, line := range lines {
			verse_id := i + 1
			writer.Write([]string{fmt.Sprint(surah_id), fmt.Sprint(verse_id), line})
			writer.Flush()
		}
	}

	file.Close()
}

func scrape(surah_id int) []string {
	url := fmt.Sprintf("https://www.clearquran.com/%03d.html", surah_id)
	response, err := http.Get(url)
	if err != nil {
		kifu.Fatal(err.Error())
	}

	chunks, err := io.ReadAll(response.Body)
	if err != nil {
		kifu.Fatal(err.Error())
	}

	regexVerses, err := regexp.Compile(`<p>[\n\w\.\s\/<>:,;!“”?—'’-]*<\/p>`)
	if err != nil {
		kifu.Fatal(err.Error())
	}

	match := regexVerses.FindString(string(chunks))

	regexNum, err := regexp.Compile(`\d+\.`)
	if err != nil {
		kifu.Fatal(err.Error())
	}

	verses := regexNum.Split(match, -1)
	results := []string{}

	i := 1
	if surah_id == 1 || surah_id == 9 {
		i = 0
	}

	for ; i < len(verses); i++ {
		line := verses[i]
		line = strings.ReplaceAll(line, "<p>", "")
		line = strings.ReplaceAll(line, "</p>", "")
		line = strings.ReplaceAll(line, "<span>", "")
		line = strings.ReplaceAll(line, "</span>", "")
		line = strings.TrimSpace(line)
		if line != "" {
			results = append(results, line)
		}
	}

	return results
}
