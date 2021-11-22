# IPNSGoServer

## This repo is currently still under development and does not work as intended.

## About 

This project is meant to bring dynamic NFTs to any NFT-supporting Blockchain. As the desire for NFT utility is on the rise, we have found a way to provide that in a decentralized way!

You will be able to do this by storing IPNS records on the NFT contract! Then when any changes are needed, you can publish to your IPNS record, and BOOM update your data on IPFS!

Unfortunately you can not creat IPNS records through js-ipfs (the browser) and that is why this is necessary. We provide a go-ipfs implementation that can publish to IPNS for you! If you're running in the browser or another server, you can request Ipfs keys from us and publish your IPNS record through us.

This let's you update your data without having to update your URI on the chain, thus raising gas prices or requiring constant updating from your owners. Instead you can do it completely off chain! All you need is a NFT supporting blockchain and infrastructure on your end to handle IPNS record retrieval.



## To Run
Make sure you have golang installed and go-ipfs. You can find instructions on their respective websites.

After that is done, to execute type:

Go mod tidy

This is similar to npm install. It will download all the packages needed, specified in go.mod

From there you can type:

Go run server.go

Now you can request keys and publish to IPNS!

