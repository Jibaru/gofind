package find

import (
	"encoding/json"
	"io"
	"os"
	"sync"
)

type JsonWriter struct {
	filename string
	result   <-chan Result
	errors   <-chan error
}

func NewJsonWriter(filename string, b *Broadcast) *JsonWriter {
	return &JsonWriter{
		filename: filename,
		result:   b.Results(),
		errors:   b.Errors(),
	}
}

func (jw *JsonWriter) StartWriting(wg *sync.WaitGroup) error {
	file, err := os.OpenFile(jw.filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	go func() {
		defer file.Close()

		var err error
		var result Result
		resultsOpen, errorsOpen := true, true

		file.WriteString("[" + "\n")

		for resultsOpen || errorsOpen {
			select {
			case err, errorsOpen = <-jw.errors:
				if errorsOpen {
					bytes, err2 := json.Marshal(err)
					if err2 != nil {
						panic(err2)
					}

					file.WriteString(string(bytes) + ",\n")
				}
			case result, resultsOpen = <-jw.result:
				if resultsOpen {
					bytes, err2 := json.Marshal(result)
					if err2 != nil {
						panic(err2)
					}

					file.WriteString(string(bytes) + ",\n")
				}
			}
		}

		// Remove last comma
		_, _ = file.Seek(-2, io.SeekEnd)
		_, _ = file.Write([]byte(" "))

		file.WriteString("]" + "\n")
		wg.Done()
	}()

	return nil
}
