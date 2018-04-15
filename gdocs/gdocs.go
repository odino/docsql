package gdocs

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Download Simply download the doc from the given URL and run a bunch of sanity checks
// before dumping it into the filesystem.
func Download(url string, filename string, timeout int64) error {
	log.Println("Downloading", url, "...")

	client := http.Client{
		Timeout: time.Duration(timeout * int64(time.Second)),
	}
	resp, err := client.Get(url)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Response from the Google Docs URL was %d, expecting 200", resp.StatusCode)
	}

	if resp.Header["Content-Type"][0] != "text/tab-separated-values" {
		return fmt.Errorf("The file we downloaded has content type '%s', while we expected 'text/tab-separated-values'. Are you sure you entered the right URL?", resp.Header["Content-Type"])
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		return err
	}

	log.Printf("Doc downloaded in %s", filename)

	return nil
}
