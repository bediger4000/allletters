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
	key     string
	word    string
	letters []rune
}

type guess struct {
	wordMap map[string]*singleWord
	words   [5]*singleWord
	depth   int
	letters [26]int
	debug   bool
}

func main() {
	dictionaryFileName := flag.String("i", "", "5-letter words file name (default stdin)")
	debug := flag.Bool("d", false, "debug output")
	flag.Parse()

	if *dictionaryFileName == "" {
		fmt.Fprintf(os.Stderr, "Need dictionary file name, -i <filename>\n")
		os.Exit(1)
	}

	words := readDictionary(*dictionaryFileName)
	fmt.Fprintf(os.Stderr, "Found %d 5-letter words\n", len(words))

	wm, uniqueWords := convertWords(words)
	fmt.Fprintf(os.Stderr, "Made %d unique-key words\n", len(wm))
	if *debug {
		for key, sw := range wm {
			fmt.Fprintf(os.Stderr, "key %q: word %q, letters %#c\n", key, sw.word, sw.letters)
		}
	}

	startGuessing(wm, uniqueWords, *debug)
}

func startGuessing(wm map[string]*singleWord, words []*singleWord, debug bool) {
	g := &guess{wordMap: wm, debug: debug}
	for i := range words {
		fmt.Fprintf(os.Stderr, "word: %s\n", words[i].word)
		/*
			type guess struct {
				words   [5]*singleWord
				depth   int
				letters [26]int
			}
		*/
		g.markLetters(words[i])
		nextGuess(g)
		g.unmarkLetters()
		delete(g.wordMap, words[i].key)
	}

	fmt.Fprintf(os.Stderr, "word map length: %d\n", len(g.wordMap))
}

func nextGuess(g *guess) {
	if g.debug {
		fmt.Fprintf(os.Stderr, "enter nextGuess\n")
		fmt.Fprintf(os.Stderr, "guess depth %d\n", g.depth)
		fmt.Fprintf(os.Stderr, "Letters: ")
		for i := 0; i < 26; i++ {
			if g.letters[i] == 0 {
				continue
			}
			fmt.Fprintf(os.Stderr, "%c", i+'a')
		}
		fmt.Fprintf(os.Stderr, "\n")
	}

	if g.depth > 4 {
		lettersUsed := 0
		for i := 0; i < 26; i++ {
			if g.letters[i] > 0 {
				lettersUsed++
			}
		}
		if lettersUsed > 23 {
			var letters []rune
			for i := 'a'; i <= 'z'; i++ {
				if g.letters[i-'a'] > 0 {
					letters = append(letters, i)
				}
			}
			fmt.Printf("guess %q:", string(letters))
			for i := range g.words {
				fmt.Printf(" %s", g.words[i].word)
			}
			fmt.Println()
		}
		return
	}

	var keyrunes [5]rune

	for i := 0; i < 26; i++ {
		if g.letters[i] > 0 {
			continue
		}
		for j := i + 1; j < 26; j++ {
			if g.letters[j] > 0 {
				continue
			}
			for k := j + 1; k < 26; k++ {
				if g.letters[k] > 0 {
					continue
				}
				for l := k + 1; l < 26; l++ {
					if g.letters[l] > 0 {
						continue
					}
					for m := l + 1; m < 26; m++ {
						if g.letters[m] > 0 {
							continue
						}

						keyrunes[0] = rune('a' + i)
						keyrunes[1] = rune('a' + j)
						keyrunes[2] = rune('a' + k)
						keyrunes[3] = rune('a' + l)
						keyrunes[4] = rune('a' + m)
						key := string(keyrunes[:])
						if candidate, ok := g.wordMap[key]; ok {
							if g.debug {
								fmt.Fprintf(os.Stderr, "key %s, word %s works\n", key, candidate.word)
							}
							g.markLetters(candidate)
							nextGuess(g)
							g.unmarkLetters()
						}
					}
				}
			}
		}
	}
}

func (gp *guess) reset() {
	gp.depth = 0
	for i := range gp.letters {
		gp.letters[i] = 0
	}
	for i := range gp.words {
		gp.words[i] = nil
	}
}

func (gp *guess) markLetters(w *singleWord) {
	gp.words[gp.depth] = w
	for i := 0; i < 5; i++ {
		gp.letters[w.letters[i]-'a']++
	}
	gp.depth++
}

func (gp *guess) unmarkLetters() {
	gp.depth--
	w := gp.words[gp.depth]
	gp.words[gp.depth] = nil
	for i := 0; i < 5; i++ {
		gp.letters[w.letters[i]-'a']--
	}
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

func convertWords(words []string) (map[string]*singleWord, []*singleWord) {

	wordMap := make(map[string]*singleWord)

	for i := range words {
		s := &singleWord{
			word: words[i],
		}

		runes := []rune(s.word)
		sort.Sort(RuneSlice(runes))
		s.key = string(runes)

		s.letters = runes

		wordMap[s.key] = s
	}

	uw := make([]*singleWord, len(wordMap))
	idx := 0
	for _, sw := range wordMap {
		uw[idx] = sw
		idx++
	}
	sort.Sort(wordsArray(uw))

	return wordMap, uw
}

type wordsArray []*singleWord

func (wa wordsArray) Len() int           { return len(wa) }
func (wa wordsArray) Less(i, j int) bool { return wa[i].word < wa[j].word }
func (wa wordsArray) Swap(i, j int)      { wa[i], wa[j] = wa[j], wa[i] }

type RuneSlice []rune

func (rs RuneSlice) Len() int           { return len(rs) }
func (rs RuneSlice) Less(i, j int) bool { return rs[i] < rs[j] }
func (rs RuneSlice) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }