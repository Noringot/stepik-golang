package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

type Row struct {
	Id            int32  `xml:"id"`
	GUID          string `xml:"guid"`
	IsActive      bool   `xml:"isActive"`
	Balance       string `xml:"balance"`
	Picture       string `xml:"picture"`
	Age           string `xml:"age"`
	EyeColor      string `xml:"eyeColor"`
	FirstName     string `xml:"first_name"`
	LastName      string `xml:"last_name"`
	Gender        string `xml:"gender"`
	Company       string `xml:"company"`
	Email         string `xml:"email"`
	Phone         string `xml:"phone"`
	Address       string `xml:"address"`
	About         string `xml:"about"`
	Registered    string `xml:"registered"`
	FavoriteFruit string `xml:"favoriteFruit"`
}

func (r *Row) IsNameContain(value string) bool {
	if value == "" {
		return false
	}

	return strings.Contains(r.FirstName, value) || strings.Contains(r.LastName, value)
}

func SearchServer() {
	file, err := os.Open("test.xml")
	var findName string
	var findAbout string
	findName = "L"

	if err != nil {
		panic(err)
	}

	d := xml.NewDecoder(file)
	rows := make([]Row, 0, 2)

	for {
		token, tokenErr := d.Token()

		if tokenErr != nil {
			if tokenErr == io.EOF {
				break
			} else {
				panic(tokenErr)
			}
		}

		row := &Row{}
		var decodeErr error

		switch tok := token.(type) {
		case xml.StartElement:
			if tok.Name.Local == "row" {
				decodeErr = d.DecodeElement(row, &tok)
				isRowCorrect := false
				if row.IsNameContain(findName) {
				}

				if isRowCorrect {
					rows = append(rows, *row)
				}
			}
		}

		if decodeErr != nil {
			panic(decodeErr)
		}
	}

	fmt.Println(len(rows))
}

func main() {
	SearchServer()
}
