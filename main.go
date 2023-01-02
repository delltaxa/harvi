package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
)

var emails []string

var domain_requests map[string]int = map[string]int {
	"domain": 0,
}

func ghost(input string) string {
	url, err := url.Parse(input)

    if err != nil {
        return "??"
    }

    hostname := strings.TrimPrefix(url.Hostname(), "www.")

	return hostname
}

func ExtractLinks(url string) ([]string, error) {
	if domain_requests[ghost(url)] >= 100 {
		return []string{}, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ref := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	rmatches := ref.FindAllString(string(html), -1)
	emails = append(emails, rmatches...)

	domain_requests[ghost(url)] = domain_requests[ghost(url)] + 1

	var fclr string = Fore["GREEN"]

	if domain_requests[ghost(url)] >= 5 {
		fclr = Fore["YELLOW"]
	} 
	if domain_requests[ghost(url)] > 10 {
		fclr = Fore["RED"]
	}

	fmt.Println(fclr+"[+] " + strconv.Itoa(domain_requests[ghost(url)]) + " FROM " + ghost(url) +Fore["RESET"]+ "                                ")

	re := regexp.MustCompile(`(https?:\/\/)([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?`)
	matches := re.FindAllString(string(html), -1)

	fmt.Printf("[+] [INFO] Total of "+strconv.Itoa(len(emails))+"\r")


	emails = RemoveDuplicates(emails)

	return matches, nil
}

func extract(url string) []string {
	urls, _ := ExtractLinks(url)
	return urls
}

func us_us(urls []string) []string {
	var result []string
	for url := range urls {
		result = append(result, extract(urls[url])...)
	}

	return result
}

func RemoveDuplicates(input []string) []string {
	output := make([]string, 0)

	seen := make(map[string]bool)

	for _, value := range input {
		if !seen[value] {
			output = append(output, value)
			seen[value] = true
		}
		
	}

	return output
}


func main() {
	fmt.Printf(logo)
	
	var target string
	if len(os.Args) > 1 {
		target = os.Args[1]
	} else {
		fmt.Println("No url specfied.\nUsage: harvi <url>")
		os.Exit(0)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for sig := range c {
			_ = sig

			fmt.Println("-----------------------------------")
			for i := range emails {
				fmt.Println(emails[i])
			}

			os.Exit(0)
		}
	}()


	var urls []string = RemoveDuplicates(us_us(extract(target)))
	_ = urls

	var last int = len(urls)

	for {
		urls = RemoveDuplicates(us_us(urls))

		if len(urls) == last {
			break
		}

		last = len(urls)
	}

	for i:=0;i<len(emails); i++ {
		fmt.Println(emails[i])
	}
}