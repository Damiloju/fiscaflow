# Workflow Rules

Rules for development workflow and process management.

## Rules

### conventional-commits.mdc
Enforces conventional commit message format for consistent and meaningful commit history.

**Features:**
- Automatically detects change types from commit descriptions
- Extracts scope from file paths
- Formats commits as `<type>(<scope>): <description>`
- Supports all conventional commit types (feat, fix, docs, style, refactor, perf, test, chore)

**Usage:**
- Automatically applied when committing changes
- Provides suggestions for commit message format
- Helps maintain clean git history

### cursor-rules-location.mdc
Ensures proper placement and organization of Cursor rule files.

**Features:**
- Enforces `.cursor/rules/` directory structure
- Validates rule file naming conventions
- Prevents rule files from being placed in incorrect locations
- Maintains consistent rule organization

**Usage:**
- Applied when creating new rule files
- Ensures rules are properly organized
- Maintains project structure standards

## Best Practices

1. **Commit Messages**: Always use conventional commit format
2. **Rule Organization**: Keep rules in appropriate category folders
3. **Naming**: Use kebab-case for rule filenames
4. **Documentation**: Include clear descriptions and examples in rules 