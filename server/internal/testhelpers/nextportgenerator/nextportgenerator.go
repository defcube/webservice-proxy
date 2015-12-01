package nextportgenerator

import "fmt"

var DefaultNextPortGenerator = NewNextPortGenerator()

type nextPortGenerator struct {
	nextCh chan int
}

func (np *nextPortGenerator) Next() int {
	return <-np.nextCh
}

func NextAsColonString() string {
	return DefaultNextPortGenerator.NextAsColonString()
}

func (np *nextPortGenerator) NextAsColonString() string {
	return fmt.Sprintf(":%v", np.Next())
}

func (np *nextPortGenerator) start() {
	np.nextCh = make(chan int)
	go func() {
		i := 8050
		for {
			np.nextCh <- i
			i++
		}
	}()
}

func NewNextPortGenerator() *nextPortGenerator {
	np := &nextPortGenerator{}
	np.start()
	return np
}
