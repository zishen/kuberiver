package config

const (
	LiftedPath = "C:\\06-temp\\karmada\\pkg\\util\\lifted"
	TmpPath    = "C:\\01-workspace\\haha"
	// K8sSynVersion must be branch release version.
	K8sSynVersion   = "1.26"
	GitHubPreRawURL = "raw.githubusercontent.com"
	LogFile         = "C:\\01-workspace\\haha.log"
)

var SkipUrlRulers = []string{
	"www.apache.org", // license used
	//"#L",             // quote code line
	"/issues/",
}
