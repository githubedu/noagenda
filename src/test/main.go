// Copyright 2017 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build example

/* housecleaning
sudo apt-get install -y libasound2-dev ; sudo apt-get install -y libwebkit2gtk-4.0
go get github.com/hajimehoshi/go-mp3
go get github.com/hajimehoshi/oto

 */

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"github.com/hajimehoshi/oto"
	"github.com/hajimehoshi/go-mp3"
  "net/http"
  "strings"
  "github.com/dustin/go-humanize"
   "time"
)

func AsyncDownload(){
  fmt.Println("starting download")
  fileUrl := "https://mp3s.nashownotes.com/NA-1211-2020-01-26-Final.mp3"

  if err := DownloadFile("noagenda.mp3", fileUrl); err != nil {
      panic(err)
  }
  fmt.Println("end download")

}

func run() error {

  go AsyncDownload()
  fmt.Println("move to next")

  time.Sleep(2 * time.Second)
  fmt.Println("move to next after 2 seconds")

	f, err := os.Open("noagenda.mp3")
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}
  fmt.Println("after decoder")

	p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer p.Close()
  fmt.Println("after player")

	fmt.Printf("Length: %d[bytes]\n", d.Length())

	if _, err := io.Copy(p, d); err != nil {
		return err
	}
  fmt.Println("fim")
	return nil
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

type WriteCounter struct {
	Total uint64
}

func DownloadFile(filepath string, url string) error {

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
  //out, err := os.Create(filepath + ".tmp")
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	// Close the file without defer so it can happen before Rename()
	out.Close()

/*	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
*/	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
