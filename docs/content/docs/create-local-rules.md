---
title: "Create Local Rules"
weight: 1
# bookFlatSection: false
bookToc: true
bookHidden: false
# bookCollapseSection: false
# bookComments: true
---

# Create local rules

The goal of LinTeX is to provide its users with the ability to check their own style.
That's done by providing a mechanism, that allows users to define their own linting
rules. This page explains how you can create your own rules.

## Discovery and File Format

LinTeX searches for rules in a directory called `.lintex/rules` relative to its working
directory. LinTeX' working directory would typically be the root directory of your LaTeX
project. To create your first rule, make sure that the directories `.lintex` and
`.lintex/rules` exist and create a file called `.lintex/rules/my-first-rule.toml`.

Rules are written in the [TOML](https://toml.io/en/) format. If you're new to TOML, just
treat it as a fancy way to denote key-value-pairs.

## First rule

Let's use this TOML code to define our first rule. This rule will prevent us from using
the sequence "foo" in our document:

```toml
name = "No 'foo'"
description = "The character sequence 'foo' is forbidden."
capture = "word"
patterns = ['((word) @word (#match? @word "foo"))']
```

This rule contains the most important keys:

`name`
: The name of the rule. This acts as an identifier for humans.

`description`
: A description, that briefly explains the rule and what's wrong, if it applies. This
  will be printed with the linting output, so it should be kept short (~ 1 line).

`patterns`
: A list of [Tree-Sitter Query
  Pattern](https://tree-sitter.github.io/tree-sitter/using-parsers#pattern-matching-with-queries)s.
  This part is responsible for finding rule violations in the document.

`capture`
: As a Tree-Sitter Query Pattern can define multiple captures, this field defines the
  capture, that matches the violation the closest. This is used to indicate the location
  of the rule violation to the user as precise as possible.

This rule has the name "No 'foo'". When it applies, the user sees the description "The
character sequence 'foo' is forbidden.", which allows them to easily find and correct
the violation in their document.

The pattern identifies parts of the document that violate the rule. It does that by
first finding all words in the document using the `(word)` query. This word is then
stored in the reference `@word`. Finally, the `#match?` predicate determines whether the
captured word is the word "foo", by checking if it matches the given regular expression.

Finally, `capture` tells LinTeX where exactly the violation is. This is useful, if you
need to capture a larger chunk of code in order to determine whether a small piece of
code contains an error.  
An example of this would be for a rule, that applies to captions
inside `figure` environments: If the caption has an error, only the caption should be
highlighted, but not the entire `figure`.  
For rules that work with no more than one capture, that capture is fine to be used here.

## Writing queries

{{< hint notice >}}

**Under construction**

This section is coming soon.

{{< /hint >}}

## Testing rules

LinTeX includes a test framework to validate that your rules work as you expect them to.
Test themselves are defined in the same TOML file as lists of violations and obediences:

```toml
[[tests.obediences]]
name = "There is no foo."
input = 'This text does not contain the forbidden word.'

[[tests.violations]]
name = "Tere is a foo."
input = 'This text does contain the word foo.'
```

Each mention of `[[tests.obediences]]` or `[[tests.violations]]` represents a test case.
LinTeX will find those test cases and run your rule against the given `input`. If it's
supposed to be an obedience, it checks that there's no error; if it's supposed to be a
violation, it checks that there is an error.  
Note that the keys `[[tests.obediences]]` and `[[tests.violations]]` can be repeated as
often as needed.

To run the tests, run `lintex test`. The output should look like this:

```shell-session
$ lintex test
Checking rule local/my-first-rule: No foo
All checks passed.
```

Let's break the rule by using `#eq?` instead of `#match?` and see what happens:

```diff
-patterns = ['((word) @word (#match? @word "foo"))']
+patterns = ['((word) @word (#eq? @word "foo"))']
```

```shell-session
$ lintex test
Checking rule local/my-first-rule: No foo
Error at violation #0: Wrong number of violations: want=1, got=0
Errors in local rules detected.
```

Replacing `#match?` with `#eq?`, breaks the rule, which is correctly reported by LinTeX.

{{< hint notice >}}

**Why does this break the rule?**

The LaTeX Tree-Sitter grammar parses interpunctuation as parts of words, s.th. in this
case the node would have the value `foo.`, which fails the `#eq?` check against `foo`.

{{< /hint >}}
