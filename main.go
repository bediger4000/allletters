package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
)

type singleWord struct {
	word    string
	letters [5]byte
	hash    uint32
}

type wordsArray []*singleWord

type guess struct {
	words   [5]*singleWord
	depth   int
	letters [26]int
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
	fmt.Fprintf(os.Stderr, "Found %d unique hashed words\n", uniqueHashCount)

	byLetter := organizeWordsByLetter(ws)

	for idx := 0; idx < 26; idx++ {
		fmt.Fprintf(os.Stderr, "%c\t%d\n", 'a'+idx, len(byLetter[idx]))
	}

	startGuessing(ws, byLetter)
}

func startGuessing(ws []*singleWord, byLetter [][]*singleWord) {
	for i := range ws {
		g := &guess{}
		g.words[0] = ws[i]
		for _, l := range ws[i].letters {
			g.letters[l-'a'] = 1
		}
		//fmt.Printf("word: %q\n", ws[i].word)
		nextGuess(g, byLetter)
	}
}

func nextGuess(g *guess, byLetter [][]*singleWord) {
	if g.depth > 3 {
		lettersUsed := 0
		for i := 0; i < 26; i++ {
			if g.letters[i] > 0 {
				lettersUsed++
			}
		}
		/*
			fmt.Printf("Used %d letters:\n", lettersUsed)
			for i := range g.words {
				fmt.Printf("\t%s\n", g.words[i].word)
			}
		*/
		if lettersUsed > 23 {
			var letters []rune
			for i := 'a'; i <= 'z'; i++ {
				if g.letters[i-'a'] > 0 {
					letters = append(letters, i)
				}
			}
			fmt.Printf("letters: %q\n", string(letters))
			for i := range g.words {
				fmt.Printf("\t%s\n", g.words[i].word)
			}
			fmt.Println()
		}
		return
	}

	for i := 0; i < 26; i++ {
		if g.letters[i] > 0 {
			continue
		}

		// Find a candidate additional word for this letter
		possibleWords := byLetter[i]
		for j := range possibleWords {

			candidate := possibleWords[j]
			foundOne := true

			for _, letter := range candidate.letters {
				if g.letters[letter-'a'] > 0 {
					foundOne = false
					break
				}
			}

			if !foundOne {
				continue
			}

			// add candidate additional word to current working guess
			g.depth++
			// fmt.Printf("candidate[%d/%c] @%d: %q\n", j, 'a'+i, g.depth, candidate.word)
			g.words[g.depth] = candidate
			for k := 0; k < 5; k++ {
				letter := candidate.letters[k]
				g.letters[letter-'a']++
			}

			// recursively call nextGuess
			nextGuess(g, byLetter)

			// remove candidate word from current working guess
			g.words[g.depth] = nil
			for k := 0; k < 5; k++ {
				letter := candidate.letters[k]
				g.letters[letter-'a']--
			}
			g.depth--
		}
	}
}

func organizeWordsByLetter(ws []*singleWord) [][]*singleWord {
	byLetter := make([][]*singleWord, 26)
	for i := range ws {
		for j := 0; j < 5; j++ {
			idx := ws[i].letters[j] - 'a'
			byLetter[idx] = append(byLetter[idx], ws[i])
		}
	}
	for i := range byLetter {
		sort.Sort(wordsArray(byLetter[i]))
	}
	return byLetter
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
		runes := []rune(word)
		foundDuplicateLetter := false
	DUPES:
		for i := 0; i < 4; i++ {
			for j := i + 1; j < 5; j++ {
				if runes[i] == runes[j] {
					foundDuplicateLetter = true
					break DUPES
				}
			}
		}
		if foundDuplicateLetter {
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

func (wa wordsArray) Len() int           { return len(wa) }
func (wa wordsArray) Less(i, j int) bool { return wa[i].word < wa[j].word }
func (wa wordsArray) Swap(i, j int)      { wa[i], wa[j] = wa[j], wa[i] }
