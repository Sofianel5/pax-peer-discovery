#!/bin/bash
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi
apt upgrade -y
wget https://go.dev/dl/go1.19.4.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.19.4.linux-amd64.tar.gz
rm go1.19.4.linux-amd64.tar.gz
# Add to path - cannot be done from sudo
#echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
#. ~/.profile