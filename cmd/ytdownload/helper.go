package ytdownload

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

func onSigInt(path string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		err := os.Remove(path)
		if err != nil {
			fmt.Printf("\r%s", strings.Repeat(" ", 36))
			fmt.Println("\rUnable to delete the file created")
		}
		fmt.Printf("\r%s", strings.Repeat(" ", 36))
		fmt.Println("\rDownload cancelled")
		os.Exit(1)
	}()
}

//GetID extracts the id from a given yotube url
func GetID(link string) (string, error) {
	u, err := url.Parse(link)
	if err != nil {
		fmt.Println(err)
		return "nil", err
	}

	var id string
	if u.Host == "youtu.be" {
		num := strings.LastIndex(link, "/")
		id = link[num+1:]
	} else {
		par, _ := url.ParseQuery(u.RawQuery)
		id = par["v"][0]
	}
	return id, nil
}

//CreateFile is used to create a file in the users downloads directory and call the onsigint func
func CreateFile(filename string) (*os.File, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}
	dir := strings.Split(home, "/")
	home = "/" + dir[1] + "/" + dir[2] + "/Downloads/"

	//to delete partially written files on ^C
	onSigInt(home + filename)

	if _, err := os.Stat(home + filename); err == nil {
		fmt.Println("A file with the given name already exists")
		return nil, "", errors.New("file doesnt exist")
	}

	out, err := os.Create(filepath.Join(home, filepath.Base(filename)))
	if err != nil {
		//fmt.Println("Unable to create file")
		fmt.Println(err)
		return nil, "", err
	}

	return out, (home + filename), nil
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.downloaded += uint64(n)
	wc.printProgress()
	return n, nil
}

func (wc writeCounter) printProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 36))

	progress := (wc.downloaded * 100) / wc.Total

	fmt.Print("\rDownloading... " + fmt.Sprint(progress) + "% complete")
}
