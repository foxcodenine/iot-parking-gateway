#!/bin/sh

SHARED_DIR="/root/dist/public"
echo "123"

# Check if the shared directory is empty
if [ ! "$(ls -A $SHARED_DIR)" ]; then
    echo "Initializing shared/public directory..."
    mkdir -p $SHARED_DIR
    cp -r /root/dist/* $SHARED_DIR
else
    echo "Shared directory already initialized."
fi

exec "$@"