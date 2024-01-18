# gomar

## Emulator for CoCo &amp; Motorola 6809, specializing in system programming and NitrOS-9 OS.

Work in progress.  See HINT at end of gomar.go.

Portions of https://github.com/strickyak/frobio
are designed to take advantage of this emulator.

But I really need to clean up some more, before documenting and advertizing.

Previously this code was in https://github.com/strickyak/doing_os9 .

## Hints (notes to self)

`go run -x --tags=coco3,level2,trace,d,cocoio gomar.go  --rom_a000 /home/strick/glap/6809/ROMS/color64bas.rom --rom_8000 /home/strick/glap/6809/ROMS/color64extbas.rom  -t 1 --basic_text --cart ~/coco-shelf/build-frobio/axiom4-whole.rom  2>_log`
