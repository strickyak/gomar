set -ex

(
  cat ../toolshed/cocoroms/extbas11.rom
  cat ../toolshed/cocoroms/bas13.rom 
) > /tmp/coco2rom

(
  cat ../toolshed/cocoroms/extbas11.rom.list
  cat ../toolshed/cocoroms/bas13.rom.list
) > /tmp/coco2rom.list

cp -vf ../nitros9/level1/coco1/NOS9_6809_L1_coco1_40d_1.dsk /tmp/NOS9_6809_L1_coco1_40d_1.dsk


go run --tags=coco1,level1,trace,d,vdg gomar.go -rom_8000 /tmp/coco2rom -internal_rom_listing  /tmp/coco2rom.list   --cart ../toolshed/cocoroms/disk11.rom -external_rom_listing ../toolshed/cocoroms/disk11.rom.list -show_vdg_screen=1  -disk /tmp/NOS9_6809_L1_coco1_40d_1.dsk
