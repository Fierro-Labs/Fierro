#!/bin/bash 
PORT=":8082"
LOCALHOST="localhost"
host="$LOCALHOST$PORT"
# echo $host

# curl --request GET http://$host/pins?ipnskey=k51qzi5uqu5dmfs21tga7t45wltgilzu6d6krek7fcvzlyhn7x2wxp4rweyla0 # get IPFS Path of record 

# curl --request POST http://$host/follow/k51qzi5uqu5diir8lcwn6n9o4k2dhohp0e16ur82vw55abv5xfi91mp00ie0ml # add IPNS ID to queue for tracking

curl --request DELETE http://$host/following/k51qzi5uqu5diir8lcwn6n9o4k2dhohp0e16ur82vw55abv5xfi91mp00ie0ml # delete IPNS ID from queue to stop tracking