package suite

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	Binary  string
	TmpHome string
	cleanDB []byte

	testHomes sync.Map // keyed by test function name → workspace path
)

func init() {}

func binaryName() string {
	if runtime.GOOS == "windows" {
		return "yolo-test.exe"
	}
	return "yolo-test"
}

const (
	cfgDir  = ".yolo-team"
	cfgFile = "settings.json"
)

func workspaceDir() string {
	cwd, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(cwd, "yolo-cli", "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			return filepath.Join(os.TempDir(), "yolo-test-workspace")
		}
		cwd = parent
	}
	return filepath.Join(cwd, "yolo-cli-test", "yolo-cli-test-workspace")
}

func Setup() {
	bin := filepath.Join(workspaceDir(), binaryName())

	cwd, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(cwd, "yolo-cli", "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			fmt.Fprintln(os.Stderr, "cannot find yolo-cli directory")
			os.Exit(1)
		}
		cwd = parent
	}

	cmd := exec.Command("go", "build", "-o", bin, "./cmd/yolo")
	cmd.Dir = filepath.Join(cwd, "yolo-cli")
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %s\n%s\n", err, out)
		os.Exit(1)
	}
	Binary = bin

	TmpHome = filepath.Join(workspaceDir(), "home")
	os.RemoveAll(TmpHome)
	os.MkdirAll(TmpHome, 0755)

	ws := filepath.Join(TmpHome, "workspace")
	os.MkdirAll(ws, 0755)
	os.MkdirAll(filepath.Join(ws, "documents"), 0755)
	writeSettings(ws)

	// Trigger DB creation once for the clean template
	runBinary("--json", "project", "list")
	cleanDB, _ = os.ReadFile(filepath.Join(ws, "yolo.db"))
}

func Teardown() {
	testHomes.Range(func(key, value interface{}) bool {
		os.RemoveAll(value.(string))
		return true
	})
	os.Remove(Binary)
	os.RemoveAll(TmpHome)
}

// ensureWorkspace lazily creates an isolated workspace per test function.
// Scans the call stack for the calling test function name.
func ensureWorkspace() string {
	key := "unknown"
	for depth := 2; depth <= 8; depth++ {
		pc, _, _, ok := runtime.Caller(depth)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		name := fn.Name()
		if strings.Contains(name, ".Test") || strings.Contains(name, ".Benchmark") {
			key = name
			break
		}
	}

	if ws, ok := testHomes.Load(key); ok {
		writeSettings(ws.(string))
		return ws.(string)
	}

	ws := filepath.Join(workspaceDir(), "tests", filepath.Base(key))
	os.RemoveAll(ws)
	os.MkdirAll(ws, 0755)
	os.MkdirAll(filepath.Join(ws, "documents"), 0755)
	if cleanDB != nil {
		os.WriteFile(filepath.Join(ws, "yolo.db"), cleanDB, 0644)
	}
	testHomes.Store(key, ws)
	writeSettings(ws)
	return ws
}

func Run(args ...string) string {
	return runBinary(args...)
}

func RunFail(args ...string) string {
	cmd := exec.Command(Binary, args...)
	cmd.Env = append(os.Environ(),
		"USERPROFILE="+TmpHome,
		"HOME="+TmpHome,
	)
	out, _ := cmd.CombinedOutput()
	return string(out)
}

func runBinary(args ...string) string {
	ws := ensureWorkspace()
	cmd := exec.Command(Binary, args...)
	cmd.Env = append(os.Environ(),
		"USERPROFILE="+TmpHome,
		"HOME="+TmpHome,
	)
	_ = ws // workspace is in effect via settings.json already written to TmpHome
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("yolo %v: %s\n%s", args, err, out))
	}
	return string(out)
}

func JSONDecode(t *testing.T, raw string) []map[string]interface{} {
	t.Helper()
	jsonStr := extractJSON(raw)
	var wrapper struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	require.NoError(t, json.Unmarshal([]byte(jsonStr), &wrapper), jsonStr)
	require.Equal(t, 0, wrapper.Code, "code != 0: %s", jsonStr)
	var items []map[string]interface{}
	require.NoError(t, json.Unmarshal(wrapper.Data, &items))
	return items
}

func extractJSON(raw string) string {
	start := strings.Index(raw, "{")
	if start < 0 {
		return raw
	}
	raw = raw[start:]
	depth := 0
	for i, ch := range raw {
		if ch == '{' {
			depth++
		} else if ch == '}' {
			depth--
			if depth == 0 {
				return raw[:i+1]
			}
		}
	}
	return raw
}

func writeSettings(ws string) {
	cfgPath := filepath.Join(TmpHome, cfgDir, cfgFile)
	os.MkdirAll(filepath.Dir(cfgPath), 0755)
	data, _ := json.MarshalIndent(map[string]string{"workspace": ws}, "", "  ")
	os.WriteFile(cfgPath, data, 0644)
}
