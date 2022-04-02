package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type singleWord struct {
	key     string // letters of .word, alphabetized
	word    string
	letters []rune // individual letters of .word, Unicode code points
}

type guess struct {
	wordMap map[string]*singleWord // avoid a global
	words   [5]*singleWord
	depth   int
	letters [26]int // count of letters in .words[] up to .depth-1
	debug   bool
}

// The letters[] element of struct guess is a count of letters: letters[0] is
// count of 'a' code points in the words[] array, using index up to, but not
// including value of depth element. It's 1 ahead.  wordMap is keyed by
// 5-letter strings, the alphabetized characters in the .word elements of the
// *singleWord structs that are values of wordMap.

func main() {
	dictionaryFileName := flag.String("i", "", "5-letter words file name (default stdin)")
	debug := flag.Bool("d", false, "debug output")
	allowDupes := flag.Bool("D", false, "allow duplicate letters in words")
	flag.Parse()

	if *dictionaryFileName == "" {
		fmt.Fprintf(os.Stderr, "Need dictionary file name, -i <filename>\n")
		os.Exit(1)
	}

	words := readDictionary(*dictionaryFileName, !*allowDupes)
	fmt.Fprintf(os.Stderr, "Found %d 5-letter words\n", len(words))

	wm, uniqueWords := convertWords(words)
	fmt.Fprintf(os.Stderr, "Made %d unique-key words\n", len(wm))
	if *debug {
		for key, sw := range wm {
			fmt.Fprintf(os.Stderr, "key %q: word %q, letters %#c\n", key, sw.word, sw.letters)
		}
	}

	alreadyOutput = make(map[string]bool)

	startGuessing(wm, uniqueWords, *debug)
}

func startGuessing(wm map[string]*singleWord, words []*singleWord, debug bool) {
	g := &guess{wordMap: wm, debug: debug}
	for i := range words {
		fmt.Fprintf(os.Stderr, "word: %s\n", words[i].word)
		g.markLetters(words[i])
		nextGuess(g)
		g.unmarkLetters()
		delete(g.wordMap, words[i].key)
	}

	fmt.Fprintf(os.Stderr, "word map length: %d\n", len(g.wordMap))
}

// Track 5 word sets that we've already encountered.  Keys are 5, 5-letter
// words in alphabetical order, concatenated, so it's possible to track
// wether a given 5 word set has already been discovered.
var alreadyOutput map[string]bool

// nextGuess called recursively: find all 5-letter keys that could lead to a
// word (via g.wordMap). Try them all as keys to g.wordMap.  If a word matches
// a 5-letter key, add it as a member of a potential 5-word set. Call
// nextGuess() with g and its additional word, to try to fill in another word.
// If g is at depth 5, we've found a 5-word set that has 25 different letters
// in it.
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
		// 5 words in guess g. See if this set of words has been encountered in
		// a different order.
		var wordsUsed [5]string
		for i := 0; i < 5; i++ {
			wordsUsed[i] = g.words[i].word
		}
		sort.Strings(wordsUsed[:])
		bigString := strings.Join(wordsUsed[:], "")
		if alreadyOutput[bigString] {
			return
		}
		alreadyOutput[bigString] = true

		fmt.Fprintf(os.Stderr, "found one\n")

		// output a string consisting of the 25 letters in the 5 words,
		// and those 5 words.
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
		return
	}

	var keyrunes [5]rune

	// Construct all 5-letter wordmap keys that are left for this
	// depth of guess. Because guess.letters[] is indexed by 'a':0,
	// 'b':1, 'c':2... marching through the letters with a 5-level-deep
	// set of for-loops gets all of the 5-letter keys possible for the
	// guess to use, and nothing more. This also imposes a "no duplicate
	// letters in the key" criteria.
	for i := 0; i < 26; i++ {
		if g.letters[i] > 0 {
			continue
		}
		keyrunes[0] = rune('a' + i)
		for j := i + 1; j < 26; j++ {
			if g.letters[j] > 0 {
				continue
			}
			keyrunes[1] = rune('a' + j)
			for k := j + 1; k < 26; k++ {
				if g.letters[k] > 0 {
					continue
				}
				keyrunes[2] = rune('a' + k)
				for l := k + 1; l < 26; l++ {
					if g.letters[l] > 0 {
						continue
					}
					keyrunes[3] = rune('a' + l)
					for m := l + 1; m < 26; m++ {
						if g.letters[m] > 0 {
							continue
						}
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

// markLetters pushes a word on the list of words
// maintained by a struct guess.
func (gp *guess) markLetters(w *singleWord) {
	gp.words[gp.depth] = w
	for i := 0; i < 5; i++ {
		gp.letters[w.letters[i]-'a']++
	}
	gp.depth++
}

// unmarkLetters pops the last word from the struct guess.
func (gp *guess) unmarkLetters() {
	gp.depth--
	w := gp.words[gp.depth]
	gp.words[gp.depth] = nil
	for i := 0; i < 5; i++ {
		gp.letters[w.letters[i]-'a']--
	}
}

// readDictionary reads contents of a file one line at a time.
// Throws away anything other than 5-letter words, and also
// 5-letter words that contain 2 or more of the same letter.
// "poppy" does not make it into the final slice of strings.
func readDictionary(dictionaryFileName string, weedOutDupes bool) []string {
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

	words := make([]string, 0, 8000) // 8000 word capacity

	lineCount := 0
	for scanner.Scan() {
		word := scanner.Text()
		lineCount++
		if len(word) != 5 {
			log.Printf("line %d, %q not length 5\n", lineCount, word)
			continue
		}
		if weedOutDupes {
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
		}
		words = append(words, word)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return words
}

// convertWords turns a list of strings into a map of *singleWord structs,
// and an array of *singleWord structs, both pointing to the same set
// of structs singleWord.
// The key to the map is a 5-letter string, where all the letters are in
// alphabetical order: the key for "cloud" is "cdlou" for example.
// This means that only the last word in the input list of strings that
// have a given key ends up in the map. Only one fo "team", "meat" and "meta"
// end up in the map. For my purposes, this is good enough.
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

		// last word with a given key stays in map
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

// Helper methods to sort slices-of-runes, and slices of *singleWord

type wordsArray []*singleWord

func (wa wordsArray) Len() int           { return len(wa) }
func (wa wordsArray) Less(i, j int) bool { return wa[i].word < wa[j].word }
func (wa wordsArray) Swap(i, j int)      { wa[i], wa[j] = wa[j], wa[i] }

type RuneSlice []rune

func (rs RuneSlice) Len() int           { return len(rs) }
func (rs RuneSlice) Less(i, j int) bool { return rs[i] < rs[j] }
func (rs RuneSlice) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
