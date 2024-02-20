package scraper

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type JobListing struct {
	Company   string    `json:"company"`
	Role      string    `json:"role"`
	Locations []string  `json:"locations"`
	Link      string    `json:"link"`
	Date      time.Time `json:"date"`
}

var tableStart = `
| Company | Role | Location | Application/Link | Date Posted |
| ------- | ---- | -------- | ---------------- | ----------- |
`
var tableEnd = `

<!-- Please leave a one line gap between this and the table TABLE_END (DO NOT CHANGE THIS LINE) -->`

func Scrape() []JobListing {
	resp, err := http.Get("https://raw.githubusercontent.com/SimplifyJobs/Summer2024-Internships/dev/README.md")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	scrapedTableData := cleanResponseBody(string(body))

	var listings []JobListing

	lastCompany := ""
	for _, line := range strings.Split(scrapedTableData, "\n") {
		lineSplit := strings.Split(line, "|")
		if lineSplit[4] == " 🔒 " {
			continue
		}
		date, lessThan2Months := parseDateFromLine(lineSplit[5])
		if !lessThan2Months {
			continue
		}
		rawLocations := strings.TrimSpace(lineSplit[3])
		if strings.Contains(rawLocations, "<details>") {
			rawLocations = strings.Split(strings.Split(rawLocations, "</summary>")[1], "</details>")[0]
		}
		locations := strings.Split(strings.TrimSpace(rawLocations), "</br>")
		listing := JobListing{
			Company:   parseCompanyFromLine(lineSplit[1], lastCompany),
			Role:      strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(lineSplit[2], "🇺🇸", ""), "🛂", "")),
			Locations: locations,
			Link:      strings.Split(strings.Split(lineSplit[4], `<a href="`)[1], `">`)[0],
			Date:      date,
		}
		lastCompany = listing.Company
		listings = append(listings, listing)
	}
	return listings
}

func cleanResponseBody(body string) string {
	return strings.Split(strings.Split(body, tableStart)[1], tableEnd)[0]
}

func parseCompanyFromLine(rawCompany, lastCompany string) string {
	company := rawCompany
	if strings.Contains(company, "**[") {
		company = strings.Split(strings.Split(company, "**[")[1], "]")[0]
	} else if strings.Contains(company, "↳") {
		company = lastCompany
	}
	return strings.TrimSpace(company)
}

func parseDateFromLine(rawDate string) (time.Time, bool) {
	rawDateSplit := strings.Split(strings.TrimSpace(rawDate), " ")
	date, err := time.Parse("02-Jan-2006", fmt.Sprintf("%s-%s-2024", rawDateSplit[1], rawDateSplit[0]))
	if err != nil {
		return date, false
	}
	if date.After(time.Now()) {
		date, _ = time.Parse("02-Jan-2006", fmt.Sprintf("%s-%s-2023", rawDateSplit[1], rawDateSplit[0]))
	}
	return date, date.After(time.Now().AddDate(0, -2, 0))
}
