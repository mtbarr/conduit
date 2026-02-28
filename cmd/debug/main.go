package main

import (
	"conduit/i18n"
	"fmt"
)

func main() {
	keys := []string{
		"reportbug_command_name",
		"requestfeature_command_name",
		"issues_command_name",
		"requestfeature_modal_title",
		"issues_command_desc",
		"issues_header",
		"issue_format",
	}
	fmt.Println("Language:", i18n.CurrentLanguage())
	for _, k := range keys {
		fmt.Printf("  %s = %q\n", k, i18n.T(k))
	}
}
