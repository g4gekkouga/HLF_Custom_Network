# HLF_Custom_Network

This reposiroty includes docker files and scripts to deploy a custom Hyperledger Fabric network and deploy chaincodes to the network.

## Network Configuration

No of Orgs : 3  
No of Peers : 1 Peer per Org  
No of Orderers : 3  
Consensus Mechanism : Raft Consensus Protocol  
State Database : CouchDB  

## Steps to Deploy

1. Deploy CA nodes using the docker-compose file in artifacts/channel/create-certificate-with-ca folder
2. Create crypto material for peer and orderer nodes to join the network using the create-certificate-with-ca.sh script in the same folder
3. Create necessary channel artifacts using create-artifacts.sh in artifacts/channel folder
4. Deploy the peer and orderer nodes using the docker-compose file in artifacts folder
5. Create a channel and join all the peers to this channel using createChannel.sh
