package compute

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

type fakeStorage struct {
	setErr, getErr, delErr error
	getValue               string

	setCalls int
	lastSet  [2]string

	getCalls int
	lastGet  string

	delCalls int
	lastDel  string
}

func (f *fakeStorage) Set(key, value string) error {
	f.setCalls++
	f.lastSet = [2]string{key, value}
	return f.setErr
}
func (f *fakeStorage) Get(key string) (string, error) {
	f.getCalls++
	f.lastGet = key
	if f.getErr != nil {
		return "", f.getErr
	}
	return f.getValue, nil
}
func (f *fakeStorage) Delete(key string) error {
	f.delCalls++
	f.lastDel = key
	return f.delErr
}

func captureStdout(fn func()) string {
	prev := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		_ = w.Close()
		os.Stdout = prev
	}()

	fn()

	_ = w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()
	return buf.String()
}

func TestComputer_Execute(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		fake            *fakeStorage
		wantOut         string
		wantErrContains string
		assertCalls     func(t *testing.T, f *fakeStorage)
	}{
		{
			name:  "set ok",
			input: "SET a b",
			fake:  &fakeStorage{},
			assertCalls: func(t *testing.T, f *fakeStorage) {
				if f.setCalls != 1 || f.lastSet != [2]string{"a", "b"} {
					t.Fatalf("unexpected set calls: calls=%d last=%v", f.setCalls, f.lastSet)
				}
				if f.getCalls != 0 || f.delCalls != 0 {
					t.Fatalf("unexpected other calls: get=%d del=%d", f.getCalls, f.delCalls)
				}
			},
		},
		{
			name:            "set error wraps",
			input:           "SET a b",
			fake:            &fakeStorage{setErr: errors.New("boom")},
			wantErrContains: "set \"a\": boom",
		},
		{
			name:    "get ok prints value",
			input:   "GET k",
			fake:    &fakeStorage{getValue: "val"},
			wantOut: "val\n",
			assertCalls: func(t *testing.T, f *fakeStorage) {
				if f.getCalls != 1 || f.lastGet != "k" {
					t.Fatalf("unexpected get calls: calls=%d last=%q", f.getCalls, f.lastGet)
				}
				if f.setCalls != 0 || f.delCalls != 0 {
					t.Fatalf("unexpected other calls: set=%d del=%d", f.setCalls, f.delCalls)
				}
			},
		},
		{
			name:            "get error wraps",
			input:           "GET x",
			fake:            &fakeStorage{getErr: errors.New("nope")},
			wantErrContains: "get \"x\": nope",
		},
		{
			name:  "del ok",
			input: "DEL d",
			fake:  &fakeStorage{},
			assertCalls: func(t *testing.T, f *fakeStorage) {
				if f.delCalls != 1 || f.lastDel != "d" {
					t.Fatalf("unexpected del calls: calls=%d last=%q", f.delCalls, f.lastDel)
				}
				if f.setCalls != 0 || f.getCalls != 0 {
					t.Fatalf("unexpected other calls: set=%d get=%d", f.setCalls, f.getCalls)
				}
			},
		},
		{
			name:            "del error wraps",
			input:           "DEL d",
			fake:            &fakeStorage{delErr: errors.New("missing")},
			wantErrContains: "del \"d\": missing",
		},
		{
			name:            "parse error surfaced",
			input:           "UNKNOWN",
			fake:            &fakeStorage{},
			wantErrContains: "unknown command",
			assertCalls: func(t *testing.T, f *fakeStorage) {
				if f.setCalls+f.getCalls+f.delCalls != 0 {
					t.Fatalf("storage should not be called on parse error")
				}
			},
		},
		{
			name:    "get with tabs/whitespace",
			input:   "GET\tKEY1",
			fake:    &fakeStorage{getValue: "X"},
			wantOut: "X\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := NewComputer(tt.fake)

			var out string
			var err error
			if tt.wantOut != "" {
				out = captureStdout(func() { err = c.Execute(tt.input) })
			} else {
				err = c.Execute(tt.input)
			}

			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErrContains, err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantOut != "" && out != tt.wantOut {
				t.Fatalf("stdout mismatch:\nwant: %q\n got: %q", tt.wantOut, out)
			}

			if tt.assertCalls != nil {
				tt.assertCalls(t, tt.fake)
			}
		})
	}
}
