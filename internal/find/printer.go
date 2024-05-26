package find

import (
	"fmt"
	"sync"

	"github.com/jibaru/gofind/internal/utils"
)

type Printer struct {
	results <-chan Result
	errors  <-chan error
}

func NewPrinter(b *Broadcast) Printer {
	return Printer{
		results: b.Results(),
		errors:  b.Errors(),
	}
}

func (fp *Printer) StartPrinting(wg *sync.WaitGroup) {
	go func() {
		var err error
		var result Result
		resultsOpen, errorsOpen := true, true

		for resultsOpen || errorsOpen {
			select {
			case err, errorsOpen = <-fp.errors:
				if errorsOpen {
					fmt.Println(utils.Red, "[❌]", err, utils.Reset)
				}
			case result, resultsOpen = <-fp.results:
				if resultsOpen {
					fmt.Println(utils.Green, "[✔️]", utils.Yellow, result, utils.Reset)
				}
			}
		}

		wg.Done()
	}()
}
