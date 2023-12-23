+++
title = 'LinTeX'
draft = false
bookHidden = false
+++

# LinTeX - A LaTeX Linter Powered by Tree-Sitter

## What is LinTeX?

LinTeX is a linter that helps you to maintain a consistent style throughout your entire
document. A linter is a tool for static analysis, which means it analyses your LaTeX
source code to find places, where the actual style is not the desired style.

{{< hint warning >}}

**LinTeX is not stable**

LinTeX is currently being developed and cannot be considered as stable software yet.
Read [State of this Project](#state-of-this-project) to find out more.

{{< /hint >}}

## How does it work?

By defaul LinTeX recusrively scans the project directory to find all `.tex` files. It
then parses them into Abstract Syntax Trees (ASTs) using
[Tree-Sitter](https://tree-sitter.github.io).

Every file is checked against a set of rules. Every rule is a Tree-Sitter query combined
with optinoal additional logic, to identify violations. After this check, a list of
style violations throughout the document is generated, that is presented to the user for
them to fix.

## Why LinTeX?

### Cross-Platform Zero Dependency

LinTeX is shipped as a single binary with all dependencies included. No need to have any
libraries or interpreters for scripting languages available.

### Tree-Sitter Backend

Tree-Sitter operates on the syntax level and therefore overcomes limitations of RegEx
based approaches. It is able to capture complex structures independent of the layout
that's used in the actual source file. And it's really fast.

### Customization as a Design Paradigm

Of course everyone has their own style and very strong opinions about it. LinTeX is
designed to allow users to define their own rules, to make sure that _your_ style is
made consistent.

## State of this Project

I started this project a while ago because

- I was writing my bachelor thesis and really would have liked to have such a tool, and
  I'm planning to use LaTeX in the future.
- I wanted to learn about Go, and it turned out that Go is quite suitable for building
  cross platform CLI applications
- I wanted to learn about Tree-Sitter

I wrote a PoC to prove that the general approach works, and since then I'm working to
make an MVP. Currently, nothing that's there should be considered stable in any way.

### Features that are there

- [x] Parse LaTeX files using Tree-Sitter.
- [x] Read rules from TOML files at compile time.
- [x] Find violations for those rules in the parsed LaTeX files.
- [x] Print out a report to the console.

### Features that are not there

- [ ] Reading rules at runtime from the user's filesystem.
- [ ] LaTeX Tree-Sitter Playground here in the docs.
- [ ] Advanced logic through an embedded programming language.

    Some rules hit limitations of Tree-Sitter or the LaTeX grammar that's used. For
    those, there should be a way to implement those without having to touch the Go
    source. My favourite is to use Lua, e.g. through
    [`gopher-lua`](https://github.com/yuin/gopher-lua).

- [ ] Control-comments to influence LinTeX' behavior in LaTeX files
    
    E.g.: `% lintex: ignore_file` for non-text files.

- [ ] Customization (command line, config files, ...)
- [ ] Apply rules to `.bib` files.
- [ ] Watchmode, servermode, etc.

    Make use of Tree-Sitter's most important performance feature, which is incremental
    parsing, by using a long running process to keep the trees around and provide a
    sensible interface.

- [ ] Editor integration? LSP implementation?
- [ ] Consider supporting [Typst](https://github.com/typst/typst)? Consider supporting
  other types of markup?
- [ ] What feature would you like to see in LinTeX? Tell me by [creating an issue with
  your suggestion](https://github.com/tobb10001/lintex/issues/new).

### How to get started?

If you want to start using LinTeX, find out [how to install](docs/installation) and [how
to use](docs/usage) LinTeX.
