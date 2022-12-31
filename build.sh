# Check max buffer size for UDP is 2.5MB, if not, set it
if [ $(sysctl -n net.core.rmem_max) -lt 2621440 ]; then
    echo "Setting max buffer size for UDP to 2.5MB"
    sudo sysctl -w net.core.rmem_max=2621440
fi
go build -o build/darkpool
echo "Darkpool built successfully"