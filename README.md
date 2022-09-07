<p align="center">
 <img width="200px" src="https://res.cloudinary.com/doy7gfxuc/image/upload/c_crop,h_800,w_900,g_west/Transparent_Logo_pt3s7z.png" align="center" alt="GitHub Readme Stats" />
 <h2 align="center">Fierro API</h2>
 <p align="center">An IPNS Pinning Service</p>
</p>
  <p align="center">
    <a href="https://github.com/Fierro-Labs/Fierro/actions">
      <img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/Fierro-Labs/Fierro/Go">
    </a>
    <a href="https://codecov.io/gh/Fierro-Labs/Fierro">
      <img src="https://codecov.io/gh/Fierro-Labs/Fierro/branch/main/graph/badge.svg?token=1IRKRN16IC"/>
    </a>
    <a href="https://github.com/Fierro-Labs/Fierro/issues">
      <img alt="Issues" src="https://img.shields.io/github/issues/Fierro-Labs/Fierro?color=0088ff" />
    </a>
    <a href="https://github.com/Fierro-Labs/Fierro/pulls">
      <img alt="GitHub pull requests" src="https://img.shields.io/github/issues-pr/Fierro-Labs/Fierro?color=0088ff" />
    </a>
    <a href="https://github.com/Fierro-Labs/Fierro/graphs/contributors" alt="Contributors">
        <img src="https://img.shields.io/github/contributors/Fierro-Labs/Fierro" />
    </a>
    <img alt="GitHub Sponsors" src="https://img.shields.io/github/sponsors/Fierro-Labs">
    <img alt="GitHub" src="https://img.shields.io/github/license/Fierro-Labs/Fierro">
    <br />
  </p>
</p>

# About 

This project is meant to bring dynamic, decentralized data to any project. It started in the Web3 Jam 2021 hackathon, out of the founders need for customizeable decentralized content. While searching for resources, there were more creators asking the same question. "What is the best way to have elements in a NFT's metadata change based on certain events using IPFS?". This api is meant for devs looking to create dapps or services that give their users customizeable content on the InterPlanetary File System (IPFS)

# Problem

IPFS is a peer-to-peer file sharing protocol that has enabled users to host a wide variety of content. If someone wants to publish a website on IPFS they will have to direct users to the specific Content Identifier (CID) where they can download the site. If the Original Poster (OP) makes any changes to that site, they will get a completely different CID that they need to redirect their users to. To combat this, IPNS can offer one an address that can point to different CIDs. The InterPlanetary Name System (IPNS) was created to bring mutability and human readable addresses to IPFS. The problem today with IPNS is the revival of IPNS records. The records have to be republished by their expiration date to maintain discoverability and security. But it is not always possible for the node to be online to republish the record.  

Unfortunately, it is not possible to create IPNS records through js-ipfs (the browser or node.js) and that is why this project is important. This project is a go-ipfs implementation that can create, publish, and even maintain IPNS records. A client browser or server, can request keys and publish IPNS records using this project on the backend. That way devs can create dapps and not be bogged down by browser or version limitations.

# How it works

The content of IPFS hashes are immutable, once you change a single letter and resubmit to IPFS. You will get back a completely different CID. But IPNS identifiers can point to different CIDs. Once you have your IPNS record, you can change what it points to as many times as you want.

The project serves keys upon request and publishes content on a users behalf to IPFS. The idea is to enable access to the IPFS network and IPNS functionality without having to run your own IPFS node. A user can only publish IPNS content if they have keys to sign the records. To enable that, the api only stores keys for the needed operations. Each key is deleted from the local node keystore and temp_key storage before returning.

# Deep dive

As previously mentioned, IPNS records need to be republished after a certain period of time. The reason for this is because this "world wide record" called the Distributed Hash Table (DHT) keeps track of IPNS records and what they point to. You can read on why they expire here [go-ipfs/issue#1958](https://github.com/ipfs/go-ipfs/issues/1958#issuecomment-410860667)

# API 

There are ten endpoints currently:

getKey() - This will generate and return a private/public key pair in .key format. *The private keys are used to create, update, and republish IPNS records.*

postKey(<key_name>.key) - This will import the private key to node.

deleteKey(<keyName> string) - This will delete private key from node.

getRecord(<IPNS_key> string) - This will resolve what IPNS record points to and return IPFS Path. *Does not do continuous resolution aka IPNS Following*

postRecord(<IPFS_CID> string, <key_name>.key) - This will publish a brand new IPNS record and return IPNS path. Saves private key to allow for republishing.

putRecord(<IPNS_key> string, <key_name>.key) - This will resolve IPNS record and returns IPFS path. Saves private key to allow for republishing.

startFollowing(<IPNS_key> string) - This will add IPNS record to queue to allow for continuous resolution aka IPNS Following. Returns status 200 & IPNS_key. *Will not resolve record immediately use GetRecord() to resolve upon request.*

stopFollowing(<IPNS_key> string) - This will remove a key from the queue. Returns status 200 & IPNS_key.

add(file) - This will add content to IPFS. Will return CID.

add(dir) - This will add directory to IPFS. Will return CID.



# To Run
Make sure you have [golang](https://go.dev/doc/install) installed and [go-ipfs](https://github.com/ipfs/go-ipfs). You can find instructions on their respective websites.

After that is done, remember to fork, clone, & init ipfs:

`ipfs init` 

Then spin up local ipfs node:

`ipfs daemon`
*this command will hold up your terminal*

Then in a separate terminal window change directory into the repo:

`cd Fierro`

Now, you need to pull in the dependencies specified in server.go.

`go get`

From there you can run the server by:

`cd src/`

`go run .`

Now you can request keys and publish to IPNS! Checkout _Tests/requests.sh_ to see example _curl_ commands