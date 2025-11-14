package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	easyjson "github.com/mailru/easyjson"
)

// вам надо написать более быструю оптимальную этой функции
/*
	!!! !!! !!!
	обратите внимание - в задании обязательно нужен отчет
	делать его лучше в самом начале, когда вы видите уже узкие места, но еще не оптимизировалм их
	так же обратите внимание на команду в параметром -http
	перечитайте еще раз задание
	!!! !!! !!!
*/

type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Browsers []string `json:"browsers"`
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer func() {
		if file.Close() != nil {
			fmt.Printf("Ошибка закрытия файла: %v \n", err)

		}
	}()

	seenBrowsers := make(map[string]int, 100)
	uniqueBrowsers := 0
	foundUsers := ""
	i := 0

	reader := bufio.NewReader(file)
	u := &User{}

	for {
		i++
		line, err := reader.ReadBytes('\n')

		if err != nil {
			if err.Error() != "EOF" {
				panic(err)
			}

			break
		}

		isHaveAndroid := false
		isHaveMSIE := false

		// fmt.Printf("%v %v\n", err, line)
		err = easyjson.Unmarshal(line, u)

		if err != nil {
			panic(err)
		}

		for _, browser := range u.Browsers {

			if strings.Contains(browser, "Android") {
				isHaveAndroid = true
				if _, ok := seenBrowsers[browser]; ok {
					seenBrowsers[browser]++
				} else {
					seenBrowsers[browser] = 1
					uniqueBrowsers++
				}
			}

			if strings.Contains(browser, "MSIE") {
				isHaveMSIE = true
				if _, ok := seenBrowsers[browser]; ok {
					seenBrowsers[browser]++
				} else {
					seenBrowsers[browser] = 1
					uniqueBrowsers++
				}
			}
		}

		if !(isHaveAndroid && isHaveMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := strings.Replace(u.Email, "@", " [at] ", 1)
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i-1, u.Name, email)

	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
