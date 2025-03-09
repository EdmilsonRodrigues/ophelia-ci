VERSION="1.0.0"
sed -i "s/Version: .*/Version: $VERSION/" deb-packaging/DEBIAN/control
dpkg-deb --build deb-packaging "dist/ophelia-ci-server_${VERSION}_amd64.deb"
