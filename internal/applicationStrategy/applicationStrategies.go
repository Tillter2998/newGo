package applicationStrategy

import "fmt"

type ApplicationStrategy interface {
	Execute(name string, directory string) error
}

type Empty struct{}

func (e *Empty) Execute(name string, directory string) error {
	fmt.Println("IT WORKED YAY")
	return nil
}

func init() {
	registry := GetRegistry()
	registry.RegisterStrategy("empty", &Empty{})
}
