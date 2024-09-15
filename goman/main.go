package main

import (
	"Goose47/goman/internal/parser"
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
		fmt.Println(fmt.Sprintf("Module %s is not found in standard library. :(", moduleName))
		return
	}
	fmt.Println(fmt.Sprintf("Package %s", module.Name))
	fmt.Println(fmt.Sprintf("Documentation: %s%s", parser.BASE_URL, module.Uri))
	fmt.Println(module.Description)
	fmt.Println("---")

	if !withItem {
		return
	}

	item, err := parser.GetItem(module, itemName)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fmt.Sprintf("%s %s", item.Type, item.Name))
	fmt.Println(item.Signature)
	fmt.Println("---")
	fmt.Println(item.Description)
	fmt.Println(item.Example)
	fmt.Println("---")
}
