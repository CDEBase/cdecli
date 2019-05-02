set -ex
git tag $VERSION -a -m "v$VERSION"
git tag latest -f -a -m "v$VERSION"
git push -f --tags
