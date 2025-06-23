# Cursor Rules

This directory contains Cursor rules that help maintain code quality, security, and development standards across the project.

## Organization

Rules are organized into the following categories:

### 📋 [Workflow](./workflow/)
Rules for development workflow and process management:
- **conventional-commits.mdc** - Enforces conventional commit message format
- **cursor-rules-location.mdc** - Ensures proper placement of Cursor rule files

### 🛠️ [Development](./development/)
Rules for development best practices and code standards:
- **code-quality.mdc** - Enforces code quality standards and best practices
- **testing-standards.mdc** - Maintains testing standards and coverage requirements
- **performance-guidelines.mdc** - Optimizes application performance

### 🔒 [Security](./security/)
Rules for security best practices and vulnerability prevention:
- **security-best-practices.mdc** - Prevents common security vulnerabilities

### 📚 [Quality](./quality/)
Rules for documentation and quality assurance:
- **documentation-standards.mdc** - Maintains comprehensive documentation standards

## Usage

These rules are automatically applied by Cursor when working on files that match their filters. They provide suggestions, enforce standards, and help maintain consistency across the codebase.

## Adding New Rules

When adding new rules:
1. Place them in the appropriate category folder
2. Follow the naming convention: `kebab-case.mdc`
3. Use the proper rule format with `<rule>` tags
4. Include clear descriptions and examples

## Rule Priority Levels

- **critical** - Must be followed (security, critical functionality)
- **high** - Strongly recommended (quality, performance)
- **medium** - Recommended (workflow, documentation)
- **low** - Optional (style, preferences) 