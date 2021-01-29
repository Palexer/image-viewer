#!/bin/sh

echo "starting compilation"

echo "linux-amd64"
fyne-cross linux -icon data/icon.png

echo "darwin-amd64"
fyne-cross darwin -icon data/icon.png

echo "windows-amd64"
fyne-cross windows -icon data/icon.png

echo

echo "creating .dmg from .app-bundle"

dd if=/dev/zero of=/tmp/ImageViewer.dmg bs=1M count=16 status=progress 
mkfs.hfsplus -v ImageViewer /tmp/ImageViewer.dmg

sudo mkdir -pv /mnt/tmp
sudo mount -o loop /tmp/ImageViewer.dmg /mnt/tmp
sudo cp -av fyne-cross/dist/darwin-amd64/ImageViewer.app /mnt/tmp

sudo umount /mnt/tmp

cp /tmp/ImageViewer.dmg fyne-cross/dist/darwin-amd64/ImageViewer.dmg

echo "dmg available at fyne-cross/dist/darwin-amd64/ImageViewer.dmg"
echo "Done"
