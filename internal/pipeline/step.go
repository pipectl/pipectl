package pipeline

type Step interface {
	Name() string
	Run(input []byte) ([]byte, error)
}
