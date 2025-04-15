# Delete existing releases
gh release delete v0.5.3 --yes
gh release delete latest --yes

# Delete existing tags
git tag -d v0.5.3
git tag -d latest
git push origin :refs/tags/v0.5.3
git push origin :refs/tags/latest

# Create new tags
git tag -a v0.5.3 -m "release v0.5.3"
git tag -f -a latest -m "Latest release v0.5.3"
git push -f origin v0.5.3
git push -f origin latest


# Install GitHub CLI if you haven't already
brew install gh  # for macOS

# Login to GitHub
gh auth login

# Create the versioned release
gh release create v0.5.3 \
  --title "v0.5.3" \
  --notes "release v0.5.3" \
  ./release/cli_darwin_amd64 \
  ./release/cli_linux_amd64 \
  ./release/cli_windows_amd64.exe \
  ./release/"Source code.zip" \
  ./release/"Source code.tar.gz"

# Create/update the latest release
gh release create latest \
  --title "Latest Release" \
  --notes "Latest release (v0.5.3)" \
  --force \
  ./release/cli_darwin_amd64 \
  ./release/cli_linux_amd64 \
  ./release/cli_windows_amd64.exe \
  ./release/"Source code.zip" \
  ./release/"Source code.tar.gz"
