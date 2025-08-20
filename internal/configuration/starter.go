package configuration

import "log"

type Starter interface {
	Run() error
	Close() error
}

type NewObject func(path string) (Starter, error)

func Runner(path string, fn ...NewObject) error {
	starters := make([]Starter, len(fn))
	for index, item := range fn {
		object, err := item(path)
		if err != nil {
			return err
		}
		if err := object.Run(); err != nil {
			return err
		}
		starters[index] = object
	}

	defer func() {
		for _, item := range starters {
			if err := item.Close(); err != nil {
				log.Panic(err)
			}
		}
	}()
	return nil
}
