# Hint: EXEC &HC000

cd $(dirname $0)

cat ../toolshed/cocoroms/extbas11.rom ../toolshed/cocoroms/bas13.rom > /tmp/rom1
cat ../toolshed/cocoroms/extbas11.rom.list ../toolshed/cocoroms/bas13.rom.list > /tmp/list1

set -x
go run --tags=coco1,noos,vdg,cocoio,$TRACE gomar.go \
    -rom_8000 /tmp/rom1 \
    -internal_rom_listing /tmp/list1  \
    --show_vdg_screen=1 --bracket_terminal \
    --cart ../collapseos-coco2/os.bin \
    "$@"
