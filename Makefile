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
	wget -O /tmp/lz4_win64_v1_9_4.zip https://github.com/lz4/lz4/releases/download/v1.9.4/lz4_win64_v1_9_4.zip
	cd /tmp/ && mkdir lib_ && unzip /tmp/lz4_win64_v1_9_4.zip -d lib_
	mv /tmp/lib_/static/liblz4_static.lib deps_win/deps-x86_64/lib/liblz4.a
	rm -rf deps_win/deps-x86_64/lib/liblz4.dll.a /tmp/lz4_win64_v1_9_4.zip /tmp/lib_
	
build: build_library_linux build_library_windows patch_ovpncli patch_wrap copy_source_code copy_header_lib patch_header patch_bind_go

integr_test:
	sudo go test -race

format:
	gofmt -w .
