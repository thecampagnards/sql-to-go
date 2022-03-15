package main

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/parser/test_driver"
	"github.com/pingcap/tidb/types"
	_ "github.com/pingcap/tidb/types/parser_driver"
)

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
