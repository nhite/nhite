package main

import "log"

type grpcUI struct {
	stdout []byte
	stderr []byte
}

// Ask asks the user for input using the given query. The response is
// returned as the given string, or an error.
func (g *grpcUI) Ask(string) (string, error) {
	return "", nil
}

// AskSecret asks the user for input using the given query, but does not echo
// the keystrokes to the terminal.
func (g *grpcUI) AskSecret(string) (string, error) {
	return "", nil
}

// Output is called for normal standard output.
func (g *grpcUI) Output(msg string) {
	log.Println("==>", msg)
	g.stdout = append(g.stdout, []byte(msg)...)
}

// Info is called for information related to the previous output.
// In general this may be the exact same as Output, but this gives
// Ui implementors some flexibility with output formats.
func (g *grpcUI) Info(msg string) {
	g.stdout = append(g.stdout, []byte(msg)...)
}

// Error is used for any error messages that might appear on standard
// error.
func (g *grpcUI) Error(msg string) {
	g.stderr = append(g.stderr, []byte(msg)...)
}

// Warn is used for any warning messages that might appear on standard
// error.
func (g *grpcUI) Warn(msg string) {
	g.stdout = append(g.stdout, []byte(msg)...)
}
