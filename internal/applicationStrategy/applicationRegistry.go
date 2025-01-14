package applicationStrategy

import (
	"errors"
	"fmt"
	"strings"
)

type ApplicationRegistry struct {
	strategies map[string]ApplicationStrategy
}

func (ar *ApplicationRegistry) RegisterStrategy(key string, strategy ApplicationStrategy) {
	ar.strategies[strings.ToLower(key)] = strategy
}

func (ar *ApplicationRegistry) GetStrategy(key string) (ApplicationStrategy, error) {
	if strategy, ok := ar.strategies[strings.ToLower(key)]; ok {
		return strategy, nil
	}
	return nil, errors.New(fmt.Sprintf("no such application type: %s", key))
}

var registry = new(ApplicationRegistry)

func init() {
	registry = &ApplicationRegistry{
		strategies: make(map[string]ApplicationStrategy),
	}
}

func GetRegistry() *ApplicationRegistry {
	return registry
}
