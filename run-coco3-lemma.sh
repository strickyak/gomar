go run --tags=coco3,level2,vdg,cocoio,gime,"$TAGS"  gomar.go -rom_8000 ../toolshed/cocoroms/coco3.rom  -internal_rom_listing ../toolshed/cocoroms/coco3.rom.list  --cart ../build-frobio/axiom41.rom --show_vdg_screen=1 --bracket_terminal "$@"
