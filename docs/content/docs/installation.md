---
title: "Installation"
weight: 1
# bookFlatSection: false
bookToc: true
bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---

# Installation

## Binary Installation

_RECOMMENDED_

1. Go to [the project's Releases page](https://github.com/tobb10001/lintex/releases).
2. Choose the release you want to install.
3. Download the binary, that corresponds to your platform.
4. Save the binary to your PATH, rename to `lintex` / `lintex.exe` if desired.
5. Ensure you have permission to execute the binary.

## Source Installation

- Download the source by cloning the Git repo.

  ```console
  $ git clone https://github.com/tobb10001/lintex
  $ cd lintex
  ```

  Alternatively the source can be downloaded as ZIP- or TAR-Archive from the releases
  page.

- Run `go build` to build your binary.

  ```console
  $ go build
  ```

- Copy the binary to your `PATH`

  ```console
  $ cp lintex /somewhere/in/your/path
  ```

- Run LinTeX

  ```console
  $ lintex
  ```
