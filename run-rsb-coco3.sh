## + go build --tags=coco3,level2,vdg,cocoio,gime, -o /tmp/for-hasty-20250202-161523/gomar.coco3.level2.vdg gomar.go
## + /usr/bin/go build -o /tmp/for-hasty-20250202-161523/waiter -x lemma-waiter.go

# mkdir -p /tmp/for-rsb
# cp ../toolshed/cocoroms/coco3.rom ../toolshed/cocoroms/coco3.rom.list /tmp/for-rsb/
#    -disk /dev/null \

TRACE=,
WHAT=coco3,level2,gime

while true
do
    case "$1" in
        +TRACE )
            TRACE=,trace
            ;;
        +* )
            echo "UNKNOWN + OPTION:  $1" >&2
            exit 13
            ;;
        * )
            break
            ;;
    esac
    shift
done

set -x
go run --tags="$WHAT,vdg,cocoio,$TRACE," gomar.go \
    -rom_8000 ../toolshed/cocoroms/coco3.rom \
    -internal_rom_listing ../toolshed/cocoroms/coco3.rom.list \
    --show_vdg_screen=1 \
    --bracket_terminal \
    "$@"



