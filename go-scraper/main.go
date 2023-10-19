package main

import (
	"fmt"
	"os"
	"bufio"
	colly "github.com/gocolly/colly"
	"strings"
	"passive/server"
)


func helpOption() string{
	helpMessageWelcome := "\nWelcome to passive v1.0.0\n\nOPTIONS:\n"
	helpMessageOptionFn := " -fn \tSearch with full-name\n"
	helpMessageOptionIp := " -ip \tSearch with ip address\n"
	helpMessageOptionU := " -u \tSearch with username\n"	
	return helpMessageWelcome+helpMessageOptionFn+helpMessageOptionIp+helpMessageOptionU
}

func main(){
	scanner := bufio.NewScanner(os.Stdin);
	for scanner.Scan(){
		userInput := strings.Split(scanner.Text(), " ");
		switch userInput[1] {
		case "--help":
		fmt.Println(helpOption());
		break
		case "-fn":
		fmt.Println("fn option");
		break
		case "-ip":
		fmt.Println("ip option");
		case "-u":
		username := userInput[2];
		
		c := colly.NewCollector();

		c.OnHTML("div#resultsHTML", func(e *colly.HTMLElement) {
			fmt.Printf("e = %v\n", e)
			innerText := e.Text
			fmt.Printf("this is the value := %v\n", innerText)
		})

		c.OnRequest(func(r *colly.Request) { 
			fmt.Println("Visiting: ", r.URL) 
		}) 
		c.OnError(func(_ *colly.Response, err error) { 
			fmt.Println("Something went wrong: ", err) 
		}) 	 
		c.OnResponse(func(r *colly.Response){
			fmt.Println("this is the response from the visit", r.StatusCode)
		});
		c.OnHTML("div#resultsHTML", func(e *colly.HTMLElement) {
			fmt.Printf("e = %v\n", e);
			innerText := e.Text;
			fmt.Printf("this is the value := %v\n", innerText);
		})
		c.OnScraped(func(r *colly.Response) { 
			fmt.Println(r.Request.URL, " scraped!") 
		})

		server.ReceiveAndSendDynamicallyLoadedPage(username);
		// c.Visit("https://whatsmyname.app/?q="+username);
		c.Visit("http://localhost:8080/send-html");
		c.Wait()
		}
	}
}



