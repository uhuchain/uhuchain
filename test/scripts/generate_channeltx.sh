#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

CONFIGTXGEN_CMD="${CONFIGTXGEN_CMD:-configtxgen}"

declare -a channels=("car-ledger")
declare -a orgs=("InsuranceAMSP" "InsuranceBMSP" "InsuranceCMSP")

export CHANNEL_DIR="/opt/gopath/src/github.com/uhuchain/uhuchain/test/uhuchain-network-dev"
export FABRIC_CFG_PATH=${CHANNEL_DIR}

echo "Generating Orderer Genesis block"
$CONFIGTXGEN_CMD -profile ThreeOrgsOrdererGenesis -outputBlock ${CHANNEL_DIR}/channel-artifacts/car-ledger.genesis.block

for i in "${channels[@]}"
do
   echo "Generating artifacts for channel: $i"

   echo "Generating channel configuration transaction"
   $CONFIGTXGEN_CMD -profile ThreeOrgsChannel -outputCreateChannelTx .${CHANNEL_DIR}/${i}.tx -channelID $i

   for j in "${orgs[@]}"
   do
     echo "Generating anchor peer update for org $j"
     $CONFIGTXGEN_CMD -profile ThreeOrgsChannel -outputAnchorPeersUpdate ${CHANNEL_DIR}/${i}${j}anchors.tx -channelID $i -asOrg $j
   done
done
