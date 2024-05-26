package find

import "sync"

type Broadcast struct {
	finder  *Finder
	results []chan Result
	errors  []chan error
}

func NewBroadcast(finder *Finder) *Broadcast {
	return &Broadcast{
		finder:  finder,
		results: make([]chan Result, 0),
		errors:  make([]chan error, 0),
	}
}

func (b *Broadcast) Results() <-chan Result {
	newResults := make(chan Result)
	b.results = append(b.results, newResults)
	return newResults
}

func (b *Broadcast) Errors() <-chan error {
	newErrors := make(chan error)
	b.errors = append(b.errors, newErrors)
	return newErrors
}

func (b *Broadcast) Broadcast(wg *sync.WaitGroup) {
	go func() {
		results := b.finder.Results()
		errors := b.finder.Errors()

		var err error
		var result Result
		resultsOpen, errorsOpen := true, true
		for resultsOpen || errorsOpen {
			select {
			case err, errorsOpen = <-errors:
				if errorsOpen {
					for _, ch := range b.errors {
						ch <- err
					}
				}
			case result, resultsOpen = <-results:
				if resultsOpen {
					for _, ch := range b.results {
						ch <- result
					}
				}
			}
		}

		for _, ch := range b.results {
			close(ch)
		}
		for _, ch := range b.errors {
			close(ch)
		}

		wg.Done()
	}()
}
