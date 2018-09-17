package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// CheckType identifies the type of a check, i.e. liveness or readiness
type CheckType int

const (
	// LIVENESS probe
	LIVENESS CheckType = iota
	// READINESS probe
	READINESS
)

// ToCheckType parses CheckType from string
func ToCheckType(s string) CheckType {
	switch strings.ToLower(s) {
	case "liveness":
		return LIVENESS
	case "readiness":
		return READINESS
	default:
		fmt.Println("Not a type: ", s)
		os.Exit(1)
	}

	return -1
}

// CheckState identifies state of a check, i.e. success or failure
type CheckState int

const (
	// SUCCESS corresponds to successful probe
	SUCCESS CheckState = iota
	// FAILURE corresponds to failed probe
	FAILURE
)

// ToCheckState parses CheckState from string
func ToCheckState(s string) CheckState {
	switch strings.ToLower(s) {
	case "success":
		return SUCCESS
	case "failure":
		return FAILURE
	default:
		fmt.Println("Not a state: ", s, " expecting: (failure, success)")
		os.Exit(1)
	}

	return -1
}

func (it CheckState) String() string {
	switch it {
	case SUCCESS:
		return "SUCCESS"
	case FAILURE:
		return "FAILURE"
	default:
		return fmt.Sprintf("%d", int(it))
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing subcommand (run, lock or unlock)")
		os.Exit(1)
	}

	run := flag.NewFlagSet("run", flag.ExitOnError)
	pRunType := run.String("type", "", "liveness or readiness")

	lock := flag.NewFlagSet("lock", flag.ExitOnError)
	pLockType := lock.String("type", "", "liveness or readiness")
	pLockState := lock.String("state", "", "success or failure")

	unlock := flag.NewFlagSet("unlock", flag.ExitOnError)
	pUnlockType := unlock.String("type", "", "liveness or readiness")

	switch strings.ToLower(os.Args[1]) {
	case "run":
		run.Parse(os.Args[2:])
		if *pRunType == "" {
			fmt.Println("-type must be set")
			os.Exit(1)
		}

		handleRun(ToCheckType(*pRunType))

	case "lock":
		lock.Parse(os.Args[2:])
		if *pLockType == "" || *pLockState == "" {
			fmt.Println("-type and -state must be set")
			os.Exit(1)
		}

		handleLock(ToCheckType(*pLockType), ToCheckState(*pLockState))

	case "unlock":
		unlock.Parse(os.Args[2:])
		if *pUnlockType == "" {
			fmt.Println("-type must be set")
			os.Exit(1)
		}

		handleUnlock(ToCheckType(*pUnlockType))

	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}

const livenessCommandEnvName = "LIVENESS_CHECK"
const readinessCommandEnvName = "READINESS_CHECK"

const livenessStateFileName = "/tmp/liveness.lock"
const readinessStateFileName = "/tmp/readiness.lock"

func handleRun(checkType CheckType) {
	switch checkType {
	case LIVENESS:
		envLiveLock, err := ioutil.ReadFile(livenessStateFileName)
		if err == nil {
			state := ToCheckState(string(envLiveLock))
			switch state {
			case SUCCESS:
				fmt.Println("Liveness check locked at success: 0")
				os.Exit(0)
			case FAILURE:
				fmt.Println("Liveness check locked at failure: 1")
				os.Exit(1)
			default:
				fmt.Println("Invalid value in file ", livenessStateFileName)
				os.Exit(254)
			}

		} else {
			envLiveValue := os.Getenv(livenessCommandEnvName)
			if len(envLiveValue) == 0 {
				fmt.Println("env ", livenessCommandEnvName, " must be set")
				os.Exit(1)
			}
			runCommandAndPipeExitCode(envLiveValue)
		}

	case READINESS:
		envReadyLock, err := ioutil.ReadFile(readinessStateFileName)
		if err == nil {
			state := ToCheckState(string(envReadyLock))
			switch state {
			case SUCCESS:
				fmt.Println("Readiness check locked at success: 0")
				os.Exit(0)
			case FAILURE:
				fmt.Println("Readiness check locked at failure: 1")
				os.Exit(1)
			default:
				fmt.Println("Invalid value in file ", readinessStateFileName)
				os.Exit(254)
			}

		} else {
			envReadyValue := os.Getenv(readinessCommandEnvName)
			if len(envReadyValue) == 0 {
				fmt.Println("env ", readinessCommandEnvName, " must be set")
				os.Exit(1)
			}
			runCommandAndPipeExitCode(envReadyValue)
		}
	default:
		fmt.Println("Unexpected error")
		os.Exit(255)
	}
}

func handleLock(checkType CheckType, state CheckState) {
	switch checkType {
	case LIVENESS:
		ioutil.WriteFile(livenessStateFileName, []byte(state.String()), 0644)
	case READINESS:
		ioutil.WriteFile(readinessStateFileName, []byte(state.String()), 0644)
	default:
		fmt.Println("Unexpected error")
		os.Exit(255)
	}
}

func handleUnlock(checkType CheckType) {
	switch checkType {
	case LIVENESS:
		os.Remove(livenessStateFileName)
	case READINESS:
		os.Remove(readinessStateFileName)
	default:
		fmt.Println("Unexpected error")
		os.Exit(255)
	}
}

func runCommandAndPipeExitCode(commandString string) {
	fmt.Println("Executing ", commandString)

	// Copy & Paste from https://stackoverflow.com/questions/10385551/get-exit-code-go
	cmd := exec.Command("sh", "-c", commandString)
	if err := cmd.Start(); err != nil {
		fmt.Println("cmd.Start: ", err)
		os.Exit(2)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				fmt.Println("FAILURE Piping Exit Status: ", status.ExitStatus())
				os.Exit(status.ExitStatus())
			}
		} else {
			fmt.Println("cmd.Wait: ", err)
			os.Exit(2)
		}
	}
	fmt.Println("SUCCESS Piping Exit Status: ", 0)
	os.Exit(0)
}
