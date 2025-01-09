package applicationStrategy

import "errors"

type ApplicationContext struct {
	strategy ApplicationStrategy
}

func (ac *ApplicationContext) SetStrategy(strategy ApplicationStrategy) {
	ac.strategy = strategy
}

func (ac *ApplicationContext) CreateApplication(name string, directory string) error {
	if ac.strategy == nil {
		return errors.New("No application strategy set")
	}

	err := ac.strategy.Execute(name, directory)
	if err != nil {
		return err
	}

	return nil
}
