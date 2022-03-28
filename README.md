# Solve Wordle the Interesting Way

Make 5, 5-letter guesses,
where 25 different letters appear.
Leave out 'Z' or 'B' or 'Q'.
Which one?

Based on the letters remaining, solve the Wordle puzzle.

Are there 5, 5-letter words that contain 25 distinct letters?

This repo attempts to solve that problem.

Under Arch Linux, relatively up-to-date on 2022-02-27,
`/usr/share/dict/words` is a symbolic link to
`/usr/share/dict/american-english`.
`/usr/share/dict/american-english` is owned by package "words 2.1-6".


```sh
tr -cd '[A-Za-z\n]' < /usr/share/dict/words | tr '[A-Z]' '[a-z]' |
grep '^.....$' | sort | uniq > words.5
```
7393 unique 5-letter words exist in that dictionary.
