#!/bin/bash
# Watch proto files and auto-generate when changed
echo "Watching proto files for changes..."

# Install fswatch if not available
# brew install fswatch (macOS)
# apt-get install fswatch (Ubuntu)

# Watch proto files and regenerate
fswatch -o ./user_service/proto/ ./proto/ | while read f; do
    echo "Proto file changed, regenerating..."
    ./scripts/gen_proto.sh
done 