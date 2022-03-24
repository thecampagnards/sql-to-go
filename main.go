package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"

	"github.com/thecampagnards/sql-to-go/models"
)

func main() {
	var generateFuncs = flag.Bool("generate-funcs", true, "Generate functions to request models")
	var modelType = flag.String("model-type", "bun", "Model output type: bun, ...")
	var outputFolder = flag.String("output-folder", "out", "Output folder")
	var packageName = flag.String("package-name", "db", "Package name of the generated files")
	flag.Parse()
	sqlFiles := flag.Args()

	log.Info().
		Bool("generate-func", *generateFuncs).
		Str("model-type", *modelType).
		Str("output-folder", *outputFolder).
		Str("package-name", *packageName).
		Strs("sql-files", sqlFiles).
		Msg("run with flags")

	log.Info().Msg("create output folder")
	if err := os.MkdirAll(*outputFolder, os.ModePerm); err != nil {
		log.Fatal().Err(err).Msg("failed to create output folder")
	}

	sql := ""
	for _, sqlFile := range sqlFiles {
		log.Info().Msg("read sql files")
		tmp, err := ioutil.ReadFile(sqlFile)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to read sql file")
		}
		sql += string(strings.Split(string(tmp), "-- migrate:down")[0])
	}

	log.Info().Msg("parse sql file")
	astNode, err := parse(sql)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse sql file")
	}

	result := extract(astNode)

	log.Info().Msg("create db models")
	for _, table := range result.Tables {
		log.Info().Str("table", table.Name).Msg("create go template for table")
		tmpl, err := template.New("template").
			Funcs(sprig.TxtFuncMap()).
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
			GenerateFuncs bool
			PackageName   string
			Result        *Result
			Table
		}{
			GenerateFuncs: *generateFuncs,
			PackageName:   *packageName,
			Result:        result,
			Table:         table,
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
