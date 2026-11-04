---
name: commit-per-fix
description: Split pending changes into one commit per fix, including several fixes mixed in the same file. Use when asked to commit work, make atomic commits, or "un commit par fix".
---

# One commit per fix

Each commit holds one fix and builds on its own, so it can be reverted alone.

## 1. Map the changes

```sh
git status --short
git diff -U0
git log -8 --format=%s
```

List the fixes, and assign every changed file and hunk to one of them, untracked
files included. Nothing is left over at the end.

- A locale key belongs to the fix that uses it, not to a "translations" commit.
- A test belongs to the fix it covers.
- Work you did not write gets its own commit. Say so in your report, since the
  user may have wanted it kept aside.

## 2. Order the commits

Put each fix after anything it depends on, so every commit builds.

## 3. Commit each fix

Use the repository's message style (`fix: ...`, `feat: ...`, `chore: ...`), with
one short line describing that fix, plus the attribution lines the session
requires. Add files by name, never with `git add .` or `git add -A`.

```sh
git add -- <files of this fix>
git commit -m "fix: <what this fix does>"
go build ./...
```

Stop at the first build failure and fix the split rather than the next commit.

### Several fixes in one file

Interactive staging (`git add -p`) is not available, and the lefthook `gofmt`
hook uses `stage_fixed`, which stages the whole working copy of each Go file. A
hunk staged alone with `git apply --cached` would therefore pull the other fixes
of that file into the commit.

Put the file itself in the intermediate state:

1. Back up first: `git diff --binary > <scratchpad>/backup.patch`, and archive the
   untracked files.
2. Keep the final content of the file.
3. Write the file with the later fixes undone, then add and commit it.
4. Repeat for each fix. After the last commit, check the file matches the content
   kept in step 2.

## 4. Verify

```sh
git status --short
git log --oneline -<count>
```

- The working tree is clean, unless the user kept something aside.
- Run the tests with `$go-test-run` on the last commit.
- Report the commits, whether tests ran on each one or only the last, and any
  commit holding work that was not yours.

## Do not

- Push, amend, or rebase unless asked.
- Commit on `main`: branch first.
- Bundle an unrelated change because it sits in a file you are already committing.
