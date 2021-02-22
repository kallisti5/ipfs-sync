default:
	go fmt
	go build
	
run:
	go fmt
	go run main.go

rel:
	go fmt
	mkdir rel/

	CGO_ENABLED=0 GOOS=linux go build -o ipfs-sync
	upx ipfs-sync
	tar -caf ipfs-sync-linux64.tar.xz ipfs-sync LICENSE README.md systemd
	mv ipfs-sync-linux64.tar.xz rel/

	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o ipfs-sync
	upx ipfs-sync
	tar -caf ipfs-sync-linuxARM.tar.xz ipfs-sync LICENSE README.md systemd
	mv ipfs-sync-linuxARM.tar.xz rel/

	CGO_ENABLED=0 GOOS=darwin go build -o ipfs-sync
	upx ipfs-sync
	tar -caf ipfs-sync-darwin64.tar.gz ipfs-sync LICENSE README.md
	mv ipfs-sync-darwin64.tar.gz rel/

	CGO_ENABLED=0 GOOS=windows go build -o ipfs-sync.exe
	upx ipfs-sync.exe
	zip ipfs-sync-win64.zip ipfs-sync.exe LICENSE README.md
	mv ipfs-sync-win64.zip rel/