# IPNSGoServer

# About 

This project is meant to bring dynamic, decentralized data to any project. It started in the Web3 Jam 2021 hackathon, out of the founders need for customizeable decentralized content. While searching for resources, there were more creators asking the same question. "What is the best way to have elements in a NFT's metadata change based on certain events using IPFS?". This api is meant for devs looking to create dapps or services that give their users customizeable content on the InterPlanetary File System (IPFS)

# Problem

IPFS is a peer-to-peer file sharing protocol that has enabled users to host a wide variety of content. If someone wants to publish a website on IPFS they will have to direct users to the specific Content Identifier (CID) where they can download the site. If the Original Poster (OP) makes any changes to that site, they will get a completely different CID that they need to redirect their users to. To combat this, IPNS can offer one an address that can point to different CIDs. The InterPlanetary Name System (IPNS) was created to bring mutability and human readable addresses to IPFS. The problem today with IPNS is the revival of IPNS records. The records have to be republished by their expiration date to maintain discoverability and security. But it is not always possible for the node to be online to republish the record.  

Unfortunately, it is not possible to create IPNS records through js-ipfs (the browser or node.js) and that is why this project is important. This project is a go-ipfs implementation that can create, publish, and even maintain IPNS records. A client browser or server, can request keys and publish IPNS records using this project on the backend. That way devs can create dapps and not be bogged down by browser or version limitations.

# How it works

The content of IPFS hashes are immutable, once you change a single letter and resubmit to IPFS. You will get back a completely different CID. But IPNS identifiers can point to different CIDs. Once you have your IPNS record, you can change what it points to as many times as you want.

The project serves keys upon request and publishes content on a users behalf to IPFS. The idea is to enable access to the IPFS network and IPNS functionality without having to run your own IPFS node. A user can only publish IPNS content if they have keys to sign the records. To enable that, the api only stores keys for the needed operations. Each key is deleted from the local node keystore and temp_key storage before returning.

- *As of now, republishing is not implemented. But this will allow devs to come up with their own republishing mechanisms on top of the api. Currently, a user will have to be online to send us their private key to be able to republish. Or send their private key along with a new CID to update their record.*

# Deep dive

As previously mentioned, IPNS records need to be republished after a certain period of time. The reason for this is because this "world wide record" called the Distributed Hash Table (DHT) keeps track of IPNS records and what they point to. You can read on why they expire here [go-ipfs/issue#1958](https://github.com/ipfs/go-ipfs/issues/1958#issuecomment-410860667)

# API 

There are two endpoints currently:

getKey() - This will generate and return a private/public key pair in .key format. These keys are used to create and update IPNS records.

postKey(<IPFS_CID>, <key_name>.key) - This will publish content to IPFS and return IPNS address.

### To Do
add(file) - This will add content to IPFS. Will return CID.

add(dir) - This will add directory to IPFS. Will return CID.

putKey(IPFSHash, publicKey, privateKey) - This will update your record if it exists.


# To Run
Make sure you have [golang](https://go.dev/doc/install) installed and [go-ipfs](https://github.com/ipfs/go-ipfs). You can find instructions on their respective websites.

After that is done, remember to fork, clone, &:

`cd IPNSGoServer`

Then you need to build the project:

`go build`

From there you can run the server by:

`go run .`

Now you can request keys and publish to IPNS!
