package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/hooligram/kifu"
)

func main() {
	surahsCsv, err := os.Open("surahs.csv")
	if err != nil {
		kifu.Fatal(err.Error())
	}

	surahsReader := csv.NewReader(surahsCsv)
	lines, err := surahsReader.ReadAll()
	if err != nil {
		kifu.Fatal(err.Error())
	}

	surahSzMap := map[int]int{}

	for _, line := range lines {
		surah_id, err := strconv.Atoi(line[0])
		if err != nil {
			kifu.Fatal(err.Error())
		}

		surah_sz, err := strconv.Atoi(line[1])
		if err != nil {
			kifu.Fatal(err.Error())
		}

		surahSzMap[surah_id] = surah_sz
	}

	file, err := os.Create("clearquran.csv")
	if err != nil {
		kifu.Fatal(err.Error())
	}

	file.Write([]byte(""))
	writer := csv.NewWriter(file)

	for surah_id := 1; surah_id <= 114; surah_id++ {
		lines := scrape(surah_id)
		kifu.Info("Scraping: surah_id=%v", surah_id)

		want := surahSzMap[surah_id]
		actual := len(lines)
		if actual != want {
			kifu.Fatal("Surah size mismatch: surah_id=%v, want=%v, actual=%v", surah_id, want, actual)
		}

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
