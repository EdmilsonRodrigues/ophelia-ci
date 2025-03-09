VERSION="1.0.0"
sed -i "s/Version: .*/Version: $VERSION/" deb-packaging/DEBIAN/control
sed -i "s/version: .*/version: $VERSION/" snap/snapcraft.yaml
dpkg-deb --build deb-packaging "dist/ophelia-ci-server_${VERSION}_amd64.deb"
