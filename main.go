package fakelog

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/adamliesko/fakelog/generator"
)

const (
	// Apache Common log format
	CommonFormat = "common"
	// Apache Combined log format
	CombinedFormat = "combined"
)

func main() {
	path := flag.String("path", "", "path to the log file (default os.Stdout)")
	rate := flag.Int("rate", 50, "approximate rate of requests per second, max of [rate, 10^6]")
	dur := flag.Int("duration", 0, "duration [s] to produce logs, 0 is inf")
	format := flag.String("format", CommonFormat, "log format")
	flag.Parse()

	rand.Seed(time.Now().Unix())
	err := run(*path, *format, *rate, *dur)
	if err != nil {
		log.Fatal(err)
	}
}

func run(path, format string, rate, dur int) error {
	f, err := getLogFile(path)
	if err != nil {
		return fmt.Errorf("unable to open file: %s : %v", path, err)
	}
	defer f.Close()

	lg, err := pickLineGenerator(format)
	if err != nil {
		return err
	}
	fl := generator.NewFakeLogger(lg, f, rate)
	// setting our stop condition - elapsed time
	if dur > 0 {
		go func(g *generator.FakeLogger) {
			time.Sleep(time.Duration(dur) * time.Second)
			g.Stop()
		}(fl)
	}
	return fl.GenerateLogs()
}

func pickLineGenerator(format string) (generator.LineGenerator, error) {
	switch format {
	case "", CommonFormat:
		return generator.ApacheCommonLine, nil
	case CombinedFormat:
		return generator.ApacheCombinedLine, nil
	default:
		return nil, fmt.Errorf("unknown log format: %s", format)
	}
}

func getLogFile(path string) (*os.File, error) {
	if path == "" {
		return os.Stdout, nil
	}
	return os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
}
