BUILD=../build-frobio/
INT=../build-frobio/pizga/Internal/

rm -f /tmp/dump/dump-[0-9][0-9][0-9][0-9][0-9][0-9]

exec - gorun --tags=coco3,level2,vdg,cocoio,gime          \
	gomar.go                                                  \
		-disk $INT/OS9DISKS/NOS9_6809_L2_coco3_80d.dsk    \
		-rom_8000 ../toolshed/cocoroms/coco3.rom          \
		--cart $BUILD/axiom41.rom                         \
		    2>&1 | tee _tty                               \
		##
		##
		##




#########################################

# - gorun --tags=coco3,level2,trace,d,vdg,cocoio,gime gomar.go -disk $INT/OS9DISKS/NOS9_6809_L2_coco3_80d.dsk -rom_8000 ../toolshed/cocoroms/coco3.rom -internal_rom_listing ../toolshed/cocoroms/coco3.rom.list --cart $BUILD/axiom41.rom -external_rom_listing $BUILD/axiom41.rom.list --borges /home/strick/borges/  --trigger_os9="(?i:F.Fork.*file='dir')"  "$@"

# - gorun --tags=coco3,level2,trace,d,vdg,cocoio,gime gomar.go -disk $INT/OS9DISKS/NOS9_6809_L2_coco3_80d.dsk -rom_8000 ../toolshed/cocoroms/coco3.rom -internal_rom_listing ../toolshed/cocoroms/coco3.rom.list --cart $BUILD/axiom41.rom -external_rom_listing $BUILD/axiom41.rom.list --borges /home/strick/borges/  --trigger_os9='(?i:Chain.*Module/file=.Shell.)'  "$@"

# - gorun --tags=coco3,level2,trace,d,vdg,cocoio,gime gomar.go -disk ../build-frobio/results/OS9DISKS/NOS9_6809_L2_coco3_80d.dsk -rom_8000 ../toolshed/cocoroms/coco3.rom -internal_rom_listing ../toolshed/cocoroms/coco3.rom.list --cart ../build-frobio/axiom4-whole.rom -external_rom_listing ~/coco-shelf/frobio/frob3/hdbdos/hdbdos.list --borges /home/strick/borges/  --trigger_pc=0xC002 --trigger_count=3 


# - gorun --tags=coco3,level2,trace,d,vdg,cocoio,gime gomar.go -disk ../build-frobio/results/OS9DISKS/NOS9_6809_L2_coco3_80d.dsk -rom_8000 ../toolshed/cocoroms/coco3.rom -internal_rom_listing ../toolshed/cocoroms/coco3.rom.list --cart ../build-frobio/axiom4-whole.rom -external_rom_listing /dev/null --borges /home/strick/borges/  --t=1

# - gorun --tags=coco3,level2,trace,d,vdg,cocoio,gime gomar.go -disk ../build-frobio/results/OS9DISKS/NOS9_6809_L2_coco3_80d.dsk -rom_8000 ../toolshed/cocoroms/coco3.rom -internal_rom_listing ../toolshed/cocoroms/coco3.rom.list --cart ../build-frobio/axiom4-whole.rom -external_rom_listing /dev/null --borges /home/strick/borges/  --trigger_pc=0xDA88  --trigger_count=3

