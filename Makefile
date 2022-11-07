SHELL := /bin/bash

dir_linux_build := "third_party/openvpn3/build/linux/client/"
dir_windows_build := "third_party/openvpn3/build/windows/client/"

patch_ovpncli:
	patch ovpncli.go patches/ovpncli_build.patch

patch_wrap:
	patch ovpncli_wrap.cxx patches/ovpncli_wrap_warnings.patch
	patch ovpncli_wrap.cxx patches/add_license_wrap_cxx.patch

patch_header:
	patch ovpncli_wrap.h patches/add_license_wrap_header.patch

patch_bind_go:		
	patch ovpncli.go patches/add_license_ovpncli_go.patch
	
build_library_linux:
	bash ./third_party/openvpn3/docker-build_linux.sh
	cp -R  $(dir_linux_build)/* .
	rm -rf Readme
	tar -xvzpf build_linux_deps.tar.gz && rm -rf build_linux_deps.tar.gz

copy_header_lib:
	cp ./third_party/openvpn3/client/ovpncli.hpp .
	
copy_source_code:
	cp -R ./third_party/openvpn3/openvpn .
		
build_library_windows:
	bash ./third_party/openvpn3/docker-build_win.sh
	cp -R  $(dir_windows_build)/* .
	rm -rf Readme
	tar -xvzpf build_windows_deps.tar.gz && rm -rf build_windows_deps.tar.gz
	
build: build_library_linux build_library_windows patch_ovpncli patch_wrap copy_source_code copy_header_lib patch_header patch_bind_go

integr_test:
	sudo go test -race

format:
	gofmt -w .
