package formatter

import (
	"fmt"
	"github.com/docker/docker/api/types/events"
	"regexp"
	"testing"
)

func TestEventFilter_Match_ImageRegexp(t *testing.T) {
	for _, c := range []struct {
		image  string
		regexp *regexp.Regexp
		match  bool
	}{
		{
			image:  "alpine:latest",
			regexp: nil,
			match:  true,
		},
		{
			image:  "alpine:latest",
			regexp: regexp.MustCompile("alpine"),
			match:  true,
		},
		{
			image:  "alpine:latest",
			regexp: regexp.MustCompile("^alpine$"),
			match:  false,
		},
		{
			image:  "alpine:latest",
			regexp: regexp.MustCompile("ubuntu"),
			match:  false,
		},
	} {
		t.Run(fmt.Sprintf("image=%s/regexp=%s", c.image, c.regexp), func(t *testing.T) {
			msg := events.Message{
				Type: "container",
				Actor: events.Actor{
					Attributes: map[string]string{"image": c.image},
				},
			}
			filter := EventFilter{ImageRegexp: c.regexp}
			if m := filter.Match(msg); m != c.match {
				t.Errorf("match wants %v but %v", c.match, m)
			}
		})
	}
}
