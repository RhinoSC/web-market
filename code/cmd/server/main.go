package main

import (
	"fmt"

	"github.com/rhinosc/web-market/code/internal/application"
)

func main() {

	// server := application.NewServerChi(":8080")

	// if err := server.Run(); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	app := application.NewDefaultHTTP(":8080")
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
