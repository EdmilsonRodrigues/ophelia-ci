VERSION=$(cat .version)
sed -i "s/version: .*/version: $VERSION/" snap/snapcraft.yaml
snapcraft