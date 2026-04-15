package pipeline_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/pipeline"
	"github.com/yourorg/logpipe/internal/config"
)

func TestRun_BasicPipeline(t *testing.T) {
	input := `{"level":"info","msg":"hello world","ts":"2024-01-01T00:00:00Z"}
{"level":"error","msg":"something failed","ts":"2024-01-01T00:00:01Z"}
`

	var out strings.Builder
	cfg := &config.Config{
		Sources: []config.Source{
			{Reader: strings.NewReader(input), Label: "test"},
		},
		Pattern:       "",
		Level:         "",
		NoColor:       true,
		Output:        &out,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := pipeline.Run(ctx, cfg)
	if err != nil && err != context.DeadlineExceeded {
		t.Fatalf("unexpected error: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "hello world") {
		t.Errorf("expected output to contain 'hello world', got: %s", result)
	}
	if !strings.Contains(result, "something failed") {
		t.Errorf("expected output to contain 'something failed', got: %s", result)
	}
}

func TestRun_WithPatternFilter(t *testing.T) {
	input := `{"level":"info","msg":"hello world","ts":"2024-01-01T00:00:00Z"}
{"level":"error","msg":"something failed","ts":"2024-01-01T00:00:01Z"}
`

	var out strings.Builder
	cfg := &config.Config{
		Sources: []config.Source{
			{Reader: strings.NewReader(input), Label: "test"},
		},
		Pattern: "hello",
		Level:   "",
		NoColor: true,
		Output:  &out,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_ = pipeline.Run(ctx, cfg)

	result := out.String()
	if !strings.Contains(result, "hello world") {
		t.Errorf("expected output to contain 'hello world', got: %s", result)
	}
	if strings.Contains(result, "something failed") {
		t.Errorf("expected output NOT to contain 'something failed', got: %s", result)
	}
}

func TestRun_WithLevelFilter(t *testing.T) {
	input := `{"level":"info","msg":"info message","ts":"2024-01-01T00:00:00Z"}
{"level":"error","msg":"error message","ts":"2024-01-01T00:00:01Z"}
{"level":"debug","msg":"debug message","ts":"2024-01-01T00:00:02Z"}
`

	var out strings.Builder
	cfg := &config.Config{
		Sources: []config.Source{
			{Reader: strings.NewReader(input), Label: "test"},
		},
		Pattern: "",
		Level:   "error",
		NoColor: true,
		Output:  &out,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_ = pipeline.Run(ctx, cfg)

	result := out.String()
	if !strings.Contains(result, "error message") {
		t.Errorf("expected output to contain 'error message', got: %s", result)
	}
	if strings.Contains(result, "info message") {
		t.Errorf("expected output NOT to contain 'info message', got: %s", result)
	}
	if strings.Contains(result, "debug message") {
		t.Errorf("expected output NOT to contain 'debug message', got: %s", result)
	}
}
