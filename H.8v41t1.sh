# - gorun --tags=coco3,level2,trace,d,vdg,cocoio,gime gomar.go -disk ../build-frobio/results/OS9DISKS/NOS9_6809_L2_coco3_80d.dsk -rom_8000 ../toolshed/cocoroms/coco3.rom -internal_rom_listing ../toolshed/cocoroms/coco3.rom.list --cart ../build-frobio/axiom4-whole.rom -external_rom_listing ~/coco-shelf/frobio/frob3/hdbdos/hdbdos.list --borges /home/strick/borges/  --trigger_pc=0xC002 --trigger_count=3 


# - gorun --tags=coco3,level2,trace,d,vdg,cocoio,gime gomar.go -disk ../build-frobio/results/OS9DISKS/NOS9_6809_L2_coco3_80d.dsk -rom_8000 ../toolshed/cocoroms/coco3.rom -internal_rom_listing ../toolshed/cocoroms/coco3.rom.list --cart ../build-frobio/axiom4-whole.rom -external_rom_listing /dev/null --borges /home/strick/borges/  --t=1

# - gorun --tags=coco3,level2,trace,d,vdg,cocoio,gime gomar.go -disk ../build-frobio/results/OS9DISKS/NOS9_6809_L2_coco3_80d.dsk -rom_8000 ../toolshed/cocoroms/coco3.rom -internal_rom_listing ../toolshed/cocoroms/coco3.rom.list --cart ../build-frobio/axiom4-whole.rom -external_rom_listing /dev/null --borges /home/strick/borges/  --trigger_pc=0xDA88  --trigger_count=3

- gorun --tags=coco3,level2,trace,d,vdg,cocoio,gime gomar.go -disk ../build-frobio/results/OS9DISKS/NOS9_6809_L2_coco3_80d.dsk -rom_8000 ../toolshed/cocoroms/coco3.rom -internal_rom_listing ../toolshed/cocoroms/coco3.rom.list --cart ../frobio/frob3/axiom41/axiom41.rom -external_rom_listing ../frobio/frob3/axiom41/axiom41.list --borges /home/strick/borges/  --t=1
