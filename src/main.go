package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func showMessage(message, title string) {
	dialog.Message("%s", message).Title(title).Info()
}

func DownloadFileInternal(url string, outputPath string) error {
	// Send GET request
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Create the output file
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func DownloadFile(url string, outputPath string) bool {
	err := DownloadFileInternal(url, outputPath)
	if err != nil {
		fmt.Println("[-] Error downloading file:", err)
		return false
	}

	fmt.Printf("[+] Downloaded %s successfully!\n", url)
	return true
}

func ConcatenateFiles(sourceFilePath string, destinationFilePath string) error {
	// Read the source file
	sourceData, err := ioutil.ReadFile(sourceFilePath)
	if err != nil {
		return err
	}

	// Open the destination file in append mode
	destinationFile, err := os.OpenFile(destinationFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Write the source data to the destination file
	_, err = destinationFile.Write(sourceData)
	if err != nil {
		return err
	}

	return nil
}

var malwareBlockListUrls = []string{
	"https://malware-filter.gitlab.io/malware-filter/phishing-filter-hosts.txt",
	"https://v.firebog.net/hosts/Prigent-Malware.txt",
	"https://v.firebog.net/hosts/RPiList-Malware.txt",
}

var adsBlockListUrls = []string{
	"https://s3.amazonaws.com/lists.disconnect.me/simple_ad.txt",
	"https://raw.githubusercontent.com/jerryn70/GoodbyeAds/master/Hosts/GoodbyeAds.txt",
	"https://raw.githubusercontent.com/ilpl/ad-hosts/master/hosts",
	"https://raw.githubusercontent.com/jdlingyu/ad-wars/master/hosts",
	"https://raw.githubusercontent.com/jerryn70/GoodbyeAds/master/Extension/GoodbyeAds-Spotify-AdBlock.txt",
}

var trackingBlockListUrls = []string{
	"https://raw.githubusercontent.com/blocklistproject/Lists/master/tracking.txt",
	"https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt",
}

const HostLocation string = "C:\\Windows\\System32\\Drivers\\etc\\hosts"

func AddToHosts(checkValue bool, urlList []string) {
	if checkValue == true {
		for _, url := range urlList {
			DownloadFile(url, "hosts.txt")
			err := ConcatenateFiles("hosts.txt", HostLocation)
			if err != nil {
				fmt.Println("[-] Error concatenating files:", err)
				return
			}
		}
	}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Menu")

	checkboxMalware := widget.NewCheck("Malware", func(checked bool) {

	})

	checkboxAds := widget.NewCheck("Ads", func(checked bool) {

	})

	checkboxTracking := widget.NewCheck("Tracking", func(checked bool) {

	})

	buttonStart := widget.NewButton("Start", func() {
		AddToHosts(checkboxMalware.Checked, malwareBlockListUrls)
		AddToHosts(checkboxAds.Checked, adsBlockListUrls)
		AddToHosts(checkboxTracking.Checked, trackingBlockListUrls)

		showMessage("Done!", "HostProtect")
	})

	buttonExit := widget.NewButton("Exit", func() {
		myApp.Quit()
	})

	menu := container.NewVBox(
		widget.NewLabel("Select what you want to block and click start"),
		checkboxMalware,
		checkboxAds,
		checkboxTracking,
		buttonStart,
		buttonExit,
	)

	myWindow.SetContent(menu)
	myWindow.ShowAndRun()
}
