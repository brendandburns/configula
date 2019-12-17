package configula

type Processor interface {
	Process(sections []Section) error
}