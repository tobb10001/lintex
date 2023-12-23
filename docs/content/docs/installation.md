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

## Binary installation

{{< hint warning >}}

**No releases so far**

LinTeX will be made available to download through GitHub Releases. However, that's not
configured yet, so currently the only way to install LinTeX is to build it from source.

{{< /hint >}}


## Source installation

{{< hint warning >}}

This section is incomplete. It only shows instructions for Unix-like systems.

{{< /hint >}}

- Download the source by cloning the Git repo.

  ```console
  $ git clone https://github.com/tobb10001/lintex
  $ cd lintex
  ```

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
