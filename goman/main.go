package main

import (
	"Goose47/goman/internal/parser"
	"Goose47/goman/internal/presenter"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: package.FunctionName package/subpackage.TypeName")
		return
	}
	modules, err := parser.GetStandardModules()
	if err != nil {
		log.Fatal(err)
	}

	next := os.Args[1]
	moduleName, itemName, withItem := strings.Cut(next, ".")

	module, ok := modules[moduleName]
	if !ok {
		fmt.Println(presenter.GetModuleNotFound(moduleName))
		return
	}
	fmt.Println(presenter.GetModuleInfo(module))

	if !withItem {
		return
	}

	item, err := parser.GetItem(module, itemName)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(presenter.GetItemInfo(*item))
}
