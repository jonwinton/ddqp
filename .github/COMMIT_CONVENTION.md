# Commit Message Convention

This project follows [Conventional Commits](https://www.conventionalcommits.org/) for commit messages.

## Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

## Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code (white-space, formatting, etc)
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **chore**: Changes to the build process or auxiliary tools

## Examples

### Simple commit
```
feat: add support for regex filters in metric queries
```

### Commit with scope
```
fix(parser): handle escaped characters in filter values
```

### Commit with body
```
refactor: simplify MetricExpression AST structure

The previous implementation had unnecessary nesting that made
traversal more complex. This change flattens the structure while
maintaining backward compatibility.
```

### Breaking change
```
feat!: change MetricQuery.String() return format

BREAKING CHANGE: The String() method now includes spaces around
operators for better readability. This may break exact string
comparisons in tests.
```

or

```
feat(parser): change default aggregation behavior

BREAKING CHANGE: queries without explicit aggregators now default
to 'avg' instead of 'sum'
```

## Scopes

Common scopes for this project:
- `parser` - Core parsing logic
- `filter` - Filter-related code
- `expression` - Expression parsing
- `monitor` - Monitor query parsing
- `generic` - Generic parser
- `test` - Test-related changes
- `ci` - CI/CD configuration
- `docs` - Documentation

## Why Conventional Commits?

1. **Automatic changelog generation** - Tools like git-cliff can generate changelogs from commit history
2. **Semantic versioning** - Commit types determine version bumps (feat = minor, fix = patch, breaking = major)
3. **Better collaboration** - Clear commit messages help reviewers understand changes
4. **Git history search** - Easy to filter by type (`git log --grep="^feat"`)

## Enforcing the Convention

This project enforces conventional commits through:
- **PR Title Validation** - GitHub Actions validates PR titles follow the format
- **svu** - Semantic version detection validates commits during release

Individual developers can optionally set up local tooling:
- [commitlint](https://commitlint.js.org/) - Lint commit messages
- [pre-commit hooks](https://pre-commit.com/) - Validate commits before pushing
- Git commit templates - Provide structure when committing

## Tips

- Keep the subject line under 72 characters
- Use the imperative mood ("add" not "added" or "adds")
- Don't end the subject line with a period
- Use the body to explain *what* and *why*, not *how*
- Reference issues in the footer: `Fixes #123` or `Closes #456`
