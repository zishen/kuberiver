package config

const (
	LiftedPath = "C:\\01-workspace\\karmada\\pkg\\util\\lifted"
	TmpPath    = "C:\\01-workspace\\haha"
	// K8sSynVersion must be branch release version.
	K8sSynVersion   = "1.27"
	GitHubPreRawURL = "raw.githubusercontent.com"
	LogFile         = "C:\\01-workspace\\haha.log"
)

var SkipUrlRulers = []string{
	"www.apache.org", // license used
	"#L",             // quote code line
	"/issues/",
}
