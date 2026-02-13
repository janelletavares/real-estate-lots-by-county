package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	listings "github.com/janelletavares/real-estate-lots-by-county/helpers"
)

func fetch(myURL string, key string) (string, error) {
	// Create client
	client := &http.Client{}
	// Create request
	// "pagination":{"currentPage":2}

	// for sale state = {"pagination":{},"isMapVisible":true,"mapBounds":{"west":-77.83640398828125,"east":-74.62290301171875,"south":39.75647040602409,"north":41.675645297572785},"regionSelection":[{"regionId":3028}],"filterState":{"sort":{"value":"globalrelevanceex"},"price":{"min":30000,"max":300000},"apa":{"value":false},"manu":{"value":false},"con":{"value":false},"apco":{"value":false},"mf":{"value":false},"tow":{"value":false},"sf":{"value":false}},"isListVisible":true,"mapZoom":9,"usersSearchTerm":"Schuylkill County PA"}

	// for sale with date https://www.zillow.com/schuylkill-county-pa/land/?searchQueryState={"pagination":{},"isMapVisible":true,"mapBounds":{"west":-77.83640398828125,"east":-74.62290301171875,"south":39.75647040602409,"north":41.675645297572785},"mapZoom":9,"usersSearchTerm":"Schuylkill County PA","regionSelection":[{"regionId":3028}],"filterState":{"sort":{"value":"globalrelevanceex"},"price":{"min":30000,"max":300000},"apa":{"value":false},"apco":{"value":false},					  "con":{"value":false},"manu":{"value":false},"mf":{"value":false},							"sf":{"value":false},"tow":{"value":false},"isListVisible":true,}

	//sold within last yr https://www.zillow.com/schuylkill-county-pa/land/?searchQueryState={"pagination":{},"isMapVisible":true,"mapBounds":{"west":-79.4431544765625,"east":-73.0161525234375,"south":38.77621078215194,"north":42.61442086067325},   "mapZoom":8,"usersSearchTerm":"Schuylkill County PA","regionSelection":[{"regionId":3028}],"filterState":{"sort":{"value":"globalrelevanceex"},"price":{"min":30000,"max":300000},"apa":{"value":false},"apco":{"value":false},"doz":{"value":"12m"},"con":{"value":false},"manu":{"value":false},"mf":{"value":false},"mp":{"min":150,"max":1500},"sf":{"value":false},"tow":{"value":false},"isListVisible":true}

	//https://www.zillow.com/schuylkill-county-pa/land/?searchQueryState=%7B%22pagination%22%3A%7B%7D%2C%22isMapVisible%22%3Atrue%2C%22mapBounds%22%3A%7B%22west%22%3A-77.83640398828125%2C%22east%22%3A-74.62290301171875%2C%22south%22%3A39.75647040602409%2C%22north%22%3A41.675645297572785%7D%2C%22regionSelection%22%3A%5B%7B%22regionId%22%3A3028%7D%5D%2C%22filterState%22%3A%7B%22sort%22%3A%7B%22value%22%3A%22globalrelevanceex%22%7D%2C%22price%22%3A%7B%22min%22%3A30000%2C%22max%22%3A300000%7D%2C%22apa%22%3A%7B%22value%22%3Afalse%7D%2C%22manu%22%3A%7B%22value%22%3Afalse%7D%2C%22con%22%3A%7B%22value%22%3Afalse%7D%2C%22apco%22%3A%7B%22value%22%3Afalse%7D%2C%22mf%22%3A%7B%22value%22%3Afalse%7D%2C%22tow%22%3A%7B%22value%22%3Afalse%7D%2C%22sf%22%3A%7B%22value%22%3Afalse%7D%7D%2C%22isListVisible%22%3Atrue%2C%22mapZoom%22%3A9%2C%22usersSearchTerm%22%3A%22Schuylkill%20County%20PA%22%7D
	//myURL := "https://www.zillow.com/Schuylkill-County-PA/land"
	// @TODO add pagination
	//@TODO pull API key from the environment
	formattedURL := fmt.Sprintf("https://app.scrapingbee.com/api/v1/?api_key=FGWTJZWBOID0XSMI3XWEKS2C7HCAXOPRD9ATXQG2GVQW64P7UESOSO4JFVARIFL8CJUOHVU3BFA4VM0I&url=%s", myURL)
	req, err := http.NewRequest("GET", formattedURL, nil)
	if err != nil {
		fmt.Println("Failure : ", err)
		return "", err
	}
	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		fmt.Println(parseFormErr)
		return "", parseFormErr
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failure : ", err)
		return "", err
	}

	// Read Response Body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failure : ", err)
		return "", err
	}

	// Display Results
	fmt.Println("response Status : ", resp.Status)
	fmt.Println("response Headers : ", resp.Header)
	fmt.Println("response Body : ", string(respBody))
	err = os.WriteFile("initial.html", []byte(respBody), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return "", err
	}
	// @TODO transform body to csv
	// listings.ExtractListings(respBody)
	return string(respBody), nil
}

