ORIGINAL_PYPROJECT="../pyproject.toml"
DEPENDENCIES=$(sed -n '/^\s*]/q; s/^\s*"\([^"]*\).*/\1/p' $ORIGINAL_PYPROJECT)

echo "$DEPENDENCIES" > src/requirements.txt
