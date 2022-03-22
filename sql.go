package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/types"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"github.com/rs/zerolog/log"
)

type Result struct {
	Tables map[string]Table
}

type Table struct {
	Columns map[string]Column
	Name    string
}

func (t Table) GetReferenceColumn(reference string) string {
	for n, c := range t.Columns {
		if c.Reference == reference {
			return n
		}
	}
	return ""
}

type Column struct {
	NotNull   bool
	Options   []string
	Type      string
	Reference string
}

func parse(sql string) ([]ast.StmtNode, error) {
	p := parser.New()

	// replace to handle postgres
	sql = strings.ReplaceAll(sql, "SERIAL", "INT NOT NULL AUTO_INCREMENT")
	r := regexp.MustCompile(`ALTER COLUMN (.+) TYPE`)
	sql = r.ReplaceAllString(sql, "MODIFY $1")

	stmtNodes, warns, err := p.ParseSQL(sql)
	if err != nil {
		return nil, err
	}
	if warns != nil {
		log.Warn().Errs("warns", warns).Msg("warning on sql parsing")
	}

	return stmtNodes, nil
}

func getColumns(cols []*ast.ColumnDef) map[string]Column {
	columns := make(map[string]Column)
	for _, col := range cols {
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
				column.Reference = option.Refer.Table.Name.O
			case ast.ColumnOptionUniqKey:
				column.Options = append(column.Options, "unique")
			}
		}
		columns[col.Name.Name.O] = column
	}
	return columns
}

func (v *Result) Enter(in ast.Node) (ast.Node, bool) {
	switch t := in.(type) {
	case *ast.CreateTableStmt:
		table := Table{
			Columns: make(map[string]Column),
			Name:    t.Table.Name.O,
		}
		table.Columns = getColumns(t.Cols)
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

	case *ast.AlterTableStmt:
		for _, s := range t.Specs {
			switch s.Tp {
			case ast.AlterTableAddColumns, ast.AlterTableModifyColumn, ast.AlterTableChangeColumn:
				columns := getColumns(s.NewColumns)
				table := v.Tables[t.Table.Name.O]
				for name, value := range columns {
					table.Columns[name] = value
				}
				v.Tables[t.Table.Name.O] = table
			case ast.AlterTableRenameColumn:
				v.Tables[t.Table.Name.O].Columns[s.NewColumnName.Name.O] = v.Tables[t.Table.Name.O].Columns[s.OldColumnName.Name.O]
				delete(v.Tables[t.Table.Name.O].Columns, s.OldColumnName.Name.O)
			case ast.AlterTableDropColumn:
				delete(v.Tables[t.Table.Name.O].Columns, s.OldColumnName.Name.O)
			case ast.AlterTableAddConstraint:
				column := v.Tables[t.Table.Name.O].Columns[s.FromKey.O]
				column.Reference = s.Constraint.Refer.Table.Name.O
				v.Tables[t.Table.Name.O].Columns[s.FromKey.O] = column
			case ast.AlterTableDropForeignKey:
				// postgres
				r := regexp.MustCompile(fmt.Sprintf(`%s_(.+)_fkey`, t.Table.Name.O))
				if !r.MatchString(s.Name) {
					// mysql
					r = regexp.MustCompile(fmt.Sprintf(`%s_ibfk_(.+)`, t.Table.Name.O))
				}
				results := r.FindStringSubmatch(s.Name)
				if len(results) > 1 {
					// delete by regerated id
					if i, err := strconv.Atoi(results[1]); err == nil {
						index := 1
						for key, column := range v.Tables[t.Table.Name.O].Columns {
							if column.Reference != "" {
								if i == index {
									column.Reference = ""
									v.Tables[t.Table.Name.O].Columns[key] = column
									break
								}
								index++
							}
						}
						// delete by column name
					} else {
						if column, ok := v.Tables[t.Table.Name.O].Columns[results[1]]; ok {
							column.Reference = ""
							v.Tables[t.Table.Name.O].Columns[results[1]] = column
						}
					}
				}
			}
		}
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
