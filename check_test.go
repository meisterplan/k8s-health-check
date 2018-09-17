package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	os.Unsetenv(livenessCommandEnvName)
	os.Unsetenv(readinessCommandEnvName)

	os.Remove(livenessStateFileName)
	os.Remove(readinessStateFileName)

	code := m.Run()
	os.Exit(code)
}

func runAndExpectExit(t *testing.T, command string, exitCodeShouldBeZero bool) string {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.Output()
	if e, ok := err.(*exec.ExitError); ok {
		if exitCodeShouldBeZero != e.Success() {
			t.Fatal("$", command, " execution success = ", e.Success(), " but expected = ", exitCodeShouldBeZero)
		}
	} else if !exitCodeShouldBeZero {
		t.Fatal("$", command, " should not return exit code 0, but did anyway")
	}
	return string(output)
}

func runAndExpectOutputAndExit(t *testing.T, command string, exitCodeShouldBeZero bool, stringToContain string) {
	actualOutput := runAndExpectExit(t, command, exitCodeShouldBeZero)
	if !strings.Contains(actualOutput, stringToContain) {
		t.Fatalf(command, " did not contain '", stringToContain, "', actual output was '", actualOutput, "'")
	}
}

func TestPlainCall(t *testing.T) {
	runAndExpectOutputAndExit(t, "./check", false, "Missing subcommand")
}
func TestLivenessNotDefined(t *testing.T) {
	runAndExpectOutputAndExit(t, "./check run -type liveness", false, "must be set")
}

func TestLivenessFails(t *testing.T) {
	os.Setenv(livenessCommandEnvName, "false")
	runAndExpectOutputAndExit(t, "./check run -type liveness", false, "FAILURE")
}

func TestLivenessSucceeds(t *testing.T) {
	os.Setenv(livenessCommandEnvName, "true")
	runAndExpectOutputAndExit(t, "./check run -type liveness", true, "SUCCESS")
}

func TestLivenessLockedFailure(t *testing.T) {
	os.Setenv(livenessCommandEnvName, "true")
	runAndExpectOutputAndExit(t, "./check lock -type liveness -state failure", true, "")
	runAndExpectOutputAndExit(t, "./check run -type liveness", false, "locked")
}

func TestLivenessLockedSuccess(t *testing.T) {
	os.Setenv(livenessCommandEnvName, "false")
	runAndExpectOutputAndExit(t, "./check lock -type liveness -state success", true, "")
	runAndExpectOutputAndExit(t, "./check run -type liveness", true, "locked")
}

func TestLivenessUnlockedFailure(t *testing.T) {
	os.Setenv(livenessCommandEnvName, "true")
	runAndExpectOutputAndExit(t, "./check lock -type liveness -state failure", true, "")
	runAndExpectOutputAndExit(t, "./check unlock -type liveness", true, "")
	runAndExpectOutputAndExit(t, "./check run -type liveness", true, "SUCCESS")
}

func TestLivenessUnlockedSuccess(t *testing.T) {
	os.Setenv(livenessCommandEnvName, "false")
	runAndExpectOutputAndExit(t, "./check lock -type liveness -state success", true, "")
	runAndExpectOutputAndExit(t, "./check unlock -type liveness", true, "")
	runAndExpectOutputAndExit(t, "./check run -type liveness", false, "FAILURE")
}

func TestReadinessNotDefined(t *testing.T) {
	runAndExpectOutputAndExit(t, "./check run -type Readiness", false, "must be set")
}

func TestReadinessFails(t *testing.T) {
	os.Setenv(readinessCommandEnvName, "false")
	runAndExpectOutputAndExit(t, "./check run -type Readiness", false, "FAILURE")
}

func TestReadinessSucceeds(t *testing.T) {
	os.Setenv(readinessCommandEnvName, "true")
	runAndExpectOutputAndExit(t, "./check run -type Readiness", true, "SUCCESS")
}

func TestReadinessLockedFailure(t *testing.T) {
	os.Setenv(readinessCommandEnvName, "true")
	runAndExpectOutputAndExit(t, "./check lock -type Readiness -state failure", true, "")
	runAndExpectOutputAndExit(t, "./check run -type Readiness", false, "locked")
}

func TestReadinessLockedSuccess(t *testing.T) {
	os.Setenv(readinessCommandEnvName, "false")
	runAndExpectOutputAndExit(t, "./check lock -type Readiness -state success", true, "")
	runAndExpectOutputAndExit(t, "./check run -type Readiness", true, "locked")
}

func TestReadinessUnlockedFailure(t *testing.T) {
	os.Setenv(readinessCommandEnvName, "true")
	runAndExpectOutputAndExit(t, "./check lock -type Readiness -state failure", true, "")
	runAndExpectOutputAndExit(t, "./check unlock -type Readiness", true, "")
	runAndExpectOutputAndExit(t, "./check run -type Readiness", true, "SUCCESS")
}

func TestReadinessUnlockedSuccess(t *testing.T) {
	os.Setenv(readinessCommandEnvName, "false")
	runAndExpectOutputAndExit(t, "./check lock -type Readiness -state success", true, "")
	runAndExpectOutputAndExit(t, "./check unlock -type Readiness", true, "")
	runAndExpectOutputAndExit(t, "./check run -type Readiness", false, "FAILURE")
}
