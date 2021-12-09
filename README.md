# IPNSGoServer

## This repo is currently still under development and does not work as intended.

## About 

This project is meant to bring dynamic, decentralized data to any project. It started in the Web3 Jam 2021 hackathon, out of the founders need for customizeable decentralized content. While searching for resources, there were more creators asking the same question. "What is the best way to have elements in a NFT's metadata change based on certain events using IPFS?". This api is meant for devs looking to create dapps or services that give their users customizeable content in a web3 way. 

# Problem

The idea of NFTs is changing. People are finding use cases for having dynamic metadata. This is possible currently by changing the data off chain then updating the URI in a token. This can be costly and slow because of the need to edit storage on the blockchain. That is why my project aimed to lower those costs. By using IPNS and storing the record/entry on the contract you can update your NFT off-chain. 
 

# How it works
We will harness the power of IPNS, the InterPlanetary Name System. IPNS was created to bring mutability and human readable URIs to IPFS. The content of IPFS hashes ARE immutable, but IPNS hashes can point to different IPFS entries. Once you have your IPNS record, you can change what it points to as many times as you want.

My system works by creating an api that accepts requests for private/public keys, then using that to create an IPNS record that gets returned to the user. From there the client will mint an NFT and create the contract with the IPNS record. As long as the dev has the ability to customize the URI in their smart contract they can take advantage of dynamic metadata and even customizeable images.

# Deep dive

Unfortunately, you can not create IPNS records through js-ipfs (the browser) and that is why this is necessary. We provide a go-ipfs implementation that can create, publish, and even maintain IPNS records for you! If you're running in the browser or another server, you can request keys from us and publish your IPNS record through us. That way devs can create dapps and not be bogged down by browser limitations.

There is a nother limitation to IPNS and that is that the records, similar to IPFS, need to be repinned every so often. The reason for this is similar to IPFS as well, where if no one is accessing it, then there is less priority to keep a log of it in the world wide record. This "world wide record" in IPNS is called the DHT, Distributed Hash Table. This table keeps track of IPNS keys and what IPFS hash they point to. It does more than just that, but with regards to this project, that is the most important function for the DHT.

A user needs a private key to be able to publish IPNS records, and those generally are the unique peerIDs generated when you initialize your ipfs node. This is easy to do in the command line, but that limitation means you would need multiple node instances to get many peerIDs to use as keys to distribute to users to create this same kind of API. For obvious reasons, that is not feasible. So we generate private/public key pairs, give to user to keep safe, then ask for it when they want to create and publish their records. 

## API 

There are two endpoints currently:

getKey() - This will generate and return a private/public key pair. These keys are used to create and update IPNS records.

postKey(IPFSHash, publicKey, privateKey) - This will embed your public key to your IPNS record, then publish to IPNS.

*to complete*
putKey(IPFSHash, publicKey, privateKey) - This will update your record if it exists. 


## To Run
Make sure you have golang installed and go-ipfs. You can find instructions on their respective websites.

After that is done, to execute type:

`go mod tidy`

This is similar to npm install. It will download all the packages needed, specified in go.mod

From there you can type:

`go run server.go`

Now you can request keys and publish to IPNS!
