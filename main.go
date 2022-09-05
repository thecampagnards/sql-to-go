package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/thecampagnards/sql-to-go/run"
)

func main() {
	var generateFuncs = flag.Bool("generate-funcs", true, "Generate functions to request models")
	var modelType = flag.String("model-type", "bun", "Model output type: bun, ...")
	var outputFolder = flag.String("output-folder", "out", "Output folder")
	var packageName = flag.String("package-name", "db", "Package name of the generated files")
	flag.Parse()
	sqlFiles := flag.Args()

	sql := ""
	for _, sqlFile := range sqlFiles {
		log.Info().Msg("read sql files")
		tmp, err := os.ReadFile(sqlFile)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to read sql file")
		}
		sql += string(strings.Split(string(tmp), "-- migrate:down")[0])
	}

	files, err := run.Run(run.RunParams{
		GenerateFuncs: *generateFuncs,
		ModelType:     *modelType,
		PackageName:   *packageName,
		SQL:           sql,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run")
	}

	log.Info().Msg("create output folder")
	if err := os.MkdirAll(*outputFolder, os.ModePerm); err != nil {
		log.Fatal().Err(err).Msg("failed to create output folder")
	}

	for file, content := range files {
		log.Info().Str("table", file).Msg("create file for table")
		f, err := os.OpenFile(filepath.Join(*outputFolder, fmt.Sprintf("%s.go", file)), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create file for table")
		}
		defer f.Close()

		_, err = f.WriteString(content)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create file for table")
		}
	}
}
