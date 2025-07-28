/*
 * GitHubber - CLI Input Handler
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: User input handling utilities for CLI interface
 */

package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GetInput prompts the user for input and returns the trimmed response
func GetInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
