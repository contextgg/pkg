package version

import (
	"bytes"
	"html/template"
	"runtime"
	"strings"
)

// Build information. Populated at build-time.
var (
	Version    string
	BuildID    string
	RevisionID string
	ShortSHA   string
	BranchName string
	RepoName   string
	GoVersion  = runtime.Version()
)

// versionInfoTmpl contains the template used by Info.
var versionInfoTmpl = `
{{.program}}, version {{.version}} (branch: {{.branch}}, revision: {{.revision}})
	build:            {{.build}}
	sha:              {{.sha}}
  repo:             {{.repo}}
  go version:       {{.goVersion}}
`

// Print returns version information.
func Print(program string) string {
	m := map[string]string{
		"program":   program,
		"version":   Version,
		"branch":    BranchName,
		"revision":  RevisionID,
		"build":     BuildID,
		"sha":       ShortSHA,
		"repo":      RepoName,
		"goVersion": GoVersion,
	}
	t := template.Must(template.New("version").Parse(versionInfoTmpl))

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "version", m); err != nil {
		panic(err)
	}
	return strings.TrimSpace(buf.String())
}
