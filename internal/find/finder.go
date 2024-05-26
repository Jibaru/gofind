package find

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode/utf8"
)

type Finder struct {
	results        chan Result
	errors         chan error
	closeWaitGroup *sync.WaitGroup
	maxResults     int
	resultsCount   int
	resultsMutex   *sync.Mutex
	stopSearching  chan struct{}
}

func NewFinder(maxResults int) Finder {
	return Finder{
		results:        make(chan Result),
		errors:         make(chan error),
		closeWaitGroup: &sync.WaitGroup{},
		resultsMutex:   &sync.Mutex{},
		maxResults:     maxResults,
		resultsCount:   0,
		stopSearching:  make(chan struct{}),
	}
}

func (f *Finder) StartFinding(dir string, search string) {
	f.closeWaitGroup.Add(1)
	go f.findFile(dir, search)

	go func() {
		f.closeWaitGroup.Wait()
		close(f.results)
		close(f.errors)
	}()
}

func (f *Finder) Results() <-chan Result {
	return f.results
}

func (f *Finder) Errors() <-chan error {
	return f.errors
}

func (f *Finder) findFile(dir string, search string) {
	defer f.closeWaitGroup.Done()

	select {
	case <-f.stopSearching:
		return
	default:
		contents, err := os.ReadDir(dir)
		if err != nil {
			f.errors <- CanNotReadDirErr{Dir: dir, Err: err}
			return
		}

		for _, content := range contents {
			select {
			case <-f.stopSearching:
				return
			default:
				if content.IsDir() {
					f.closeWaitGroup.Add(1)
					go f.findFile(filepath.Join(dir, content.Name()), search)
				} else {
					filePath := filepath.Join(dir, content.Name())
					if strings.Contains(strings.ToLower(content.Name()), strings.ToLower(search)) {
						if f.addResult(Result{Name: filePath}) {
							return
						}
					}

					if isTextFile(filePath) {
						f.closeWaitGroup.Add(1)
						go f.searchInFile(filePath, search)
					}
				}
			}
		}
	}
}

func (f *Finder) searchInFile(filePath string, search string) {
	defer f.closeWaitGroup.Done()

	select {
	case <-f.stopSearching:
		return
	default:
		file, err := os.Open(filePath)
		if err != nil {
			f.errors <- CanNotReadFileErr{File: filePath, Err: err}
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			select {
			case <-f.stopSearching:
				return
			default:
				if strings.Contains(strings.ToLower(scanner.Text()), strings.ToLower(search)) {
					if f.addResult(Result{Name: filePath}) {
						return
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			f.errors <- CanNotReadFileErr{File: filePath, Err: err}
		}
	}
}

func (f *Finder) addResult(result Result) bool {
	f.resultsMutex.Lock()
	defer f.resultsMutex.Unlock()

	if f.resultsCount < f.maxResults {
		f.results <- result
		f.resultsCount++
		if f.resultsCount >= f.maxResults {
			close(f.stopSearching)
			return true
		}
	}
	return false
}

func isTextFile(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := make([]byte, 512)
	n, err := reader.Read(buffer)
	if err != nil {
		return false
	}

	if utf8.Valid(buffer[:n]) {
		return true
	}

	return false
}
