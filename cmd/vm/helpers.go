package vm

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func resolveVMArg(args []string) (string, error) {
	if len(args) > 0 && args[0] != "" {
		return args[0], nil
	}
	if v := os.Getenv("GOCT_VM"); v != "" {
		return v, nil
	}
	if stat, _ := os.Stdin.Stat(); stat.Mode()&os.ModeCharDevice == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			if line := strings.TrimSpace(scanner.Text()); line != "" {
				return line, nil
			}
		}
	}
	return "", errors.New("VM not specified: use positional arg, set GOCT_VM, or pipe via stdin")
}

func parseSize(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("size cannot be empty")
	}
	mult := int64(1)
	if len(s) > 1 {
		switch s[len(s)-1] {
		case 'g', 'G':
			mult = 1 << 30
			s = s[:len(s)-1]
		case 'm', 'M':
			mult = 1 << 20
			s = s[:len(s)-1]
		case 't', 'T':
			mult = 1 << 40
			s = s[:len(s)-1]
		}
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return v * mult, nil
}

func parseIndex(s string, out *int32) (bool, error) {
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return false, fmt.Errorf("invalid index %q: %w", s, err)
	}
	*out = int32(v)
	return true, nil
}
