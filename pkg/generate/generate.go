package generate

import (
	"context"
)

type (
	Generator interface {
		Evaluate(ctx context.Context) ([]string, error)
	}

	CombinationStream interface {
		Next(ctx context.Context) (map[string]string, error)
	}
)

type generatorImp struct {
	combinations CombinationStream
	template     template
}

func NewGenerator(file string, combinations CombinationStream) Generator {
	return &generatorImp{
		template:     newTemplate(file),
		combinations: combinations,
	}
}

// Evaluate uses specified template and combination stream to build/return the combinations of
// documents built together
func (g *generatorImp) Evaluate(ctx context.Context) ([]string, error) {
	var result template

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
				return result.documents, nil
			}

			result = g.template.with(combination, result)
		}
	}
}
