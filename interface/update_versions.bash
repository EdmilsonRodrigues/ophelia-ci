PYPROJECT="../pyproject.toml"
VERSION=$(sed -n 's/^\s*version = "\([^"]*\).*/\1/p' $PYPROJECT)
ROCKCRAFT_PATH="src/rockcraft.yaml"
SNAPCRAFT_PATH="src/snap/snapcraft.yaml"
CHARMCRAFT_PATH="src/charm/charmcraft.yaml"


sed -i "s/version: .*/version: '$VERSION'/" $ROCKCRAFT_PATH
sed -i "s/version: .*/version: $VERSION/" $SNAPCRAFT_PATH
# sed -i "s/version: .*/version: $VERSION/" $CHARMCRAFT_PATH