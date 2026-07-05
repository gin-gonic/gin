---
name: "Resolve PR Comments"
description: "Fetch all review comments from a GitHub PR and resolve them by applying the suggested fixes in the PR's source branch. Use when: address PR feedback, resolve review comments, apply Copilot suggestions, fix PR review issues."
argument-hint: "PR number or full GitHub PR URL"
agent: "agent"
tools: [fetch_webpage, githubRepo]
---

You are helping resolve open review comments on a GitHub Pull Request.

## Steps

### 1. Fetch the PR page
Use `#fetch_webpage` to retrieve the PR at:
`https://github.com/$REPO_OWNER/$REPO_NAME/pull/$PR_NUMBER`

Where `$REPO_OWNER`, `$REPO_NAME`, and `$PR_NUMBER` come from the argument (URL or number). If only a number is given, infer the repo from the current workspace's git remote.

### 2. Extract review comments
From the fetched page identify:
- **Each review comment**: reviewer, file, line(s), and the concern raised
- **Suggested code changes**: diff blocks inside comment bodies
- **Inline code comments** on specific files/lines

### 3. Identify the source branch
From the PR page, find the source branch (e.g., `owner:branch-name`). Check out that branch locally:
```
git checkout -b <branch> origin/<branch>   # or: git checkout <branch>
```

### 4. Resolve each comment
For every open comment:
1. Read the affected file at the referenced lines
2. Apply the fix — prefer the reviewer's suggested changeset when provided; otherwise implement the intent of the comment
3. Run any relevant tests to verify correctness
4. Note which comment is resolved

### 5. Commit and push
Stage all changed files, write a clear commit message that references which review concerns were addressed, and push to the source branch:
```
git add <files>
git commit -m "fix: resolve PR review comments\n\n<bullet list of addressed concerns>"
git push origin <branch>
```

### 6. Report back
Summarize what was changed per comment:
- ✅ Comment from `@reviewer` on `file:line` → what was done
- ⚠️ Any comment that could not be automatically resolved and needs manual attention

## Notes
- Always run tests before committing
- If a suggested changeset is outdated (marked as such on GitHub), implement the intent rather than applying the diff verbatim
- Keep each fix minimal — only change what the review comment targets
