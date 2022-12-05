package models

// nolint: lll
const bun = `
// Code generated by sql-to-go; DO NOT EDIT.

package {{ $.PackageName }}

import (
{{- if $.GenerateFuncs }}
	"context"
	"errors"
{{- end }}
{{- range $column := $.Columns }}
	{{- if hasPrefix "time" $column.Type }}
	"time"
	{{- end }}
{{- end }}

	"github.com/uptrace/bun"
)

type {{ camelcase $.Name }} struct {
	bun.BaseModel ` + "`" + `bun:"table:{{ $.Name }}"` + "`" + `

	{{- range $key, $column := $.Columns }}
		{{ camelcase $key | replace "Id" "ID" }} {{ if not $column.NotNull }}*{{ end }}{{ $column.Type }}  ` + "`" + `bun:"{{ $key }},{{ $column.Options | uniq | join "," }}"` + "`" + `
	{{- end }}

	{{- range $key, $column := $.Columns }}
		{{- if $column.Reference }}
			{{ camelcase $key | trimSuffix "Id" }} *{{ camelcase $column.Reference }} ` + "`" + `bun:"rel:belongs-to"` + "`" + `
		{{- end }}
	{{- end }}
	{{- range $table := $.Result.Tables }}
		{{- range $key, $column := $table.Columns }}
			{{- if eq $column.Reference $.Name }}
				{{- if $.Table.GetReferenceColumn $key }}
					{{ camelcase $table.Name }} *{{ camelcase $table.Name }} ` + "`" + `bun:"rel:has-one,join:id={{ $key }}"` + "`" + `
				{{- else }}
					{{ trimSuffix $.Name (trimSuffix "_id" $key) | camelcase }}{{ $table.Name | camelcase }}s []*{{ camelcase $table.Name }} ` + "`" + `bun:"rel:has-many,join:id={{ $key }}"` + "`" + `
				{{- end }}
			{{- end }}
		{{- end }}
	{{- end }}
}

{{- if $.GenerateFuncs }}
// Select{{ camelcase $.Name }}
func Select{{ camelcase $.Name }}(ctx context.Context, db *bun.DB, {{ camelcase $.Name | untitle }} {{ camelcase $.Name }}) (*{{ camelcase $.Name }}, error) {
	err := db.NewSelect().
		Model(&{{ camelcase $.Name | untitle }}).
		WherePK().
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &{{ camelcase $.Name | untitle }}, nil
}

// Select{{ camelcase $.Name }}sParams
type Select{{ camelcase $.Name }}sParams struct {
	Page *int
	PerPage *int
}

// Select{{ camelcase $.Name }}s
func Select{{ camelcase $.Name }}s(ctx context.Context, db *bun.DB, {{ camelcase $.Name | untitle }}s []{{ camelcase $.Name }}, params *Select{{ camelcase $.Name }}sParams) ([]{{ camelcase $.Name }}, int, error) {
	query := db.NewSelect().Model(&{{ camelcase $.Name | untitle }}s)

	if params != nil {
		if params.Page != nil {
			if *params.Page < 1 {
				return nil, 0, errors.New("page must be > 1")
			}
			query.Offset(*params.Page)
		}
		if params.PerPage != nil {
			if *params.PerPage < 1 {
				return nil, 0, errors.New("per page must be > 1")
			}
			query.Limit(*params.Page)
		}
	}

	count, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return {{ camelcase $.Name | untitle }}s, count, nil
}

// Create{{ camelcase $.Name }}
func Create{{ camelcase $.Name }}(ctx context.Context, db *bun.DB, {{ camelcase $.Name | untitle }} {{ camelcase $.Name }}) (*{{ camelcase $.Name }}, error) {
	_, err := db.NewInsert().
		Model(&{{ camelcase $.Name | untitle }}).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &{{ camelcase $.Name | untitle }}, nil
}

// Update{{ camelcase $.Name }}
func Update{{ camelcase $.Name }}(ctx context.Context, db *bun.DB, {{ camelcase $.Name | untitle }} {{ camelcase $.Name }}) (*{{ camelcase $.Name }}, error) {
	_, err := db.NewUpdate().
		Model(&{{ camelcase $.Name | untitle }}).
		WherePK().
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &{{ camelcase $.Name | untitle }}, nil
}

// Delete{{ camelcase $.Name }}
func Delete{{ camelcase $.Name }}(ctx context.Context, db *bun.DB, {{ camelcase $.Name | untitle }} {{ camelcase $.Name }}) error {
	_, err := db.NewDelete().
		Model(&{{ camelcase $.Name | untitle }}).
		WherePK().
		Exec(ctx)
	return err
}
{{- end }}
`
