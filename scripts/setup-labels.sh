#!/bin/bash
# Setup GitHub labels for go-webgpu/webgpu
# Run: bash scripts/setup-labels.sh

set -e

echo "Creating labels for go-webgpu/webgpu..."

# Priority
gh label create "priority: critical" --color "b60205" --description "Release blocker, security, data loss" --force
gh label create "priority: high" --color "d93f0b" --description "Important for next release" --force
gh label create "priority: medium" --color "fbca04" --description "Normal priority" --force
gh label create "priority: low" --color "fef2c0" --description "Backlog, nice to have" --force

# Type
gh label create "type: bug" --color "d73a4a" --description "Something isn't working" --force
gh label create "type: feature" --color "a2eeef" --description "New capability" --force
gh label create "type: enhancement" --color "1d76db" --description "Improve existing feature" --force
gh label create "type: docs" --color "0e8a16" --description "Documentation only" --force
gh label create "type: refactor" --color "5319e7" --description "Code cleanup, no behavior change" --force
gh label create "type: performance" --color "ff7619" --description "Speed/memory improvement" --force
gh label create "type: security" --color "b60205" --description "Security vulnerability" --force
gh label create "type: test" --color "c2e0c6" --description "Test coverage" --force

# Status
gh label create "status: triage" --color "d4c5f9" --description "Needs initial review" --force
gh label create "status: confirmed" --color "0052cc" --description "Verified, ready for work" --force
gh label create "status: in-progress" --color "fbca04" --description "Actively being worked on" --force
gh label create "status: blocked" --color "e11d21" --description "Waiting on external dependency" --force
gh label create "status: needs-info" --color "E2A1C2" --description "Awaiting reporter response" --force
gh label create "status: review" --color "6f42c1" --description "In code review" --force

# Resolution
gh label create "~duplicate" --color "cfd3d7" --description "Duplicate of another issue" --force
gh label create "~wontfix" --color "ffffff" --description "Will not implement" --force
gh label create "~invalid" --color "e4e669" --description "Not reproducible or invalid" --force
gh label create "~stale" --color "ededed" --description "Closed due to inactivity" --force

# Effort
gh label create "effort: 1" --color "c2e0c6" --description "Trivial, < 1 hour" --force
gh label create "effort: 2" --color "bfdadc" --description "Small, 1-4 hours" --force
gh label create "effort: 3" --color "c5def5" --description "Medium, ~1 day" --force
gh label create "effort: 5" --color "d4c5f9" --description "Large, 2-3 days" --force
gh label create "effort: 8" --color "f9d0c4" --description "Very large, ~1 week" --force
gh label create "effort: 13" --color "e99695" --description "Epic, 2+ weeks (split it)" --force

# Contributor
gh label create "good first issue" --color "7057ff" --description "Good for newcomers" --force
gh label create "help wanted" --color "008672" --description "Extra attention needed" --force
gh label create "mentor available" --color "0e8a16" --description "Maintainer will mentor" --force

# Area (repo-specific)
gh label create "area: wgpu" --color "0052cc" --description "wgpu/ package (core bindings)" --force
gh label create "area: examples" --color "0052cc" --description "examples/ directory" --force
gh label create "area: ffi" --color "0052cc" --description "goffi integration, FFI issues" --force
gh label create "area: types" --color "0052cc" --description "Type definitions, gputypes" --force
gh label create "area: build" --color "0052cc" --description "Build, CI/CD, releases" --force

# Platform (repo-specific)
gh label create "os: windows" --color "006b75" --description "Windows-specific" --force
gh label create "os: linux" --color "006b75" --description "Linux-specific" --force
gh label create "os: macos" --color "006b75" --description "macOS-specific" --force
gh label create "upstream" --color "d93f0b" --description "Depends on wgpu-native/goffi" --force

echo "Done! Created 40 labels."
