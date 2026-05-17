// Package shell locates command-line executables on the user's machine.
//
// It searches PATH first, then a caller-supplied list of well-known install
// locations. Use it whenever an application has to wrap a third-party CLI and
// the binary is not guaranteed to be on PATH.
package shell
