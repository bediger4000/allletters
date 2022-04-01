# Solve Wordle the Interesting Way

### Wordle Solving Method

Make 5, 5-letter guesses,
where 25 different letters appear.
Leave out 'Z', 'J', 'B' or 'Q'.

Based on the letters remaining, solve the Wordle puzzle.

---

## Are there 5, 5-letter words that contain 25 distinct letters?

This repo attempts to solve that problem.

## Dictionary of 5-letter words

Under Arch Linux, relatively up-to-date on 2022-03-27,
`/usr/share/dict/words` is a symbolic link to
`/usr/share/dict/american-english`.
`/usr/share/dict/american-english` is owned by package "words 2.1-6".


```sh
tr -cd '[A-Za-z\n]' < /usr/share/dict/words | tr '[A-Z]' '[a-z]' |
grep '^.....$' | sort | uniq > words.5
```
7393 unique 5-letter words exist in that dictionary.
Some of these words have duplicate letters ("teeth"),
some don't make sense: this dictionary has common Roman Numeral strings,
like "clxiv",
because this dictionary is intended for use in spell checking.
Why alert on Roman Numerals used to denote page order of introductions?

This dictionary is different than the dictionary used internally by
the Wordle JavaScript app.

## Solutions

### Using /usr/share/dict/words

```
benzs   clxiv   fjord   gawky   thump 
benzs   clxvi   fjord   gawky   thump 
bjork   clxiv   fazed   gwyns   thump 
bjork   clxvi   fazed   gwyns   thump 
bumph   frock   gyved   jinxs   waltz 
```

"clxiv" and "clxvi" are just Roman Numerals.

Wordle-the-app doesn't allow "jinxs" as a guess.

### Using Wordle Dictionary

```
bemix waqfs clunk vozhd grypt
bling waqfs treck vozhd jumpy
blunk waqfs cimex vozhd grypt
brick waqfs vozhd glent jumpy
brung waqfs xylic vozhd kempt
chunk waltz vibex fjord gymps
clipt waqfs jumby vozhd kreng
fjord waltz vibex gucks nymph
glent waqfs jumby prick vozhd
jumby waqfs treck vozhd pling
```

