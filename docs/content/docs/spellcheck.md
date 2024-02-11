---
title: "Spellcheck"
weight: 1
# bookFlatSection: false
bookToc: true
bookHidden: false
# bookCollapseSection: false
# bookComments: true
---

# Consistent spelling

## What is consistent spelling?

You can use LinTeX to validate that your spelling is consistent throughout your
document. This is done by defining a list of words with their correct spelling, as well
as a regular expression, that captures all alternative spellings.

{{< hint warning >}}

**LinTeX is not a spell checker!**

Do not confuse this with a regular spell checker, that checks all words for orthography!
Use another tool for this, e.g. [LTeX](https://valentjn.github.io/ltex/).

{{< /hint >}}

As an example we'll use the word `colour`, which is well known to be spelled `color` in
American English, but `colour` in British English.\* In a paper that's supposed to be
written in British English, we would want to make sure to always use `colour` and never
`color`.

\* This is not a great example, because ideally you'll have your actual spell checker
configured to either American or British English and identify this for you, but it's
just an example.

## Configuration

First, make sure that your configuration file exists as `.lintex/config.toml`. If it
doesn't, create an empty file at this location.

To activate spell checking, add a list under the key `spellcheck`:

```toml
spellchecks = [
	{ correct = "colour", regex = "colou?r" },
]
```

Each list element is a mapping with two keys: `correct` denotes the correct spelling,
that is supposed to be used; `regex` denotes a pattern, that matches all possible
spellings that might be used.

After saving the file, LinTeX will check for the correct spelling on the next run. To
test it, let's create a LaTeX file, that contains both spellings.

```tex
color colour
```

Running LinTeX on this file shows

```shell-session
$ lintex
test.tex:1:0
spelling/colour: Spelling: colour
   1: color colour
   2:
Check that 'colour' is spelled correct.
```
