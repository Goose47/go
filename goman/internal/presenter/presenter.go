package presenter

import (
	"Goose47/goman/internal/client"
	"Goose47/goman/internal/types"
	"fmt"
)

func GetModuleInfo(module types.Module) string {
	return fmt.Sprintf(
		"Package %s\nDocumentation: %s%s\n%s\n---",
		module.Name,
		client.BASE_URL,
		module.Uri,
		module.Description,
	)
}

func GetModuleNotFound(name string) string {
	return fmt.Sprintf("Module %s is not found in standard library. :(", name)
}
func GetItemInfo(item types.Item) string {
	message := fmt.Sprintf("%s %s \n%s", item.Type, item.Name, item.Signature)
	if item.Description != "" || item.Example != "" {
		message += fmt.Sprintf("\n---\n%s\n%s\n---", item.Description, item.Example)
	}
	return message
}
