package distinct_test

import (
	"testing"

	"github.com/your-org/logpipe/internal/distinct"
)

func base(msg string) distinct.Entry {
	return distinct.Entry{"message": msg, "level": "info"}
}

func TestApply_FirstOccurrence_Passes(t *testing.T) {
	p := distinct.New(distinct.Config{})
	_, ok := p.Apply(base("hello"))
	if !ok {
		t.Fatal("expected first occurrence to pass")
	}
}

func TestApply_SecondOccurrence_Dropped(t *testing.T) {
	p := distinct.New(distinct.Config{})
	p.Apply(base("hello"))
	_, ok := p.Apply(base("hello"))
	if ok {
		t.Fatal("expected duplicate to be dropped")
	}
}

func TestApply_DifferentValues_BothPass(t *testing.T) {
	p := distinct.New(distinct.Config{})
	_, ok1 := p.Apply(base("foo"))
	_, ok2 := p.Apply(base("bar"))
	if !ok1 || !ok2 {
		t.Fatal("expected distinct values to both pass")
	}
}

func TestApply_CustomField_UsedForDedup(t *testing.T) {
	p := distinct.New(distinct.Config{Field: "request_id"})
	e1 := distinct.Entry{"request_id": "abc", "message": "x"}
	e2 := distinct.Entry{"request_id": "abc", "message": "y"}
	p.Apply(e1)
	_, ok := p.Apply(e2)
	if ok {
		t.Fatal("expected same request_id to be dropped")
	}
}

func TestApply_MissingField_AlwaysPasses(t *testing.T) {
	p := distinct.New(distinct.Config{Field: "trace_id"})
	e := distinct.Entry{"message": "no trace"}
	_, ok1 := p.Apply(e)
	_, ok2 := p.Apply(e)
	if !ok1 || !ok2 {
		t.Fatal("entries without the key field should always pass")
	}
}

func TestReset_ClearsSeen(t *testing.T) {
	p := distinct.New(distinct.Config{})
	p.Apply(base("hello"))
	p.Reset()
	_, ok := p.Apply(base("hello"))
	if !ok {
		t.Fatal("expected entry to pass after reset")
	}
}

func TestLen_TracksDistinctCount(t *testing.T) {
	p := distinct.New(distinct.Config{})
	p.Apply(base("a"))
	p.Apply(base("b"))
	p.Apply(base("a")) // duplicate
	if got := p.Len(); got != 2 {
		t.Fatalf("expected Len 2, got %d", got)
	}
}