func fetchAll() {
	//	sendClassic()
	//key := os.Getenv("SCRAPING_KEY")
	var lots string
	nextPage := 0
	headers, err := listings.GetHeaders()
	if err != nil {
		panic(err)
	}
	lots = headers

	data, err := os.ReadFile("input/one.json")
	if err != nil {
		panic(err)
	}

	var rows [][]string
	err = json.Unmarshal(data, &rows)
	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		fmt.Printf("row: county: %+v state: %+v\n", row[0], row[1])
		county := row[0]
		state := row[1]

		forSaleOrSold := [2]bool{true, false}

		for _, forSale := range forSaleOrSold {
			fmt.Printf("doing search for %+v\n", forSale)
			for nextPage > -1 {
				formattedURL, _ := generateURL(county, state, forSale, 0)
				fmt.Println(formattedURL)
				// @TODO replace with scraper call
				b, err := os.ReadFile("html/initial.html")
				if err != nil {
					fmt.Println("Error reading file:", err)
					return
				}
				_, nextPage, err = listings.ExtractListings(string(b))
				if err != nil {
					fmt.Println("Error reading file:", err)
					return
				}

			}
			nextPage = 0
			fmt.Println(lots) // @TODO write to either for sale or sold file
			tag := "for-sale"
			if !forSale {
				tag = "sold"
			}
			err = os.WriteFile(fmt.Sprintf("output/wednesday-%s-%s-%s.csv", county, state, tag), []byte(headers), 0600)
		}
	}
}

func main() {
	/*
		action := os.Args[1]
		switch action {
		case "fetch":
			fetchAll()
		case "report":
			fallthrough
		default:
			fmt.Println("not implemented")
		}

	*/
	// create concurrency groups
	concurrency := make(chan string, 2)
	output := make(chan string)
	go startScraper(concurrency, output)
	go startScraper(concurrency, output)

	states, err := os.ReadDir("../not_done")
	if err != nil {
		log.Fatal(err)
	}

	for _, state := range states {
		fmt.Println(state.Name(), state.IsDir())
		counties, err := os.ReadDir(fmt.Sprintf("../not_done/%s", state.Name()))
		if err != nil {
			log.Fatal(err)
		}

		for _, county := range counties {
			// add to workers
			select {
			case msg := <-output:
				fmt.Println("received message: ", msg)
			default:
				fmt.Println("no message received")
			}

			concurrency <- fmt.Sprintf("python3 ../entrypoint.py %s %s", state.Name(), county.Name())

			//fmt.Println(county.Name(), county.IsDir())
			//d, err := os.ReadFile(fmt.Sprintf("../not_done/%s/%s", state.Name(), county.Name()))
			//if err != nil {
			//	log.Fatal(err)
			//}
			//zipcodes, err := json.Marshal(d)
			//if err != nil {
			//	log.Fatal(err)
			//}
			//fmt.Printf("%+v\n", zipcodes)
		}

	}
}

func startScraper(input chan string, output chan string) {
	command := <-input
	items := strings.Split(command, " ")
	console, err := exec.Command(items[0], items[1:]...).Output()
	if err != nil {
		var execErr *exec.Error
		var exitErr *exec.ExitError
		switch {
		case errors.As(err, &execErr):
			output <- fmt.Sprintf("failed executing:", err)
		case errors.As(err, &exitErr):
			exitCode := exitErr.ExitCode()
			output <- fmt.Sprintf("command exit rc =", exitCode)
		default:
			panic(err)
		}
	} else {
		output <- string(console)
	}
}

