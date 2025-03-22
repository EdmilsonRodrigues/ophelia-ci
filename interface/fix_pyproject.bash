#!/bin/bash

ORIGINAL_PYPROJECT="../pyproject.toml"
BRIEFCASE_PYPROJECT="pyproject.toml"
CONFIG="src/ophelia_ci_interface/config.py"

DEPENDENCIES=$(sed -n '/^\s*]/q; s/^\s*"\([^"]*\).*/\1/p' $ORIGINAL_PYPROJECT)
VERSION=$(sed -n 's/^\s*version = "\([^"]*\).*/\1/p' $ORIGINAL_PYPROJECT)
FORMATTED_DEPS=$(echo "$DEPENDENCIES" | awk '{print "    \"" $0 "\","}')

BRIEFCASE_START=$(sed "/^\s*requires/q; s/version = .*/version = \"${VERSION}\"/" $BRIEFCASE_PYPROJECT)
# BRIEFCASE_START=$(echo $BRIEFCASE_START | sed "s/version = .*/version = \"${VERSION}\"/")
BRIEFCASE_START_WITH_DEPENDENCIES=$(echo "${BRIEFCASE_START}
${FORMATTED_DEPS}
]")
BRIEFCASE_ENDS=$(sed '1,/^\s*]/d' $BRIEFCASE_PYPROJECT)
sed -i "s/VERSION = .*/VERSION = '$VERSION'/" $CONFIG

printf "%s\n" "$BRIEFCASE_START_WITH_DEPENDENCIES" "$BRIEFCASE_ENDS" > $BRIEFCASE_PYPROJECT
