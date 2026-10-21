package agent

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Autumn-27/norma/permission"
	actool "github.com/Autumn-27/norma/tool"
)

// decodeRemoteScript extracts and decodes the base64 payload from an
// ssh-wrapped command, so tests can assert on the exact remote script.
func decodeRemoteScript(t *testing.T, wrapped string) string {
	t.Helper()
	const open = `bash -c "echo `
	i := strings.Index(wrapped, open)
	if i < 0 {
		t.Fatalf("no bash -c wrapper found in %q", wrapped)
	}
	rest := wrapped[i+len(open):]
	// The payload is standard base64 (no spaces); the wrapper continues with
	// " | base64 -d | bash" after it, so the first space ends the payload.
	end := strings.IndexByte(rest, ' ')
	if end < 0 {
		t.Fatalf("payload not terminated in %q", wrapped)
	}
	script, err := base64.StdEncoding.DecodeString(rest[:end])
	if err != nil {
		t.Fatalf("payload not valid base64: %v", err)
	}
	return string(script)
}

func TestRemoteSSHCommandShape(t *testing.T) {
	wrapped := remoteSSHCommand("macewindu@192.168.3.140", "/Users/macewindu/sw-exec", "nmap -v target")
	if !strings.HasPrefix(wrapped, "ssh macewindu@192.168.3.140 'bash -c \"") {
		t.Errorf("missing ssh host prefix: %q", wrapped)
	}
	if !strings.HasSuffix(wrapped, `"'`) {
		t.Errorf("missing closing quotes: %q", wrapped)
	}
	script := decodeRemoteScript(t, wrapped)
	wantPrefix := `export PATH="/opt/homebrew/bin:/opt/homebrew/sbin:$PATH"; cd /Users/macewindu/sw-exec && `
	if !strings.HasPrefix(script, wantPrefix) {
		t.Errorf("script prefix = %q, want prefix %q", script, wantPrefix)
	}
	if !strings.HasSuffix(script, "nmap -v target") {
		t.Errorf("script lost the original command: %q", script)
	}
}

func TestRemoteSSHCommandNoWorkDir(t *testing.T) {
	script := decodeRemoteScript(t, remoteSSHCommand("u@h", "", "ls"))
	if strings.Contains(script, "cd ") {
		t.Errorf("empty workDir must not emit a cd: %q", script)
	}
	if !strings.HasSuffix(script, "ls") {
		t.Errorf("script lost the command: %q", script)
	}
}

// The whole point of the base64 wrapper: hostile quoting, newlines and
// expansions in the model's command must survive three shell re-parses.
func TestRemoteSSHCommandSurvivesQuoting(t *testing.T) {
	nasty := "echo 'single' \"double $HOME `bt`\"\nline2\ttab; rm -rf /tmp/x && echo $'\\x41'"
	script := decodeRemoteScript(t, remoteSSHCommand("u@h", "/tmp/wd", nasty))
	if !strings.HasSuffix(script, nasty) {
		t.Errorf("round-trip mismatch:\n got suffix %q\nwant suffix %q", script, nasty)
	}
}

func TestWrapBashRemoteLocalPassthrough(t *testing.T) {
	tools := []actool.CoreTool{fakeTool("Bash"), fakeTool("Read")}
	out := wrapBashRemote(tools, "", "")
	if len(out) != len(tools) || out[0] != tools[0] || out[1] != tools[1] {
		t.Errorf("empty host must return the tools unchanged")
	}
}

func TestWrapBashRemoteOnlyBashWrapped(t *testing.T) {
	read := fakeTool("Read")
	bash := fakeTool("Bash")
	out := wrapBashRemote([]actool.CoreTool{bash, read}, "u@h", "/wd")
	if out[1] != read {
		t.Errorf("non-Bash tools must pass through untouched")
	}
	if out[0] == bash {
		t.Errorf("Bash tool must be replaced by the wrapper")
	}
	if out[0].Name() != "Bash" {
		t.Errorf("wrapped tool must keep the name Bash, got %q", out[0].Name())
	}
	if !strings.Contains(out[0].Description(), "u@h") || !strings.Contains(out[0].Description(), "/wd") {
		t.Errorf("wrapped description must disclose host and cwd: %q", out[0].Description())
	}
}

