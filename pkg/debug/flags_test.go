package debug

import (
	"os"
	"testing"
)

func cleanupDebugEnv() {
	os.Unsetenv("GOCT_TRACE")
	os.Unsetenv("GOCT_VERBOSE")
	os.Unsetenv("GOCT_DUMP")
}

func TestDebugFlags_Resolve_EnvVar(t *testing.T) {
	tests := []struct {
		name     string
		envKey   string
		envValue string
		expected TraceLevel
	}{
		{"GOCT_TRACE=true", "GOCT_TRACE", "true", TraceLevelTrace},
		{"GOCT_VERBOSE=true", "GOCT_VERBOSE", "true", TraceLevelVerbose},
		{"GOCT_DUMP=true", "GOCT_DUMP", "true", TraceLevelDump},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean all env vars before each test
			cleanupDebugEnv()
			defer cleanupDebugEnv()

			// Set the env var under test
			os.Setenv(tt.envKey, tt.envValue)

			flags := &DebugFlags{}
			result := flags.Resolve()

			if result.TraceLevel != tt.expected {
				t.Errorf("Resolve() = %v, want %v", result.TraceLevel, tt.expected)
			}
		})
	}
}

func TestDebugFlags_Resolve_FlagPriority(t *testing.T) {
	cleanupDebugEnv()
	defer cleanupDebugEnv()

	// Test that CLI flags take priority over env vars
	// When GOCT_TRACE=true but flag has Verbose=true, should return Verbose
	os.Setenv("GOCT_TRACE", "true")

	flags := &DebugFlags{Verbose: true}
	result := flags.Resolve()

	if result.TraceLevel != TraceLevelVerbose {
		t.Errorf("Resolve() with Verbose flag = %v, want %v (flag should override env)",
			result.TraceLevel, TraceLevelVerbose)
	}
}

func TestDebugFlags_Resolve_AllFlags(t *testing.T) {
	cleanupDebugEnv()
	defer cleanupDebugEnv()

	// Test multiple flags set - last one wins
	flags := &DebugFlags{Trace: true, Verbose: true, Dump: true}
	result := flags.Resolve()

	if result.TraceLevel != TraceLevelDump {
		t.Errorf("Resolve() with all flags = %v, want %v", result.TraceLevel, TraceLevelDump)
	}
}

func TestDebugFlags_Resolve_NoneSet(t *testing.T) {
	cleanupDebugEnv()
	defer cleanupDebugEnv()

	// No flags, no env vars
	flags := &DebugFlags{}
	result := flags.Resolve()

	if result.TraceLevel != TraceLevelOff {
		t.Errorf("Resolve() with no flags = %v, want %v", result.TraceLevel, TraceLevelOff)
	}
}