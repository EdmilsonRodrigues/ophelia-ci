OPERATION="i"
CHANGELOG="deb-packaging/DEBIAN/changelog"

if [ "$#" -eq 1 ]; then
    OPERATION="$1"
fi

case "$OPERATION" in
    "i")
        echo "inserting"
        read -p "Enter your message: " MESSAGE        
        dch -c $CHANGELOG -mi "$MESSAGE"
        ;;
    *)
        echo "Invalid operation"
        ;;
esac
