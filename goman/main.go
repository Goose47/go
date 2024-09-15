package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"strings"
)

const BASE_URL = "https://pkg.go.dev"

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: package.FunctionName package/subpackage.TypeName")
		return
	}
	modules, err := GetStandardModules()
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
	fmt.Println(fmt.Sprintf("Documentation: %s%s", BASE_URL, module.Uri))
	fmt.Println(module.Description)
	fmt.Println("---")

	if !withItem {
		return
	}

	item, err := GetItem(module, itemName)
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

type Module struct {
	Name        string
	Uri         string
	Description string
}

type Item struct {
	Type        string
	Name        string
	Signature   string
	Description string
	Example     string
}

func GetStandardModules() (map[string]Module, error) {
	res, err := http.Get(fmt.Sprintf("%s/std", BASE_URL))
	if err != nil {
		return nil, err
	}
	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	var nextModuleName string
	modules := make(map[string]Module)

	var parse func(*html.Node)
	parse = func(n *html.Node) {
		found := false
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" && (a.Val == "UnitDirectories-pathCell" || a.Val == "UnitDirectories-subdirectory") {
					found = true
					break
				}
			}
		}

		if !found {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				parse(c)
			}
			return
		}

		var nextModule Module

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}
			if c.Data == "span" && c.FirstChild != nil {
				for gc := c.FirstChild; gc != nil; gc = gc.NextSibling {
					if gc.Data == "a" {
						nextModule.Uri = gc.Attr[0].Val
						nextModule.Name = fmt.Sprintf("%s/%s", nextModuleName, gc.FirstChild.Data)
					}
				}
			}
			if c.Data == "div" {
				descriptionBlock := false
				for _, a := range c.Attr {
					if a.Key == "class" && a.Val == "UnitDirectories-mobileSynopsis" {
						descriptionBlock = true
						break
					}
				}
				if descriptionBlock && c.FirstChild != nil {
					nextModule.Description = c.FirstChild.Data
					continue
				}
				for gc := c.FirstChild; gc != nil; gc = gc.NextSibling {
					switch gc.Data {
					case "a":
						nextModule.Uri = gc.Attr[0].Val
						nextModule.Name = gc.FirstChild.Data
						nextModuleName = gc.FirstChild.Data
					case "span":
						nextModuleName = gc.FirstChild.Data
					}
				}
			}
		}
		modules[nextModule.Name] = nextModule
	}
	parse(doc)
	return modules, nil
}

func GetItem(module Module, name string) (*Item, error) {
	res, err := http.Get(fmt.Sprintf("%s/%s", BASE_URL, module.Uri))
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	var item Item
	item.Name = name

	var parse func(*html.Node)
	parse = func(n *html.Node) {
		found := false
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == name {
					found = true
					break
				}
			}
		}

		if !found {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				parse(c)
			}
			return
		}

		for _, a := range n.Attr {
			if a.Key == "data-kind" {
				item.Type = a.Val
				break
			}
		}

		var signature string
		for c := n.NextSibling; c != nil; c = c.NextSibling {
			switch c.Data {
			case "div":
				for gc := c.FirstChild; gc != nil; gc = gc.NextSibling {
					if gc.Data != "pre" {
						continue
					}
					for gc = gc.FirstChild; gc != nil; gc = gc.NextSibling {
						if gc.Data == "a" {
							signature += gc.FirstChild.Data
						} else {
							signature += gc.Data
						}
					}
					break
				}
				item.Signature = signature
			case "p":
				item.Description = c.FirstChild.Data
			case "details":
				for gc := c.FirstChild; gc != nil; gc = gc.NextSibling {
					if gc.Data != "div" {
						continue
					}
					for gc := gc.FirstChild; gc != nil; gc = gc.NextSibling {
						if gc.Data != "pre" {
							continue
						}
						item.Example = gc.FirstChild.Data
						break
					}
					break
				}
			}
		}
	}
	parse(doc)

	return &item, nil
}
