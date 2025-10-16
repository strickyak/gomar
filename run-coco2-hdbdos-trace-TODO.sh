cat ../toolshed/cocoroms/extbas11.rom ../toolshed/cocoroms/bas13.rom > /tmp/rom1
cat ../toolshed/cocoroms/extbas11.rom.list ../toolshed/cocoroms/bas13.rom.list > /tmp/list1.list

go run --tags=coco1,noos,vdg,cocoio,trace  gomar.go \
        --cart ~/modoc/coco-shelf/toolshed/hdbdos/hdbdw3cc2.rom  \
        --external_rom_listing ~/modoc/coco-shelf/toolshed/hdbdos/hdbdw3cc2.rom.list   \
        -rom_8000 /tmp/rom1 \
        -internal_rom_listing /tmp/list1.list \
        --show_vdg_screen=1 \
        --bracket_terminal \
        -t=1 -max=10000000 \
        -v=amit  ## ?-v
