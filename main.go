package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Masterminds/sprig"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/parser/test_driver"
	"github.com/pingcap/tidb/types"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"github.com/rs/zerolog/log"

	"github.com/thecampagnards/sql-to-go/models"
)

func main() {
	var sqlFile = flag.String("sql-file", "example.sql", "SQL file to parse")
	var modelType = flag.String("model-type", "bun", "Model output type: bun, ...")
	var outputFolder = flag.String("output-folder", "out", "Output folder")

	flag.Parse()

	log.Info().
		Str("sql-file", *sqlFile).
		Str("model-type", *modelType).
		Str("output-folder", *outputFolder).
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
		err = tmpl.Execute(file, table)
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

type Result struct {
	Tables []Table
}

type Table struct {
	Name    string
	Columns []Column
}

type Column struct {
	Name    string
	Type    string
	Null    bool
	Options []string
}

func parse(sql string) (*ast.StmtNode, error) {
	p := parser.New()

	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	return &stmtNodes[0], nil
}

func (v *Result) Enter(in ast.Node) (ast.Node, bool) {
	if t, ok := in.(*ast.CreateTableStmt); ok {
		table := Table{Name: t.Table.Name.O}
		for _, col := range t.Cols {
			column := Column{
				Name: col.Name.Name.O,
				Null: true,
			}
			switch col.Tp.EvalType() {
			case types.ETInt:
				column.Type = "int64"
			case types.ETReal, types.ETDecimal:
				column.Type = "float64"
			case types.ETDatetime, types.ETTimestamp:
				column.Type = "time.Time"
			case types.ETDuration:
				column.Type = "time.Duration"
			case types.ETJson:
				column.Type = "interface{}"
			case types.ETString:
				column.Type = "string"
			}
			for _, option := range col.Options {
				if option.Tp == ast.ColumnOptionNotNull {
					column.Null = false
					column.Options = append(column.Options, "not null")
				}
				if option.Tp == ast.ColumnOptionUniqKey {
					column.Null = false
					column.Options = append(column.Options, "unique")
				}
			}
			table.Columns = append(table.Columns, column)
		}
		v.Tables = append(v.Tables, table)
	}
	return in, false
}

func (v *Result) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func extract(rootNode *ast.StmtNode) *Result {
	v := &Result{}
	(*rootNode).Accept(v)
	return v
}
