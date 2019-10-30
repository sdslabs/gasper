# Contribution Guidelines

## Contents
- [Coding Style Guide](#coding-style-guide)
- [Pre-Commit](#pre-commit)
- [Git Commit Message Format](#git-commit-message-format)

## Coding Style Guide

Gasper follows standard effective go guidelines. You can refer to [this guide](https://github.com/golang/go/wiki/CodeReviewComments)
for more information.

## Pre-Commit

Before committing any changes, make sure you perform the following tasks:-

1. Formatting the codebase
    ```bash
    $ make fmt
    🔨 Formatting
    👍 Done
    ```

2. Vetting the codebase and resolving any errors if found
    ```bash
    $ make vet
    🔨 Vetting
    👍 Done
    ```

3. Linting the codebase
    ```bash
    $ make lint
    🔨 Linting
    👍 Done
    ```

## Git Commit Message Format

Taken from https://github.com/angular/angular.js/blob/master/CONTRIBUTING.md and modified as required.
Each commit message consists of a **header**, a **body** and a **footer**. The header has a special
format that includes a **type**, a **scope** and a **subject**:

```
<type>(<scope>): <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

Any line of the commit message cannot be longer 100 characters! This allows the message to be easier
to read on github as well as in various git tools.

### Type

Must be one of the following:

- **feat**: A new feature
- **fix**: A bug fix
- **style**: CSS Changes
- **cleanup**: Changes that do not affect the meaning of the code (white-space, formatting, missing
  semi-colons, dead code removal etc.)
- **refactor**: A code change that neither fixes a bug or adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or fixing them
- **chore**: Changes to the build process or auxiliary tools and libraries such as documentation
  generation
- **tracking**: Any kind of tracking which includes Bug Tracking, User Tracking, Anyalytics, AB-Testing etc
- **docs**: Documentation only changes

### Scope

The scope could be anything specifying place/context of the commit change. For example `editor`,
`reports`, `FileName`, `ServiceName`, `DirectiveName`, `FunctionName` etc..

### Subject

The subject contains succinct description of the change:

- use the imperative, present tense: "change" not "changed" nor "changes"
- don't capitalize first letter
- no dot (.) at the end

### Body

Just as in the **subject**, use the imperative, present tense: "change" not "changed" nor "changes"
The body should include the motivation for the change and contrast this with previous behavior.

### Footer

The footer should contain any information about **Breaking Changes** and is also the place to
reference issues that this commit **Closes**.
