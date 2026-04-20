package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const defaultOrderField = "Name"

var allowedOrderField = []string{"Id", "Age", "Name"}

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
	Name          string
	Gender        string `xml:"gender"`
	Company       string `xml:"company"`
	Email         string `xml:"email"`
	Phone         string `xml:"phone"`
	Address       string `xml:"address"`
	About         string `xml:"about"`
	Registered    string `xml:"registered"`
	FavoriteFruit string `xml:"favoriteFruit"`
}

func (r *Row) IsNameOrAboutContain(value string) bool {
	if value == "" {
		return false
	}

	return strings.Contains(r.FirstName, value) || strings.Contains(r.LastName, value)
}

func _SearchServer(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	orderField := r.URL.Query().Get("order_field")
	orderBy := r.URL.Query().Get("order_by")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	_, err := strconv.Atoi(offset)

	if err != nil {
		http.Error(w, "cannot get offset as int", http.StatusBadRequest)
	}

	fmt.Println(query, orderField, orderBy, limit, offset)
}

// query, orderField string, orderBy, limit, offset int
func SearchServer(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	offset := r.URL.Query().Get("offset")
	lim := r.URL.Query().Get("limit")
	orderField := r.URL.Query().Get("order_field")
	orderByUrl := r.URL.Query().Get("order_by")

	if orderField == "" {
		orderField = defaultOrderField
	}

	if !slices.Contains(allowedOrderField, orderField) {
		http.Error(w, ErrorBadOrderField, http.StatusBadRequest)
	}

	fmt.Printf("query: %s,\noffset: %s,\nlim: %s,\norderField: %s,\norderByUrl: %s", query, offset, lim, orderField, orderByUrl)

	return

	file, err := os.Open("test.xml")

	if err != nil {
		panic(err)
	}

	d := xml.NewDecoder(file)
	rows := make([]Row, 0, 2)

	handled := 0
	toSkip, err := strconv.Atoi(offset)

	if err != nil {
		http.Error(w, "cannot get offset as int", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(lim)

	if err != nil {
		http.Error(w, "cannot get limit as int", http.StatusBadRequest)
		return
	}

	orderBy, err := strconv.Atoi(orderByUrl)

	if err != nil {
		http.Error(w, "cannot get orderBy as int", http.StatusBadRequest)
		return
	}

	fmt.Printf("query: %s\norderField: %s\norderBy: %s\nlimit: %s\noffset: %s\n", query, orderField, orderByUrl, offset, lim)
XML:
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
				isRowCorrect := true
				if query != "" {
					isRowCorrect = row.IsNameOrAboutContain(query)
				}

				if isRowCorrect {
					if toSkip > 0 {
						toSkip--
					} else {
						if handled >= limit {
							break XML
						}

						row.Name = row.FirstName + " " + row.LastName
						rows = append(rows, *row)
						handled++
					}
				}
			}
		}

		if decodeErr != nil {
			panic(decodeErr)
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		switch orderField {
		case "id":
			if orderBy == 1 {
				return rows[i].Id > rows[j].Id
			} else {
				return rows[i].Id < rows[j].Id
			}
		case "age":
			if orderBy == 1 {
				return rows[i].Age > rows[j].Age
			} else {
				return rows[i].Age < rows[j].Age
			}
		case "name", "":
			if orderBy == 1 {
				return rows[i].Name > rows[j].Name
			} else {
				return rows[i].Name < rows[j].Name
			}
		default:
			return false
		}
	})

	res, err := json.Marshal(rows)

	if err != nil {
		http.Error(w, "cannot pack result into json", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func main() {
	http.HandleFunc("/search", SearchServer)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
	// res, _ := SearchServer("", "id", -1, 4, 1)

	// PrintResult(*res)
}
