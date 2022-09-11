package models

type Provider interface {
	Name() string
	Units() ([]Unit, error)
}
