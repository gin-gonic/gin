# How to Contribute

## Your First Pull Request
We use github for our codebase. You can start by reading [How To Pull Request](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/about-pull-requests).

## Branch Organization
We use [git-flow](https://nvie.com/posts/a-successful-git-branching-model/) as our branch organization, as known as [FDD](https://en.wikipedia.org/wiki/Feature-driven_development)

## Bugs
### 1. How to Find Known Issues
We are using [Github Issues](https://github.com/cloudwego/kitex/issues) for our public bugs. We keep a close eye on this and try to make it clear when we have an internal fix in progress. Before filing a new task, try to make sure your problem doesn’t already exist.

### 2. Reporting New Issues
Providing a reduced test code is a recommended way for reporting issues. Then can placed in:
- Just in issues
- [Golang Playground](https://play.golang.org/)

### 3. Security Bugs
Please do not report the safe disclosure of bugs to public issues. Contact us by [Support Email](mailto:conduct@cloudwego.io)

## How to Get in Touch
- [Email](mailto:conduct@cloudwego.io)

## Submit a Pull Request
Before you submit your Pull Request (PR) consider the following guidelines:
1. Search [GitHub](https://github.com/cloudwego/kitex/pulls) for an open or closed PR that relates to your submission. You don't want to duplicate existing efforts.
2. Be sure that an issue describes the problem you're fixing, or documents the design for the feature you'd like to add. Discussing the design upfront helps to ensure that we're ready to accept your work.
3. [Fork](https://docs.github.com/en/github/getting-started-with-github/fork-a-repo) the cloudwego/kitex repo.
4. In your forked repository, make your changes in a new git branch:
    ```
    git checkout -b my-fix-branch develop
    ```
5. Create your patch, including appropriate test cases.
6. Follow our [Style Guides](#code-style-guides).
7. Commit your changes using a descriptive commit message that follows [AngularJS Git Commit Message Conventions](https://docs.google.com/document/d/1QrDFcIiPjSLDn3EL15IJygNPiHORgU1_OOAqWjiDU5Y/edit).
   Adherence to these conventions is necessary because release notes are automatically generated from these messages.
8. Push your branch to GitHub:
    ```
    git push origin my-fix-branch
    ```
9. In GitHub, send a pull request to `kitex:develop`

## Contribution Prerequisites
- Our development environment keeps up with [Go Official](https://golang.org/project/).
- You need fully checking with lint tools before submit your pull request. [gofmt](https://golang.org/pkg/cmd/gofmt/) and [golangci-lint](https://github.com/golangci/golangci-lint)
- You are familiar with [Github](https://github.com)
- Maybe you need familiar with [Actions](https://github.com/features/actions)(our default workflow tool).

## Code Style Guides
Also see [Pingcap General advice](https://pingcap.github.io/style-guide/general.html).

Good resources:
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
