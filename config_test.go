package main

import (
	"bytes"
	"testing"
)

func TestConfig(t *testing.T) {
	type result struct {
		backend string
		proxy   bool
	}

	cases := []struct {
		Config string
		Tests  map[string]result
	}{
		{
			Config: `
# Comment
go.universe.tf 1.2.3.4
*.universe.tf 2.3.4.5
# Comment
google.* 3.4.5.6
/gooo+gle\.com/ 4.5.6.7
foobar.net 6.7.8.9 PROXY

/(.*)\.d1\.tyd\.us/ /$1.internal.nakagawa/
/(?P<name>.*)\.d2\.tyd\.us/ /internal.nakagawa.${name}/
/(?P<first>.*)\.(?P<second>.*)\.(?P<third>.*)\.nakagawa/ /nakagawa.${third}.${second}.${first}/
/((.*)\.d3\.tyd\.us)/ /$2.internal.nakagawa/
/(?P<port>^[0-9]{1,4})\.(?P<host>.*)\.d4\.tyd\.us/ /internal.nakagawa.${host}:${port}/
`,
			Tests: map[string]result{
				"go.universe.tf":     result{"1.2.3.4", false},
				"foo.universe.tf":    result{"2.3.4.5", false},
				"bar.universe.tf":    result{"2.3.4.5", false},
				"google.com":         result{"3.4.5.6", false},
				"google.fr":          result{"3.4.5.6", false},
				"goooooooooogle.com": result{"4.5.6.7", false},
				"foobar.net":         result{"6.7.8.9", true},

				"blah.com":            result{"", false},
				"google.com.br":       result{"", false},
				"foo.bar.universe.tf": result{"", false},
				"goooooglexcom":       result{"", false},

				"lucus.d1.tyd.us":     result{"lucus.internal.nakagawa", false},
				"lucus.d2.tyd.us":     result{"internal.nakagawa.lucus", false},
				"1.2.3.nakagawa":      result{"nakagawa.3.2.1", false},
				"lucus.d3.tyd.us":     result{"lucus.internal.nakagawa", false},
				"443.lucus.d4.tyd.us": result{"internal.nakagawa.lucus:443", false},
			},
		},
	}

	for _, test := range cases {
		var cfg Config
		if err := cfg.Read(bytes.NewBufferString(test.Config)); err != nil {
			t.Fatalf("Failed to read config (%s):\n%q", err, test.Config)
		}

		for hostname, expected := range test.Tests {
			backend, proxy := cfg.Match(hostname)
			if expected.backend != backend {
				t.Errorf("cfg.Match(%q) is %q, want %q", hostname, backend, expected.backend)
			}
			if expected.proxy != proxy {
				t.Errorf("cfg.Match(%q).proxy is %v, want %v", hostname, proxy, expected.proxy)
			}
		}
	}
}