func TestNewRemoteBashRewritesCommand(t *testing.T) {
	var captured string
	inner := actool.Build(actool.Spec{
		Name:        "Bash",
		Description: "bash",
		Schema:      map[string]any{},
		Run: func(ctx context.Context, input json.RawMessage, tc *actool.ToolContext) (actool.Result, error) {
			captured = string(input)
			return actool.Text("ok"), nil
		},
	})
	wrapped := newRemoteBash(inner, "macewindu@192.168.3.140", "/Users/macewindu/sw-exec")
	res, err := wrapped.Call(context.Background(), json.RawMessage(`{"command":"sw-egress curl -s https://example.invalid","timeout":30}`), &actool.ToolContext{})
	if err != nil {
		t.Fatalf("call failed: %v", err)
	}
	if res.Flatten() != "ok" {
		t.Errorf("inner result not propagated: %+v", res)
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(captured), &m); err != nil {
		t.Fatalf("inner input not valid JSON: %v", err)
	}
	cmd, _ := m["command"].(string)
	if !strings.HasPrefix(cmd, "ssh macewindu@192.168.3.140 ") {
		t.Errorf("inner command not ssh-rewritten: %q", cmd)
	}
	script := decodeRemoteScript(t, cmd)
	if !strings.HasSuffix(script, "sw-egress curl -s https://example.invalid") {
		t.Errorf("remote script lost the original command: %q", script)
	}
	if got, _ := m["timeout"].(float64); got != 30 {
		t.Errorf("unrelated input fields must round-trip, timeout = %v", m["timeout"])
	}
}

func TestNewRemoteBashPassthroughMalformed(t *testing.T) {
	for _, in := range []string{
		`not json`,
		`{"timeout":5}`,
		`{"command":42}`,
		`{"command":null}`,
	} {
		var captured string
		inner := actool.Build(actool.Spec{
			Name:        "Bash",
			Description: "bash",
			Schema:      map[string]any{},
			Run: func(ctx context.Context, input json.RawMessage, tc *actool.ToolContext) (actool.Result, error) {
				captured = string(input)
				return actool.Text("inner"), nil
			},
		})
		wrapped := newRemoteBash(inner, "u@h", "/wd")
		if _, err := wrapped.Call(context.Background(), json.RawMessage(in), &actool.ToolContext{}); err != nil {
			t.Fatalf("input %q: call failed: %v", in, err)
		}
		if captured != in {
			t.Errorf("input %q: inner saw %q, want byte-identical passthrough", in, captured)
		}
	}
}

// The catastrophic-command floor must evaluate the model's original command,
// not the ssh wrapper around it.
func TestNewRemoteBashPermissionsSeesOriginalCommand(t *testing.T) {
	var seen string
	inner := actool.Build(actool.Spec{
		Name:        "Bash",
		Description: "bash",
		Schema:      map[string]any{},
		Permissions: func(ctx context.Context, input json.RawMessage, pc permission.Context) permission.Decision {
			seen = string(input)
			return permission.Denied("test deny")
		},
		Run: func(ctx context.Context, input json.RawMessage, tc *actool.ToolContext) (actool.Result, error) {
			return actool.Text("should not run"), nil
		},
	})
	wrapped := newRemoteBash(inner, "u@h", "")
	original := `{"command":"mkfs /dev/sda1"}`
	if d := wrapped.CheckPermissions(context.Background(), json.RawMessage(original), permission.Context{}); d.Behavior != permission.Deny {
		t.Errorf("inner permission decision not propagated: %+v", d)
	}
	if seen != original {
		t.Errorf("CheckPermissions saw %q, want the original unwrapped command", seen)
	}
}

func fakeTool(name string) actool.CoreTool {
	return actool.Build(actool.Spec{
		Name:        name,
		Description: name + " tool",
		Schema:      map[string]any{},
		Run: func(ctx context.Context, input json.RawMessage, tc *actool.ToolContext) (actool.Result, error) {
			return actool.Text(name), nil
		},
	})
}
