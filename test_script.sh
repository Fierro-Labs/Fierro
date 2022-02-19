#!/bin/bash

key1=$(ipfs key gen temp1)
key2=$(ipfs key gen temp2)
key3=$(ipfs key gen temp3)


ipfs name publish --key=temp1 /ipfs/QmSVjCYjy4jYZynyC2i5GeFgjhq1bLCK2vrkRz5ffnssqo
ipfs name publish --key=temp2 /ipfs/QmSVjCYjy4jYZynyC2i5GeFgjhq1bLCK2vrkRz5ffnssqo
ipfs name publish --key=temp3 /ipfs/QmSVjCYjy4jYZynyC2i5GeFgjhq1bLCK2vrkRz5ffnssqo

curl --request GET http://localhost:8082/getRecord?ipnskey=$key1
curl --request GET http://localhost:8082/getRecord?ipnskey=$key2
curl --request GET http://localhost:8082/getRecord?ipnskey=$key3

ipfs key rm temp1
ipfs key rm temp2
ipfs key rm temp3
