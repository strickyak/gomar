cd $(dirname $0)

#  go run --tags=coco1,noos,vdg,cocoio,trace  gomar.go -rom_8000 /tmp/rom1 -internal_rom_listing /tmp/list1 --cart ../build-frobio/axiom41.rom --show_vdg_screen=1 --bracket_terminal  -global_map ~/coco-shelf/nekot-coco-microkernel/kernel/_nekot.map -global_listing ~/coco-shelf/nekot-coco-microkernel/kernel/_nekot.o.list  -t=1 

echo "
  Suggestion:  in ../build-frobio do this:
      make run-lemma LAN=127.0.0.1 FORCE=--force_page=98
  or (if you configure your wired ethernet) this:
      make run-lemma LAN=10.23.23.23 FORCE=--force_page=98

  But first, if you are forcing page 98,
      cd ~/coco-shelf/nekot-coco-microkernel/nekot1
      make
" >&2

cat ../toolshed/cocoroms/extbas11.rom ../toolshed/cocoroms/bas13.rom > /tmp/rom1
cat ../toolshed/cocoroms/extbas11.rom.list ../toolshed/cocoroms/bas13.rom.list > /tmp/list1

set -x
go run --tags=coco1,noos,vdg,cocoio,$TRACE  gomar.go \
    -rom_8000 /tmp/rom1 \
    -internal_rom_listing /tmp/list1  \
    --cart ../build-frobio/axiom41.rom \
    --show_vdg_screen=1 --bracket_terminal \
    --global_map ../nekot-coco-microkernel/build-for-16k-cocoio/_nekot1.decb.map     \
    --global_listing ../nekot-coco-microkernel/build-for-16k-cocoio/_nekot1.o.list   \
    "$@"
