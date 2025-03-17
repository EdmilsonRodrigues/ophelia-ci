VERSION=$(cat .version)
sed -i "s/Version: .*/Version: $VERSION/" deb-packaging/DEBIAN/control
dpkg-deb --build deb-packaging "../dist/ophelia-ci-client_${VERSION}_amd64.deb"
