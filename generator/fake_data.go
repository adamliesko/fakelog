package generator

var (
	codes = newWeightedChoice(
		map[string]int{
			"200": 5,
			"404": 1,
			"500": 1,
			"503": 1,
			"401": 1,
			"403": 1,
			"301": 1,
		},
	)
	methods = newWeightedChoice(
		map[string]int{
			"GET":    5,
			"DELETE": 1,
			"POST":   1,
			"PUT":    1,
			"PATCH":  1,
		},
	)
	userAgents = []string{
		"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2224.3 Safari/537.36",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; rv:2.2) Gecko/20110201",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.75.14 (KHTML, like Gecko) Version/7.0.3 Safari/7046A194A",
		"Opera/9.80 (X11; Linux i686; Ubuntu/14.10) Presto/2.12.388 Version/12.16",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5376e Safari/8536.25",
		"Mozilla/5.0 (iPad; CPU OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5376e Safari/8536.25",
		"filthy crawler",
		"Mozilla/5.0 (Android 4.4; Mobile; rv:41.0) Gecko/41.0 Firefox/41.0",
	}
	referrers = []string{"google.com", "http://yahoo.com/landing", "https://twitter.com/janice43", "http://lost.at", "http://oc.dk", "http://loopy.ch/dash33?=34"}
	userNames = []string{"john_doe", "leet_coder", "abrv", "sarah_cooper", "grace_hooper", "ada_l"}
	endpoints = []string{"/articles", "/blog/ethan", "/users", "/login", "/signup", "/trending", "/popular", "/article/", "/subscribe/monthly?userID="}
)
