BUILD=../build-frobio/
INT=../build-frobio/pizga/Internal/

rm -f /tmp/dump/dump-[0-9][0-9][0-9][0-9][0-9][0-9]

exec - gorun --tags=coco3,level2,trace,d,vdg,cocoio,gime,copico  \
	gomar.go                                              \
		-disk $INT/OS9DISKS/NOS9_6809_L2_coco3_80d.dsk    \
		-rom_8000 ../toolshed/cocoroms/coco3.rom          \
		--cart $BUILD/axiom41.rom                         \
		--borges /home/strick/borges/                     \
		--t=1                                             \
		--copico_server='127.0.0.1:2321'                  \
		    2>&1 | tee _tty                               \
            ##
