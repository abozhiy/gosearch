package main

import (
	"os"

	"gosearch/internal/app"
	"gosearch/internal/global"
	"gosearch/internal/use_cases/scan"
)

func main() {
	// url := os.Args
	// runtime.Breakpoint()
	// fmt.Println("Args:", args)

	global.Container = app.InitContainer()
	scan.Perform(os.Getenv("TEST_SEARCH_URL"))
}
