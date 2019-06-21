package fileserver

import (
	"net/http"
	"os"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(caddy.Module{
		Name: "http.matchers.file",
		New:  func() interface{} { return new(FileMatcher) },
	})
}

// FileMatcher is a matcher that can match requests
// based on the local file system.
// TODO: Not sure how to do this well; we'd need the ability to
// hide files, etc...
// TODO: Also consider a feature to match directory that
// contains a certain filename (use filepath.Glob), useful
// if wanting to map directory-URI requests where the dir
// has index.php to PHP backends, for example (although this
// can effectively be done with rehandling already)
type FileMatcher struct {
	Root  string   `json:"root"`
	Path  string   `json:"path"`
	Flags []string `json:"flags"`
}

// Match matches the request r against m.
func (m FileMatcher) Match(r *http.Request) bool {
	fullPath := sanitizedPathJoin(m.Root, m.Path)
	var match bool
	if len(m.Flags) > 0 {
		match = true
		fi, err := os.Stat(fullPath)
		for _, f := range m.Flags {
			switch f {
			case "EXIST":
				match = match && os.IsNotExist(err)
			case "DIR":
				match = match && err == nil && fi.IsDir()
			default:
				match = false
			}
		}
	}
	return match
}

// Interface guard
var _ caddyhttp.RequestMatcher = (*FileMatcher)(nil)