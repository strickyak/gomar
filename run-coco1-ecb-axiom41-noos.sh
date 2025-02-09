cd $(dirname $0)

echo "
  Suggestion:  in ../build-frobio do this:
      make run-lemma LAN=127.0.0.1 FORCE=--force_page=98
  or (if you configure your wired ethernet) this:
      make run-lemma LAN=10.23.23.23 FORCE=--force_page=98

  But first, if you are forcing page 98,
      cd ~/coco-shelf/nekot-coco-microkernel/kernel
      make
" >&2

cat ../toolshed/cocoroms/extbas11.rom ../toolshed/cocoroms/bas13.rom > /tmp/rom1
cat ../toolshed/cocoroms/extbas11.rom.list ../toolshed/cocoroms/bas13.rom.list > /tmp/list1

set -x
go run --tags=coco1,noos,vdg,cocoio, gomar.go \
    -rom_8000 /tmp/rom1 \
    -internal_rom_listing /tmp/list1  \
    --cart ../build-frobio/axiom41.rom \
    --show_vdg_screen=1 --bracket_terminal \
    "$@"
