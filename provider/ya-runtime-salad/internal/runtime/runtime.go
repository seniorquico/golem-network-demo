package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/akamensky/argparse"
)

type runCommandArgs struct {
	// Duration of run (in seconds). Must be a positive multiple of frequency. Defaults to 300.
	Duration int `json:"duration"`
	// Frequency of updates (in seconds). Must be a positive number. Defaults to 60.
	Frequency int `json:"frequency"`
}

func HandleInvocation(ctx context.Context, argv []string, stdout, stderr io.Writer) int {
	parser := argparse.NewParser("ya-runtime-salad", "SaladCloud Runtime")
	command := parser.StringPositional(nil)
	_ = parser.String("", "cpu-cores", nil)
	_ = parser.String("", "entrypoint", nil)
	_ = parser.String("", "hostname", nil)
	_ = parser.String("", "inet-endpoint", nil)
	_ = parser.String("", "mem-gib", nil)
	_ = parser.String("", "storage-gib", nil)
	_ = parser.String("", "task-package", nil)
	_ = parser.String("", "volume-override", nil)
	_ = parser.String("", "vpn-endpoint", nil)
	_ = parser.String("", "workdir", nil)

	cut := len(argv)
	for i, v := range argv {
		if v == "--" {
			cut = i
			break
		}
	}
	if err := parser.Parse(argv[:cut]); err != nil {
		fmt.Fprintf(stderr, "failed to parse arguments: %v\n", err)
		return 1
	}

	var c string
	if command == nil {
		c = "version"
	} else {
		c = *command
	}

	switch c {
	case "version":
		fmt.Fprintln(stdout, "0.0.1")
		return 0
	case "test":
		return 0
	case "offer-template":
		template := map[string]any{
			"constraints": "",
			"properties": map[string]any{
				"golem.!exp.gap-35.v1.inf.gpu.clocks.graphics.mhz":     2100,
				"golem.!exp.gap-35.v1.inf.gpu.clocks.memory.mhz":       7001,
				"golem.!exp.gap-35.v1.inf.gpu.clocks.sm.mhz":           2100,
				"golem.!exp.gap-35.v1.inf.gpu.clocks.video.mhz":        1950,
				"golem.!exp.gap-35.v1.inf.gpu.cuda.compute-capability": "8.6",
				"golem.!exp.gap-35.v1.inf.gpu.cuda.cores":              4864,
				"golem.!exp.gap-35.v1.inf.gpu.cuda.enabled":            true,
				"golem.!exp.gap-35.v1.inf.gpu.cuda.version":            "13.0",
				"golem.!exp.gap-35.v1.inf.gpu.memory.bandwidth.gib":    448.0,
				"golem.!exp.gap-35.v1.inf.gpu.memory.total.gib":        8.0,
				"golem.!exp.gap-35.v1.inf.gpu.model":                   "NVIDIA GeForce RTX 3060 Ti",
				"golem.inf.cpu.brand":                                  "Intel(R) Core(TM) i9-14900K",
				"golem.inf.cpu.model":                                  "Stepping 1 Family 6 Model 183",
				"golem.inf.cpu.vendor":                                 "GenuineIntel",
				"golem.runtime.capabilities":                           []string{"!exp:gpu"},
			},
		}
		enc := json.NewEncoder(stdout)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(template); err != nil {
			fmt.Fprintf(stderr, "failed to encode offer template response: %v\n", err)
			return 1
		}
		return 0
	case "deploy":
		template := map[string]any{
			"startMode": "empty",
			"valid": map[string]any{
				"Ok": "",
			},
			"vols": []any{},
		}
		enc := json.NewEncoder(stdout)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(template); err != nil {
			fmt.Fprintf(stderr, "failed to encode deploy response: %v\n", err)
			return 1
		}
		return 0
	case "start":
		return 0
	case "run":
		durationSeconds := 300
		frequencySeconds := 60
		if len(argv) > (cut + 1) {
			var args runCommandArgs
			if err := json.Unmarshal([]byte(argv[cut+1]), &args); err == nil {
				durationSeconds = args.Duration
				if durationSeconds <= 0 {
					durationSeconds = 1
				}

				frequencySeconds = args.Frequency
				if frequencySeconds <= 0 {
					frequencySeconds = 1
				}
			}
		}

		deadline := time.Now().Add(time.Duration(durationSeconds) * time.Second)
		ticker := time.NewTicker(time.Duration(frequencySeconds) * time.Second)
		defer ticker.Stop()
	loop:
		for {
			select {
			case <-ctx.Done():
				fmt.Fprintln(stdout, "interrupted")
				break loop
			case t := <-ticker.C:
				fmt.Fprintln(stdout, t.Format(time.RFC3339))
				if t.After(deadline) {
					fmt.Fprintln(stdout, "completed")
					break loop
				}
			}
		}

		return 0
	case "stop":
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", c)
		return 1
	}
}
