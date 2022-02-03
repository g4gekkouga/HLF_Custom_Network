export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${PWD}/artifacts/channel/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export PEER0_ORG1_CA=${PWD}/artifacts/channel/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export PEER0_ORG2_CA=${PWD}/artifacts/channel/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export PEER0_ORG3_CA=${PWD}/artifacts/channel/crypto-config/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt
export FABRIC_CFG_PATH=${PWD}/artifacts/channel/config/

export CHANNEL_NAME=mychannel

setGlobalsForPeer0Org1() {
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
}

setGlobalsForPeer0Org2() {
    export CORE_PEER_LOCALMSPID="Org2MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/artifacts/channel/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
    export CORE_PEER_ADDRESS=localhost:9051
}

CHANNEL_NAME="mychannel"
CC_NAME="adcc"


chaincodeQuery() {
    setGlobalsForPeer0Org2

    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "GetResource","Args":["subject1", "object1", "read"]}'

    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "GetResource","Args":["subject2", "object2", "read"]}'

    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "GetResource","Args":["subject1", "object3", "read"]}'

    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"function": "GetResource","Args":["subject3", "object1", "read"]}'
}

chaincodeQuery
