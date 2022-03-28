package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

type singleWord struct {
	word    string
	letters [5]byte
	hash    uint32
}

type guess struct {
	words   [5]*singleWord
	depth   int
	letters [26]byte
}

func main() {
	dictionaryFileName := flag.String("i", "", "5-letter words file name (default stdin)")
	flag.Parse()

	words := readDictionary(*dictionaryFileName)
	fmt.Fprintf(os.Stderr, "Found %d 5-letter words\n", len(words))

	ws := covertToWords(words)

	m := make(map[uint32]int)
	for _, s := range ws {
		m[s.hash]++
	}

	uniqueHashCount := 0
	for hash, count := range m {
		if count > 1 {
			fmt.Printf("%x count %d\n", hash, count)
			continue
		}
		uniqueHashCount++
	}
	fmt.Printf("Found %d unique hashed words\n", uniqueHashCount)
}

func readDictionary(dictionaryFileName string) []string {
	fin := os.Stdin
	if dictionaryFileName != "" {
		var err error
		fin, err = os.Open(dictionaryFileName)
		if err != nil {
			log.Fatal(err)
		}
		defer fin.Close()
	}

	scanner := bufio.NewScanner(fin)

	words := make([]string, 0, 8000)

	lineCount := 0
	for scanner.Scan() {
		word := scanner.Text()
		lineCount++
		if len(word) != 5 {
			log.Printf("line %d, %q not length 5\n", lineCount, word)
			continue
		}
		words = append(words, word)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return words
}

func covertToWords(words []string) []*singleWord {
	singles := make([]*singleWord, len(words))

	for i := range words {
		singles[i] = &singleWord{
			word: words[i],
		}
		var hash uint32 = 5381
		for j := 0; j < 5; j++ {
			singles[i].letters[j] = words[i][j]
			hash = ((hash << 5) + hash) + uint32(words[i][j])
		}
		singles[i].hash = hash
	}

	return singles
}
