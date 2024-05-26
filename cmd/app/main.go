package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/jibaru/gofind/internal/find"
	"github.com/jibaru/gofind/internal/utils"
)

func printLogo() {
	fmt.Println(utils.Yellow, "   ____         ___  _             _  ")
	fmt.Println("  / ___| _____ / __|(_) __  __ ___| | ")
	fmt.Println(" | |  _ /  _  \\| |_ | ||  \\| |/  _  | ")
	fmt.Println(" | |_| || |_| ||  _|| ||     || |_| | ")
	fmt.Println("  \\____|\\_____/|_|  |_||_|\\__|\\_____| ")
	fmt.Println("                                      ", utils.Reset)
}

func main() {
	var dir string
	var search string
	var maxCoincidences int
	var logFile string
	var seeHelp bool

	flag.StringVar(&dir, "d", "./", "directory to search")
	flag.StringVar(&search, "s", "", "search string")
	flag.IntVar(&maxCoincidences, "m", 100, "max number of coincidences")
	flag.StringVar(&logFile, "l", "coincidences.json", "log file")
	flag.BoolVar(&seeHelp, "h", false, "see help")
	flag.Parse()

	if seeHelp {
		printLogo()
		fmt.Println("Usage: gofind -d <directory> -s <search string> -m <max number of coincidences>")
		return
	}

	if search == "" {
		printLogo()
		fmt.Println("search string is required")
		return
	}

	fmt.Println(
		utils.Yellow,
		fmt.Sprintf("searching for '%s' in '%s' with a max of %v coincidences", search, dir, maxCoincidences),
		utils.Reset,
	)

	wg := &sync.WaitGroup{}
	wg.Add(3)

	finder := find.NewFinder(maxCoincidences)
	broadcast := find.NewBroadcast(&finder)
	writer := find.NewJsonWriter(logFile, broadcast)
	printer := find.NewPrinter(broadcast)

	broadcast.Broadcast(wg)

	finder.StartFinding(dir, search)
	printer.StartPrinting(wg)
	err := writer.StartWriting(wg)
	if err != nil {
		panic(err)
	}

	wg.Wait()
}
