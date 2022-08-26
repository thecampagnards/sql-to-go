package main

import (
	"fmt"
	"regexp"
	"syscall/js"

	"mvdan.cc/gofumpt/format"

	"github.com/thecampagnards/sql-to-go/run"
)

func main() {
	js.Global().Set("parse", parse())
	<-make(chan bool)
}

func parse() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		input := args[0].String()
		files, err := run.Run(run.RunParams{
			GenerateFuncs: false,
			ModelType:     "bun",
			PackageName:   "db",
			SQL:           input,
		})
		if err != nil {
			fmt.Printf("unable to parse SQL: %s\n", err)
			return err.Error()
		}

		result := ""
		for _, content := range files {
			if result != "" {
				re := regexp.MustCompile(`(?m)((.|\n)*)\)`)
				content = re.ReplaceAllString(content, "")
			}
			result += fmt.Sprintf("%s\n", content)
		}
		bytes, err := format.Source([]byte(result), format.Options{})
		if err != nil {
			fmt.Printf("unable to format code: %s\n", err)
			return err.Error()
		}
		return string(bytes)
	})
}
