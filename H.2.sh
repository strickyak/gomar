(
  cat ../toolshed/cocoroms/extbas11.rom
  cat ../toolshed/cocoroms/bas13.rom 
) > /tmp/coco2rom

(
  cat ../toolshed/cocoroms/extbas11.rom.list
  cat ../toolshed/cocoroms/bas13.rom.list
) > /tmp/coco2rom.list


- gorun --tags=coco1,level1,trace,d,vdg gomar.go -t=1 -rom_8000 /tmp/coco2rom -internal_rom_listing  /tmp/coco2rom.list   --cart ../toolshed/cocoroms/disk11.rom -external_rom_listing ../toolshed/cocoroms/disk11.rom.list
