#!/bin/sh

echo "Start the Network"
minifab netup -s couchdb -e true -i 2.4.8 -o insuranceCompany.insuranceClaimPostAccident.com

sleep 5

echo "create a channel"
minifab create -c autochannel

sleep 2

echo "join the peers to this channel"
minifab join -c autochannel

sleep 2

echo "Anchor update"
minifab anchorupdate

sleep 2

echo "Profile Generation"
minifab profilegen
