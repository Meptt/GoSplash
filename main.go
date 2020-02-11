package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/cavaliercoder/grab"
	"github.com/reujab/wallpaper"
)

const apiKEY = "11424b845fa388c2191aefe2d41bfff9dfe6a47cc51d0023463e5783871b2c4d"
const apiURL = "https://api.unsplash.com/photos/random?client_id=" + apiKEY

// Information struct to store json informations
type Information struct {
	Width       int
	Height      int
	Description string
	Urls        struct {
		Full  string
		Small string
	}
	Location struct {
		Title string
	}
}

func main() {
	data := ReadFile("details.json")
	var input string
	for {
		fmt.Println(`
 ██████╗  ██████╗ ███████╗██████╗ ██╗      █████╗ ███████╗██╗  ██╗
██╔════╝ ██╔═══██╗██╔════╝██╔══██╗██║     ██╔══██╗██╔════╝██║  ██║
██║  ███╗██║   ██║███████╗██████╔╝██║     ███████║███████╗███████║
██║   ██║██║   ██║╚════██║██╔═══╝ ██║     ██╔══██║╚════██║██╔══██║
╚██████╔╝╚██████╔╝███████║██║     ███████╗██║  ██║███████║██║  ██║
 ╚═════╝  ╚═════╝ ╚══════╝╚═╝     ╚══════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝
		`)
		fmt.Println("1 - Download random image")
		fmt.Println("2 - Get Wallpaper info")
		fmt.Println("3 - Change Wallpaper")
		fmt.Println("You can use 'q' or 'x' to quit")
		fmt.Print(": ")
		fmt.Scanln(&input)
		if input == "1" {
			if err := DownloadFile("details.json", apiURL); err != nil {
				panic(err)
			}
		} else if input == "2" {
			fmt.Printf("Photo Width : %d, Photo Height %d, Description : %s, Location : %s", data.Width, data.Height, data.Description, data.Location.Title)
		} else if input == "3" {
			go ChangeWallpaper()
		} else if input == "x" || input == "q" {
			break
		}
	}
}

// ChangeWallpaper walllpaper changer function. Get data from json file
// and downloader
func ChangeWallpaper() {
	data := ReadFile("details.json")
	dir, _ := os.Getwd()
	resp, err := grab.Get(".", data.Urls.Full)
	if err != nil {
		log.Fatal(err)
	}
	src := &resp.Filename
	out := "image.jpg"
	os.Rename(*src, out)
	background, err := wallpaper.Get()
	if err != nil {
		panic(err)
	}
	wallpaper.SetFromFile(dir + "\\image.jpg")
	fmt.Println("Wallpaper changed to :", background)
}

// ReadFile read details.json and print to screen some information
func ReadFile(filepath string) Information {
	jsonfile, _ := ioutil.ReadFile(filepath)
	jsonveri := []byte(jsonfile)
	var veri Information
	err := json.Unmarshal(jsonveri, &veri)
	if err != nil {
		log.Fatalln(err)
	}
	return veri
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
