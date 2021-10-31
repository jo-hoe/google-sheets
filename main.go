package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/jo-hoe/google-sheets/client"
	"github.com/jo-hoe/google-sheets/reader"
)

func main() {
	b, err := ioutil.ReadFile("C:\\Users\\johan\\Downloads\\airbnbnotifications-bf7c8720e0fe.json") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	client, err := client.NewServiceAccountClient(context.Background(), string(b))
	//srv, err := sheets.New(client)
	//resp, err := srv.Spreadsheets.Values.Get("1Ytu0Y6UbKewdoXlukrs1lQersBI4ynqILTuksgJAnFU", "Sheet2").Do()

	readerCloser, err := reader.NewSheetReader(client, "1Ytu0Y6UbKewdoXlukrs1lQersBI4ynqILTuksgJAnFU", "Sheet2")

	buf := new(strings.Builder)
	io.Copy(buf, readerCloser)
	// check errors
	fmt.Println(buf.String())
	//csv := csv.NewReader(readerCloser)
	//csvResult, err := csv.ReadAll()
	fmt.Printf("Error: %v", err)
}
