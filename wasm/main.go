package main

import (
	"fmt"
	"syscall/js"
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
		fmt.Printf("input: %s\n", input)
		/*	files, err := run.Run(run.RunParams{
				GenerateFuncs: false,
				ModelType:     "bun",
				PackageName:   "db",
				SQL:           input,
			})
			if err != nil {
				fmt.Printf("unable to parse SQL: %s\n", err)
				return err.Error()
			}
		*/
		files := map[string]string{"test": "test"}

		result := ""
		for file, content := range files {
			result += fmt.Sprintf("// %s\n%s\n", file, content)
		}
		fmt.Printf("result: %s\n", result)
		return result
	})
}
