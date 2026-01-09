package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const agentInstructions = `
## grepai - Semantic Code Search

This project uses grepai for semantic code search. Instead of guessing file locations
or grepping for exact text, use grepai to find relevant code by intent.

### Usage

Search for code semantically:
` + "```bash" + `
grepai search "function that handles user authentication"
grepai search "error handling for API requests"
grepai search "database connection setup"
` + "```" + `

### Tips

- Use natural language queries describing what the code does
- Results include file paths, line numbers, and relevance scores
- The index is maintained in real-time when ` + "`grepai watch`" + ` is running

`

const agentMarker = "## grepai - Semantic Code Search"

var agentSetupCmd = &cobra.Command{
	Use:   "agent-setup",
	Short: "Configure AI agents to use grepai",
	Long: `Configure AI agent environments to leverage grepai for context retrieval.

This command will:
- Detect agent configuration files (.cursorrules, CLAUDE.md)
- Append instructions for using grepai search
- Ensure idempotence (won't add duplicate instructions)`,
	RunE: runAgentSetup,
}

func runAgentSetup(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	agentFiles := []string{
		".cursorrules",
		"CLAUDE.md",
		".claude/settings.md",
	}

	found := false
	modified := 0

	for _, file := range agentFiles {
		path := filepath.Join(cwd, file)

		// Check if file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		found = true
		fmt.Printf("Found: %s\n", file)

		// Read existing content
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("  Warning: could not read %s: %v\n", file, err)
			continue
		}

		// Check if already configured
		if strings.Contains(string(content), agentMarker) {
			fmt.Printf("  Already configured, skipping\n")
			continue
		}

		// Append instructions
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("  Warning: could not open %s for writing: %v\n", file, err)
			continue
		}

		// Add newlines if needed
		if len(content) > 0 && content[len(content)-1] != '\n' {
			f.WriteString("\n")
		}
		f.WriteString("\n")
		f.WriteString(agentInstructions)
		f.Close()

		fmt.Printf("  Added grepai instructions\n")
		modified++
	}

	if !found {
		fmt.Println("No agent configuration files found.")
		fmt.Println("\nSupported files:")
		for _, file := range agentFiles {
			fmt.Printf("  - %s\n", file)
		}
		fmt.Println("\nCreate one of these files and run 'grepai agent-setup' again,")
		fmt.Println("or manually add instructions for using 'grepai search'.")
		return nil
	}

	if modified > 0 {
		fmt.Printf("\nUpdated %d file(s).\n", modified)
	} else {
		fmt.Println("\nAll files already configured.")
	}

	return nil
}
