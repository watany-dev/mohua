#!/bin/bash

# SageMakerが利用可能なリージョンを取得
regions=$(aws ec2 describe-regions --query "Regions[].RegionName" --output text)

# すべてのリージョンでmofuaを実行
for region in $regions; do
    echo "Running mofua in region: $region"
    ./mohua -r "$region"
done

echo "All regions processed."
