package generator

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	maxRate          = 1000
	apacheTimeFormat = "02/01/2006:15:04:05 -0700" // dd/MMM/yyyy:hh:mm:ss +-hhmm
)

// LineGenerator is the signature of fake log lines generating function.
type LineGenerator func() string

// FakeLogger generates fake and random log lines.
type FakeLogger struct {
	lineFn LineGenerator
	out    io.WriteCloser
	rate   int
	cancel chan bool
}

// NewFakeLogger returns configured FakeLogger, enforcing sensible rate.
func NewFakeLogger(lineFn LineGenerator, out io.WriteCloser, rate int) *FakeLogger {
	if rate > maxRate || rate == 0 {
		rate = maxRate
	}

	return &FakeLogger{
		lineFn: lineFn,
		out:    out,
		rate:   rate,
		cancel: make(chan bool),
	}
}

// GenerateLogs generates W3-c common log lines with rps of rate to the output out.
func (g *FakeLogger) GenerateLogs() error {

	// we want a bit randomized rate, therefore avoiding to use a predictive ticker
	sleepBase := int64((1.0 / float64(g.rate)) * float64(maxRate))

	for {
		select {
		case <-g.cancel:
			return nil
		default:
			line := g.lineFn()
			if _, err := g.out.Write([]byte(line + "\n")); err != nil {
				return fmt.Errorf("unable to write log line to out: %v", err)
			}

			// if we want to get close to maxRate wa can't sleep at all
			if sleepBase != 0 {
				add := rand.Int63n((sleepBase / 5) + 1)
				var d time.Duration
				if add%2 == 0 {
					d = time.Duration(sleepBase - add)
				} else {
					d = time.Duration(sleepBase + add)
				}
				time.Sleep(d * time.Millisecond)
			}
		}
	}
}

// Stop notifies the cancel chanel of a FakeLogger, resulting into stopping it's log generation in GenerateLogs().
func (g *FakeLogger) Stop() {
	g.cancel <- true
}

// ApacheCommonLine generates an apache common access log line. It implements the Line generator func signature.
// format source: https://httpd.apache.org/docs/2.4/logs.html#accesslog
func ApacheCommonLine() string {
	ip := ipV4Address()
	u := randomEle(userNames)
	dt := time.Now().Format(apacheTimeFormat)
	code := codes.next()
	size := rand.Intn(30000)
	method := methods.next()

	ep := endpoints[rand.Intn(len(endpoints)-1)]
	if strings.HasSuffix(ep, "/") || strings.HasSuffix(ep, "=") {
		ep = ep + strconv.Itoa(rand.Intn(10000))
	}

	return fmt.Sprintf(`%s %s - [%s] "%s %s HTTP/1.1" %s %d`, ip, u, dt, method, ep, code, size)
}

// ApacheCombinedLine generates an apache combined access log line. It implements the Line generator func signature.
// format source: https://httpd.apache.org/docs/2.4/logs.html#accesslog
func ApacheCombinedLine() string {
	ua := randomEle(userAgents)
	ref := randomEle(referrers)

	return ApacheCommonLine() + fmt.Sprintf(` "%s" "%s"`, ref, ua)
}

func randomEle(in []string) string {
	return in[rand.Intn(len(userNames)-1)]
}

func ipV4Address() string {
	blocks := []string{}
	for i := 0; i < 4; i++ {
		number := rand.Intn(255)
		blocks = append(blocks, strconv.Itoa(number))
	}

	return strings.Join(blocks, ".")
}
