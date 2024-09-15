package parser

import (
	"Goose47/goman/internal/types"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
)

const BASE_URL = "https://pkg.go.dev"

func GetStandardModules() (map[string]types.Module, error) {
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
	modules := make(map[string]types.Module)

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

		var nextModule types.Module

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

func GetItem(module types.Module, name string) (*types.Item, error) {
	res, err := http.Get(fmt.Sprintf("%s/%s", BASE_URL, module.Uri))
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	var item types.Item
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
