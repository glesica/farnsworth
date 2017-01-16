#!/usr/bin/env bash

set -e

SCRIPT_DIR=`dirname $BASH_SOURCE`
TEMP_DIR="$SCRIPT_DIR/_temp"
PROJECT_DIR=java_project

rm -rf "$TEMP_DIR"
mkdir -p "$TEMP_DIR"

"$SCRIPT_DIR"/../farnsworth archive --project "$SCRIPT_DIR/$PROJECT_DIR" "$TEMP_DIR/$PROJECT_DIR.zip"
unzip "$TEMP_DIR/$PROJECT_DIR.zip" -d "$TEMP_DIR" > /dev/null

# Gradle file made the trip
if [ ! -f "$TEMP_DIR/$PROJECT_DIR/build.gradle" ]; then
    rm -rf "$TEMP_DIR"
    echo "FAILURE - archive"
    exit 1
fi

# Deleting a file on the merge path and merging the origin replaces it
testfile="$TEMP_DIR/$PROJECT_DIR/src/test/java/AppTest.java"
rm "$testfile"
"$SCRIPT_DIR"/../farnsworth merge --project "$TEMP_DIR/$PROJECT_DIR" "$SCRIPT_DIR/$PROJECT_DIR"
if [ ! -f "$testfile" ]; then
    rm -rf "$TEMP_DIR"
    echo "FAILURE - merge"
    exit 1
fi

rm -rf "$TEMP_DIR"
echo "SUCCESS"
exit 0

