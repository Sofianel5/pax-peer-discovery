# Test for sudo
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi
sysctl -w net.core.rmem_max=2500000
go build -o build/darkpool
echo "Darkpool built successfully"