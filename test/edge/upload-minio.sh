#!/bin/bash

# Define the S3 bucket and folder
S3_BUCKET="aon-production-mongodb-backups"
S3_FOLDER="uber-production/"

# Define the MinIO instance details
MINIO_ALIAS="myminio"
MINIO_BUCKET="neo9-testing-edge-mongodb-backups"
MINIO_HOST="http://localhost:9000"
MINIO_ACCESS_KEY="minioadmin"
MINIO_SECRET_KEY="minioadmin"

# Set up the MinIO alias (if not already done)
mc alias set $MINIO_ALIAS $MINIO_HOST $MINIO_ACCESS_KEY $MINIO_SECRET_KEY || exit 1

# Get the latest .gz file from S3
LATEST_FILE=$(aws s3api list-objects-v2 \
    --bucket "$S3_BUCKET" \
    --prefix "$S3_FOLDER" \
    --query 'Contents[?ends_with(Key, `.gz`)] | sort_by(@, &LastModified)[-1].Key' \
    --output text)

if [ "$LATEST_FILE" == "None" ] || [ -z "$LATEST_FILE" ]; then
    exit 1
fi

echo "Latest .gz file: s3://$S3_BUCKET/$LATEST_FILE"

mc mb "$MINIO_ALIAS/$MINIO_BUCKET" || echo "Bucket $MINIO_BUCKET already exists"

# Stream the file directly from S3 to MinIO
echo "Uploading to MinIO..."
aws s3 cp "s3://$S3_BUCKET/$LATEST_FILE" - | mc pipe "$MINIO_ALIAS/$MINIO_BUCKET/$(basename "$LATEST_FILE")"

if [ $? -eq 0 ]; then
    echo "Successfully uploaded $LATEST_FILE to MinIO bucket $MINIO_BUCKET"
else
    echo "Failed to upload $LATEST_FILE to MinIO"
    exit 1
fi