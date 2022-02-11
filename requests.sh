#!/bin/bash 

curl --request GET http://localhost:8082/getKey 

curl --request POST http://localhost:8082/postKey?CID=QmSVjCYjy4jYZynyC2i5GeFgjhq1bLCK2vrkRz5ffnssqo -F "file=@temp.key" -vvv

