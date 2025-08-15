package configuration

type Starter interface {
	Run() error
}

type NewObject func() (Starter, error)

func Runner(fn ...NewObject) error {
	for _, item := range fn {
		object, err := item()
		if err != nil {
			return err
		}
		if err := object.Run(); err != nil {
			return err
		}
	}
	return nil
}
