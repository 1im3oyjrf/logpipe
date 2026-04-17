package filter_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/filter"
)

func TestFilter_NoPattern_MatchesAll(t *testing.T) {
	f := filter.New("", nil, false)
	entry := filter.Entry{"level": "info", "msg": "hello world"}
	if !f.Match(entry) {
		t.Error("expected empty pattern to match all entries")
	}
}

func TestFilter_MatchesSubstring(t *testing.T) {
	f := filter.New("error", nil, false)

	matching := filter.Entry{"level": "error", "msg": "something failed"}
	if !f.Match(matching) {
		t.Error("expected entry with 'error' level to match")
	}

	nonMatching := filter.Entry{"level": "info", "msg": "all good"}
	if f.Match(nonMatching) {
		t.Error("expected entry without 'error' to not match")
	}
}

func TestFilter_CaseInsensitive(t *testing.T) {
	f := filter.New("ERROR", nil, false)
	entry := filter.Entry{"level": "error", "msg": "failed"}
	if !f.Match(entry) {
		t.Error("expected case-insensitive match")
	}
}

func TestFilter_CaseSensitive_NoMatch(t *testing.T) {
	f := filter.New("ERROR", nil, true)
	entry := filter.Entry{"level": "error", "msg": "failed"}
	if f.Match(entry) {
		t.Error("expected case-sensitive filter to not match lowercase 'error'")
	}
}

func TestFilter_CaseSensitive_Match(t *testing.T) {
	f := filter.New("ERROR", nil, true)
	entry := filter.Entry{"level": "ERROR", "msg": "failed"}
	if !f.Match(entry) {
		t.Error("expected case-sensitive filter to match exact uppercase 'ERROR'")
	}
}

func TestFilter_SpecificFields(t *testing.T) {
	f := filter.New("timeout", []string{"msg"}, false)

	matching := filter.Entry{"level": "error", "msg": "request timeout occurred"}
	if !f.Match(matching) {
		t.Error("expected match on 'msg' field")
	}

	// 'timeout' is in a non-targeted field
	nonMatching := filter.Entry{"level": "timeout", "msg": "all good"}
	if f.Match(nonMatching) {
		t.Error("expected no match when pattern is only in non-targeted field")
	}
}

func TestFilter_MissingField(t *testing.T) {
	f := filter.New("debug", []string{"level"}, false)
	entry := filter.Entry{"msg": "some message"}
	if f.Match(entry) {
		t.Error("expected no match when targeted field is absent")
	}
}
