Here’s a compact Git cheat sheet with the commands most worth memorizing for day-to-day work.

# Setup

```bash
git config --global user.name "Your Name"
git config --global user.email "you@example.com"
git config --global init.defaultBranch main
```

# Start & Clone

```bash
git init                 # Start a new repo
git clone <repo-url>     # Copy existing repo
```

# Check Status

```bash
git status               # See changed files
git log                  # Commit history
git log --oneline --graph --all
```

# Add & Commit

```bash
git add file.txt         # Add one file
git add .                # Add everything
git commit -m "message"  # Commit changes
```

# Branching

```bash
git branch               # List branches
git branch feature-x     # Create branch
git checkout feature-x   # Switch branch
git checkout -b feature-x # Create + switch

# Modern alternative
git switch feature-x
git switch -c feature-x
```

# Pull & Push

```bash
git pull origin main     # Get latest changes
git push origin main     # Push commits
```

# Remote Repositories

```bash
git remote -v            # Show remotes
git remote add origin <repo-url>
```

# Undo Changes

```bash
git restore file.txt     # Discard local file changes
git reset HEAD file.txt  # Unstage file

git reset --soft HEAD~1  # Undo commit, keep changes
git reset --hard HEAD~1  # Remove commit + changes
```

# Stashing

```bash
git stash                # Save unfinished work
git stash pop            # Restore stash
```

# Merge & Rebase

```bash
git merge feature-x      # Merge branch
git rebase main          # Rebase onto main
```

# Sync Branch

```bash
git fetch                # Download latest refs
git pull --rebase origin main
```

# See Differences

```bash
git diff                 # Unstaged changes
git diff --staged        # Staged changes
```

# Delete Branches

```bash
git branch -d feature-x   # Delete local branch
git push origin --delete feature-x
```

# Handy Everyday Workflow

```bash
git status
git add .
git commit -m "Describe changes"
git pull --rebase origin main
git push origin main
```

# Emergency Commands

```bash
git reflog               # Recover lost commits
git stash                # Save work quickly
git checkout -- .        # Discard all local changes (older syntax)
```

# Most Important 10 to Memorize

```bash
git status
git add .
git commit -m "msg"
git pull
git push
git branch
git checkout -b new-branch
git switch branch-name
git merge branch-name
git log --oneline
```

If you want, I can also give:

* a **1-page printable Git cheat sheet**
* **Git commands for beginners vs advanced**
* **most common Git mistakes + fixes**
* **GitHub workflow commands**
* **interactive rebase cheat sheet**
