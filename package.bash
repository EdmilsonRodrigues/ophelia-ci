VERSION="1.0.0"
sed -i "s/Version: .*/Version: $VERSION/" packaging/DEBIAN/control
dpkg-deb --build packaging "dist/ophelia-ci-server_${VERSION}_amd64.deb"
