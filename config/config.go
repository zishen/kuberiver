package config

const (
	LiftedPath    = "C:\\01-workspace\\karmada\\pkg\\util\\lifted"
	TmpPath       = "C:\\haha"
	K8sSynVersion = "1.26"
)

var SkipUrlRulers = []string{
	"www.apache.org", // license used
	"#L",             // quote code line
	"/issues/",
}
