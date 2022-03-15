package models

type ModelType string

const (
	BunType ModelType = "bun"
)

var Models = map[ModelType]string{
	BunType: bun,
}
