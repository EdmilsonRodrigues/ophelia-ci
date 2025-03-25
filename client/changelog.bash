OPERATION="a"
CHANGELOG="deb-packaging/DEBIAN/changelog"

if [ "$#" -eq 1 ]; then
    OPERATION="$1"
fi

case "$OPERATION" in
    "i")
        echo "incrementing version"
        read -p "Enter your message: " MESSAGE        
        dch -c $CHANGELOG -mi "$MESSAGE"
        ;;
    "a")
        echo "appending to current"
        read -p "Enter your message: " MESSAGE
        dch -c $CHANGELOG -ma "$MESSAGE"
        ;;
    "v")
        echo "changing to new version"
        read -p "Enter your version: " VERSION
        read -p "Enter your message: " MESSAGE
        dch -c $CHANGELOG -v $VERSION -m "$MESSAGE"
        echo $VERSION > .version
        ;;
    *)
        echo "Invalid operation"
        ;;
esac
