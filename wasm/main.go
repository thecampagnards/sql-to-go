package main

import (
	"fmt"
	"regexp"
	"syscall/js"

	"github.com/thecampagnards/sql-to-go/run"
)

func main() {
	js.Global().Set("parse", parse())
	<-make(chan bool)
}

func parse() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid parameters"
		}
		input := args[0].String()
		files, err := run.Run(run.RunParams{
			GenerateFuncs: false,
			ModelType:     "bun",
			PackageName:   "db",
			SQL:           input,
		})
		if err != nil {
			return fmt.Sprintf("Unable to convert: %s\n", err)
		}

		result := ""
		for _, content := range files {
			if result != "" {
				// remove package part
				re := regexp.MustCompile(`(?m)((.|\n)*)\)\n`)
				content = re.ReplaceAllString(content, "")
			}
			result += fmt.Sprintf("%s", content)
		}
		return result
	})
}
