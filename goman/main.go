package main

import (
	"Goose47/goman/internal/client"
	"Goose47/goman/internal/parser"
	"Goose47/goman/internal/presenter"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: package.FunctionName package/subpackage.TypeName")
		return
	}

	reader, err := client.CachedFetch("/std")
	if err != nil {
		fmt.Println(err)
		return
	}

	modules, err := parser.GetStandardModules(reader)
	if err != nil {
		log.Fatal(err)
	}

	var messages = make([]string, len(os.Args)-1)
	var wg sync.WaitGroup
	wg.Add(len(os.Args) - 1)

	for i := 0; i < len(os.Args)-1; i++ {
		go func() {
			defer wg.Done()
			next := os.Args[i+1]
			moduleName, itemName, withItem := strings.Cut(next, ".")

			module, ok := modules[moduleName]
			if !ok {
				messages[i] = presenter.GetModuleNotFound(moduleName)
				return
			}

			if len(os.Args) == 2 {
				messages[i] = presenter.GetModuleInfo(module)
			}

			if !withItem {
				return
			}

			reader, err := client.CachedFetch(module.Uri)
			if err != nil {
				fmt.Println(err)
				return
			}

			item, err := parser.GetItem(reader, itemName)
			if err != nil {
				fmt.Println(err)
				return
			}

			messages[i] += "\n" + presenter.GetItemInfo(*item)
		}()
	}

	wg.Wait()

	for _, message := range messages {
		fmt.Println(message)
	}
}
