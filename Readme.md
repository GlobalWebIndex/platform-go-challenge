# Ownify Crm API

1. Testing - start Local Memsql
   `docker run -i --init \
    --name singlestore-ciab \
    -e LICENSE_KEY="BDc0N2Y2YTQzZjM2MTQwYTM4MmFlNTJhOGI4Y2I2MGI1AAAAAAAAAAAEAAAAAAAAACgwNAIYMWL4bAtDSFFwpdt5xsVJwuo3/NM6n0ilAhgVC+he8d+tSaNqHwRtQZzvMyDlKwLn6ZYAAA==" \
    -e ROOT_PASSWORD="helloworld" \
    -p 3306:3306 -p 8080:8080 \
    singlestore/cluster-in-a-box
docker start singlestore-ciab`
