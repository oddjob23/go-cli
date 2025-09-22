package utils

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	Success = color.New(color.FgGreen).SprintFunc()
	Error   = color.New(color.FgRed).SprintFunc()
	Warning = color.New(color.FgYellow).SprintFunc()
	Info    = color.New(color.FgCyan).SprintFunc()
	Bold    = color.New(color.Bold).SprintFunc()
	Gray    = color.New(color.FgHiBlack).SprintFunc()
)

// CliOutput provides a clean CLI output interface
type CliOutput struct {
	verbose bool
}

// NewCliOutput creates a new CLI output handler
func NewCliOutput(verbose bool) *CliOutput {
	return &CliOutput{
		verbose: verbose,
	}
}

// Legacy functions for backward compatibility
func PrintSuccess(format string, args ...interface{}) {
	fmt.Printf("✅ %s\n", Success(fmt.Sprintf(format, args...)))
}

func PrintError(format string, args ...interface{}) {
	fmt.Printf("❌ %s\n", Error(fmt.Sprintf(format, args...)))
}

func PrintWarning(format string, args ...interface{}) {
	fmt.Printf("⚠️  %s\n", Warning(fmt.Sprintf(format, args...)))
}

func PrintInfo(format string, args ...interface{}) {
	fmt.Printf("ℹ️  %s\n", Info(fmt.Sprintf(format, args...)))
}

func PrintBold(format string, args ...interface{}) {
	fmt.Printf("%s\n", Bold(fmt.Sprintf(format, args...)))
}

// Simple CLI output methods
func (c *CliOutput) Info(format string, args ...interface{}) {
	fmt.Printf("ℹ️  %s\n", Info(fmt.Sprintf(format, args...)))
}

func (c *CliOutput) Success(format string, args ...interface{}) {
	fmt.Printf("✅ %s\n", Success(fmt.Sprintf(format, args...)))
}

func (c *CliOutput) Warning(format string, args ...interface{}) {
	fmt.Printf("⚠️  %s\n", Warning(fmt.Sprintf(format, args...)))
}

func (c *CliOutput) Error(format string, args ...interface{}) {
	fmt.Printf("❌ %s\n", Error(fmt.Sprintf(format, args...)))
}

func (c *CliOutput) Debug(format string, args ...interface{}) {
	if c.verbose {
		fmt.Printf("🔍 %s\n", Gray(fmt.Sprintf(format, args...)))
	}
}

func (c *CliOutput) Plain(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func (c *CliOutput) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
