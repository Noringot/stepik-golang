package main

import "fmt"

// <id>0</id>
// <guid>1a6fa827-62f1-45f6-b579-aaead2b47169</guid>
// <first_name>Boyd</first_name>
// <last_name>Wolf</last_name>

func PrintResult(rows []Row) {

	for _, row := range rows {
		fmt.Printf(
			"Id: %d\nGUID: %s\nName: %s\nAge: %s\n-------------\n",
			row.Id, row.GUID, row.Name, row.Age,
		)
	}
}
