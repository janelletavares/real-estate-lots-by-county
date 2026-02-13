package listings

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ExtractListing parses the provided HTML string and returns
// extracted listing data as a CSV-formatted string.
// example:
// price;address;link;badge
// $11,000;130 Spring Mountain Dr, Zion Grove, PA 17985;https://www.zillow.com/homedetails/...;On Zillow for 30 days
func ExtractListings(html string) (string, int, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", 0, err
	}

	var records [][]string
	records = append(records, []string{
		"price",
		"address",
		"link",
		"badge",
	})

	// Each property card
	doc.Find(".photo-cards > li").Each(func(_ int, card *goquery.Selection) {

		// Badge
		f := func(i int, s *goquery.Selection) bool {
			link, _ := s.Attr("class")
			return strings.HasPrefix(link, "StyledPropertyCardBadge-")
		}

		badge := strings.TrimSpace(
			card.Find("span").FilterFunction(f).First().Text(),
		)
		//fmt.Printf("full badge text: %s\n", badge)

		// Price (first text containing $)
		price := ""
		card.Find("a > div > span").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			t := strings.TrimSpace(s.Text())
			if strings.Contains(t, "$") {
				price = t
				return false
			}
			return true
		})
		//fmt.Printf("price: %s\n", price)

		// Address
		address := strings.TrimSpace(
			card.Find("address, [data-testid='address'], a").First().Text(),
		)
		//fmt.Printf("address: %s\n", address)

		// Deep link
		link := ""
		if href, exists := card.Find("a[href]").First().Attr("href"); exists {
			link = strings.TrimSpace(href)
		}
		//fmt.Printf("link: %s\n", link)

		// Only emit rows that actually look like listings
		if price != "" && address != "" && link != "" {
			records = append(records, []string{
				price,
				address,
				link,
				badge,
			})
			//		} else {
			//			fmt.Println("skip: ", price, address, link, badge)
		}
		//os.Exit(0)
	})

	nav := doc.Find(".search-pagination").Find("span").First().Text()
	nextPage, err := processNavigation(nav)
	if err != nil {
		fmt.Println("warn: ", err)
	}
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = ';'
	writer.WriteAll(records[1:]) // skip headers

	if err := writer.Error(); err != nil {
		return "", 0, err
	}

	return buf.String(), nextPage, nil
}

func processNavigation(nav string) (int, error) {
	re, err := regexp.Compile(`Page (\w+) of (\w+)`)
	if err != nil {
		return 0, err
	}
	return -1, nil // @TODO remove
	matches := re.FindStringSubmatch(nav)
	first := matches[1]
	second := matches[2]
	if first != second {
		i, err := strconv.Atoi(first)
		if err != nil {
			return 0, err
		}
		i++
		return i, nil
	}
	return -1, err
}

func GetHeaders() (string, error) {
	var records [][]string
	records = append(records, []string{
		"price",
		"address",
		"link",
		"badge",
	})

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = ';'
	writer.WriteAll(records)

	if err := writer.Error(); err != nil {
		return "", err
	}

	return buf.String(), nil
}
