#!/bin/sh

PUBLIC_DIR="/root/dist/public"

# Check if the shared directory is empty
if [ ! "$(ls -A $PUBLIC_DIR)" ]; then
    echo "Initializing shared/public directory..."
    mkdir -p $PUBLIC_DIR
    cp -r /root/dist/staging_images/* $PUBLIC_DIR
else
    echo "Shared directory already initialized."
fi

exec "$@"