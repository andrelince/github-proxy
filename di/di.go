package di

import (
	"go.uber.org/dig"
)

func NewDI() (*dig.Container, error) {
	c := dig.New()

	if err := buildConfig(c); err != nil {
		return nil, err
	}

	return c, nil
}
