package generate

import (
	"context"
	"fmt"
	"io"
)

type (
	Generator interface {
		Generate(ctx context.Context) ([]string, error)
	}

	CombinationStream interface {
		Next(ctx context.Context) (map[string]string, error)
	}
)

type generatorImp struct {
	combinations CombinationStream
	template     template
}

func NewGenerator(file io.Reader, combinations CombinationStream) (Generator, error) {
	compiledTemplate, err := buildTemplate(file)
	if err != nil {
		return nil, fmt.Errorf("failed to build template: %w", err)
	}

	return &generatorImp{
		template:     compiledTemplate,
		combinations: combinations,
	}, nil
}

// Genreate uses specified template and combination stream to build/return the combinations of
// documents built together
func (g *generatorImp) Generate(ctx context.Context) ([]string, error) {
	// Wait for the context to end or the combinations to be done
	for {
		select {
		case <-ctx.Done():
			return []string{}, ctx.Err()
		default:
			combination, err := g.combinations.Next(ctx)
			if err != nil {
				return []string{}, err
			}

			if combination == nil {
				return g.template.processedDocuments, nil
			}

			g.template.with(combination)
		}
	}
}
