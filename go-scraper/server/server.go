package server

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"sync"
)

var htmlChan = make(chan string)
var wg sync.WaitGroup

func ReceiveAndSendDynamicallyLoadedPage(username string) {
	http.HandleFunc("/get-username", func(w http.ResponseWriter, r *http.Request){
		if r.Method == http.MethodGet{
			_, err := w.Write([]byte(username))
			if  err != nil{
				http.Error(w, "Failed to write username", http.StatusInternalServerError)
				return 
			}
		}else{
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})



	http.HandleFunc("/recieve-html", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusBadRequest)
				return
			}
			html := string(body)
			fmt.Println("HTML is:", html)

			htmlChan <- html
			w.Write([]byte("HTML received and processed"))
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/send-html", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			htmlData, ok := <-htmlChan
			if !ok {
				http.Error(w, "HTML data not available", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "text/html")

			_, err := w.Write([]byte(htmlData))
			if err != nil {
				http.Error(w, "Failed to write HTML data to response", http.StatusInternalServerError)
				return
			}

			// Close htmlChan to indicate that HTML is sent
			close(htmlChan)
			fmt.Println("HTML sent. Closing server...")
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		//once htmlChan is closed it will run the deferred statement
		//and close the server
		select {
		case <-htmlChan:
		}
	}()

	wg.Wait()
}

