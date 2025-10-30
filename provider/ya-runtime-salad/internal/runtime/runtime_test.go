package runtime

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/afero"
)

func TestVersionCommand(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	exitCode := HandleInvocation(
		context.Background(),
		[]string{"/root/.local/lib/yagna/plugins/ya-runtime-test", "version"},
		fs,
		&stdout,
		&stderr,
	)

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d", exitCode)
	}

	expectedOutput := "0.0.1\n"
	if stdout.String() != expectedOutput {
		t.Fatalf("Expected stdout %q, got %q", expectedOutput, stdout.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("Expected empty stderr, got %q", stderr.String())
	}
}

func TestTestCommand(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	exitCode := HandleInvocation(
		context.Background(),
		[]string{"/root/.local/lib/yagna/plugins/ya-runtime-test", "test"},
		fs,
		&stdout,
		&stderr,
	)

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("Expected empty stdout, got %q", stdout.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("Expected empty stderr, got %q", stderr.String())
	}
}

func TestOfferTemplateCommand(t *testing.T) {
	testCases := []struct {
		name   string
		config string
		output string
	}{
		{
			name:   "cpu",
			config: `{"golem.inf.cpu.brand":"Intel(R) Core(TM) i9-14900K","golem.inf.cpu.model":"Stepping 1 Family 6 Model 183","golem.inf.cpu.vendor":"GenuineIntel"}`,
			output: `{"constraints":"","properties":{"golem.inf.cpu.brand":"Intel(R) Core(TM) i9-14900K","golem.inf.cpu.model":"Stepping 1 Family 6 Model 183","golem.inf.cpu.vendor":"GenuineIntel","golem.runtime.capabilities":[]}}`,
		},
		{
			name:   "gpu",
			config: `{"golem.!exp.gap-35.v1.inf.gpu.clocks.graphics.mhz":2100,"golem.!exp.gap-35.v1.inf.gpu.clocks.memory.mhz":7001,"golem.!exp.gap-35.v1.inf.gpu.clocks.sm.mhz":2100,"golem.!exp.gap-35.v1.inf.gpu.clocks.video.mhz":1950,"golem.!exp.gap-35.v1.inf.gpu.cuda.compute-capability":"8.6","golem.!exp.gap-35.v1.inf.gpu.cuda.cores":4864,"golem.!exp.gap-35.v1.inf.gpu.cuda.enabled":true,"golem.!exp.gap-35.v1.inf.gpu.cuda.version":"13.0","golem.!exp.gap-35.v1.inf.gpu.memory.bandwidth.gib":448.0,"golem.!exp.gap-35.v1.inf.gpu.memory.total.gib":8.0,"golem.!exp.gap-35.v1.inf.gpu.model":"NVIDIA GeForce RTX 3060 Ti","golem.inf.cpu.brand":"Intel(R) Core(TM) i9-14900K","golem.inf.cpu.model":"Stepping 1 Family 6 Model 183","golem.inf.cpu.vendor":"GenuineIntel"}`,
			output: `{"constraints":"","properties":{"golem.!exp.gap-35.v1.inf.gpu.clocks.graphics.mhz":2100,"golem.!exp.gap-35.v1.inf.gpu.clocks.memory.mhz":7001,"golem.!exp.gap-35.v1.inf.gpu.clocks.sm.mhz":2100,"golem.!exp.gap-35.v1.inf.gpu.clocks.video.mhz":1950,"golem.!exp.gap-35.v1.inf.gpu.cuda.compute-capability":"8.6","golem.!exp.gap-35.v1.inf.gpu.cuda.cores":4864,"golem.!exp.gap-35.v1.inf.gpu.cuda.enabled":true,"golem.!exp.gap-35.v1.inf.gpu.cuda.version":"13.0","golem.!exp.gap-35.v1.inf.gpu.memory.bandwidth.gib":448,"golem.!exp.gap-35.v1.inf.gpu.memory.total.gib":8,"golem.!exp.gap-35.v1.inf.gpu.model":"NVIDIA GeForce RTX 3060 Ti","golem.inf.cpu.brand":"Intel(R) Core(TM) i9-14900K","golem.inf.cpu.model":"Stepping 1 Family 6 Model 183","golem.inf.cpu.vendor":"GenuineIntel","golem.runtime.capabilities":["!exp:gpu"]}}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewMemMapFs()
			userConfigDir, err := os.UserConfigDir()
			if err != nil {
				t.Fatalf("Failed to get user config dir: %v", err)
			}

			path := filepath.Join(userConfigDir, "ya-runtime-salad", "template.json")
			afero.WriteFile(fs, path, []byte(tc.config), 0o600)

			stdout := bytes.Buffer{}
			stderr := bytes.Buffer{}
			exitCode := HandleInvocation(
				context.Background(),
				[]string{"/root/.local/lib/yagna/plugins/ya-runtime-test", "offer-template"},
				fs,
				&stdout,
				&stderr,
			)

			if exitCode != 0 {
				t.Fatalf("Expected exit code 0, got %d", exitCode)
			}

			expectedOutput := tc.output + "\n"
			if stdout.String() != expectedOutput {
				t.Fatalf("Expected stdout %q, got %q", expectedOutput, stdout.String())
			}

			if stderr.Len() != 0 {
				t.Fatalf("Expected empty stderr, got %q", stderr.String())
			}
		})
	}
}

func TestDeployCommand(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	exitCode := HandleInvocation(
		context.Background(),
		[]string{
			"/root/.local/lib/yagna/plugins/ya-runtime-test",
			"--workdir",
			"/root/.local/share/ya-provider/exe-unit/work/eb7c05f95a0ac7f67d217f75155864172c740988c179da3d46aa13b33fb1d9a0/5cf174d64b8a493f82bebd3e9307db7c",
			"--task-package",
			"/root/.local/share/ya-provider/exe-unit/cache/54200484f0115776b5ed339a30811bc621479ccc4169181aeabca4ad8cb07ab2_64c5a5548cea45177ba89dcc13bea00bd7e8d6db5bbf81872fa462f3",
			"deploy",
			"--",
		},
		fs,
		&stdout,
		&stderr,
	)

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d", exitCode)
	}

	expectedOutput := `{"startMode":"empty","valid":{"Ok":""},"vols":[]}` + "\n"
	if stdout.String() != expectedOutput {
		t.Fatalf("Expected stdout %q, got %q", expectedOutput, stdout.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("Expected empty stderr, got %q", stderr.String())
	}
}

func TestStartCommand(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	exitCode := HandleInvocation(
		context.Background(),
		[]string{
			"/root/.local/lib/yagna/plugins/ya-runtime-test",
			"--workdir",
			"/root/.local/share/ya-provider/exe-unit/work/eb7c05f95a0ac7f67d217f75155864172c740988c179da3d46aa13b33fb1d9a0/5cf174d64b8a493f82bebd3e9307db7c",
			"--task-package",
			"/root/.local/share/ya-provider/exe-unit/cache/54200484f0115776b5ed339a30811bc621479ccc4169181aeabca4ad8cb07ab2_64c5a5548cea45177ba89dcc13bea00bd7e8d6db5bbf81872fa462f3",
			"start",
			"--",
		},
		fs,
		&stdout,
		&stderr,
	)

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("Expected empty stdout, got %q", stdout.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("Expected empty stderr, got %q", stderr.String())
	}
}

func TestRunCommand(t *testing.T) {
	testCases := []struct {
		name      string
		args      string
		duration  int
		frequency int
	}{
		{
			name:      "fast",
			args:      `{"duration":5,"frequency":2}`,
			duration:  5,
			frequency: 2,
		},
		{
			name:      "slow",
			args:      `{"duration":10,"frequency":3}`,
			duration:  10,
			frequency: 3,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewMemMapFs()
			stdout := bytes.Buffer{}
			stderr := bytes.Buffer{}
			exitCode := HandleInvocation(
				context.Background(),
				[]string{
					"/root/.local/lib/yagna/plugins/ya-runtime-test",
					"--workdir",
					"/root/.local/share/ya-provider/exe-unit/work/eb7c05f95a0ac7f67d217f75155864172c740988c179da3d46aa13b33fb1d9a0/5cf174d64b8a493f82bebd3e9307db7c",
					"--task-package",
					"/root/.local/share/ya-provider/exe-unit/cache/54200484f0115776b5ed339a30811bc621479ccc4169181aeabca4ad8cb07ab2_64c5a5548cea45177ba89dcc13bea00bd7e8d6db5bbf81872fa462f3",
					"run",
					"--entrypoint",
					"command",
					"--",
					tc.args,
				},
				fs,
				&stdout,
				&stderr,
			)

			if exitCode != 0 {
				t.Fatalf("Expected exit code 0, got %d", exitCode)
			}

			lines := bytes.Split(stdout.Bytes(), []byte("\n"))
			if len(lines) > 0 && len(lines[len(lines)-1]) == 0 {
				lines = lines[:len(lines)-1]
			}

			expectedLines := (tc.duration+tc.frequency-1)/tc.frequency + 1
			if len(lines) != expectedLines {
				t.Fatalf("Expected %d stdout lines, got %d", expectedLines, len(lines))
			}
			if string(lines[len(lines)-1]) != "completed" {
				t.Fatalf("Expected last stdout line to be 'completed', got %q", string(lines[len(lines)-1]))
			}

			if stderr.Len() != 0 {
				t.Fatalf("Expected empty stderr, got %q", stderr.String())
			}
		})
	}
}

func TestRunCommandCancel(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(2500 * time.Millisecond)
		cancel()
	}()
	exitCode := HandleInvocation(
		ctx,
		[]string{
			"/root/.local/lib/yagna/plugins/ya-runtime-test",
			"--workdir",
			"/root/.local/share/ya-provider/exe-unit/work/eb7c05f95a0ac7f67d217f75155864172c740988c179da3d46aa13b33fb1d9a0/5cf174d64b8a493f82bebd3e9307db7c",
			"--task-package",
			"/root/.local/share/ya-provider/exe-unit/cache/54200484f0115776b5ed339a30811bc621479ccc4169181aeabca4ad8cb07ab2_64c5a5548cea45177ba89dcc13bea00bd7e8d6db5bbf81872fa462f3",
			"run",
			"--entrypoint",
			"command",
			"--",
			`{"duration":5,"frequency":2}`,
		},
		fs,
		&stdout,
		&stderr,
	)

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d", exitCode)
	}

	lines := bytes.Split(stdout.Bytes(), []byte("\n"))
	if len(lines) > 0 && len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}

	if len(lines) != 2 {
		t.Fatalf("Expected 2 stdout lines, got %d", len(lines))
	}
	if string(lines[len(lines)-1]) != "interrupted" {
		t.Fatalf("Expected last stdout line to be 'interrupted', got %q", string(lines[len(lines)-1]))
	}

	if stderr.Len() != 0 {
		t.Fatalf("Expected empty stderr, got %q", stderr.String())
	}
}