type ValueBlock struct {
	Value bool `json:"value"`
}
type ValueString struct {
	Value string `json:"value"`
}
type forSaleFilterState struct {
	Sort ValueString `json:"sort"`
	APA  ValueBlock  `json:"apa"`
	APCO ValueBlock  `json:"apco"`
	CON  ValueBlock  `json:"con"`
	MANU ValueBlock  `json:"manu"`
	MF   ValueBlock  `json:"mf"`
	TOW  ValueBlock  `json:"tow"`
}
type SoldFilterState struct {
	Sort ValueString `json:"sort"`
	APA  ValueBlock  `json:"apa"`
	APCO ValueBlock  `json:"apco"`
	CON  ValueBlock  `json:"con"`
	MANU ValueBlock  `json:"manu"`
	MF   ValueBlock  `json:"mf"`
	TOW  ValueBlock  `json:"tow"`
	Doz  Timeframe   `json:"doz,omitempty"`
	Mp   Range       `json:"mp,omitempty"`
	Rs   ValueBlock  `json:"rs,omitempty"`
	Fsba ValueBlock  `json:"fsba,omitempty"`
	Fsbo ValueBlock  `json:"fsbo,omitempty"`
	Nc   ValueBlock  `json:"nc,omitempty"`
	Cmsn ValueBlock  `json:"cmsn,omitempty"`
	Auc  ValueBlock  `json:"auc,omitempty"`
	Fore ValueBlock  `json:"fore,omitempty"`
}

type Timeframe struct {
	Value string `json:"value"`
}

type Range struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type Page struct {
	CurrentPage int `json:"currentPage,omitempty"`
}

type ForSaleSearchQueryState struct {
	Pagination    Page               `json:"pagination"`
	UST           string             `json:"usersSearchTerm"`
	IsMapVisible  bool               `json:"isMapVisible"`
	MapZoom       int                `json:"mapZoom"`
	FilterState   forSaleFilterState `json:"filterState"`
	IsListVisible bool               `json:"isListVisible"`
}

type SoldSearchQueryState struct {
	Pagination    Page            `json:"pagination"`
	UST           string          `json:"usersSearchTerm"`
	IsMapVisible  bool            `json:"isMapVisible"`
	MapZoom       int             `json:"mapZoom"`
	FilterState   SoldFilterState `json:"filterState"`
	IsListVisible bool            `json:"isListVisible"`
}

// category?

func generateURL(county string, state string, forSale bool, pageID int) (string, error) {
	var query any
	if forSale {
		var fs = ForSaleSearchQueryState{
			IsMapVisible: true,
			UST:          fmt.Sprintf("%s, %s", county, state),
			FilterState: forSaleFilterState{
				Sort: ValueString{Value: "globalrelevanceex"},
				APA:  ValueBlock{Value: false},
				APCO: ValueBlock{Value: false},
				CON:  ValueBlock{Value: false},
				MANU: ValueBlock{Value: false},
				MF:   ValueBlock{Value: false},
				TOW:  ValueBlock{Value: false},
			},
			IsListVisible: true,
			MapZoom:       9,
		}
		query = fs
	} else {
		var s = SoldSearchQueryState{
			IsMapVisible: true,
			UST:          fmt.Sprintf("%s, %s", county, state),
			FilterState: SoldFilterState{
				Sort: ValueString{Value: "globalrelevanceex"},
				APA:  ValueBlock{Value: false},
				APCO: ValueBlock{Value: false},
				CON:  ValueBlock{Value: false},
				MANU: ValueBlock{Value: false},
				MF:   ValueBlock{Value: false},
				TOW:  ValueBlock{Value: false},
				Doz:  Timeframe{Value: "12m"}, // longer?
				Mp:   Range{Min: 150, Max: 1500},
				Rs:   ValueBlock{Value: false},
				Fsba: ValueBlock{Value: false},
				Fsbo: ValueBlock{Value: false},
				Nc:   ValueBlock{Value: false},
				Cmsn: ValueBlock{Value: false},
				Auc:  ValueBlock{Value: false},
				Fore: ValueBlock{Value: false},
			},
			IsListVisible: true,
			MapZoom:       9,
		}
		query = s
	}

	bytes, err := json.Marshal(query)
	if err != nil {
		return "", err
	}
	searchQueryState := url.PathEscape(string(bytes))

	// assemble URL
	forSaleString := "land/"
	if !forSale {
		forSaleString = "sold/"
	}
	if pageID > 0 {
		forSaleString += fmt.Sprintf("/%i_p/", pageID)
	}
	location := strings.ReplaceAll(county, " ", "-")
	location = fmt.Sprintf("%s-%s", strings.ToLower(location), strings.ToLower(state))
	myURL := fmt.Sprintf("https://www.zillow.com/%s/%s?searchQueryState=%v", location, forSaleString, searchQueryState)
	return myURL, nil
}
