package shell

import (
	"errors"
	"os"
	"os/exec"
)

var ErrBinaryNotFound = errors.New("shell: binary not found")

type ResolutionSource int

const (
	SourceUnknown ResolutionSource = iota
	SourcePath
	SourceCandidate
)

func (s ResolutionSource) String() string {
	switch s {
	case SourcePath:
		return "path"
	case SourceCandidate:
		return "candidate"
	default:
		return "unknown"
	}
}

type ResolvedExecutable struct {
	absolutePath string
	source       ResolutionSource
}

func (r ResolvedExecutable) AbsolutePath() string {
	return r.absolutePath
}

func (r ResolvedExecutable) Source() ResolutionSource {
	return r.source
}

type Candidates struct {
	names      []string
	candidates []string
}

func NewCandidates() Candidates {
	return Candidates{}
}

func (c Candidates) WithName(name string) Candidates {
	if name == "" {
		return c
	}
	c.names = append(c.names, name)
	return c
}

func (c Candidates) WithCandidate(path string) Candidates {
	if path == "" {
		return c
	}
	c.candidates = append(c.candidates, path)
	return c
}

func (c Candidates) Resolve() (ResolvedExecutable, error) {
	for _, name := range c.names {
		if absolute, err := exec.LookPath(name); err == nil {
			return ResolvedExecutable{absolutePath: absolute, source: SourcePath}, nil
		}
	}
	for _, candidate := range c.candidates {
		if isExecutableFile(candidate) {
			return ResolvedExecutable{absolutePath: candidate, source: SourceCandidate}, nil
		}
	}
	return ResolvedExecutable{}, ErrBinaryNotFound
}

func isExecutableFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}
