buildall:
	fyne-cross linux -icon data/icon.png
	fyne-cross darwin -icon data/icon.png
	fyne-cross windows -icon data/icon.png

	# turn macOS .app to dmg
	dd if=/dev/zero of=/tmp/ImageViewer.dmg bs=1M count=16 status=progress 
	mkfs.hfsplus -v ImageViewer /tmp/ImageViewer.dmg

	sudo mkdir -pv /mnt/tmp
	sudo mount -o loop /tmp/ImageViewer.dmg /mnt/tmp
	sudo cp -av fyne-cross/dist/darwin-amd64/ImageViewer.app /mnt/tmp

	sudo umount /mnt/tmp

	cp /tmp/ImageViewer.dmg fyne-cross/dist/darwin-amd64/ImageViewer.dmg

clean:
	if [ -d "fyne-cross" ]; then rm -r fyne-cross; fi
