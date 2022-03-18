package main

import (
	"strings"

	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/types"
	_ "github.com/pingcap/tidb/types/parser_driver"
)

type Result struct {
	Tables map[string]Table
}

type References []string

func (rs References) Contain(table string) bool {
	for _, r := range rs {
		if r == table {
			return true
		}
	}
	return false
}

type Table struct {
	Columns    map[string]Column
	Name       string
	References References
}

type Column struct {
	NotNull bool
	Options []string
	Type    string
}

func parse(sql string) ([]ast.StmtNode, error) {
	p := parser.New()

	// replace postgres serials
	sql = strings.ReplaceAll(sql, "SERIAL", "INT NOT NULL AUTO_INCREMENT")

	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	return stmtNodes, nil
}

func (v *Result) Enter(in ast.Node) (ast.Node, bool) {
	if t, ok := in.(*ast.CreateTableStmt); ok {
		table := Table{
			Columns: make(map[string]Column),
			Name:    t.Table.Name.O,
		}
		for _, col := range t.Cols {
			column := Column{}
			switch col.Tp.EvalType() {
			case types.ETDatetime, types.ETTimestamp:
				column.Type = "time.Time"
			case types.ETDecimal, types.ETReal:
				column.Type = "float64"
			case types.ETDuration:
				column.Type = "time.Duration"
			case types.ETInt:
				column.Type = "int64"
			case types.ETJson:
				column.Type = "interface{}"
			case types.ETString:
				column.Type = "string"
			}
			for _, option := range col.Options {
				switch option.Tp {
				case ast.ColumnOptionAutoIncrement:
					column.Options = append(column.Options, "autoincrement")
				case ast.ColumnOptionNotNull:
					column.NotNull = true
					column.Options = append(column.Options, "notnull")
				case ast.ColumnOptionPrimaryKey:
					column.Options = append(column.Options, "pk")
				case ast.ColumnOptionReference:
					table.References = append(table.References, option.Refer.Table.Name.O)
				case ast.ColumnOptionUniqKey:
					column.Options = append(column.Options, "unique")
				}
			}
			table.Columns[col.Name.Name.O] = column
		}
		for _, cons := range t.Constraints {
			// nolint: gocritic
			switch cons.Tp {
			case ast.ConstraintPrimaryKey:
				for _, key := range cons.Keys {
					column := table.Columns[key.Column.Name.O]
					column.Options = append(column.Options, "pk")
					table.Columns[key.Column.Name.O] = column
				}
			}
		}
		v.Tables[table.Name] = table
	}
	return in, false
}

func (v *Result) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func extract(rootNode []ast.StmtNode) *Result {
	v := &Result{
		Tables: make(map[string]Table),
	}
	for _, node := range rootNode {
		node.Accept(v)
	}
	return v
}
