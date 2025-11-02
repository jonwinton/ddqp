#!/bin/bash
# Automated release script using svu for version calculation
# Usage: ./scripts/create-release.sh [major|minor|patch]

set -e

# Ensure we're in a git repository
if [ ! -d .git ]; then
    echo "? Error: Not in a git repository"
    exit 1
fi

# Check for uncommitted changes
if [ -n "$(git status --porcelain)" ]; then
    echo "? Error: You have uncommitted changes"
    git status --short
    exit 1
fi

# Activate Hermit environment
if [ ! -f ./bin/activate-hermit ]; then
    echo "? Error: Hermit not found"
    exit 1
fi

echo "?? Analyzing commits to determine next version..."

# Get current version
current_version=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
echo "?? Current version: $current_version"

# Calculate next version based on conventional commits
if [ -n "$1" ]; then
    # Manual override
    case $1 in
        major)
            next_version=$(svu major)
            ;;
        minor)
            next_version=$(svu minor)
            ;;
        patch)
            next_version=$(svu patch)
            ;;
        *)
            echo "? Invalid argument. Use: major, minor, or patch"
            exit 1
            ;;
    esac
    echo "?? Manual version bump: $next_version ($1)"
else
    # Automatic detection based on conventional commits
    next_version=$(svu next)
    echo "?? Auto-detected next version: $next_version"
fi

# Preview changelog
echo ""
echo "?? Preview of changes for $next_version:"
echo "??????????????????????????????????????????????????????"
git-cliff --config cliff.toml --unreleased --strip all
echo "??????????????????????????????????????????????????????"
echo ""

# Confirm
read -p "?? Create release $next_version? [y/N] " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "? Release cancelled"
    exit 1
fi

# Create and push tag
echo "???  Creating tag $next_version..."
git tag "$next_version"

echo "??  Pushing tag to origin..."
git push origin "$next_version"

echo ""
echo "? Release $next_version initiated!"
echo ""
echo "?? GitHub Actions will now:"
echo "   1. Generate changelog"
echo "   2. Create GitHub release"
echo "   3. Publish release notes"
echo ""
echo "?? Monitor progress: https://github.com/jonwinton/ddqp/actions"
