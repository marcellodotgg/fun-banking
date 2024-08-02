# Pull Latest
git pull

# Build Styles and Hash
npm install -D tailwindcss
npx tailwindcss -i ./public/styles.css -o ./public/styles.min.css --minify

# Generate hash
HASH=$(md5sum public/styles.min.css | cut -d ' ' -f 1)
NEW_CSS_FILE_NAME="styles-$HASH.min.css"

mv ./public/styles.min.css "./public/$NEW_CSS_FILE_NAME"

# Update the .env file with the new hash
if grep -q "BUILD_HASH=" .env; then
    if [ "$OS" = "Darwin" ]; then
        # macOS version using awk
        awk -v hash="$HASH" '/BUILD_HASH=/ { $0="BUILD_HASH=" hash } { print }' .env > .env.tmp && mv .env.tmp .env
    else
        # Linux version using awk
        awk -v hash="$HASH" '/BUILD_HASH=/ { $0="BUILD_HASH=" hash } { print }' .env > .env.tmp && mv .env.tmp .env
    fi
else
    echo "BUILD_HASH=$HASH" >> .env
fi

echo "Updated build hash to: $HASH"

# Build docker
docker build -t fun_banking .
docker kill fun_banking_container
docker container prune -f
docker run -d \
       -p 8082:8080 \
       -v $(pwd)/fun_banking.db:/app/fun_banking.db \
       -e DATABASE_URL=/app/fun_banking.db \
       -e GIN_MODE=release \
       --name fun_banking_container fun_banking
