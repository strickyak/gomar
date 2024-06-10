# - gorun --tags=coco3,level2,trace,d,vdg,cocoio,gime gomar.go -disk ../build-frobio/results/OS9DISKS/NOS9_6809_L2_coco3_80d.dsk -rom_8000 ../toolshed/cocoroms/coco3.rom -internal_rom_listing ../toolshed/cocoroms/coco3.rom.list --cart ../build-frobio/axiom4-whole.rom -external_rom_listing ~/coco-shelf/frobio/frob3/hdbdos/hdbdos.list --borges /home/strick/borges/  --trigger_pc=0xC002 --trigger_count=3 


# - gorun --tags=coco3,level2,trace,d,vdg,cocoio,gime gomar.go -disk ../build-frobio/results/OS9DISKS/NOS9_6809_L2_coco3_80d.dsk -rom_8000 ../toolshed/cocoroms/coco3.rom -internal_rom_listing ../toolshed/cocoroms/coco3.rom.list --cart ../build-frobio/axiom4-whole.rom -external_rom_listing /dev/null --borges /home/strick/borges/  --t=1

# - gorun --tags=coco3,level2,trace,d,vdg,cocoio,gime gomar.go -disk ../build-frobio/results/OS9DISKS/NOS9_6809_L2_coco3_80d.dsk -rom_8000 ../toolshed/cocoroms/coco3.rom -internal_rom_listing ../toolshed/cocoroms/coco3.rom.list --cart ../build-frobio/axiom4-whole.rom -external_rom_listing /dev/null --borges /home/strick/borges/  --trigger_pc=0xDA88  --trigger_count=3

INKEY=--inkey='~~~~~~~~~~~~~~~~~~~~~~~~252
~~~~~~~~~~~~~~~@
'
SHOW='--show_vdg_screen=0'
SHOW='--show_vdg_screen=1 --show_gime_screen=1'

case "$1" in
	h1)
		INKEY='--inkey_file=inkey-test-data/decb_many_dirs.txt'
		SHOW='--show_vdg_screen=1'
		;;
	h2)
		INKEY='--inkey_file=inkey-test-data/decb_many_saves.txt'
		SHOW='--show_vdg_screen=1'
		;;
esac

- gorun \
  --tags=coco3,level2,trace,d,vdg,cocoio,gime gomar.go \
    -disk ../build-frobio/lemma/LEMMINGS/Rogue512.dsk \
      -rom_8000 ../toolshed/cocoroms/coco3.rom \
      -internal_rom_listing ../toolshed/cocoroms/coco3.rom.list \
          --cart ../build-frobio/axiom41.rom \
	  -external_rom_listing ../build-frobio/hdbdos.rom.list \
  --borges /home/strick/borges/  \
  --t=1  \
		  $SHOW \
		  "$INKEY"
#   --trigger_count=1 \
#    -disk ../build-frobio/lemma/OS9DISKS/NOS9_6809_L2_coco3_80d.dsk \
exit $?
