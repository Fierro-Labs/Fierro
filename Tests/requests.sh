#!/bin/bash 
PORT=":8082"
LOCALHOST="localhost"
host="$LOCALHOST$PORT"
# echo $host

# curl --request POST http://$host/addFile -F "file=@Hello" # add a file from current directory called Hello, to IPFS

# curl --request GET http://$host/getKey --output temp.key # returns key to user, you specify as temp.key

# curl --request POST http://$host/postKey -F "file=@temp.key" -v # send file to API

# curl --request DELETE http://$host/deleteKey?keyName=temp # delete key from remote node

# curl --request POST http://$host/postRecord?CID=QmWEzjhLRjaJdUboeZkc7Cy9H7vDynUCxm52Dn5Grev2J4 -F "file=@temp.key" -v # publish brand new record to IPNS

# curl --request GET http://$host/getRecord?ipnskey=k51qzi5uqu5dmfs21tga7t45wltgilzu6d6krek7fcvzlyhn7x2wxp4rweyla0 # get IPFS Path of record 

# curl --request POST http://$host/startFollowing?ipnskey=k51qzi5uqu5diir8lcwn6n9o4k2dhohp0e16ur82vw55abv5xfi91mp00ie0ml # add IPNS ID to queue for tracking

# curl --request DELETE http://$host/stopFollowing?ipnskey=k51qzi5uqu5diir8lcwn6n9o4k2dhohp0e16ur82vw55abv5xfi91mp00ie0ml # delete IPNS ID from queue to stop tracking

# for ease of use, just using local commands
# ipfspath=$(ipfs add Hello)
# ipfsoutput=$(ipfs key gen temp1)
# ipfspath=$(echo ipfsoutput | awk '{split($0,a); print a[1]}')
# echo $ipfspath
# ipfs key export temp1
# ipfs name publish --key=temp1 "$ipfspath"
# ipfs key rm temp1

# curl --request PUT http://$host/putRecord?ipnskey="$ipnskey" -F "file=@temp1.key" -v # hand over IPNS record to allow republishing.
