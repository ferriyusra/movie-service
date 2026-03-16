#!/bin/bash

# Hook wrapper for OpenAPI update
# This is called automatically when API files change

# Read the tool input from stdin
TOOL_INPUT=$(cat)

# Extract the file path from the tool input
FILE_PATH=$(echo "$TOOL_INPUT" | grep -o '"file_path":"[^"]*"' | cut -d'"' -f4)

# Check if the modified file is an API-related file
if [[ "$FILE_PATH" =~ controllers/.*\.go$ ]] || \
   [[ "$FILE_PATH" =~ routers/.*\.go$ ]] || \
   [[ "$FILE_PATH" =~ models/.*\.go$ ]] || \
   [[ "$FILE_PATH" =~ services/.*\.go$ ]] || \
   [[ "$FILE_PATH" == *"/main.go" ]]; then
    
    # Get the directory of this script
    HOOK_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    SCRIPT="$HOOK_DIR/../scripts/update-openapi.sh"
    
    # Execute the main script
    exec "$SCRIPT" --force
fi

# Exit silently if not an API file
exit 0