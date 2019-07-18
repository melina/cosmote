// Get the balance of the Cosmote Prepaid Mobile Internet

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	port := ":8800"
	http.HandleFunc("/", checkRemainingMBs)
	log.Println("Listening on", port)
	log.Fatal(http.ListenAndServe(port, nil))
}


func checkRemainingMBs(w http.ResponseWriter, r *http.Request) {
	URL := "http://ciotgprepaid.cosmote.gr/iotg/iotg.portal?_nfpb=true&_pageLabel=RemainderPage&iotgRemainder_right_actionOverride=%2FFlows%2FRemainder%2Fbegin&_windowLabel=iotgRemainder_right"
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		panic(err)
	}

	// Create HTTP client with timeout
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Get the response body as a string
	dataInBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	pageContent := string(dataInBytes)

	// Find the substring
	MBStartIndex := strings.Index(pageContent, "<span class=\"green\">Χρήση που απομένει:</span>&nbsp;")
	if MBStartIndex == -1 {
		fmt.Fprintf(w, "No title element found")
		os.Exit(0)
	}
	// Get the start index
	MBStartIndex += 68

	// find the index of the MB
	MBEndIndex := strings.Index(pageContent, "MB")
	if MBEndIndex == -1 {
		fmt.Fprintf(w, "No closing tag for title found.")
		os.Exit(0)
	}

	MBLeft := []byte(pageContent[MBStartIndex:MBEndIndex])

	fmt.Fprintf(w, "Remaining: %s MB \n", MBLeft)

}
