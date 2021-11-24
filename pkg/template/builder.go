package template

import (
	"context"
	"fmt"
	"io"
)

type (
	Builder interface {
		Build(ctx context.Context) ([]string, error)
	}

	CombinationStream interface {
		Next(ctx context.Context) (map[string]string, error)
	}
)

type builder struct {
	combinations CombinationStream
	template     template
}

func NewBuilder(file io.Reader, combinations CombinationStream) (Builder, error) {
	compiledTemplate, err := newTemplate(file)
	if err != nil {
		return nil, fmt.Errorf("failed to build template: %w", err)
	}

	return &builder{
		template:     compiledTemplate,
		combinations: combinations,
	}, nil
}

// Build uses the current builder's template and combination stream to
// construct the combinations of manifests built together
func (g *builder) Build(ctx context.Context) ([]string, error) {
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
				return g.template.processedManifests, nil
			}

			g.template.with(combination)
		}
	}
}
