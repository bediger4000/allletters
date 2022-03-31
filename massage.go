package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	inMatch := false
	words := make([]string, 5)
	var idx int

	for scanner.Scan() {
		line := scanner.Text()
		strings.TrimSpace(line)
		if strings.HasPrefix(line, "letters:") {
			inMatch = true
			continue
		}
		if line == "" {
			sort.Strings(words)
			for i := 0; i < 5; i++ {
				fmt.Printf("%s ", words[i])
			}
			fmt.Println()
			inMatch = false
			idx = 0
		}
		if inMatch {
			words[idx] = line
			idx++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
