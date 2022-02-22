#!/bin/bash 
PORT=":8082"
LOCALHOST="localhost"
host="$LOCALHOST$PORT"
echo $host

# curl --request POST http://$host/addFile -F "file=@Hello" -v

# curl --request GET http://$host/getKey

# curl --request POST http://$host/postKey? -F "file=@temp.key" -v

# curl --request DELETE http://$host/deleteKey?keyName=temp

# curl --request GET http://$host/getRecord?ipnskey=k51qzi5uqu5dm876hw4kh2mn58rnajofhoohohymt9bui38q6ogsa0rrct6fnh

# curl --request POST http://$host/startFollowing?ipnskey=k51qzi5uqu5dm876hw4kh2mn58rnajofhoohohymt9bui38q6ogsa0rrct6fnh

curl --request DELETE http://$host/stopFollowing?ipnskey=k51qzi5uqu5dm876hw4kh2mn58rnajofhoohohymt9bui38q6ogsa0rrct6fnh

# curl --request POST http://$host/postRecord?CID=QmSVjCYjy4jYZynyC2i5GeFgjhq1bLCK2vrkRz5ffnssqo -F "file=@temp.key" -v

# for ease of use, just using local commands
# ipfspath=$(ipfs add Hello)
# ipnskey=$(ipfs key gen temp1)
# ipfs key export temp1
# ipfs publish --key=temp1 "$ipfspath"
# ipfs key rm temp1

# curl --request PUT http://$host/putRecord?ipnskey="$ipnskey" -F "file=@temp1.key" -v