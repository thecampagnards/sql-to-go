package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"

	"github.com/thecampagnards/sql-to-go/models"
)

func main() {
	var modelType = flag.String("model-type", "bun", "Model output type: bun, ...")
	var outputFolder = flag.String("output-folder", "out", "Output folder")
	var sqlFile = flag.String("sql-file", "example.sql", "SQL file to parse")

	flag.Parse()

	log.Info().
		Str("model-type", *modelType).
		Str("output-folder", *outputFolder).
		Str("sql-file", *sqlFile).
		Msg("run with flags")

	log.Info().Msg("create output folder")
	if err := os.MkdirAll(*outputFolder, os.ModePerm); err != nil {
		log.Fatal().Err(err).Msg("failed to create output folder")
	}

	log.Info().Msg("read sql file")
	sql, err := ioutil.ReadFile(*sqlFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read sql file")
	}

	log.Info().Msg("parse sql file")
	astNode, err := parse(string(sql))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse sql file")
	}

	result := extract(astNode)

	log.Info().Msg("create db models")
	for _, table := range result.Tables {
		log.Info().Str("table", table.Name).Msg("create go template for table")
		tmpl, err := template.New("template").
			Funcs(sprig.FuncMap()).
			Parse(models.Models[models.ModelType(*modelType)])
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create go template for table")
		}

		log.Info().Str("table", table.Name).Msg("create file for table")
		newpath := filepath.Join(*outputFolder, fmt.Sprintf("%s.go", table.Name))
		file, err := os.OpenFile(newpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create file for table")
		}
		defer file.Close()

		log.Info().Str("table", table.Name).Msg("render template for table")
		err = tmpl.Execute(file, struct {
			Result *Result
			Table
		}{
			Result: result,
			Table:  table,
		})
		file.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to render template for table") // nolint: gocritic
		}
	}

	log.Info().Msg("format db models")
	if err := exec.Command("gofmt", "-w", *outputFolder).Run(); err != nil {
		log.Fatal().Err(err).Msg("failed to format db models")
	}

	log.Info().Msg("done")
}
