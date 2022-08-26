package run

import (
	"bytes"
	"fmt"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"

	"github.com/thecampagnards/sql-to-go/models"
)

type RunParams struct {
	GenerateFuncs bool
	ModelType     string
	PackageName   string
	SQL           string
}

// Run return a map with filename as index and its content
func Run(params RunParams) (map[string]string, error) {
	log.Info().
		Bool("generate-func", params.GenerateFuncs).
		Str("model-type", params.ModelType).
		Str("package-name", params.PackageName).
		Msg("run with options")

	log.Info().Msg("parse sql file")
	astNode, err := parse(params.SQL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sql file: %w", err)
	}

	result := extract(astNode)
	output := make(map[string]string)

	log.Info().Msg("create db models")
	for _, table := range result.Tables {
		log.Info().Str("table", table.Name).Msg("create go template for table")
		tmpl, err := template.New("template").
			Funcs(sprig.TxtFuncMap()).
			Parse(models.Models[models.ModelType(params.ModelType)])
		if err != nil {
			return nil, fmt.Errorf("failed to create go template for table: %w", err)
		}

		log.Info().Str("table", table.Name).Msg("render template for table")
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, struct {
			GenerateFuncs bool
			PackageName   string
			Result        *Result
			Table
		}{
			GenerateFuncs: params.GenerateFuncs,
			PackageName:   params.PackageName,
			Result:        result,
			Table:         table,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to render template for table: %w", err)
		}
		output[table.Name] = buf.String()
	}

	log.Info().Msg("done")
	return output, nil
}
