package parser

import (
	"Goose47/goman/internal/types"
	"strings"
	"testing"
)

func TestGetItem(t *testing.T) {
	var tests = []struct {
		in, name string
		want     types.Item
	}{
		{"", "name", types.Item{Name: "name"}},
		{
			"<div class=\"UnitDirectories-pathCell\"><div><a href=\"/bufio@go1.23.1\">bufio</a></div><div class=\"UnitDirectories-mobileSynopsis\">Package bufio implements buffered I/O. It wraps an io.Reader or io.Writer object, creating another object (Reader or Writer) that also implements the interface but provides buffering and some help for textual I/O.</div></div>",
			"name",
			types.Item{Name: "name"},
		},
		{
			"<div class=\"Documentation-function\"><h4 tabindex=\"-1\" id=\"Println\" data-kind=\"function\" class=\"Documentation-functionHeader\"><span>func <a class=\"Documentation-source\" href=\"https://cs.opensource.google/go/go/+/go1.23.1:src/fmt/print.go;l=313\">Println</a> <a class=\"Documentation-idLink\" href=\"#Println\" aria-label=\"Go to Println\">¶</a></span><span class=\"Documentation-sinceVersion\"></span></h4><div class=\"Documentation-declaration\"><pre>sig</pre></div><p>desc</p><details tabindex=\"-1\" id=\"example-Println\" class=\"Documentation-exampleDetails js-exampleContainer\"><summary class=\"Documentation-exampleDetailsHeader\">Example <a href=\"#example-Println\" aria-label=\"Go to Example\">¶</a></summary><div class=\"Documentation-exampleDetailsBody\"><textarea class=\"Documentation-exampleCode code\" spellcheck=\"false\" style=\"height: 19.625rem;\"></textarea><pre><span class=\"Documentation-exampleOutputLabel\">Output:</span><span class=\"Documentation-exampleOutput\">Kim is 22 years old.</span></pre></div></div></details></div>",
			"Println",
			types.Item{Name: "Println", Type: "function", Signature: "sig", Description: "desc", Example: "span"},
		},
	}

	for _, tc := range tests {
		res, err := GetItem(strings.NewReader(tc.in), tc.name)

		if err != nil {
			t.Fatal(err)
		}

		if res == nil {
			t.Errorf("Expected %s, got nil", tc.want)
			continue
		}

		if *res != tc.want {
			t.Errorf("Expected %s, got %s", tc.want, *res)
		}
	}
}
