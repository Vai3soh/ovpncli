# ovpncli - Golang wrapping client call openvpn3 C++ library

Library for wrapping openvpn3 (https://github.com/OpenVPN/openvpn3) functionality in go way

For bindings use Swig library: https://github.com/swig/swig

Build for linux, windows (x64) OS.

How to build (use docker):
```
git clone --recursive https://github.com/Vai3soh/ovpncli
cd go-openvpn-ng
make build 
```
How to use:
    see test in ovpn_test.go
