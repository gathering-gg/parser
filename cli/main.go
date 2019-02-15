package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gathering-gg/parser"
	"github.com/gathering-gg/parser/api"
	"github.com/gathering-gg/parser/config"
)

const fileName = "output_log.txt"

//SIFormat prints bytes in the International System of Units format
func siformat(numIn int64) string {
	suffix := "B" //just assume bytes
	num := float64(numIn)
	units := []string{"", "K", "M", "G", "T", "P", "E", "Z"}
	for _, unit := range units {
		if num < 1000.0 {
			return fmt.Sprintf("%3.1f%s%s", num, unit, suffix)
		}
		num = (num / 1000)
	}
	return fmt.Sprintf("%.1f%s%s", num, "Yi", suffix)
}

func upload(loc string) error {
	file, err := os.Open(loc)
	if err != nil {
		return err
	}
	defer file.Close()
	req, err := api.UploadFile("/upload/raw", "file", file)
	if err != nil {
		log.Printf("error creating file upload request: %v\n", err.Error())
		return err
	}
	var data interface{}
	_, err = api.Do(req, data)
	if err != nil {
		log.Printf("error uploading file: %v\n", err.Error())
		return nil
	}
	log.Printf("file upload success!")
	return nil
}

// ParseAll gets all data from a log
func ParseAll(filePath string) (gathering.UploadData, error) {
	data := gathering.UploadData{}
	f, err := os.Open(filePath)
	if err != nil {
		return data, err
	}
	alog, err := gathering.ParseLog(f)
	if err != nil {
		return data, err
	}
	col, err := alog.Collection()
	if err != nil {
		log.Printf("error getting collection: %v\n", err.Error())
	} else {
		data.Collection = col
	}
	rank, err := alog.Rank()
	if err != nil {
		log.Printf("error getting rank: %v\n", err.Error())
	} else {
		data.Rank = rank
	}
	inv, err := alog.Inventory()
	if err != nil {
		log.Printf("error getting inventory: %v\n", err.Error())
	} else {
		data.Inventory = inv
	}
	name, err := alog.Auth()
	if err != nil {
		log.Printf("error getting auth: %v\n", err.Error())
	} else {
		// TODO: Make this prettier, we only need the name
		data.Auth = &gathering.ArenaAuthRequest{
			Payload: gathering.ArenaAuthRequestPayload{
				PlayerName: string(name),
			},
		}
	}
	decks, err := alog.Decks()
	if err != nil {
		log.Printf("error getting decks: %v\n", err.Error())
	} else {
		data.Decks = decks
	}
	boosters, err := alog.Boosters()
	if err != nil {
		log.Printf("error getting boosters: %v\n", err.Error())
	} else {
		data.Boosters = boosters
	}
	matches, err := alog.Matches()
	if err != nil {
		log.Printf("error getting matches: %v\n", err.Error())
	} else {
		data.Matches = matches
	}
	events, err := alog.Events()
	if err != nil {
		log.Printf("error getting events: %v\n", err.Error())
	} else {
		data.Events = events
	}
	running, err := gathering.IsArenaRunning()
	if err != nil {
		log.Printf("error getting mtga.exe running status: %v\n", err.Error())
	} else {
		data.IsPlaying = running
	}
	return data, err
}

// onChange parse out all info and upload to the server
func onChange(f string) {
	log.Println("log file updated, parsing")
	body, err := ParseAll(f)
	if err != nil {
		log.Printf("error parsing log file: %v\n", err.Error())
	}
	log.Println("uploading body (even if error)")
	req, err := api.Upload("/upload/json", body)
	if err != nil {
		log.Printf("error creating request: %v\n", err.Error())
		return
	}
	var data map[string]interface{}
	_, err = api.Do(req, &data)
	if err != nil {
		log.Printf("error uploading data: %v\n", err.Error())
		return
	}
	log.Printf("upload success!")
}

// main
// Start the program
func main() {
	var fileFlag = flag.String("file", "", "The absolute or relative file path where the log file is located. This is useful when running on non-windows platforms where the log directory is not well known.")
	var tokenFlag = flag.String("token", "", "Required: Your authentication token")
	var uploadFlag = flag.Bool("upload", false, "Upload the log file instead of parsing. If provided, the client will not continue running, but will parse once, upload, and exit.")
	var versionFlag = flag.Bool("version", false, "Show the current running version")
	var timerFlag = flag.Int("timer", 30, "How often do you want the log file to be read in seconds? Changing this to be higher will delay updates to gathering.gg, but will increase performance. Defaults to 30 seconds")
	flag.Parse()
	if *versionFlag {
		fmt.Println(config.Version)
		return
	}
	log.Println("gathering.gg client starting")
	if *tokenFlag == "" {
		log.Fatalln("Error, need authentication token to upload data! Use `-token=TOKEN`")
	}
	api.Token = *tokenFlag
	file := *fileFlag
	if file == "" {
		user, err := user.Current()
		if err != nil {
			log.Fatalf("failed to get current user: %v\n", err.Error())
		}
		log.Printf("user home directory: '%v', using for finding log location\n", user.HomeDir)
		if gathering.LogDir == "" {
			log.Fatalf("Fatal: No log directory specified and the log location is unknown on this platform: '%v'. Please see `-help`\n", runtime.GOOS)
		}
		file = filepath.Join(user.HomeDir, gathering.LogDir, fileName)
	}
	// Finished Setup
	// At this point we should know where the log file is, what the user token
	// is and can parse the log file and begin the watch loop.
	if *uploadFlag {
		log.Println("Uploading raw log file (this may take a while)")
		onChange(file)
		upload(file)
		return
	}
	log.Printf("watching file: '%v'\n", file)
	// New Watcher
	// Checks the file size every X duration and on change will fire the
	// event. Easier to use and manage and sure to work.
	watcher := NewWatcher(file, time.Duration(*timerFlag)*time.Second)
	defer watcher.Stop()
	done := make(chan bool)
	log.Printf("adding log file location and watching: %v\n", file)
	// Main watch loop. Whenever the file is updated, reparse and upload
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("file updated, size:", siformat(event.Size))
				onChange(file)
			case err := <-watcher.Errors:
				log.Println("watcher error:", err)
				if strings.Index(err.Error(), "no such file or directory") > -1 {
					log.Println("Do you need to put quotes around the path? You may if the path has spaces in it. Use -file='/path/with spaces in it/file'")
				}
			}
		}
	}()
	log.Println("entering main watch loop")
	watcher.Start()
	<-done
}
