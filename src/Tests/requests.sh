#!/bin/bash 

# get IPFS Path of record 
# curl --header 'Content-Type: application/json' --header 'Accept: */*' --header 'Authorization: Bearer testauth' --request GET http://localhost:8082/pins/

# add IPNS ID to queue for tracking
curl --header 'Content-Type: application/json' --header 'Accept: */*' --header 'Authorization: Bearer testauth'  -d '{"cid":"k51qzi5uqu5diir8lcwn6n9o4k2dhohp0e16ur82vw55abv5xfi91mp00ie0ml"}' --request POST http://$localhost:8082/follow/

# delete IPNS ID from queue to stop tracking
# curl --header 'Content-Type: application/json' --header 'Accept: */*' --header 'Authorization: Bearer testauth' --request DELETE http://localhost:8082/follow/{requestID}