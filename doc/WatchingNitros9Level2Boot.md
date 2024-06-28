# Watching Nitros9 Level2 Boot

`<strick>` Henry Strickland `<@yak.net>` Jun 28, 2024

Here's a look at all the Kernel System Calls that Nitros9 Level2 makes
on a coco3
while booting up to an interactive shell.

## What we're looking at

Because I'm particularly interested in some directory operations at the
moment, the `STARTUP` file contains this:

```
dir /dd
dir /b1
```

where both `/dd` (the default drive) and `/b1` are instances of the
RBLemma device, a network block device with synchronous operations.
That is, it does not use interrupts; rather it waits until the data
arrives.  And it just so happens that the contents of `/dd` and `/b1`
are identical, so the output of the `dir` commands will be the same.

(That `STARTUP` file contents is defined in
`coco-shelf/frobio/frob3/lemmings/Nitros9_Coco3_M6809_Level2_40col_lem.go` )

The coco3 Level2 kernel was built from https://github.com/nitros9project/nitros9
using the Coco Shelf ( https://github.com/strickyak/coco-shelf ).
All github repositories should be accurate at the end of the day on Jun 28, 2024.

## How I ran it

We're looking at a log of Kernel System Calls and their corresponding Return
operations in the logging output of the Gomar emulator
( https://github.com/strickyak/gomar ).  Gomar was run with these
parameters:

```
BUILD=../build-frobio/
INT=../build-frobio/pizga/Internal/

gorun --tags=coco3,level2,trace,d,vdg,cocoio,gime gomar.go \
        -disk $INT/OS9DISKS/NOS9_6809_L2_coco3_80d.dsk     \
	        -rom_8000 ../toolshed/cocoroms/coco3.rom   \
		        --cart $BUILD/axiom41.rom
```

where `gorun` is a script in `coco-shelf/bin` that invokes `go run` with
some environment variables set.  (`go` is the command to invoke the compiler
for the Go Language ( https://go.dev/ ) which must be version 1.18 
or later.)

But before you run that `gorun` command, in another terminal window,
you must start a Lemma Server, which by defaults listens on your
`localhost` for connections from the emulated booting coco3's CocoIO
ethernet card (with the Wiznet W5100S chip).  In the `coco-shelf/build-frobio`
directory, type this command and let it run until the emulator is finished.
Then you can hit Control-C to kill it.

```
cd coco-shelf/build-frobio
make all run-lemma   FORCE=--force_page=3272
```

Instead of page 3272 (which has the experimental LemMan file manager
to be debugged) you could use page 327 or various other
Nitros9 kernels.

## The Output

The verbose debugging output, by default, goes into a
file named `_log`.  We run a short Python script to extract
just the Kernel System Calls, with some pretty-printing
enhancements, into the file `_kern`:

```
python3 ../frobio/frob3/helper/show_kernel_calls.py < _log > _kern
```

From `_kern` we can do some greps to create shorter versions:

```
egrep -v 'F[$](VModul|All64|Ret64|LDABX|PrsNam|SRqMem|Move|LDDDXY|Find64|SRtMem|AlHRAM)\>' _kern > _short
egrep -w t1 _short > _user
```

The return operations start with an arrow `<----`.  The words `t0` and `t1`
refer to the MMU being in task 0 or task 1, which corresponds to Nitros9 being
in system mode (t0) or in user mode (t1).  The kernel often makes lots of calls
in system mode as part of a user mode call.

## Shortest: Just the system calls from User (t1) space.

This is the file `_user`, created above.
It only shows calls made from User Space (t1), and their corresponding
return operations.

The trailing numbers like `#8411284` show the cycle
counts as counted by the emulator.  Due to the emulator doing a poor
job of having realistic interrupts, they may vary a lot from an
actual coco.

```
{proc=2"SysGo"} t1: OS9$09 <F$Icpt> {routine=da8b storage=0000} #8301674
:   <---- t1 OS9$09 <F$Icpt> {routine=da8b storage=0000} #8301674 #8303952
{proc=2"SysGo"} t1: OS9$0c <F$ID> {} #8303964
:   <---- t1 OS9$0c <F$ID> {} #8303964 #8306236
{proc=2"SysGo"} t1: OS9$0d <F$SPrior> {pid=02 priority=80} #8306252
:   <---- t1 OS9$0d <F$SPrior> {pid=02 priority=80} #8306252 #8308656
{proc=2"SysGo"} t1: OS9$00 <F$Link> {type/lang=00 module/file='Init'} #8308696
:   <---- t1 OS9$00 <F$Link> {type/lang=00 module/file='Init'} #8308696 -> addr $9e00 entry $ade0 #8352448
{proc=2"SysGo"} t1: OS9$8a <I$Write> {NitrOS-9/6809 Level 2 V3.3.0} #8353444
<---- t1 OS9$8a <I$Write> {NitrOS-9/6809 Level 2 V3.3.0} #8353444 #8411218
{proc=2"SysGo"} t1: OS9$8c <I$WritLn> {} #8411284
:   <---- t1 OS9$8c <I$WritLn> {} #8411284 #8423858
{proc=2"SysGo"} t1: OS9$8a <I$Write> {Tandy Color Computer 3} #8424676
:   <---- t1 OS9$8a <I$Write> {Tandy Color Computer 3} #8424676 #8440428
{proc=2"SysGo"} t1: OS9$8c <I$WritLn> {} #8440494
:   <---- t1 OS9$8c <I$WritLn> {} #8440494 #8453068
{proc=2"SysGo"} t1: OS9$8a <I$Write> {(C) 2014 The NitrOS-9 Project
:   <---- t1 OS9$8a <I$Write> {(C) 2014 The NitrOS-9 Project
{proc=2"SysGo"} t1: OS9$16 <F$STime> {y85 m12 d31 h23 m59 s59} #8561324
:   <---- t1 OS9$16 <F$STime> {y85 m12 d31 h23 m59 s59} #8561324 #8618608
{proc=2"SysGo"} t1: OS9$86 <I$ChgDir> {mode=04, d93a='CMDS'} #8618642
:   <---- t1 OS9$86 <I$ChgDir> {mode=04, d93a='CMDS'} #8618642 #8947230
{proc=2"SysGo"} t1: OS9$86 <I$ChgDir> {mode=01, d932='/DD'} #8947264
:   <---- t1 OS9$86 <I$ChgDir> {mode=01, d932='/DD'} #8947264 #9163556
{proc=2"SysGo"} t1: OS9$86 <I$ChgDir> {mode=04, d936='/DD/CMDS'} #9163596
:   <---- t1 OS9$86 <I$ChgDir> {mode=04, d936='/DD/CMDS'} #9163596 #9488494
{proc=2"SysGo"} t1: OS9$0c <F$ID> {} #9488524
:   <---- t1 OS9$0c <F$ID> {} #9488524 #9490796
{proc=2"SysGo"} t1: OS9$18 <F$GPrDsc> {} #9490828
:   <---- t1 OS9$18 <F$GPrDsc> {} #9490828 #9504954
{proc=2"SysGo"} t1: OS9$4f <F$MapBlk> {beginningBlock=0 numBlocks=1} #9504996
:   <---- t1 OS9$4f <F$MapBlk> {beginningBlock=0 numBlocks=1} #9504996 #9509802
{proc=2"SysGo"} t1: OS9$8d <I$GetStt> {path=1 27==SS.KySns  : Getstat/SetStat for COCO keyboard} #9510374
:   <---- t1 OS9$8d <I$GetStt> {path=1 27==SS.KySns  : Getstat/SetStat for COCO keyboard} #9510374 #9515938
{proc=2"SysGo"} t1: OS9$03 <F$Fork> {Module/file='Shell' param="startup{32}-p{13}" lang/type=1 pages=0} #9516016
:   <---- t1 OS9$03 <F$Fork> {Module/file='Shell' param="startup{32}-p{13}" lang/type=1 pages=0} #9516016 #12769468
{proc=2"SysGo"} t1: OS9$04 <F$Wait> {} #12769486
{proc=3"Shell"} t1: OS9$09 <F$Icpt> {routine=e06b storage=0000} #12780034
:   <---- t1 OS9$09 <F$Icpt> {routine=e06b storage=0000} #12780034 #12782312
{proc=3"Shell"} t1: OS9$0c <F$ID> {} #12782364
:   <---- t1 OS9$0c <F$ID> {} #12782364 #12784636
{proc=3"Shell"} t1: OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #12787256
:   <---- t1 OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #12787256 #12795290
{proc=3"Shell"} t1: OS9$84 <I$Open> {f7f4='.'} #12796032
:   <---- t1 OS9$84 <I$Open> {f7f4='.'} #12796032 -> path $3 #13143592
{proc=3"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #13143686
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #13143686 #13197682
{proc=3"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #13197802
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #13197802 #13205358
{proc=3"Shell"} t1: OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #13205626
:   <---- t1 OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #13205626 #13213688
{proc=3"Shell"} t1: OS9$8f <I$Close> {path=3} #13217092
:   <---- t1 OS9$8f <I$Close> {path=3} #13217092 #13229802
{proc=3"Shell"} t1: OS9$1c <F$SUser> {} #13232352
:   <---- t1 OS9$1c <F$SUser> {} #13232352 #13234606
{proc=3"Shell"} t1: OS9$21 <F$NMLink> {LangType=00, 0e6d='startup'} #13257724
{proc=3"Shell"} t1: OS9$84 <I$Open> {0e6d='startup'} #13315360
{proc=3"Shell"} t1: OS9$82 <I$Dup> {$0} #14392028
:   <---- t1 OS9$82 <I$Dup> {$0} #14392028 -> path $3 #14395356
{proc=3"Shell"} t1: OS9$8f <I$Close> {path=0} #14395406
:   <---- t1 OS9$8f <I$Close> {path=0} #14395406 #14399308
{proc=3"Shell"} t1: OS9$84 <I$Open> {0e6d='startup'} #14399528
:   <---- t1 OS9$84 <I$Open> {0e6d='startup'} #14399528 -> path $0 #14807940
{proc=3"Shell"} t1: OS9$21 <F$NMLink> {LangType=11, e00d='Shell'} #14808826
:   <---- t1 OS9$21 <F$NMLink> {LangType=11, e00d='Shell'} #14808826 #14832058
{proc=3"Shell"} t1: OS9$03 <F$Fork> {Module/file='Shell' param="-P{32}X{32}PATH=;-p{13}" lang/type=11 pages=1f} #14832254
:   <---- t1 OS9$03 <F$Fork> {Module/file='Shell' param="-P{32}X{32}PATH=;-p{13}" lang/type=11 pages=1f} #14832254 #14888870
{proc=3"Shell"} t1: OS9$1d <F$UnLoad> {} #14888944
:   <---- t1 OS9$1d <F$UnLoad> {} #14888944 #14909610
{proc=3"Shell"} t1: OS9$8f <I$Close> {path=0} #14909852
:   <---- t1 OS9$8f <I$Close> {path=0} #14909852 #14915110
{proc=3"Shell"} t1: OS9$82 <I$Dup> {$3} #14915132
:   <---- t1 OS9$82 <I$Dup> {$3} #14915132 -> path $0 #14918358
{proc=3"Shell"} t1: OS9$8f <I$Close> {path=3} #14918408
:   <---- t1 OS9$8f <I$Close> {path=3} #14918408 #14922310
{proc=3"Shell"} t1: OS9$04 <F$Wait> {} #14922642
{proc=4"Shell"} t1: OS9$09 <F$Icpt> {routine=e06b storage=0000} #14933190
:   <---- t1 OS9$09 <F$Icpt> {routine=e06b storage=0000} #14933190 #14935468
{proc=4"Shell"} t1: OS9$0c <F$ID> {} #14935520
:   <---- t1 OS9$0c <F$ID> {} #14935520 #14937792
{proc=4"Shell"} t1: OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #14940412
:   <---- t1 OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #14940412 #14948446
{proc=4"Shell"} t1: OS9$84 <I$Open> {f7f4='.'} #14949188
:   <---- t1 OS9$84 <I$Open> {f7f4='.'} #14949188 -> path $3 #15334764
{proc=4"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #15334858
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #15334858 #15392058
{proc=4"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #15392178
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #15392178 #15399734
{proc=4"Shell"} t1: OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #15400002
:   <---- t1 OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #15400002 #15406430
{proc=4"Shell"} t1: OS9$8f <I$Close> {path=3} #15409834
:   <---- t1 OS9$8f <I$Close> {path=3} #15409834 #15427734
{proc=4"Shell"} t1: OS9$1c <F$SUser> {} #15430386
:   <---- t1 OS9$1c <F$SUser> {} #15430386 #15432640
{proc=4"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #15461340
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #15461340 #15465484
{proc=4"Shell"} t1: OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #15465526
:   <---- t1 OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #15465526 #15471212
{proc=4"Shell"} t1: OS9$1c <F$SUser> {} #15471490
:   <---- t1 OS9$1c <F$SUser> {} #15471490 #15473744
{proc=4"Shell"} t1: OS9$8b <I$ReadLn> {} #15473886
:   <---- t1 OS9$8b <I$ReadLn> {} #15473886 #15530026
{proc=4"Shell"} t1: OS9$8b <I$ReadLn> {} #15530422
:   <---- t1 OS9$8b <I$ReadLn> {} #15530422 #15540560
{proc=4"Shell"} t1: OS9$1c <F$SUser> {} #15544186
:   <---- t1 OS9$1c <F$SUser> {} #15544186 #15546440
{proc=4"Shell"} t1: OS9$21 <F$NMLink> {LangType=00, 0e6d='dir'} #15566544
{proc=4"Shell"} t1: OS9$84 <I$Open> {0e6d='dir'} #15622434
:   <---- t1 OS9$84 <I$Open> {0e6d='dir'} #15622434 -> path $3 #16110232
{proc=4"Shell"} t1: OS9$89 <I$Read> {path=3 buf=00d6 size=4d} #16110274
:   <---- t1 OS9$89 <I$Read> {path=3 buf=00d6 size=4d} #16110274 #16166644
{proc=4"Shell"} t1: OS9$8f <I$Close> {path=3} #16170222
:   <---- t1 OS9$8f <I$Close> {path=3} #16170222 #16186920
{proc=4"Shell"} t1: OS9$21 <F$NMLink> {LangType=11, 0e6d='dir'} #16187084
{proc=4"Shell"} t1: OS9$22 <F$NMLoad> {LangType=11, 0e6d='dir'} #16242912
:   <---- t1 OS9$22 <F$NMLoad> {LangType=11, 0e6d='dir'} #16242912 #17054712
{proc=4"Shell"} t1: OS9$03 <F$Fork> {Module/file='dir' param="/dd{13}" lang/type=11 pages=1f} #17054920
:   <---- t1 OS9$03 <F$Fork> {Module/file='dir' param="/dd{13}" lang/type=11 pages=1f} #17054920 #17097596
{proc=4"Shell"} t1: OS9$1d <F$UnLoad> {} #17097660
:   <---- t1 OS9$1d <F$UnLoad> {} #17097660 #17104516
{proc=4"Shell"} t1: OS9$04 <F$Wait> {} #17105078
{proc=5"Dir"} t1: OS9$8d <I$GetStt> {path=1 26==SS.ScSiz  : Return screen size for COCO} #17108792
:   <---- t1 OS9$8d <I$GetStt> {path=1 26==SS.ScSiz  : Return screen size for COCO} #17108792 #17113298
{proc=5"Dir"} t1: OS9$84 <I$Open> {1efc='/dd'} #17113842
:   <---- t1 OS9$84 <I$Open> {1efc='/dd'} #17113842 -> path $3 #17391174
{proc=5"Dir"} t1: OS9$86 <I$ChgDir> {mode=81, 1efc='/dd'} #17391222
:   <---- t1 OS9$86 <I$ChgDir> {mode=81, 1efc='/dd'} #17391222 #17648874
{proc=5"Dir"} t1: OS9$15 <F$Time> {buf=d} #17651396
:   <---- t1 OS9$15 <F$Time> {buf=d} #17651396 #17655286
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #17658852
:   <---- t1 OS9$8c <I$WritLn> {} #17658852 #17688664
{proc=5"Dir"} t1: OS9$88 <I$Seek> {path=3 pos=00000040} #17688754
:   <---- t1 OS9$88 <I$Seek> {path=3 pos=00000040} #17688754 #17692280
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17692342
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17692342 #17752874
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17754374
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17754374 #17761930
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17763514
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17763514 #17771070
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17772582
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17772582 #17781772
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #17782536
:   <---- t1 OS9$8c <I$WritLn> {} #17782536 #17814172
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17814242
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17814242 #17821798
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17823250
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17823250 #17830946
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17832416
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17832416 #17893938
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17897902
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17897902 #17905458
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #17906466
:   <---- t1 OS9$8c <I$WritLn> {} #17906466 #17938386
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17938456
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17938456 #17946012
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17947440
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #17952932
:   <---- t1 OS9$8c <I$WritLn> {} #17952932 #17973020
{proc=5"Dir"} t1: OS9$06 <F$Exit> {status=0} #17973058
:   :   <---- t1 OS9$04 <F$Wait> {} #17105078 #18014642
{proc=4"Shell"} t1: OS9$8b <I$ReadLn> {} #18020766
:   <---- t1 OS9$8b <I$ReadLn> {} #18020766 #18029270
{proc=4"Shell"} t1: OS9$1c <F$SUser> {} #18032896
:   <---- t1 OS9$1c <F$SUser> {} #18032896 #18035150
{proc=4"Shell"} t1: OS9$21 <F$NMLink> {LangType=00, 0e6d='dir'} #18057706
{proc=4"Shell"} t1: OS9$84 <I$Open> {0e6d='dir'} #18113704
:   <---- t1 OS9$84 <I$Open> {0e6d='dir'} #18113704 -> path $3 #18611338
{proc=4"Shell"} t1: OS9$89 <I$Read> {path=3 buf=00d6 size=4d} #18613832
:   <---- t1 OS9$89 <I$Read> {path=3 buf=00d6 size=4d} #18613832 #18670086
{proc=4"Shell"} t1: OS9$8f <I$Close> {path=3} #18676496
:   <---- t1 OS9$8f <I$Close> {path=3} #18676496 #18692814
{proc=4"Shell"} t1: OS9$21 <F$NMLink> {LangType=11, 0e6d='dir'} #18692978
{proc=4"Shell"} t1: OS9$22 <F$NMLoad> {LangType=11, 0e6d='dir'} #18748966
:   <---- t1 OS9$22 <F$NMLoad> {LangType=11, 0e6d='dir'} #18748966 #19591372
{proc=4"Shell"} t1: OS9$03 <F$Fork> {Module/file='dir' param="/b1{13}" lang/type=11 pages=1f} #19591580
:   <---- t1 OS9$03 <F$Fork> {Module/file='dir' param="/b1{13}" lang/type=11 pages=1f} #19591580 #19634256
{proc=4"Shell"} t1: OS9$1d <F$UnLoad> {} #19634320
:   <---- t1 OS9$1d <F$UnLoad> {} #19634320 #19641176
{proc=4"Shell"} t1: OS9$04 <F$Wait> {} #19641738
{proc=5"Dir"} t1: OS9$8d <I$GetStt> {path=1 26==SS.ScSiz  : Return screen size for COCO} #19645452
:   <---- t1 OS9$8d <I$GetStt> {path=1 26==SS.ScSiz  : Return screen size for COCO} #19645452 #19649958
{proc=5"Dir"} t1: OS9$84 <I$Open> {1efc='/b1'} #19650502
:   <---- t1 OS9$84 <I$Open> {1efc='/b1'} #19650502 -> path $3 #19923422
{proc=5"Dir"} t1: OS9$86 <I$ChgDir> {mode=81, 1efc='/b1'} #19925922
:   <---- t1 OS9$86 <I$ChgDir> {mode=81, 1efc='/b1'} #19925922 #20180338
{proc=5"Dir"} t1: OS9$15 <F$Time> {buf=d} #20182860
:   <---- t1 OS9$15 <F$Time> {buf=d} #20182860 #20186750
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #20187864
:   <---- t1 OS9$8c <I$WritLn> {} #20187864 #20217724
{proc=5"Dir"} t1: OS9$88 <I$Seek> {path=3 pos=00000040} #20217814
:   <---- t1 OS9$88 <I$Seek> {path=3 pos=00000040} #20217814 #20221340
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20221402
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20221402 #20286108
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20287608
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20287608 #20295164
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20296748
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20296748 #20304304
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20305816
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20305816 #20315006
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #20315770
:   <---- t1 OS9$8c <I$WritLn> {} #20315770 #20347466
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20347536
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20347536 #20355092
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20356544
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20356544 #20364240
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20365710
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20365710 #20424436
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20425948
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20425948 #20437870
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #20438878
:   <---- t1 OS9$8c <I$WritLn> {} #20438878 #20470910
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20470980
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20470980 #20478536
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20479964
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #20485456
:   <---- t1 OS9$8c <I$WritLn> {} #20485456 #20505544
{proc=5"Dir"} t1: OS9$06 <F$Exit> {status=0} #20505582
:   :   <---- t1 OS9$04 <F$Wait> {} #19641738 #20547166
{proc=4"Shell"} t1: OS9$8b <I$ReadLn> {} #20548542
{proc=4"Shell"} t1: OS9$1c <F$SUser> {} #20558528
:   <---- t1 OS9$1c <F$SUser> {} #20558528 #20560782
{proc=4"Shell"} t1: OS9$06 <F$Exit> {status=0} #20560826
:   :   <---- t1 OS9$04 <F$Wait> {} #14922642 #20591368
{proc=3"Shell"} t1: OS9$1c <F$SUser> {} #20592658
:   <---- t1 OS9$1c <F$SUser> {} #20592658 #20594912
{proc=3"Shell"} t1: OS9$06 <F$Exit> {status=0} #20594956
:   :   <---- t1 OS9$04 <F$Wait> {} #12769486 #20624490
{proc=2"SysGo"} t1: OS9$03 <F$Fork> {Module/file='AutoEx' param="{13}" lang/type=1 pages=0} #20624552
{proc=2"SysGo"} t1: OS9$05 <F$Chain> {Module/file='Shell' param="i=/1{13}" lang/type=1 pages=0} #21775490
{proc=2"Shell"} t1: OS9$09 <F$Icpt> {routine=e06b storage=0000} #25036204
:   <---- t1 OS9$09 <F$Icpt> {routine=e06b storage=0000} #25036204 #25038862
{proc=2"Shell"} t1: OS9$0c <F$ID> {} #25038914
:   <---- t1 OS9$0c <F$ID> {} #25038914 #25041186
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25043806
:   <---- t1 OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25043806 #25050258
{proc=2"Shell"} t1: OS9$84 <I$Open> {f7f4='.'} #25051000
:   <---- t1 OS9$84 <I$Open> {f7f4='.'} #25051000 -> path $3 #25410946
{proc=2"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #25411040
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #25411040 #25468664
{proc=2"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #25468784
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #25468784 #25477974
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25478242
:   <---- t1 OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25478242 #25484670
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=3} #25488074
:   <---- t1 OS9$8f <I$Close> {path=3} #25488074 #25500840
{proc=2"Shell"} t1: OS9$1c <F$SUser> {} #25505638
:   <---- t1 OS9$1c <F$SUser> {} #25505638 #25507892
{proc=2"Shell"} t1: OS9$82 <I$Dup> {$0} #25511790
:   <---- t1 OS9$82 <I$Dup> {$0} #25511790 -> path $3 #25515118
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=0} #25515168
:   <---- t1 OS9$8f <I$Close> {path=0} #25515168 #25519210
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25519594
:   <---- t1 OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25519594 #25526046
{proc=2"Shell"} t1: OS9$84 <I$Open> {00b5='/Term'} #25526078
:   <---- t1 OS9$84 <I$Open> {00b5='/Term'} #25526078 -> path $0 #25872932
{proc=2"Shell"} t1: OS9$82 <I$Dup> {$1} #25873106
:   <---- t1 OS9$82 <I$Dup> {$1} #25873106 -> path $4 #25876468
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=1} #25876518
:   <---- t1 OS9$8f <I$Close> {path=1} #25876518 #25880620
{proc=2"Shell"} t1: OS9$82 <I$Dup> {$0} #25880652
:   <---- t1 OS9$82 <I$Dup> {$0} #25880652 -> path $1 #25883912
{proc=2"Shell"} t1: OS9$82 <I$Dup> {$2} #25884058
:   <---- t1 OS9$82 <I$Dup> {$2} #25884058 -> path $5 #25887454
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=2} #25887504
:   <---- t1 OS9$8f <I$Close> {path=2} #25887504 #25893300
{proc=2"Shell"} t1: OS9$82 <I$Dup> {$1} #25893332
:   <---- t1 OS9$82 <I$Dup> {$1} #25893332 -> path $2 #25896626
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=0 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25896766
:   <---- t1 OS9$8d <I$GetStt> {path=0 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25896766 #25903218
{proc=2"Shell"} t1: OS9$8a <I$Write> {{0}} #25903300
:   <---- t1 OS9$8a <I$Write> {{0}} #25903300 #25913120
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25913654
:   <---- t1 OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25913654 #25920106
{proc=2"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #25923798
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #25923798 #25928186
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #25928228
:   <---- t1 OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #25928228 #25934286
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=3} #25934526
:   <---- t1 OS9$8f <I$Close> {path=3} #25934526 #25938748
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=4} #25938882
:   <---- t1 OS9$8f <I$Close> {path=4} #25938882 #25943164
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=5} #25943298
:   <---- t1 OS9$8f <I$Close> {path=5} #25943298 #25948822
{proc=2"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #25962084
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #25962084 #25966472
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #25966514
:   <---- t1 OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #25966514 #25972572
{proc=2"Shell"} t1: OS9$8a <I$Write> {
:   <---- t1 OS9$8a <I$Write> {
{proc=2"Shell"} t1: OS9$15 <F$Time> {buf=2da} #25991154
:   <---- t1 OS9$15 <F$Time> {buf=2da} #25991154 #25995044
{proc=2"Shell"} t1: OS9$8c <I$WritLn> {} #25998350
:   <---- t1 OS9$8c <I$WritLn> {} #25998350 #26020568
{proc=2"Shell"} t1: OS9$1c <F$SUser> {} #26020730
:   <---- t1 OS9$1c <F$SUser> {} #26020730 #26022984
{proc=2"Shell"} t1: OS9$8c <I$WritLn> {} #26023194
:   <---- t1 OS9$8c <I$WritLn> {} #26023194 #26043822
{proc=2"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #26043906
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #26043906 #26048294
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #26048344
:   <---- t1 OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #26048344 #26054402
{proc=2"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.Opt    : Read/Write PD Options} #26054562
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.Opt    : Read/Write PD Options} #26054562 #26062350
{proc=2"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.SSig   : Send signal on data ready} #26062394
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.SSig   : Send signal on data ready} #26062394 #26066850
{proc=2"Shell"} t1: OS9$0a <F$Sleep> {ticks=0000} #26066940
```

# Short: Both User and Kernel calls, with noisy, less important calls removed.

This is the file `_short`
See the `egrep -v` command above, to see which system calls were filtered out.

```
{proc=1} t0: OS9$35 <F$Boot> {} #2146194
:   {proc=1} t0: OS9$00 <F$Link> {type/lang=c1 module/file='Boot'} #2146728
:   :   <---- t0 OS9$00 <F$Link> {type/lang=c1 module/file='Boot'} #2146728 -> addr $ee30 entry $ee42 #2155076
:   <---- t0 OS9$35 <F$Boot> {} #2146194 #4261348
{proc=1} t0: OS9$00 <F$Link> {type/lang=c0 module/file='Init'} #4261394
:   <---- t0 OS9$00 <F$Link> {type/lang=c0 module/file='Init'} #4261394 -> addr $9e00 entry $ade0 #4303174
{proc=1} t0: OS9$00 <F$Link> {type/lang=c0 module/file='krnp2'} #4303400
:   <---- t0 OS9$00 <F$Link> {type/lang=c0 module/file='krnp2'} #4303400 -> addr $8700 entry $8713 #4348666
{proc=1} t0: OS9$32 <F$SSvc> {table=8789} #4348844
:   <---- t0 OS9$32 <F$SSvc> {table=8789} #4348844 #4354632
{proc=1} t0: OS9$86 <I$ChgDir> {mode=05, 9e34='/DD'} #4354832
:   {proc=1} t0: OS9$00 <F$Link> {type/lang=c1 module/file='IOMan'} #4355210
:   :   <---- t0 OS9$00 <F$Link> {type/lang=c1 module/file='IOMan'} #4355210 -> addr $93db entry $93ee #4396422
:   {proc=1} t0: OS9$32 <F$SSvc> {table=943d} #4456778
:   :   <---- t0 OS9$32 <F$SSvc> {table=943d} #4456778 #4458622
:   {proc=1} t0: OS9$80 <I$Attach> {9e35='DD'} #4464064
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ 9e35 dat@ 640} #4465280
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ 9e35 dat@ 640} #4465280 #4475798
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #4475930
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #4475930 -> addr $e704 entry $e72e #4492114
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #4492200
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #4492200 -> addr $9e77 entry $9e88 #4532798
:   :   <---- t0 OS9$80 <I$Attach> {9e35='DD'} #4464064 #4552604
:   {proc=1} t0: OS9$81 <I$Detach> {8300} #4669628
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #4670190
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #4670190 #4670930
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #4670950
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #4670950 #4671690
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #4671710
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #4671710 #4672450
:   :   <---- t0 OS9$81 <I$Detach> {8300} #4669628 #4672800
:   <---- t0 OS9$86 <I$ChgDir> {mode=05, 9e34='/DD'} #4354832 #4673856
{proc=1} t0: OS9$84 <I$Open> {9e37='/Term'} #4674062
:   {proc=1} t0: OS9$80 <I$Attach> {9e38='Term'} #4680206
:   :   {proc=1} t0: OS9$34 <F$SLink> {"Term" type f0 name@ 9e38 dat@ 640} #4681422
:   :   :   <---- t0 OS9$34 <F$SLink> {"Term" type f0 name@ 9e38 dat@ 640} #4681422 #4710720
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='VTIO'} #4710852
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='VTIO'} #4710852 -> addr $b8dd entry $be5c #4747624
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='SCF'} #4747710
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='SCF'} #4747710 -> addr $b165 entry $b410 #4786072
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=c1 module/file='JoyDrv'} #4805622
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=c1 module/file='JoyDrv'} #4805622 -> addr $c554 entry $c568 #4839968
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=c1 module/file='SndDrv'} #4840124
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=c1 module/file='SndDrv'} #4840124 -> addr $c48e entry $c4a2 #4876370
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=c1 module/file='CoWin'} #4876878
:   :   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$00 <F$Link> {type/lang=c1 module/file='CoWin'} #4876878 #4923758
:   :   {proc=1} t0: OS9$01 <F$Load> {type/lang=c1 filename='CoWin'} #4923798
:   :   :   {proc=1} t0: OS9$4b <F$AllPrc> {} #4924172
:   :   :   :   <---- t0 OS9$4b <F$AllPrc> {} #4924172 #4938400
:   :   :   {proc=1} t0: OS9$84 <I$Open> {c1ee='CoWin'} #4938574
:   :   :   :   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #4944378
:   :   :   :   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #4945594
:   :   :   :   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #4945594 #4956832
:   :   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #4956964
:   :   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #4956964 -> addr $e704 entry $e72e #4973148
:   :   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #4973234
:   :   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #4973234 -> addr $9e77 entry $9e88 #5013832
:   :   :   :   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #4944378 #5016022
:   :   :   :   {proc=1} t0: OS9$81 <I$Detach> {8300} #5244628
:   :   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #5245190
:   :   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #5245190 #5245930
:   :   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #5245950
:   :   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #5245950 #5246690
:   :   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #5246710
:   :   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #5246710 #5247450
:   :   :   :   :   <---- t0 OS9$81 <I$Detach> {8300} #5244628 #5247800
:   :   :   :   <-- ERROR: $d8(E$PNNF   :Path Name Not Found): OS9KERNEL0 OS9$84 <I$Open> {c1ee='CoWin'} #4938574 #5248866
:   :   :   {proc=1} t0: OS9$4c <F$DelPrc> {} #5249028
:   :   :   :   {proc=1} t0: OS9$40 <F$DelTsk> {proc_desc=7e00} #5249480
:   :   :   :   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7e00} #5249480 #5250074
:   :   :   :   <---- t0 OS9$4c <F$DelPrc> {} #5249028 #5253188
:   :   :   <-- ERROR: $d8(E$PNNF   :Path Name Not Found): OS9KERNEL0 OS9$01 <F$Load> {type/lang=c1 filename='CoWin'} #4923798 #5253540
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=c1 module/file='CoGrf'} #5253860
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=c1 module/file='CoGrf'} #5253860 -> addr $c6d9 entry $c9ec #5287654
:   :   {proc=1} t0: OS9$21 <F$NMLink> {LangType=c1, c7cc='grfdrv'} #5288352
:   :   :   {proc=1} t0: OS9$4e <F$FModul> {"grfdrv" type c1 name@ c7cc dat@ 640} #5288742
:   :   :   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$4e <F$FModul> {"grfdrv" type c1 name@ c7cc dat@ 640} #5288742 #5334166
:   :   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$21 <F$NMLink> {LangType=c1, c7cc='grfdrv'} #5288352 #5334430
:   :   {proc=1} t0: OS9$22 <F$NMLoad> {LangType=c1, c7c4='../CMDS/grfdrv'} #5334692
:   :   :   {proc=1} t0: OS9$4b <F$AllPrc> {} #5335070
:   :   :   :   <---- t0 OS9$4b <F$AllPrc> {} #5335070 #5348962
:   :   :   {proc=1} t0: OS9$84 <I$Open> {c7c4='../CMDS/grfdrv'} #5349136
:   :   :   :   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #5354940
:   :   :   :   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #5356156
:   :   :   :   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #5356156 #5367394
:   :   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #5367526
:   :   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #5367526 -> addr $e704 entry $e72e #5383710
:   :   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #5383796
:   :   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #5383796 -> addr $9e77 entry $9e88 #5424394
:   :   :   :   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #5354940 #5426584
:   :   :   :   <---- t0 OS9$84 <I$Open> {c7c4='../CMDS/grfdrv'} #5349136 -> path $2 #6074750
:   :   :   {proc=1} t0: OS9$3f <F$AllTsk> {processDesc=7e00} #6074808
:   :   :   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7e00} #6074808 #6075672
:   :   :   {proc=2} t0: OS9$41 <F$SetTsk> {} #6076468
:   :   :   :   <---- t0 OS9$41 <F$SetTsk> {} #6076468 #6077182
:   :   :   {proc=2} t0: OS9$89 <I$Read> {path=2 buf=0000 size=9} #6077234
:   :   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=0000 size=9} #6077234 #6132272
:   :   :   {proc=2} t0: OS9$41 <F$SetTsk> {} #6134790
:   :   :   :   <---- t0 OS9$41 <F$SetTsk> {} #6134790 #6135504
:   :   :   {proc=2} t0: OS9$89 <I$Read> {path=2 buf=0009 size=228d} #6135556
:   :   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=0009 size=228d} #6135556 #7946042
:   :   :   {proc=2} t0: OS9$89 <I$Read> {path=2 buf=2296 size=9} #7990368
:   :   :   :   <-- ERROR: $d3(E$EOF    :End of File): OS9KERNEL0 OS9$89 <I$Read> {path=2 buf=2296 size=9} #7990368 #7993924
:   :   :   {proc=1} t0: OS9$8f <I$Close> {path=2} #7994024
:   :   :   :   {proc=1} t0: OS9$81 <I$Detach> {8300} #8000882
:   :   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #8001444
:   :   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #8001444 #8002184
:   :   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #8002204
:   :   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #8002204 #8002944
:   :   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #8002964
:   :   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #8002964 #8003704
:   :   :   :   :   <---- t0 OS9$81 <I$Detach> {8300} #8000882 #8004054
:   :   :   :   <---- t0 OS9$8f <I$Close> {path=2} #7994024 #8005110
:   :   :   {proc=1} t0: OS9$4c <F$DelPrc> {} #8005242
:   :   :   :   {proc=1} t0: OS9$40 <F$DelTsk> {proc_desc=7e00} #8005694
:   :   :   :   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7e00} #8005694 #8006364
:   :   :   :   <---- t0 OS9$4c <F$DelPrc> {} #8005242 #8009478
:   :   :   <---- t0 OS9$22 <F$NMLoad> {LangType=c1, c7c4='../CMDS/grfdrv'} #5334692 #8013550
:   :   {proc=1} t0: OS9$4e <F$FModul> {"grfdrv" type c1 name@ c7cc dat@ 640} #8015324
:   :   :   <---- t0 OS9$4e <F$FModul> {"grfdrv" type c1 name@ c7cc dat@ 640} #8015324 #8022466
:   :   {proc=1} t0: OS9$13 <F$AllBit> {bitmap=1071 first=0 count=1} #8058394
:   :   :   {proc=1} t0: OS9$4a <F$STABX> {} #8059772
:   :   :   :   <---- t0 OS9$4a <F$STABX> {} #8059772 #8060540
:   :   :   <---- t0 OS9$13 <F$AllBit> {bitmap=1071 first=0 count=1} #8058394 #8060794
:   :   <---- t0 OS9$80 <I$Attach> {9e38='Term'} #4680206 #8062034
:   {proc=1} t0: OS9$80 <I$Attach> {d511='Term'} #8081736
:   :   {proc=1} t0: OS9$34 <F$SLink> {"Term" type f0 name@ d511 dat@ 640} #8082952
:   :   :   <---- t0 OS9$34 <F$SLink> {"Term" type f0 name@ d511 dat@ 640} #8082952 #8115056
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='VTIO'} #8115188
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='VTIO'} #8115188 -> addr $b8dd entry $be5c #8153326
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='SCF'} #8153412
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='SCF'} #8153412 -> addr $b165 entry $b410 #8193140
:   :   <---- t0 OS9$80 <I$Attach> {d511='Term'} #8081736 #8195386
:   <---- t0 OS9$84 <I$Open> {9e37='/Term'} #4674062 -> path $1 #8199346
{proc=1} t0: OS9$82 <I$Dup> {$1} #8199384
:   <---- t0 OS9$82 <I$Dup> {$1} #8199384 -> path $1 #8200778
{proc=1} t0: OS9$82 <I$Dup> {$1} #8200800
:   <---- t0 OS9$82 <I$Dup> {$1} #8200800 -> path $1 #8202194
{proc=1} t0: OS9$00 <F$Link> {type/lang=c0 module/file='krnp3'} #8202230
:   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$00 <F$Link> {type/lang=c0 module/file='krnp3'} #8202230 #8248804
{proc=1} t0: OS9$03 <F$Fork> {Module/file='SysGo' param="" lang/type=1 pages=0} #8249016
:   {proc=1} t0: OS9$82 <I$Dup> {$1} #8264194
:   :   <---- t0 OS9$82 <I$Dup> {$1} #8264194 -> path $1 #8265588
:   {proc=1} t0: OS9$82 <I$Dup> {$1} #8265652
:   :   <---- t0 OS9$82 <I$Dup> {$1} #8265652 -> path $1 #8267046
:   {proc=1} t0: OS9$82 <I$Dup> {$1} #8267110
:   :   <---- t0 OS9$82 <I$Dup> {$1} #8267110 -> path $1 #8268504
:   {proc=2} t0: OS9$34 <F$SLink> {"SysGo" type 1 name@ 9e2f dat@ 640} #8268694
:   :   <---- t0 OS9$34 <F$SLink> {"SysGo" type 1 name@ 9e2f dat@ 640} #8268694 #8290974
:   {proc=2} t0: OS9$07 <F$Mem> {desired_size=fc} #8292254
:   :   {proc=2} t0: OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7a00} #8292772
:   :   :   <---- t0 OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7a00} #8292772 #8294362
:   :   <---- t0 OS9$07 <F$Mem> {desired_size=fc} #8292254 #8294668
:   {proc=1} t0: OS9$3f <F$AllTsk> {processDesc=7a00} #8294992
:   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7a00} #8294992 #8295856
:   {proc=1} t0: OS9$40 <F$DelTsk> {proc_desc=7a00} #8298440
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7a00} #8298440 #8299110
:   {proc=1} t0: OS9$2c <F$AProc> {proc=7a00} #8299216
:   :   <---- t0 OS9$2c <F$AProc> {proc=7a00} #8299216 #8299940
:   <---- t0 OS9$03 <F$Fork> {Module/file='SysGo' param="" lang/type=1 pages=0} #8249016 #8300180
{proc=1} t0: OS9$2d <F$NProc> {} #8300198
{proc=2"SysGo"} t1: OS9$09 <F$Icpt> {routine=da8b storage=0000} #8301674
:   <---- t1 OS9$09 <F$Icpt> {routine=da8b storage=0000} #8301674 #8303952
{proc=2"SysGo"} t1: OS9$0c <F$ID> {} #8303964
:   <---- t1 OS9$0c <F$ID> {} #8303964 #8306236
{proc=2"SysGo"} t1: OS9$0d <F$SPrior> {pid=02 priority=80} #8306252
:   <---- t1 OS9$0d <F$SPrior> {pid=02 priority=80} #8306252 #8308656
{proc=2"SysGo"} t1: OS9$00 <F$Link> {type/lang=00 module/file='Init'} #8308696
:   <---- t1 OS9$00 <F$Link> {type/lang=00 module/file='Init'} #8308696 -> addr $9e00 entry $ade0 #8352448
{proc=2"SysGo"} t1: OS9$8a <I$Write> {NitrOS-9/6809 Level 2 V3.3.0} #8353444
<---- t1 OS9$8a <I$Write> {NitrOS-9/6809 Level 2 V3.3.0} #8353444 #8411218
{proc=2"SysGo"} t1: OS9$8c <I$WritLn> {} #8411284
:   <---- t1 OS9$8c <I$WritLn> {} #8411284 #8423858
{proc=2"SysGo"} t1: OS9$8a <I$Write> {Tandy Color Computer 3} #8424676
:   <---- t1 OS9$8a <I$Write> {Tandy Color Computer 3} #8424676 #8440428
{proc=2"SysGo"} t1: OS9$8c <I$WritLn> {} #8440494
:   <---- t1 OS9$8c <I$WritLn> {} #8440494 #8453068
{proc=2"SysGo"} t1: OS9$8a <I$Write> {(C) 2014 The NitrOS-9 Project
:   <---- t1 OS9$8a <I$Write> {(C) 2014 The NitrOS-9 Project
{proc=2"SysGo"} t1: OS9$16 <F$STime> {y85 m12 d31 h23 m59 s59} #8561324
:   {proc=1} t0: OS9$00 <F$Link> {type/lang=c1 module/file='Clock'} #8564204
:   :   <---- t0 OS9$00 <F$Link> {type/lang=c1 module/file='Clock'} #8564204 -> addr $d626 entry $d7c2 #8589886
:   {proc=1} t0: OS9$00 <F$Link> {type/lang=21 module/file='Clock2'} #8589992
:   :   <---- t0 OS9$00 <F$Link> {type/lang=21 module/file='Clock2'} #8589992 -> addr $d82d entry $d841 #8614326
:   {proc=2} t0: OS9$32 <F$SSvc> {table=d639} #8614610
:   :   <---- t0 OS9$32 <F$SSvc> {table=d639} #8614610 #8615766
:   <---- t1 OS9$16 <F$STime> {y85 m12 d31 h23 m59 s59} #8561324 #8618608
{proc=2"SysGo"} t1: OS9$86 <I$ChgDir> {mode=04, d93a='CMDS'} #8618642
:   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #8625182
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #8626398
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #8626398 #8639042
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #8639174
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #8639174 -> addr $e704 entry $e72e #8661100
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #8661186
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #8661186 -> addr $9e77 entry $9e88 #8704824
:   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #8625182 #8708704
:   {proc=2} t0: OS9$81 <I$Detach> {8300} #8940124
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #8940686
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #8940686 #8941426
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #8941446
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #8941446 #8942186
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #8942206
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #8942206 #8942946
:   :   <---- t0 OS9$81 <I$Detach> {8300} #8940124 #8943296
:   <---- t1 OS9$86 <I$ChgDir> {mode=04, d93a='CMDS'} #8618642 #8947230
{proc=2"SysGo"} t1: OS9$86 <I$ChgDir> {mode=01, d932='/DD'} #8947264
:   {proc=2} t0: OS9$80 <I$Attach> {d933='DD'} #8954100
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ d933 dat@ 7a40} #8955316
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ d933 dat@ 7a40} #8955316 #8967970
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #8968102
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #8968102 -> addr $e704 entry $e72e #8987326
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #8987412
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #8987412 -> addr $9e77 entry $9e88 #9031050
:   :   <---- t0 OS9$80 <I$Attach> {d933='DD'} #8954100 #9033296
:   {proc=2} t0: OS9$81 <I$Detach> {8300} #9158084
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #9158646
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #9158646 #9159386
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #9159406
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #9159406 #9160146
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #9160166
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #9160166 #9160906
:   :   <---- t0 OS9$81 <I$Detach> {8300} #9158084 #9161256
:   <---- t1 OS9$86 <I$ChgDir> {mode=01, d932='/DD'} #8947264 #9163556
{proc=2"SysGo"} t1: OS9$86 <I$ChgDir> {mode=04, d936='/DD/CMDS'} #9163596
:   {proc=2} t0: OS9$80 <I$Attach> {d937='DD/CMDS'} #9170432
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD/CMDS" type f0 name@ d937 dat@ 7a40} #9171648
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD/CMDS" type f0 name@ d937 dat@ 7a40} #9171648 #9185936
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #9186068
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #9186068 -> addr $e704 entry $e72e #9203658
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #9203744
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #9203744 -> addr $9e77 entry $9e88 #9249016
:   :   <---- t0 OS9$80 <I$Attach> {d937='DD/CMDS'} #9170432 #9251262
:   {proc=2} t0: OS9$81 <I$Detach> {8300} #9483022
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #9483584
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #9483584 #9484324
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #9484344
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #9484344 #9485084
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #9485104
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #9485104 #9485844
:   :   <---- t0 OS9$81 <I$Detach> {8300} #9483022 #9486194
:   <---- t1 OS9$86 <I$ChgDir> {mode=04, d936='/DD/CMDS'} #9163596 #9488494
{proc=2"SysGo"} t1: OS9$0c <F$ID> {} #9488524
:   <---- t1 OS9$0c <F$ID> {} #9488524 #9490796
{proc=2"SysGo"} t1: OS9$18 <F$GPrDsc> {} #9490828
:   {proc=2} t0: OS9$37 <F$GProcP> {id=02} #9491918
:   :   <---- t0 OS9$37 <F$GProcP> {id=02} #9491918 #9492626
:   <---- t1 OS9$18 <F$GPrDsc> {} #9490828 #9504954
{proc=2"SysGo"} t1: OS9$4f <F$MapBlk> {beginningBlock=0 numBlocks=1} #9504996
:   {proc=2} t0: OS9$3e <F$FreeHB> {} #9506184
:   :   <---- t0 OS9$3e <F$FreeHB> {} #9506184 #9507302
:   {proc=2} t0: OS9$3c <F$SetImg> {} #9507392
:   :   <---- t0 OS9$3c <F$SetImg> {} #9507392 #9508160
:   <---- t1 OS9$4f <F$MapBlk> {beginningBlock=0 numBlocks=1} #9504996 #9509802
{proc=2"SysGo"} t1: OS9$8d <I$GetStt> {path=1 27==SS.KySns  : Getstat/SetStat for COCO keyboard} #9510374
:   <---- t1 OS9$8d <I$GetStt> {path=1 27==SS.KySns  : Getstat/SetStat for COCO keyboard} #9510374 #9515938
{proc=2"SysGo"} t1: OS9$03 <F$Fork> {Module/file='Shell' param="startup{32}-p{13}" lang/type=1 pages=0} #9516016
:   {proc=2} t0: OS9$82 <I$Dup> {$1} #9531804
:   :   <---- t0 OS9$82 <I$Dup> {$1} #9531804 -> path $1 #9533198
:   {proc=2} t0: OS9$82 <I$Dup> {$1} #9533262
:   :   <---- t0 OS9$82 <I$Dup> {$1} #9533262 -> path $1 #9534656
:   {proc=2} t0: OS9$82 <I$Dup> {$1} #9534720
:   :   <---- t0 OS9$82 <I$Dup> {$1} #9534720 -> path $1 #9536114
:   {proc=3} t0: OS9$34 <F$SLink> {"Shell" type 1 name@ d93f dat@ 7a40} #9536304
:   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$34 <F$SLink> {"Shell" type 1 name@ d93f dat@ 7a40} #9536304 #9587664
:   {proc=2} t0: OS9$01 <F$Load> {type/lang=7a filename='Shell'} #9587714
:   :   {proc=2} t0: OS9$4b <F$AllPrc> {} #9588088
:   :   :   <---- t0 OS9$4b <F$AllPrc> {} #9588088 #9604800
:   :   {proc=2} t0: OS9$84 <I$Open> {d93f='Shell'} #9604974
:   :   :   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #9610778
:   :   :   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #9611994
:   :   :   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #9611994 #9624638
:   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #9624770
:   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #9624770 -> addr $e704 entry $e72e #9643994
:   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #9644080
:   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #9644080 -> addr $9e77 entry $9e88 #9687718
:   :   :   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #9610778 #9691598
:   :   :   <---- t0 OS9$84 <I$Open> {d93f='Shell'} #9604974 -> path $2 #10455994
:   :   {proc=2} t0: OS9$3f <F$AllTsk> {processDesc=7600} #10456052
:   :   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7600} #10456052 #10456946
:   :   {proc=4} t0: OS9$41 <F$SetTsk> {} #10457938
:   :   :   <---- t0 OS9$41 <F$SetTsk> {} #10457938 #10458652
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=0000 size=9} #10458704
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=0000 size=9} #10458704 #10512232
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=0009 size=1b4e} #10514216
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=0009 size=1b4e} #10514216 #12020178
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1b57 size=9} #12068168
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1b57 size=9} #12068168 #12073686
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1b60 size=e8} #12075670
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1b60 size=e8} #12075670 #12139804
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1c48 size=9} #12188092
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1c48 size=9} #12188092 #12193610
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1c51 size=4a} #12195594
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1c51 size=4a} #12195594 #12203692
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1c9b size=9} #12254122
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1c9b size=9} #12254122 #12261274
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1ca4 size=19} #12263258
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1ca4 size=19} #12263258 #12269004
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1cbd size=9} #12320958
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1cbd size=9} #12320958 #12326476
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1cc6 size=3e} #12328460
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1cc6 size=3e} #12328460 #12388096
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d04 size=9} #12442924
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d04 size=9} #12442924 #12448442
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d0d size=23} #12450426
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d0d size=23} #12450426 #12456354
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d30 size=9} #12511986
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d30 size=9} #12511986 #12517504
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d39 size=1b} #12519488
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d39 size=1b} #12519488 #12526936
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d54 size=9} #12582566
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d54 size=9} #12582566 #12589718
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d5d size=5e} #12591702
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d5d size=5e} #12591702 #12598530
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1dbb size=9} #12657012
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1dbb size=9} #12657012 #12662530
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1dc4 size=1e} #12664514
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1dc4 size=1e} #12664514 #12670430
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1de2 size=9} #12728954
:   :   :   <-- ERROR: $d3(E$EOF    :End of File): OS9KERNEL0 OS9$89 <I$Read> {path=2 buf=1de2 size=9} #12728954 #12732510
:   :   {proc=2} t0: OS9$8f <I$Close> {path=2} #12732610
:   :   :   {proc=2} t0: OS9$81 <I$Detach> {8300} #12741938
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #12742500
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #12742500 #12743240
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #12743260
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #12743260 #12744000
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #12744020
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #12744020 #12744760
:   :   :   :   <---- t0 OS9$81 <I$Detach> {8300} #12741938 #12745110
:   :   :   <---- t0 OS9$8f <I$Close> {path=2} #12732610 #12746166
:   :   {proc=2} t0: OS9$4c <F$DelPrc> {} #12746298
:   :   :   {proc=2} t0: OS9$40 <F$DelTsk> {proc_desc=7600} #12746750
:   :   :   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7600} #12746750 #12747420
:   :   :   <---- t0 OS9$4c <F$DelPrc> {} #12746298 #12750272
:   :   {proc=3} t0: OS9$4d <F$ELink> {} #12753320
:   :   :   <---- t0 OS9$4d <F$ELink> {} #12753320 #12756098
:   :   <---- t0 OS9$01 <F$Load> {type/lang=7a filename='Shell'} #9587714 #12756430
:   {proc=3} t0: OS9$07 <F$Mem> {desired_size=1f00} #12757750
:   :   {proc=3} t0: OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7800} #12758268
:   :   :   <---- t0 OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7800} #12758268 #12759594
:   :   <---- t0 OS9$07 <F$Mem> {desired_size=1f00} #12757750 #12759900
:   {proc=2} t0: OS9$3f <F$AllTsk> {processDesc=7800} #12760224
:   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7800} #12760224 #12761118
:   {proc=2} t0: OS9$40 <F$DelTsk> {proc_desc=7800} #12766484
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7800} #12766484 #12767154
:   {proc=2} t0: OS9$2c <F$AProc> {proc=7800} #12767260
:   :   <---- t0 OS9$2c <F$AProc> {proc=7800} #12767260 #12767984
:   <---- t1 OS9$03 <F$Fork> {Module/file='Shell' param="startup{32}-p{13}" lang/type=1 pages=0} #9516016 #12769468
{proc=2"SysGo"} t1: OS9$04 <F$Wait> {} #12769486
:   {proc=2} t0: OS9$40 <F$DelTsk> {proc_desc=7a00} #12770920
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7a00} #12770920 #12771590
:   {proc=2} t0: OS9$2d <F$NProc> {} #12771644
{proc=3"Shell"} t1: OS9$09 <F$Icpt> {routine=e06b storage=0000} #12780034
:   <---- t1 OS9$09 <F$Icpt> {routine=e06b storage=0000} #12780034 #12782312
{proc=3"Shell"} t1: OS9$0c <F$ID> {} #12782364
:   <---- t1 OS9$0c <F$ID> {} #12782364 #12784636
{proc=3"Shell"} t1: OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #12787256
:   <---- t1 OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #12787256 #12795290
{proc=3"Shell"} t1: OS9$84 <I$Open> {f7f4='.'} #12796032
:   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #12802778
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #12803994
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #12803994 #12831734
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #12831866
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #12831866 -> addr $e704 entry $e72e #12863744
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #12863830
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #12863830 -> addr $9e77 entry $9e88 #12921798
:   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #12802778 #12924044
:   <---- t1 OS9$84 <I$Open> {f7f4='.'} #12796032 -> path $3 #13143592
{proc=3"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #13143686
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #13143686 #13197682
{proc=3"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #13197802
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #13197802 #13205358
{proc=3"Shell"} t1: OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #13205626
:   <---- t1 OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #13205626 #13213688
{proc=3"Shell"} t1: OS9$8f <I$Close> {path=3} #13217092
:   {proc=3} t0: OS9$81 <I$Detach> {8300} #13224646
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #13225208
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #13225208 #13225948
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #13225968
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #13225968 #13226708
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #13226728
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #13226728 #13227468
:   :   <---- t0 OS9$81 <I$Detach> {8300} #13224646 #13227818
:   <---- t1 OS9$8f <I$Close> {path=3} #13217092 #13229802
{proc=3"Shell"} t1: OS9$1c <F$SUser> {} #13232352
:   <---- t1 OS9$1c <F$SUser> {} #13232352 #13234606
{proc=3"Shell"} t1: OS9$21 <F$NMLink> {LangType=00, 0e6d='startup'} #13257724
:   {proc=3} t0: OS9$4e <F$FModul> {"startup" type 0 name@ e6d dat@ 7840} #13258840
:   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$4e <F$FModul> {"startup" type 0 name@ e6d dat@ 7840} #13258840 #13314140
:   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL1 OS9$21 <F$NMLink> {LangType=00, 0e6d='startup'} #13257724 #13315332
{proc=3"Shell"} t1: OS9$84 <I$Open> {0e6d='startup'} #13315360
:   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #13322024
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #13323240
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #13323240 #13350980
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #13351112
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #13351112 -> addr $e704 entry $e72e #13382990
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #13383076
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #13383076 -> addr $9e77 entry $9e88 #13441002
:   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #13322024 #13443248
:   {proc=3} t0: OS9$81 <I$Detach> {8300} #14385884
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #14386446
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #14386446 #14387186
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #14387206
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #14387206 #14387946
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #14387966
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #14387966 #14388706
:   :   <---- t0 OS9$81 <I$Detach> {8300} #14385884 #14389056
:   <-- ERROR: $d8(E$PNNF   :Path Name Not Found): OS9KERNEL1 OS9$84 <I$Open> {0e6d='startup'} #13315360 #14391782
{proc=3"Shell"} t1: OS9$82 <I$Dup> {$0} #14392028
:   <---- t1 OS9$82 <I$Dup> {$0} #14392028 -> path $3 #14395356
{proc=3"Shell"} t1: OS9$8f <I$Close> {path=0} #14395406
:   <---- t1 OS9$8f <I$Close> {path=0} #14395406 #14399308
{proc=3"Shell"} t1: OS9$84 <I$Open> {0e6d='startup'} #14399528
:   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #14407666
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #14408882
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #14408882 #14436622
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #14436754
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #14436754 -> addr $e704 entry $e72e #14468632
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #14468718
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #14468718 -> addr $9e77 entry $9e88 #14526644
:   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #14407666 #14528890
:   <---- t1 OS9$84 <I$Open> {0e6d='startup'} #14399528 -> path $0 #14807940
{proc=3"Shell"} t1: OS9$21 <F$NMLink> {LangType=11, e00d='Shell'} #14808826
:   {proc=3} t0: OS9$4e <F$FModul> {"Shell" type 11 name@ e00d dat@ 7840} #14809942
:   :   <---- t0 OS9$4e <F$FModul> {"Shell" type 11 name@ e00d dat@ 7840} #14809942 #14829876
:   <---- t1 OS9$21 <F$NMLink> {LangType=11, e00d='Shell'} #14808826 #14832058
{proc=3"Shell"} t1: OS9$03 <F$Fork> {Module/file='Shell' param="-P{32}X{32}PATH=;-p{13}" lang/type=11 pages=1f} #14832254
:   {proc=3} t0: OS9$82 <I$Dup> {$2} #14848732
:   :   <---- t0 OS9$82 <I$Dup> {$2} #14848732 -> path $2 #14850126
:   {proc=3} t0: OS9$82 <I$Dup> {$1} #14850190
:   :   <---- t0 OS9$82 <I$Dup> {$1} #14850190 -> path $1 #14853218
:   {proc=3} t0: OS9$82 <I$Dup> {$1} #14853282
:   :   <---- t0 OS9$82 <I$Dup> {$1} #14853282 -> path $1 #14854676
:   {proc=4} t0: OS9$34 <F$SLink> {"Shell" type 11 name@ e00d dat@ 7840} #14854866
:   :   <---- t0 OS9$34 <F$SLink> {"Shell" type 11 name@ e00d dat@ 7840} #14854866 #14875356
:   {proc=4} t0: OS9$07 <F$Mem> {desired_size=1f00} #14876676
:   :   {proc=4} t0: OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7500} #14877194
:   :   :   <---- t0 OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7500} #14877194 #14878982
:   :   <---- t0 OS9$07 <F$Mem> {desired_size=1f00} #14876676 #14879288
:   {proc=3} t0: OS9$3f <F$AllTsk> {processDesc=7500} #14879612
:   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7500} #14879612 #14880506
:   {proc=3} t0: OS9$40 <F$DelTsk> {proc_desc=7500} #14885886
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7500} #14885886 #14886556
:   {proc=3} t0: OS9$2c <F$AProc> {proc=7500} #14886662
:   :   <---- t0 OS9$2c <F$AProc> {proc=7500} #14886662 #14887386
:   <---- t1 OS9$03 <F$Fork> {Module/file='Shell' param="-P{32}X{32}PATH=;-p{13}" lang/type=11 pages=1f} #14832254 #14888870
{proc=3"Shell"} t1: OS9$1d <F$UnLoad> {} #14888944
:   {proc=3} t0: OS9$4e <F$FModul> {"Shell" type 11 name@ e00d dat@ 7840} #14890060
:   :   <---- t0 OS9$4e <F$FModul> {"Shell" type 11 name@ e00d dat@ 7840} #14890060 #14908360
:   <---- t1 OS9$1d <F$UnLoad> {} #14888944 #14909610
{proc=3"Shell"} t1: OS9$8f <I$Close> {path=0} #14909852
:   <---- t1 OS9$8f <I$Close> {path=0} #14909852 #14915110
{proc=3"Shell"} t1: OS9$82 <I$Dup> {$3} #14915132
:   <---- t1 OS9$82 <I$Dup> {$3} #14915132 -> path $0 #14918358
{proc=3"Shell"} t1: OS9$8f <I$Close> {path=3} #14918408
:   <---- t1 OS9$8f <I$Close> {path=3} #14918408 #14922310
{proc=3"Shell"} t1: OS9$04 <F$Wait> {} #14922642
:   {proc=3} t0: OS9$40 <F$DelTsk> {proc_desc=7800} #14924076
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7800} #14924076 #14924746
:   {proc=3} t0: OS9$2d <F$NProc> {} #14924800
{proc=4"Shell"} t1: OS9$09 <F$Icpt> {routine=e06b storage=0000} #14933190
:   <---- t1 OS9$09 <F$Icpt> {routine=e06b storage=0000} #14933190 #14935468
{proc=4"Shell"} t1: OS9$0c <F$ID> {} #14935520
:   <---- t1 OS9$0c <F$ID> {} #14935520 #14937792
{proc=4"Shell"} t1: OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #14940412
:   <---- t1 OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #14940412 #14948446
{proc=4"Shell"} t1: OS9$84 <I$Open> {f7f4='.'} #14949188
:   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #14980136
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #14981352
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #14981352 #15009092
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #15009224
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #15009224 -> addr $e704 entry $e72e #15041144
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #15041230
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #15041230 -> addr $9e77 entry $9e88 #15099156
:   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #14980136 #15101402
:   <---- t1 OS9$84 <I$Open> {f7f4='.'} #14949188 -> path $3 #15334764
{proc=4"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #15334858
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #15334858 #15392058
{proc=4"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #15392178
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #15392178 #15399734
{proc=4"Shell"} t1: OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #15400002
:   <---- t1 OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #15400002 #15406430
{proc=4"Shell"} t1: OS9$8f <I$Close> {path=3} #15409834
:   {proc=4} t0: OS9$81 <I$Detach> {8300} #15419970
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #15420532
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #15420532 #15421272
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #15421292
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #15421292 #15422032
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #15422052
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #15422052 #15422792
:   :   <---- t0 OS9$81 <I$Detach> {8300} #15419970 #15423142
:   <---- t1 OS9$8f <I$Close> {path=3} #15409834 #15427734
{proc=4"Shell"} t1: OS9$1c <F$SUser> {} #15430386
:   <---- t1 OS9$1c <F$SUser> {} #15430386 #15432640
{proc=4"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #15461340
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #15461340 #15465484
{proc=4"Shell"} t1: OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #15465526
:   <---- t1 OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #15465526 #15471212
{proc=4"Shell"} t1: OS9$1c <F$SUser> {} #15471490
:   <---- t1 OS9$1c <F$SUser> {} #15471490 #15473744
{proc=4"Shell"} t1: OS9$8b <I$ReadLn> {} #15473886
:   <---- t1 OS9$8b <I$ReadLn> {} #15473886 #15530026
{proc=4"Shell"} t1: OS9$8b <I$ReadLn> {} #15530422
:   <---- t1 OS9$8b <I$ReadLn> {} #15530422 #15540560
{proc=4"Shell"} t1: OS9$1c <F$SUser> {} #15544186
:   <---- t1 OS9$1c <F$SUser> {} #15544186 #15546440
{proc=4"Shell"} t1: OS9$21 <F$NMLink> {LangType=00, 0e6d='dir'} #15566544
:   {proc=4} t0: OS9$4e <F$FModul> {"dir" type 0 name@ e6d dat@ 7540} #15569294
:   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$4e <F$FModul> {"dir" type 0 name@ e6d dat@ 7540} #15569294 #15621214
:   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL1 OS9$21 <F$NMLink> {LangType=00, 0e6d='dir'} #15566544 #15622406
{proc=4"Shell"} t1: OS9$84 <I$Open> {0e6d='dir'} #15622434
:   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #15653782
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #15654998
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #15654998 #15682738
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #15682870
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #15682870 -> addr $e704 entry $e72e #15714748
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #15714834
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #15714834 -> addr $9e77 entry $9e88 #15772760
:   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #15653782 #15775006
:   <---- t1 OS9$84 <I$Open> {0e6d='dir'} #15622434 -> path $3 #16110232
{proc=4"Shell"} t1: OS9$89 <I$Read> {path=3 buf=00d6 size=4d} #16110274
:   <---- t1 OS9$89 <I$Read> {path=3 buf=00d6 size=4d} #16110274 #16166644
{proc=4"Shell"} t1: OS9$8f <I$Close> {path=3} #16170222
:   {proc=4} t0: OS9$81 <I$Detach> {8300} #16178776
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #16179338
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #16179338 #16180078
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #16180098
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #16180098 #16180838
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #16180858
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #16180858 #16181598
:   :   <---- t0 OS9$81 <I$Detach> {8300} #16178776 #16181948
:   <---- t1 OS9$8f <I$Close> {path=3} #16170222 #16186920
{proc=4"Shell"} t1: OS9$21 <F$NMLink> {LangType=11, 0e6d='dir'} #16187084
:   {proc=4} t0: OS9$4e <F$FModul> {"dir" type 11 name@ e6d dat@ 7540} #16188200
:   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$4e <F$FModul> {"dir" type 11 name@ e6d dat@ 7540} #16188200 #16241702
:   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL1 OS9$21 <F$NMLink> {LangType=11, 0e6d='dir'} #16187084 #16242894
{proc=4"Shell"} t1: OS9$22 <F$NMLoad> {LangType=11, 0e6d='dir'} #16242912
:   {proc=4} t0: OS9$4b <F$AllPrc> {} #16244016
:   :   <---- t0 OS9$4b <F$AllPrc> {} #16244016 #16260878
:   {proc=4} t0: OS9$84 <I$Open> {0e6d=''} #16261052
:   :   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #16291858
:   :   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #16293074
:   :   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #16293074 #16320814
:   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #16320946
:   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #16320946 -> addr $e704 entry $e72e #16352824
:   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #16352910
:   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #16352910 -> addr $9e77 entry $9e88 #16410836
:   :   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #16291858 #16413082
:   :   <---- t0 OS9$84 <I$Open> {0e6d=''} #16261052 -> path $4 #16743306
:   {proc=4} t0: OS9$3f <F$AllTsk> {processDesc=7300} #16743364
:   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7300} #16743364 #16744258
:   {proc=5} t0: OS9$41 <F$SetTsk> {} #16745400
:   :   <---- t0 OS9$41 <F$SetTsk> {} #16745400 #16746114
:   {proc=5} t0: OS9$89 <I$Read> {path=4 buf=0000 size=9} #16746166
:   :   <---- t0 OS9$89 <I$Read> {path=4 buf=0000 size=9} #16746166 #16796554
:   {proc=5} t0: OS9$89 <I$Read> {path=4 buf=0009 size=398} #16798538
:   :   <---- t0 OS9$89 <I$Read> {path=4 buf=0009 size=398} #16798538 #16965102
:   {proc=5} t0: OS9$89 <I$Read> {path=4 buf=03a1 size=9} #17025018
:   :   <-- ERROR: $d3(E$EOF    :End of File): OS9KERNEL0 OS9$89 <I$Read> {path=4 buf=03a1 size=9} #17025018 #17028574
:   {proc=4} t0: OS9$8f <I$Close> {path=4} #17028674
:   :   {proc=4} t0: OS9$81 <I$Detach> {8300} #17037942
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #17038504
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #17038504 #17039244
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #17039264
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #17039264 #17040004
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #17040024
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #17040024 #17040764
:   :   :   <---- t0 OS9$81 <I$Detach> {8300} #17037942 #17041114
:   :   <---- t0 OS9$8f <I$Close> {path=4} #17028674 #17044722
:   {proc=4} t0: OS9$4c <F$DelPrc> {} #17044854
:   :   {proc=4} t0: OS9$40 <F$DelTsk> {proc_desc=7300} #17045306
:   :   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7300} #17045306 #17045976
:   :   <---- t0 OS9$4c <F$DelPrc> {} #17044854 #17048744
:   <---- t1 OS9$22 <F$NMLoad> {LangType=11, 0e6d='dir'} #16242912 #17054712
{proc=4"Shell"} t1: OS9$03 <F$Fork> {Module/file='dir' param="/dd{13}" lang/type=11 pages=1f} #17054920
:   {proc=4} t0: OS9$82 <I$Dup> {$2} #17072916
:   :   <---- t0 OS9$82 <I$Dup> {$2} #17072916 -> path $2 #17074310
:   {proc=4} t0: OS9$82 <I$Dup> {$1} #17074374
:   :   <---- t0 OS9$82 <I$Dup> {$1} #17074374 -> path $1 #17075768
:   {proc=4} t0: OS9$82 <I$Dup> {$1} #17075832
:   :   <---- t0 OS9$82 <I$Dup> {$1} #17075832 -> path $1 #17077226
:   {proc=5} t0: OS9$34 <F$SLink> {"dir" type 11 name@ e6d dat@ 7540} #17077416
:   :   <---- t0 OS9$34 <F$SLink> {"dir" type 11 name@ e6d dat@ 7540} #17077416 #17084356
:   {proc=5} t0: OS9$07 <F$Mem> {desired_size=1f00} #17085690
:   :   {proc=5} t0: OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7300} #17086208
:   :   :   <---- t0 OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7300} #17086208 #17087930
:   :   <---- t0 OS9$07 <F$Mem> {desired_size=1f00} #17085690 #17088236
:   {proc=4} t0: OS9$3f <F$AllTsk> {processDesc=7300} #17088560
:   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7300} #17088560 #17091088
:   {proc=4} t0: OS9$40 <F$DelTsk> {proc_desc=7300} #17094612
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7300} #17094612 #17095282
:   {proc=4} t0: OS9$2c <F$AProc> {proc=7300} #17095388
:   :   <---- t0 OS9$2c <F$AProc> {proc=7300} #17095388 #17096112
:   <---- t1 OS9$03 <F$Fork> {Module/file='dir' param="/dd{13}" lang/type=11 pages=1f} #17054920 #17097596
{proc=4"Shell"} t1: OS9$1d <F$UnLoad> {} #17097660
:   {proc=4} t0: OS9$4e <F$FModul> {"Dir" type 11 name@ e3 dat@ 7540} #17098776
:   :   <---- t0 OS9$4e <F$FModul> {"Dir" type 11 name@ e3 dat@ 7540} #17098776 #17103266
:   <---- t1 OS9$1d <F$UnLoad> {} #17097660 #17104516
{proc=4"Shell"} t1: OS9$04 <F$Wait> {} #17105078
:   {proc=4} t0: OS9$40 <F$DelTsk> {proc_desc=7500} #17106512
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7500} #17106512 #17107182
:   {proc=4} t0: OS9$2d <F$NProc> {} #17107236
{proc=5"Dir"} t1: OS9$8d <I$GetStt> {path=1 26==SS.ScSiz  : Return screen size for COCO} #17108792
:   <---- t1 OS9$8d <I$GetStt> {path=1 26==SS.ScSiz  : Return screen size for COCO} #17108792 #17113298
{proc=5"Dir"} t1: OS9$84 <I$Open> {1efc='/dd'} #17113842
:   {proc=5} t0: OS9$80 <I$Attach> {1efd=''} #17145738
:   :   {proc=1} t0: OS9$34 <F$SLink> {"dd" type f0 name@ 1efd dat@ 7340} #17146954
:   :   :   <---- t0 OS9$34 <F$SLink> {"dd" type f0 name@ 1efd dat@ 7340} #17146954 #17172016
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #17172148
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #17172148 -> addr $e704 entry $e72e #17205432
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #17205518
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #17205518 -> addr $9e77 entry $9e88 #17264850
:   :   <---- t0 OS9$80 <I$Attach> {1efd=''} #17145738 #17267096
:   <---- t1 OS9$84 <I$Open> {1efc='/dd'} #17113842 -> path $3 #17391174
{proc=5"Dir"} t1: OS9$86 <I$ChgDir> {mode=81, 1efc='/dd'} #17391222
:   {proc=5} t0: OS9$80 <I$Attach> {1efd=''} #17397918
:   :   {proc=1} t0: OS9$34 <F$SLink> {"dd" type f0 name@ 1efd dat@ 7340} #17399134
:   :   :   <---- t0 OS9$34 <F$SLink> {"dd" type f0 name@ 1efd dat@ 7340} #17399134 #17424154
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #17424286
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #17424286 -> addr $e704 entry $e72e #17457570
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #17457656
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #17457656 -> addr $9e77 entry $9e88 #17516988
:   :   <---- t0 OS9$80 <I$Attach> {1efd=''} #17397918 #17519234
:   {proc=5} t0: OS9$81 <I$Detach> {8300} #17643402
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #17643964
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #17643964 #17644704
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #17644724
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #17644724 #17645464
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #17645484
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #17645484 #17646224
:   :   <---- t0 OS9$81 <I$Detach> {8300} #17643402 #17646574
:   <---- t1 OS9$86 <I$ChgDir> {mode=81, 1efc='/dd'} #17391222 #17648874
{proc=5"Dir"} t1: OS9$15 <F$Time> {buf=d} #17651396
:   <---- t1 OS9$15 <F$Time> {buf=d} #17651396 #17655286
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #17658852
:   <---- t1 OS9$8c <I$WritLn> {} #17658852 #17688664
{proc=5"Dir"} t1: OS9$88 <I$Seek> {path=3 pos=00000040} #17688754
:   <---- t1 OS9$88 <I$Seek> {path=3 pos=00000040} #17688754 #17692280
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17692342
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17692342 #17752874
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17754374
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17754374 #17761930
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17763514
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17763514 #17771070
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17772582
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17772582 #17781772
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #17782536
:   <---- t1 OS9$8c <I$WritLn> {} #17782536 #17814172
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17814242
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17814242 #17821798
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17823250
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17823250 #17830946
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17832416
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17832416 #17893938
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17897902
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17897902 #17905458
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #17906466
:   <---- t1 OS9$8c <I$WritLn> {} #17906466 #17938386
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17938456
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17938456 #17946012
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17947440
:   <-- ERROR: $d3(E$EOF    :End of File): OS9KERNEL1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17947440 #17952726
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #17952932
:   <---- t1 OS9$8c <I$WritLn> {} #17952932 #17973020
{proc=5"Dir"} t1: OS9$06 <F$Exit> {status=0} #17973058
:   {proc=5} t0: OS9$8f <I$Close> {path=2} #17974214
:   :   <---- t0 OS9$8f <I$Close> {path=2} #17974214 #17976046
:   {proc=5} t0: OS9$8f <I$Close> {path=1} #17976124
:   :   <---- t0 OS9$8f <I$Close> {path=1} #17976124 #17978434
:   {proc=5} t0: OS9$8f <I$Close> {path=1} #17978512
:   :   {proc=5} t0: OS9$37 <F$GProcP> {id=04} #17980828
:   :   :   <---- t0 OS9$37 <F$GProcP> {id=04} #17980828 #17981536
:   :   <---- t0 OS9$8f <I$Close> {path=1} #17978512 #17982274
:   {proc=5} t0: OS9$8f <I$Close> {path=4} #17982352
:   :   {proc=5} t0: OS9$81 <I$Detach> {8300} #17989468
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #17990030
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #17990030 #17990770
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #17990790
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #17990790 #17991530
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #17991550
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #17991550 #17992290
:   :   :   <---- t0 OS9$81 <I$Detach> {8300} #17989468 #17992640
:   :   <---- t0 OS9$8f <I$Close> {path=4} #17982352 #17996248
:   {proc=5} t0: OS9$3b <F$DelImg> {beginBlock=0 numBlocks=1 processDesc=7300} #17996662
:   :   <---- t0 OS9$3b <F$DelImg> {beginBlock=0 numBlocks=1 processDesc=7300} #17996662 #17997440
:   {proc=5} t0: OS9$02 <F$UnLink> {u=e000 magic=ae48 module='''} #17997498
:   :   <---- t0 OS9$02 <F$UnLink> {u=e000 magic=ae48 module='''} #17997498 #18007114
:   {proc=5} t0: OS9$40 <F$DelTsk> {proc_desc=7300} #18007154
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7300} #18007154 #18007824
:   {proc=5} t0: OS9$40 <F$DelTsk> {proc_desc=7300} #18008276
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7300} #18008276 #18008870
:   {proc=5} t0: OS9$2c <F$AProc> {proc=7500} #18011420
:   :   <---- t0 OS9$2c <F$AProc> {proc=7500} #18011420 #18012144
:   {proc=5} t0: OS9$2d <F$NProc> {} #18012156
:   :   <---- t1 OS9$04 <F$Wait> {} #17105078 #18014642
{proc=4"Shell"} t1: OS9$8b <I$ReadLn> {} #18020766
:   <---- t1 OS9$8b <I$ReadLn> {} #18020766 #18029270
{proc=4"Shell"} t1: OS9$1c <F$SUser> {} #18032896
:   <---- t1 OS9$1c <F$SUser> {} #18032896 #18035150
{proc=4"Shell"} t1: OS9$21 <F$NMLink> {LangType=00, 0e6d='dir'} #18057706
:   {proc=4} t0: OS9$4e <F$FModul> {"dir" type 0 name@ e6d dat@ 7540} #18058822
:   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$4e <F$FModul> {"dir" type 0 name@ e6d dat@ 7540} #18058822 #18112484
:   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL1 OS9$21 <F$NMLink> {LangType=00, 0e6d='dir'} #18057706 #18113676
{proc=4"Shell"} t1: OS9$84 <I$Open> {0e6d='dir'} #18113704
:   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #18145052
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #18146268
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #18146268 #18174168
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #18174300
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #18174300 -> addr $e704 entry $e72e #18206338
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #18206424
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #18206424 -> addr $9e77 entry $9e88 #18264510
:   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #18145052 #18266756
:   <---- t1 OS9$84 <I$Open> {0e6d='dir'} #18113704 -> path $3 #18611338
{proc=4"Shell"} t1: OS9$89 <I$Read> {path=3 buf=00d6 size=4d} #18613832
:   <---- t1 OS9$89 <I$Read> {path=3 buf=00d6 size=4d} #18613832 #18670086
{proc=4"Shell"} t1: OS9$8f <I$Close> {path=3} #18676496
:   {proc=4} t0: OS9$81 <I$Detach> {8300} #18685050
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #18685612
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #18685612 #18686352
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #18686372
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #18686372 #18687112
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #18687132
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #18687132 #18687872
:   :   <---- t0 OS9$81 <I$Detach> {8300} #18685050 #18688222
:   <---- t1 OS9$8f <I$Close> {path=3} #18676496 #18692814
{proc=4"Shell"} t1: OS9$21 <F$NMLink> {LangType=11, 0e6d='dir'} #18692978
:   {proc=4} t0: OS9$4e <F$FModul> {"dir" type 11 name@ e6d dat@ 7540} #18694094
:   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$4e <F$FModul> {"dir" type 11 name@ e6d dat@ 7540} #18694094 #18747756
:   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL1 OS9$21 <F$NMLink> {LangType=11, 0e6d='dir'} #18692978 #18748948
{proc=4"Shell"} t1: OS9$22 <F$NMLoad> {LangType=11, 0e6d='dir'} #18748966
:   {proc=4} t0: OS9$4b <F$AllPrc> {} #18750070
:   :   <---- t0 OS9$4b <F$AllPrc> {} #18750070 #18766932
:   {proc=4} t0: OS9$84 <I$Open> {0e6d=''} #18767106
:   :   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #18797912
:   :   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #18799128
:   :   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #18799128 #18827028
:   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #18827160
:   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #18827160 -> addr $e704 entry $e72e #18859198
:   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #18859284
:   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #18859284 -> addr $9e77 entry $9e88 #18917370
:   :   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #18797912 #18919616
:   :   <---- t0 OS9$84 <I$Open> {0e6d=''} #18767106 -> path $4 #19265582
:   {proc=4} t0: OS9$3f <F$AllTsk> {processDesc=7300} #19265640
:   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7300} #19265640 #19266534
:   {proc=5} t0: OS9$41 <F$SetTsk> {} #19269310
:   :   <---- t0 OS9$41 <F$SetTsk> {} #19269310 #19270024
:   {proc=5} t0: OS9$89 <I$Read> {path=4 buf=0000 size=9} #19270076
:   :   <---- t0 OS9$89 <I$Read> {path=4 buf=0000 size=9} #19270076 #19323604
:   {proc=5} t0: OS9$89 <I$Read> {path=4 buf=0009 size=398} #19325588
:   :   <---- t0 OS9$89 <I$Read> {path=4 buf=0009 size=398} #19325588 #19501700
:   {proc=5} t0: OS9$89 <I$Read> {path=4 buf=03a1 size=9} #19561678
:   :   <-- ERROR: $d3(E$EOF    :End of File): OS9KERNEL0 OS9$89 <I$Read> {path=4 buf=03a1 size=9} #19561678 #19565234
:   {proc=4} t0: OS9$8f <I$Close> {path=4} #19566968
:   :   {proc=4} t0: OS9$81 <I$Detach> {8300} #19574602
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #19575164
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #19575164 #19575904
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #19575924
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #19575924 #19576664
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #19576684
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #19576684 #19577424
:   :   :   <---- t0 OS9$81 <I$Detach> {8300} #19574602 #19577774
:   :   <---- t0 OS9$8f <I$Close> {path=4} #19566968 #19581382
:   {proc=4} t0: OS9$4c <F$DelPrc> {} #19581514
:   :   {proc=4} t0: OS9$40 <F$DelTsk> {proc_desc=7300} #19581966
:   :   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7300} #19581966 #19582636
:   :   <---- t0 OS9$4c <F$DelPrc> {} #19581514 #19585404
:   <---- t1 OS9$22 <F$NMLoad> {LangType=11, 0e6d='dir'} #18748966 #19591372
{proc=4"Shell"} t1: OS9$03 <F$Fork> {Module/file='dir' param="/b1{13}" lang/type=11 pages=1f} #19591580
:   {proc=4} t0: OS9$82 <I$Dup> {$2} #19609576
:   :   <---- t0 OS9$82 <I$Dup> {$2} #19609576 -> path $2 #19610970
:   {proc=4} t0: OS9$82 <I$Dup> {$1} #19611034
:   :   <---- t0 OS9$82 <I$Dup> {$1} #19611034 -> path $1 #19612428
:   {proc=4} t0: OS9$82 <I$Dup> {$1} #19612492
:   :   <---- t0 OS9$82 <I$Dup> {$1} #19612492 -> path $1 #19613886
:   {proc=5} t0: OS9$34 <F$SLink> {"dir" type 11 name@ e6d dat@ 7540} #19614076
:   :   <---- t0 OS9$34 <F$SLink> {"dir" type 11 name@ e6d dat@ 7540} #19614076 #19621016
:   {proc=5} t0: OS9$07 <F$Mem> {desired_size=1f00} #19622350
:   :   {proc=5} t0: OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7300} #19622868
:   :   :   <---- t0 OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7300} #19622868 #19624590
:   :   <---- t0 OS9$07 <F$Mem> {desired_size=1f00} #19622350 #19624896
:   {proc=4} t0: OS9$3f <F$AllTsk> {processDesc=7300} #19626854
:   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7300} #19626854 #19627748
:   {proc=4} t0: OS9$40 <F$DelTsk> {proc_desc=7300} #19631272
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7300} #19631272 #19631942
:   {proc=4} t0: OS9$2c <F$AProc> {proc=7300} #19632048
:   :   <---- t0 OS9$2c <F$AProc> {proc=7300} #19632048 #19632772
:   <---- t1 OS9$03 <F$Fork> {Module/file='dir' param="/b1{13}" lang/type=11 pages=1f} #19591580 #19634256
{proc=4"Shell"} t1: OS9$1d <F$UnLoad> {} #19634320
:   {proc=4} t0: OS9$4e <F$FModul> {"Dir" type 11 name@ e3 dat@ 7540} #19635436
:   :   <---- t0 OS9$4e <F$FModul> {"Dir" type 11 name@ e3 dat@ 7540} #19635436 #19639926
:   <---- t1 OS9$1d <F$UnLoad> {} #19634320 #19641176
{proc=4"Shell"} t1: OS9$04 <F$Wait> {} #19641738
:   {proc=4} t0: OS9$40 <F$DelTsk> {proc_desc=7500} #19643172
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7500} #19643172 #19643842
:   {proc=4} t0: OS9$2d <F$NProc> {} #19643896
{proc=5"Dir"} t1: OS9$8d <I$GetStt> {path=1 26==SS.ScSiz  : Return screen size for COCO} #19645452
:   <---- t1 OS9$8d <I$GetStt> {path=1 26==SS.ScSiz  : Return screen size for COCO} #19645452 #19649958
{proc=5"Dir"} t1: OS9$84 <I$Open> {1efc='/b1'} #19650502
:   {proc=5} t0: OS9$80 <I$Attach> {1efd=''} #19682448
:   :   {proc=1} t0: OS9$34 <F$SLink> {"b1" type f0 name@ 1efd dat@ 7340} #19683664
:   :   :   <---- t0 OS9$34 <F$SLink> {"b1" type f0 name@ 1efd dat@ 7340} #19683664 #19706984
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #19707116
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #19707116 -> addr $e704 entry $e72e #19740400
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #19740486
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #19740486 -> addr $9e77 entry $9e88 #19799818
:   :   <---- t0 OS9$80 <I$Attach> {1efd=''} #19682448 #19802538
:   <---- t1 OS9$84 <I$Open> {1efc='/b1'} #19650502 -> path $3 #19923422
{proc=5"Dir"} t1: OS9$86 <I$ChgDir> {mode=81, 1efc='/b1'} #19925922
:   {proc=5} t0: OS9$80 <I$Attach> {1efd=''} #19932668
:   :   {proc=1} t0: OS9$34 <F$SLink> {"b1" type f0 name@ 1efd dat@ 7340} #19933884
:   :   :   <---- t0 OS9$34 <F$SLink> {"b1" type f0 name@ 1efd dat@ 7340} #19933884 #19957152
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #19957284
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #19957284 -> addr $e704 entry $e72e #19990568
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #19990654
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #19990654 -> addr $9e77 entry $9e88 #20049986
:   :   <---- t0 OS9$80 <I$Attach> {1efd=''} #19932668 #20052324
:   {proc=5} t0: OS9$81 <I$Detach> {831a} #20174866
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #20175428
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #20175428 #20176168
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #20176188
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #20176188 #20176928
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec61 magic=87cd module='B1'} #20176948
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec61 magic=87cd module='B1'} #20176948 #20177688
:   :   <---- t0 OS9$81 <I$Detach> {831a} #20174866 #20178038
:   <---- t1 OS9$86 <I$ChgDir> {mode=81, 1efc='/b1'} #19925922 #20180338
{proc=5"Dir"} t1: OS9$15 <F$Time> {buf=d} #20182860
:   <---- t1 OS9$15 <F$Time> {buf=d} #20182860 #20186750
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #20187864
:   <---- t1 OS9$8c <I$WritLn> {} #20187864 #20217724
{proc=5"Dir"} t1: OS9$88 <I$Seek> {path=3 pos=00000040} #20217814
:   <---- t1 OS9$88 <I$Seek> {path=3 pos=00000040} #20217814 #20221340
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20221402
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20221402 #20286108
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20287608
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20287608 #20295164
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20296748
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20296748 #20304304
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20305816
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20305816 #20315006
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #20315770
:   <---- t1 OS9$8c <I$WritLn> {} #20315770 #20347466
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20347536
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20347536 #20355092
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20356544
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20356544 #20364240
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20365710
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20365710 #20424436
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20425948
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20425948 #20437870
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #20438878
:   <---- t1 OS9$8c <I$WritLn> {} #20438878 #20470910
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20470980
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20470980 #20478536
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20479964
:   <-- ERROR: $d3(E$EOF    :End of File): OS9KERNEL1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20479964 #20485250
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #20485456
:   <---- t1 OS9$8c <I$WritLn> {} #20485456 #20505544
{proc=5"Dir"} t1: OS9$06 <F$Exit> {status=0} #20505582
:   {proc=5} t0: OS9$8f <I$Close> {path=2} #20506738
:   :   <---- t0 OS9$8f <I$Close> {path=2} #20506738 #20508570
:   {proc=5} t0: OS9$8f <I$Close> {path=1} #20508648
:   :   <---- t0 OS9$8f <I$Close> {path=1} #20508648 #20510958
:   {proc=5} t0: OS9$8f <I$Close> {path=1} #20511036
:   :   {proc=5} t0: OS9$37 <F$GProcP> {id=04} #20513352
:   :   :   <---- t0 OS9$37 <F$GProcP> {id=04} #20513352 #20514060
:   :   <---- t0 OS9$8f <I$Close> {path=1} #20511036 #20514798
:   {proc=5} t0: OS9$8f <I$Close> {path=4} #20514876
:   :   {proc=5} t0: OS9$81 <I$Detach> {831a} #20521992
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #20522554
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #20522554 #20523294
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #20523314
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #20523314 #20524054
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec61 magic=87cd module='B1'} #20524074
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec61 magic=87cd module='B1'} #20524074 #20524814
:   :   :   <---- t0 OS9$81 <I$Detach> {831a} #20521992 #20525164
:   :   <---- t0 OS9$8f <I$Close> {path=4} #20514876 #20528772
:   {proc=5} t0: OS9$3b <F$DelImg> {beginBlock=0 numBlocks=1 processDesc=7300} #20529186
:   :   <---- t0 OS9$3b <F$DelImg> {beginBlock=0 numBlocks=1 processDesc=7300} #20529186 #20529964
:   {proc=5} t0: OS9$02 <F$UnLink> {u=e000 magic=ae48 module='''} #20530022
:   :   <---- t0 OS9$02 <F$UnLink> {u=e000 magic=ae48 module='''} #20530022 #20539638
:   {proc=5} t0: OS9$40 <F$DelTsk> {proc_desc=7300} #20539678
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7300} #20539678 #20540348
:   {proc=5} t0: OS9$40 <F$DelTsk> {proc_desc=7300} #20540800
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7300} #20540800 #20541394
:   {proc=5} t0: OS9$2c <F$AProc> {proc=7500} #20543944
:   :   <---- t0 OS9$2c <F$AProc> {proc=7500} #20543944 #20544668
:   {proc=5} t0: OS9$2d <F$NProc> {} #20544680
:   :   <---- t1 OS9$04 <F$Wait> {} #19641738 #20547166
{proc=4"Shell"} t1: OS9$8b <I$ReadLn> {} #20548542
:   <-- ERROR: $d3(E$EOF    :End of File): OS9KERNEL1 OS9$8b <I$ReadLn> {} #20548542 #20558142
{proc=4"Shell"} t1: OS9$1c <F$SUser> {} #20558528
:   <---- t1 OS9$1c <F$SUser> {} #20558528 #20560782
{proc=4"Shell"} t1: OS9$06 <F$Exit> {status=0} #20560826
:   {proc=4} t0: OS9$8f <I$Close> {path=2} #20561982
:   :   {proc=4} t0: OS9$81 <I$Detach> {8300} #20567528
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #20568090
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #20568090 #20568830
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #20568850
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #20568850 #20569590
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #20569610
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #20569610 #20570350
:   :   :   <---- t0 OS9$81 <I$Detach> {8300} #20567528 #20570700
:   :   <---- t0 OS9$8f <I$Close> {path=2} #20561982 #20571756
:   {proc=4} t0: OS9$8f <I$Close> {path=1} #20571834
:   :   <---- t0 OS9$8f <I$Close> {path=1} #20571834 #20574144
:   {proc=4} t0: OS9$8f <I$Close> {path=1} #20574222
:   :   {proc=4} t0: OS9$37 <F$GProcP> {id=03} #20576538
:   :   :   <---- t0 OS9$37 <F$GProcP> {id=03} #20576538 #20577246
:   :   <---- t0 OS9$8f <I$Close> {path=1} #20574222 #20577954
:   {proc=4} t0: OS9$3b <F$DelImg> {beginBlock=0 numBlocks=1 processDesc=7500} #20578396
:   :   <---- t0 OS9$3b <F$DelImg> {beginBlock=0 numBlocks=1 processDesc=7500} #20578396 #20579174
:   {proc=4} t0: OS9$02 <F$UnLink> {u=e000 magic=ae48 module='''} #20579232
:   :   <---- t0 OS9$02 <F$UnLink> {u=e000 magic=ae48 module='''} #20579232 #20583756
:   {proc=4} t0: OS9$40 <F$DelTsk> {proc_desc=7500} #20583796
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7500} #20583796 #20584466
:   {proc=4} t0: OS9$40 <F$DelTsk> {proc_desc=7500} #20584918
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7500} #20584918 #20585512
:   {proc=4} t0: OS9$2c <F$AProc> {proc=7800} #20588146
:   :   <---- t0 OS9$2c <F$AProc> {proc=7800} #20588146 #20588870
:   {proc=4} t0: OS9$2d <F$NProc> {} #20588882
:   :   <---- t1 OS9$04 <F$Wait> {} #14922642 #20591368
{proc=3"Shell"} t1: OS9$1c <F$SUser> {} #20592658
:   <---- t1 OS9$1c <F$SUser> {} #20592658 #20594912
{proc=3"Shell"} t1: OS9$06 <F$Exit> {status=0} #20594956
:   {proc=3} t0: OS9$8f <I$Close> {path=1} #20596112
:   :   <---- t0 OS9$8f <I$Close> {path=1} #20596112 #20598362
:   {proc=3} t0: OS9$8f <I$Close> {path=1} #20598440
:   :   <---- t0 OS9$8f <I$Close> {path=1} #20598440 #20600750
:   {proc=3} t0: OS9$8f <I$Close> {path=1} #20600828
:   :   {proc=3} t0: OS9$37 <F$GProcP> {id=02} #20603144
:   :   :   <---- t0 OS9$37 <F$GProcP> {id=02} #20603144 #20603852
:   :   <---- t0 OS9$8f <I$Close> {path=1} #20600828 #20604560
:   {proc=3} t0: OS9$3b <F$DelImg> {beginBlock=0 numBlocks=1 processDesc=7800} #20605002
:   :   <---- t0 OS9$3b <F$DelImg> {beginBlock=0 numBlocks=1 processDesc=7800} #20605002 #20605780
:   {proc=3} t0: OS9$02 <F$UnLink> {u=e000 magic=ae48 module='''} #20605838
:   :   <---- t0 OS9$02 <F$UnLink> {u=e000 magic=ae48 module='''} #20605838 #20616822
:   {proc=3} t0: OS9$40 <F$DelTsk> {proc_desc=7800} #20616862
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7800} #20616862 #20617532
:   {proc=3} t0: OS9$40 <F$DelTsk> {proc_desc=7800} #20617984
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7800} #20617984 #20618578
:   {proc=3} t0: OS9$2c <F$AProc> {proc=7a00} #20621268
:   :   <---- t0 OS9$2c <F$AProc> {proc=7a00} #20621268 #20621992
:   {proc=3} t0: OS9$2d <F$NProc> {} #20622004
:   :   <---- t1 OS9$04 <F$Wait> {} #12769486 #20624490
{proc=2"SysGo"} t1: OS9$03 <F$Fork> {Module/file='AutoEx' param="{13}" lang/type=1 pages=0} #20624552
:   {proc=2} t0: OS9$82 <I$Dup> {$1} #20641922
:   :   <---- t0 OS9$82 <I$Dup> {$1} #20641922 -> path $1 #20643316
:   {proc=2} t0: OS9$82 <I$Dup> {$1} #20643380
:   :   <---- t0 OS9$82 <I$Dup> {$1} #20643380 -> path $1 #20644774
:   {proc=2} t0: OS9$82 <I$Dup> {$1} #20644838
:   :   <---- t0 OS9$82 <I$Dup> {$1} #20644838 -> path $1 #20646232
:   {proc=3} t0: OS9$34 <F$SLink> {"AutoEx" type 1 name@ d945 dat@ 7a40} #20646422
:   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$34 <F$SLink> {"AutoEx" type 1 name@ d945 dat@ 7a40} #20646422 #20696668
:   {proc=2} t0: OS9$01 <F$Load> {type/lang=7a filename='AutoEx'} #20696718
:   :   {proc=2} t0: OS9$4b <F$AllPrc> {} #20697092
:   :   :   <---- t0 OS9$4b <F$AllPrc> {} #20697092 #20713804
:   :   {proc=2} t0: OS9$84 <I$Open> {d945='AutoEx'} #20713978
:   :   :   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #20719782
:   :   :   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #20720998
:   :   :   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #20720998 #20736768
:   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #20736900
:   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #20736900 -> addr $e704 entry $e72e #20755982
:   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #20756068
:   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #20756068 -> addr $9e77 entry $9e88 #20802832
:   :   :   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #20719782 #20805170
:   :   :   {proc=2} t0: OS9$81 <I$Detach> {8300} #21750524
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #21751086
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #21751086 #21751826
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #21751846
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #21751846 #21752586
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #21752606
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #21752606 #21753346
:   :   :   :   <---- t0 OS9$81 <I$Detach> {8300} #21750524 #21753696
:   :   :   <-- ERROR: $d8(E$PNNF   :Path Name Not Found): OS9KERNEL0 OS9$84 <I$Open> {d945='AutoEx'} #20713978 #21754762
:   :   {proc=2} t0: OS9$4c <F$DelPrc> {} #21754924
:   :   :   {proc=2} t0: OS9$40 <F$DelTsk> {proc_desc=7600} #21755376
:   :   :   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7600} #21755376 #21755970
:   :   :   <---- t0 OS9$4c <F$DelPrc> {} #21754924 #21758822
:   :   <-- ERROR: $d8(E$PNNF   :Path Name Not Found): OS9KERNEL0 OS9$01 <F$Load> {type/lang=7a filename='AutoEx'} #20696718 #21759174
:   {proc=2} t0: OS9$8f <I$Close> {path=1} #21759350
:   :   <---- t0 OS9$8f <I$Close> {path=1} #21759350 #21761540
:   {proc=2} t0: OS9$8f <I$Close> {path=1} #21761618
:   :   <---- t0 OS9$8f <I$Close> {path=1} #21761618 #21763808
:   {proc=2} t0: OS9$8f <I$Close> {path=1} #21763886
:   :   <---- t0 OS9$8f <I$Close> {path=1} #21763886 #21766076
:   {proc=3} t0: OS9$02 <F$UnLink> {u=0000 magic=0000 module=''} #21766540
:   :   <---- t0 OS9$02 <F$UnLink> {u=0000 magic=0000 module=''} #21766540 #21767208
:   {proc=2} t0: OS9$40 <F$DelTsk> {proc_desc=7800} #21767248
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7800} #21767248 #21767842
:   {proc=2} t0: OS9$40 <F$DelTsk> {proc_desc=7800} #21767992
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7800} #21767992 #21768586
:   <-- ERROR: $d8(E$PNNF   :Path Name Not Found): OS9KERNEL1 OS9$03 <F$Fork> {Module/file='AutoEx' param="{13}" lang/type=1 pages=0} #20624552 #21772772
{proc=2"SysGo"} t1: OS9$05 <F$Chain> {Module/file='Shell' param="i=/1{13}" lang/type=1 pages=0} #21775490
:   {proc=2} t0: OS9$02 <F$UnLink> {u=d893 magic=87cd module='SysGo'} #21795370
:   :   <---- t0 OS9$02 <F$UnLink> {u=d893 magic=87cd module='SysGo'} #21795370 #21796110
:   {proc=2} t0: OS9$34 <F$SLink> {"Shell" type 1 name@ d93f dat@ 7840} #21796552
:   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$34 <F$SLink> {"Shell" type 1 name@ d93f dat@ 7840} #21796552 #21849352
:   {proc=3} t0: OS9$01 <F$Load> {type/lang=78 filename='Shell'} #21849402
:   :   {proc=3} t0: OS9$4b <F$AllPrc> {} #21849776
:   :   :   <---- t0 OS9$4b <F$AllPrc> {} #21849776 #21866488
:   :   {proc=3} t0: OS9$84 <I$Open> {d93f='Shell'} #21866662
:   :   :   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #21872466
:   :   :   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #21873682
:   :   :   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #21873682 #21887818
:   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #21887950
:   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #21887950 -> addr $e704 entry $e72e #21908666
:   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #21908752
:   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #21908752 -> addr $9e77 entry $9e88 #21955516
:   :   :   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #21872466 #21957854
:   :   :   <---- t0 OS9$84 <I$Open> {d93f='Shell'} #21866662 -> path $2 #22713216
:   :   {proc=3} t0: OS9$3f <F$AllTsk> {processDesc=7600} #22713274
:   :   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7600} #22713274 #22714168
:   :   {proc=4} t0: OS9$41 <F$SetTsk> {} #22715160
:   :   :   <---- t0 OS9$41 <F$SetTsk> {} #22715160 #22715874
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=0000 size=9} #22715926
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=0000 size=9} #22715926 #22767884
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=0009 size=1b4e} #22769868
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=0009 size=1b4e} #22769868 #24269546
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1b57 size=9} #24321212
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1b57 size=9} #24321212 #24326730
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1b60 size=e8} #24328714
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1b60 size=e8} #24328714 #24389580
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1c48 size=9} #24441322
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1c48 size=9} #24441322 #24446840
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1c51 size=4a} #24448824
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1c51 size=4a} #24448824 #24455288
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1c9b size=9} #24508952
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1c9b size=9} #24508952 #24514470
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1ca4 size=19} #24516454
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1ca4 size=19} #24516454 #24523834
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1cbd size=9} #24575534
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1cbd size=9} #24575534 #24582686
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1cc6 size=3e} #24584670
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1cc6 size=3e} #24584670 #24641790
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d04 size=9} #24696144
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d04 size=9} #24696144 #24703296
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d0d size=23} #24705280
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d0d size=23} #24705280 #24711208
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d30 size=9} #24767780
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d30 size=9} #24767780 #24773298
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d39 size=1b} #24775282
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d39 size=1b} #24775282 #24781096
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d54 size=9} #24839080
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d54 size=9} #24839080 #24844598
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d5d size=5e} #24848216
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d5d size=5e} #24848216 #24855044
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1dbb size=9} #24914026
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1dbb size=9} #24914026 #24919544
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1dc4 size=1e} #24921528
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1dc4 size=1e} #24921528 #24927444
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1de2 size=9} #24986248
:   :   :   <-- ERROR: $d3(E$EOF    :End of File): OS9KERNEL0 OS9$89 <I$Read> {path=2 buf=1de2 size=9} #24986248 #24989804
:   :   {proc=3} t0: OS9$8f <I$Close> {path=2} #24989904
:   :   :   {proc=3} t0: OS9$81 <I$Detach> {8300} #24999232
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #24999794
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #24999794 #25000534
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #25000554
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #25000554 #25001294
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #25001314
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #25001314 #25002054
:   :   :   :   <---- t0 OS9$81 <I$Detach> {8300} #24999232 #25002404
:   :   :   <---- t0 OS9$8f <I$Close> {path=2} #24989904 #25003460
:   :   {proc=3} t0: OS9$4c <F$DelPrc> {} #25003592
:   :   :   {proc=3} t0: OS9$40 <F$DelTsk> {proc_desc=7600} #25004044
:   :   :   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7600} #25004044 #25004714
:   :   :   <---- t0 OS9$4c <F$DelPrc> {} #25003592 #25007566
:   :   {proc=2} t0: OS9$4d <F$ELink> {} #25011180
:   :   :   <---- t0 OS9$4d <F$ELink> {} #25011180 #25013958
:   :   <---- t0 OS9$01 <F$Load> {type/lang=78 filename='Shell'} #21849402 #25014290
:   {proc=2} t0: OS9$07 <F$Mem> {desired_size=1f00} #25015610
:   :   <---- t0 OS9$07 <F$Mem> {desired_size=1f00} #25015610 #25016410
:   {proc=3} t0: OS9$3f <F$AllTsk> {processDesc=7a00} #25016740
:   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7a00} #25016740 #25017634
:   {proc=3} t0: OS9$40 <F$DelTsk> {proc_desc=7800} #25021352
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7800} #25021352 #25022022
:   {proc=3} t0: OS9$40 <F$DelTsk> {proc_desc=7a00} #25024712
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7a00} #25024712 #25027016
:   {proc=1} t0: OS9$2c <F$AProc> {proc=7a00} #25027078
:   :   <---- t0 OS9$2c <F$AProc> {proc=7a00} #25027078 #25027802
:   {proc=1} t0: OS9$2d <F$NProc> {} #25027814
{proc=2"Shell"} t1: OS9$09 <F$Icpt> {routine=e06b storage=0000} #25036204
:   <---- t1 OS9$09 <F$Icpt> {routine=e06b storage=0000} #25036204 #25038862
{proc=2"Shell"} t1: OS9$0c <F$ID> {} #25038914
:   <---- t1 OS9$0c <F$ID> {} #25038914 #25041186
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25043806
:   <---- t1 OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25043806 #25050258
{proc=2"Shell"} t1: OS9$84 <I$Open> {f7f4='.'} #25051000
:   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #25059328
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #25060544
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #25060544 #25088444
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #25088576
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #25088576 -> addr $e704 entry $e72e #25120614
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #25120700
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #25120700 -> addr $9e77 entry $9e88 #25178786
:   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #25059328 #25181124
:   <---- t1 OS9$84 <I$Open> {f7f4='.'} #25051000 -> path $3 #25410946
{proc=2"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #25411040
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #25411040 #25468664
{proc=2"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #25468784
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #25468784 #25477974
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25478242
:   <---- t1 OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25478242 #25484670
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=3} #25488074
:   {proc=2} t0: OS9$81 <I$Detach> {8300} #25495684
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #25496246
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #25496246 #25496986
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #25497006
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #25497006 #25497746
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #25497766
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #25497766 #25498506
:   :   <---- t0 OS9$81 <I$Detach> {8300} #25495684 #25498856
:   <---- t1 OS9$8f <I$Close> {path=3} #25488074 #25500840
{proc=2"Shell"} t1: OS9$1c <F$SUser> {} #25505638
:   <---- t1 OS9$1c <F$SUser> {} #25505638 #25507892
{proc=2"Shell"} t1: OS9$82 <I$Dup> {$0} #25511790
:   <---- t1 OS9$82 <I$Dup> {$0} #25511790 -> path $3 #25515118
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=0} #25515168
:   <---- t1 OS9$8f <I$Close> {path=0} #25515168 #25519210
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25519594
:   <---- t1 OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25519594 #25526046
{proc=2"Shell"} t1: OS9$84 <I$Open> {00b5='/Term'} #25526078
:   {proc=2} t0: OS9$80 <I$Attach> {00b6='X'} #25534438
:   :   {proc=1} t0: OS9$34 <F$SLink> {"Term" type f0 name@ b6 dat@ 7a40} #25535654
:   :   :   <---- t0 OS9$34 <F$SLink> {"Term" type f0 name@ b6 dat@ 7a40} #25535654 #25575126
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='VTIO'} #25575258
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='VTIO'} #25575258 -> addr $b8dd entry $be5c #25629118
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='SCF'} #25629204
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='SCF'} #25629204 -> addr $b165 entry $b410 #25685482
:   :   <---- t0 OS9$80 <I$Attach> {00b6='X'} #25534438 #25687784
:   {proc=1} t0: OS9$80 <I$Attach> {d511='Term'} #25707530
:   :   {proc=1} t0: OS9$34 <F$SLink> {"Term" type f0 name@ d511 dat@ 640} #25708746
:   :   :   <---- t0 OS9$34 <F$SLink> {"Term" type f0 name@ d511 dat@ 640} #25708746 #25756572
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='VTIO'} #25756704
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='VTIO'} #25756704 -> addr $b8dd entry $be5c #25810564
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='SCF'} #25810650
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='SCF'} #25810650 -> addr $b165 entry $b410 #25866928
:   :   <---- t0 OS9$80 <I$Attach> {d511='Term'} #25707530 #25869230
:   <---- t1 OS9$84 <I$Open> {00b5='/Term'} #25526078 -> path $0 #25872932
{proc=2"Shell"} t1: OS9$82 <I$Dup> {$1} #25873106
:   <---- t1 OS9$82 <I$Dup> {$1} #25873106 -> path $4 #25876468
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=1} #25876518
:   <---- t1 OS9$8f <I$Close> {path=1} #25876518 #25880620
{proc=2"Shell"} t1: OS9$82 <I$Dup> {$0} #25880652
:   <---- t1 OS9$82 <I$Dup> {$0} #25880652 -> path $1 #25883912
{proc=2"Shell"} t1: OS9$82 <I$Dup> {$2} #25884058
:   <---- t1 OS9$82 <I$Dup> {$2} #25884058 -> path $5 #25887454
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=2} #25887504
:   <---- t1 OS9$8f <I$Close> {path=2} #25887504 #25893300
{proc=2"Shell"} t1: OS9$82 <I$Dup> {$1} #25893332
:   <---- t1 OS9$82 <I$Dup> {$1} #25893332 -> path $2 #25896626
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=0 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25896766
:   <---- t1 OS9$8d <I$GetStt> {path=0 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25896766 #25903218
{proc=2"Shell"} t1: OS9$8a <I$Write> {{0}} #25903300
:   <---- t1 OS9$8a <I$Write> {{0}} #25903300 #25913120
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25913654
:   <---- t1 OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25913654 #25920106
{proc=2"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #25923798
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #25923798 #25928186
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #25928228
:   <---- t1 OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #25928228 #25934286
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=3} #25934526
:   <---- t1 OS9$8f <I$Close> {path=3} #25934526 #25938748
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=4} #25938882
:   <---- t1 OS9$8f <I$Close> {path=4} #25938882 #25943164
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=5} #25943298
:   {proc=2} t0: OS9$37 <F$GProcP> {id=01} #25946478
:   :   <---- t0 OS9$37 <F$GProcP> {id=01} #25946478 #25947186
:   <---- t1 OS9$8f <I$Close> {path=5} #25943298 #25948822
{proc=2"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #25962084
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #25962084 #25966472
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #25966514
:   <---- t1 OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #25966514 #25972572
{proc=2"Shell"} t1: OS9$8a <I$Write> {
:   <---- t1 OS9$8a <I$Write> {
{proc=2"Shell"} t1: OS9$15 <F$Time> {buf=2da} #25991154
:   <---- t1 OS9$15 <F$Time> {buf=2da} #25991154 #25995044
{proc=2"Shell"} t1: OS9$8c <I$WritLn> {} #25998350
:   <---- t1 OS9$8c <I$WritLn> {} #25998350 #26020568
{proc=2"Shell"} t1: OS9$1c <F$SUser> {} #26020730
:   <---- t1 OS9$1c <F$SUser> {} #26020730 #26022984
{proc=2"Shell"} t1: OS9$8c <I$WritLn> {} #26023194
:   <---- t1 OS9$8c <I$WritLn> {} #26023194 #26043822
{proc=2"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #26043906
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #26043906 #26048294
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #26048344
:   <---- t1 OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #26048344 #26054402
{proc=2"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.Opt    : Read/Write PD Options} #26054562
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.Opt    : Read/Write PD Options} #26054562 #26062350
{proc=2"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.SSig   : Send signal on data ready} #26062394
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.SSig   : Send signal on data ready} #26062394 #26066850
{proc=2"Shell"} t1: OS9$0a <F$Sleep> {ticks=0000} #26066940
:   {proc=2} t0: OS9$40 <F$DelTsk> {proc_desc=7a00} #26068250
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7a00} #26068250 #26068920
:   {proc=2} t0: OS9$2d <F$NProc> {} #26068974
```

## Long.  All kernel calls, unfiltered.

This is the file `_kern` from above.

```
{proc=1} t0: OS9$2e <F$VModul> {addr=0d06="REL" map=[3f 0 0 0 0 0 0 0]} #2116862
:   <---- t0 OS9$2e <F$VModul> {addr=0d06="REL" map=[3f 0 0 0 0 0 0 0]} #2116862 #2122564
{proc=1} t0: OS9$2e <F$VModul> {addr=0e30="Boot" map=[3f 0 0 0 0 0 0 0]} #2123582
:   <---- t0 OS9$2e <F$VModul> {addr=0e30="Boot" map=[3f 0 0 0 0 0 0 0]} #2123582 #2130936
{proc=1} t0: OS9$2e <F$VModul> {addr=1000="Krn" map=[3f 0 0 0 0 0 0 0]} #2131800
:   <---- t0 OS9$2e <F$VModul> {addr=1000="Krn" map=[3f 0 0 0 0 0 0 0]} #2131800 #2140132
{proc=1} t0: OS9$35 <F$Boot> {} #2146194
:   {proc=1} t0: OS9$00 <F$Link> {type/lang=c1 module/file='Boot'} #2146728
:   :   <---- t0 OS9$00 <F$Link> {type/lang=c1 module/file='Boot'} #2146728 -> addr $ee30 entry $ee42 #2155076
:   {proc=1} t0: OS9$28 <F$SRqMem> {size=65d0} #2156158
:   :   <---- t0 OS9$28 <F$SRqMem> {size=65d0} #2156158 -> size $6600 addr $8700 #2168442
:   {proc=1} t0: OS9$2e <F$VModul> {addr=0700="KrnP2" map=[2 3 4 3f 0 0 0 0]} #3524094
:   :   <---- t0 OS9$2e <F$VModul> {addr=0700="KrnP2" map=[2 3 4 3f 0 0 0 0]} #3524094 #3534708
:   {proc=1} t0: OS9$2e <F$VModul> {addr=13db="IOMan" map=[2 3 4 3f 0 0 0 0]} #3535880
:   :   <---- t0 OS9$2e <F$VModul> {addr=13db="IOMan" map=[2 3 4 3f 0 0 0 0]} #3535880 #3547128
:   {proc=1} t0: OS9$2e <F$VModul> {addr=1e00="Init" map=[2 3 4 3f 0 0 0 0]} #3548146
:   :   <---- t0 OS9$2e <F$VModul> {addr=1e00="Init" map=[2 3 4 3f 0 0 0 0]} #3548146 #3560726
:   {proc=1} t0: OS9$2e <F$VModul> {addr=1e77="RBF" map=[2 3 4 3f 0 0 0 0]} #3561590
:   :   <---- t0 OS9$2e <F$VModul> {addr=1e77="RBF" map=[2 3 4 3f 0 0 0 0]} #3561590 #3575344
:   {proc=1} t0: OS9$2e <F$VModul> {addr=3165="SCF" map=[2 3 4 3f 0 0 0 0]} #3576208
:   :   <---- t0 OS9$2e <F$VModul> {addr=3165="SCF" map=[2 3 4 3f 0 0 0 0]} #3576208 #3591406
:   {proc=1} t0: OS9$2e <F$VModul> {addr=38dd="VTIO" map=[2 3 4 3f 0 0 0 0]} #3592424
:   :   <---- t0 OS9$2e <F$VModul> {addr=38dd="VTIO" map=[2 3 4 3f 0 0 0 0]} #3592424 #3609366
:   {proc=1} t0: OS9$2e <F$VModul> {addr=448e="SndDrv" map=[2 3 4 3f 0 0 0 0]} #3610692
:   :   <---- t0 OS9$2e <F$VModul> {addr=448e="SndDrv" map=[2 3 4 3f 0 0 0 0]} #3610692 #3630526
:   {proc=1} t0: OS9$2e <F$VModul> {addr=4554="JoyDrv" map=[2 3 4 3f 0 0 0 0]} #3631852
:   :   <---- t0 OS9$2e <F$VModul> {addr=4554="JoyDrv" map=[2 3 4 3f 0 0 0 0]} #3631852 #3652730
:   {proc=1} t0: OS9$2e <F$VModul> {addr=46d9="CoGrf" map=[2 3 4 3f 0 0 0 0]} #3653902
:   :   <---- t0 OS9$2e <F$VModul> {addr=46d9="CoGrf" map=[2 3 4 3f 0 0 0 0]} #3653902 #3675998
:   {proc=1} t0: OS9$2e <F$VModul> {addr=54db="Term" map=[2 3 4 3f 0 0 0 0]} #3677016
:   :   <---- t0 OS9$2e <F$VModul> {addr=54db="Term" map=[2 3 4 3f 0 0 0 0]} #3677016 #3700300
:   {proc=1} t0: OS9$2e <F$VModul> {addr=551f="W" map=[2 3 4 3f 0 0 0 0]} #3700856
:   :   <---- t0 OS9$2e <F$VModul> {addr=551f="W" map=[2 3 4 3f 0 0 0 0]} #3700856 #3724748
:   {proc=1} t0: OS9$2e <F$VModul> {addr=5560="W1" map=[2 3 4 3f 0 0 0 0]} #3725458
:   :   <---- t0 OS9$2e <F$VModul> {addr=5560="W1" map=[2 3 4 3f 0 0 0 0]} #3725458 #3751218
:   {proc=1} t0: OS9$2e <F$VModul> {addr=55a2="W2" map=[2 3 4 3f 0 0 0 0]} #3751928
:   :   <---- t0 OS9$2e <F$VModul> {addr=55a2="W2" map=[2 3 4 3f 0 0 0 0]} #3751928 #3779560
:   {proc=1} t0: OS9$2e <F$VModul> {addr=55e4="W3" map=[2 3 4 3f 0 0 0 0]} #3780270
:   :   <---- t0 OS9$2e <F$VModul> {addr=55e4="W3" map=[2 3 4 3f 0 0 0 0]} #3780270 #3809774
:   {proc=1} t0: OS9$2e <F$VModul> {addr=5626="Clock" map=[2 3 4 3f 0 0 0 0]} #3810946
:   :   <---- t0 OS9$2e <F$VModul> {addr=5626="Clock" map=[2 3 4 3f 0 0 0 0]} #3810946 #3842164
:   {proc=1} t0: OS9$2e <F$VModul> {addr=582d="Clock2" map=[2 3 4 3f 0 0 0 0]} #3843490
:   :   <---- t0 OS9$2e <F$VModul> {addr=582d="Clock2" map=[2 3 4 3f 0 0 0 0]} #3843490 #3878102
:   {proc=1} t0: OS9$2e <F$VModul> {addr=5893="SysGo" map=[2 3 4 3f 0 0 0 0]} #3879274
:   :   <---- t0 OS9$2e <F$VModul> {addr=5893="SysGo" map=[2 3 4 3f 0 0 0 0]} #3879274 #3913852
:   {proc=1} t0: OS9$2e <F$VModul> {addr=5a8f="LemMan" map=[2 3 4 3f 0 0 0 0]} #3915178
:   :   <---- t0 OS9$2e <F$VModul> {addr=5a8f="LemMan" map=[2 3 4 3f 0 0 0 0]} #3915178 #3950762
:   {proc=1} t0: OS9$2e <F$VModul> {addr=66b5="Lemmer" map=[2 3 4 3f 0 0 0 0]} #3952088
:   :   <---- t0 OS9$2e <F$VModul> {addr=66b5="Lemmer" map=[2 3 4 3f 0 0 0 0]} #3952088 #3991956
:   {proc=1} t0: OS9$2e <F$VModul> {addr=66df="Lem" map=[2 3 4 3f 0 0 0 0]} #3992820
:   :   <---- t0 OS9$2e <F$VModul> {addr=66df="Lem" map=[2 3 4 3f 0 0 0 0]} #3992820 #4033526
:   {proc=1} t0: OS9$2e <F$VModul> {addr=6704="RBLemma" map=[2 3 4 3f 0 0 0 0]} #4035006
:   :   <---- t0 OS9$2e <F$VModul> {addr=6704="RBLemma" map=[2 3 4 3f 0 0 0 0]} #4035006 #4077906
:   {proc=1} t0: OS9$2e <F$VModul> {addr=6c3c="DD" map=[2 3 4 3f 0 0 0 0]} #4078616
:   :   <---- t0 OS9$2e <F$VModul> {addr=6c3c="DD" map=[2 3 4 3f 0 0 0 0]} #4078616 #4120632
:   {proc=1} t0: OS9$2e <F$VModul> {addr=6c61="B1" map=[2 3 4 3f 0 0 0 0]} #4121342
:   :   <---- t0 OS9$2e <F$VModul> {addr=6c61="B1" map=[2 3 4 3f 0 0 0 0]} #4121342 #4165360
:   {proc=1} t0: OS9$2e <F$VModul> {addr=6c86="B2" map=[2 3 4 3f 0 0 0 0]} #4166070
:   :   <---- t0 OS9$2e <F$VModul> {addr=6c86="B2" map=[2 3 4 3f 0 0 0 0]} #4166070 #4212080
:   {proc=1} t0: OS9$2e <F$VModul> {addr=6cab="B3" map=[2 3 4 3f 0 0 0 0]} #4212790
:   :   <---- t0 OS9$2e <F$VModul> {addr=6cab="B3" map=[2 3 4 3f 0 0 0 0]} #4212790 #4260792
:   <---- t0 OS9$35 <F$Boot> {} #2146194 #4261348
{proc=1} t0: OS9$00 <F$Link> {type/lang=c0 module/file='Init'} #4261394
:   <---- t0 OS9$00 <F$Link> {type/lang=c0 module/file='Init'} #4261394 -> addr $9e00 entry $ade0 #4303174
{proc=1} t0: OS9$00 <F$Link> {type/lang=c0 module/file='krnp2'} #4303400
:   <---- t0 OS9$00 <F$Link> {type/lang=c0 module/file='krnp2'} #4303400 -> addr $8700 entry $8713 #4348666
{proc=1} t0: OS9$32 <F$SSvc> {table=8789} #4348844
:   <---- t0 OS9$32 <F$SSvc> {table=8789} #4348844 #4354632
{proc=1} t0: OS9$86 <I$ChgDir> {mode=05, 9e34='/DD'} #4354832
:   {proc=1} t0: OS9$00 <F$Link> {type/lang=c1 module/file='IOMan'} #4355210
:   :   <---- t0 OS9$00 <F$Link> {type/lang=c1 module/file='IOMan'} #4355210 -> addr $93db entry $93ee #4396422
:   {proc=1} t0: OS9$28 <F$SRqMem> {size=3c9} #4396596
:   :   <---- t0 OS9$28 <F$SRqMem> {size=3c9} #4396596 -> size $400 addr $8300 #4406190
:   {proc=1} t0: OS9$30 <F$All64> {table=0} #4437060
:   :   {proc=1} t0: OS9$28 <F$SRqMem> {size=100} #4437446
:   :   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #4437446 -> size $100 addr $8200 #4447168
:   :   <---- t0 OS9$30 <F$All64> {table=0} #4437060 -> base $8200 blocknum $1 addr $8240 #4455922
:   {proc=1} t0: OS9$31 <F$Ret64> {block_num=1 address=8200} #4455950
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=1 address=8200} #4455950 #4456728
:   {proc=1} t0: OS9$32 <F$SSvc> {table=943d} #4456778
:   :   <---- t0 OS9$32 <F$SSvc> {table=943d} #4456778 #4458622
:   {proc=1} t0: OS9$30 <F$All64> {table=8200} #4458854
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #4458854 -> base $8200 blocknum $1 addr $8240 #4461208
:   {proc=1} t0: OS9$49 <F$LDABX> {} #4461282
:   :   <---- t0 OS9$49 <F$LDABX> {} #4461282 #4462136
:   {proc=1} t0: OS9$10 <F$PrsNam> {path='/DD'} #4462214
:   :   <---- t0 OS9$10 <F$PrsNam> {path='/DD'} #4462214 #4464022
:   {proc=1} t0: OS9$80 <I$Attach> {9e35='DD'} #4464064
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ 9e35 dat@ 640} #4465280
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ 9e35 dat@ 640} #4465280 #4475798
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #4475930
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #4475930 -> addr $e704 entry $e72e #4492114
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #4492200
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #4492200 -> addr $9e77 entry $9e88 #4532798
:   :   {proc=1} t0: OS9$28 <F$SRqMem> {size=a7} #4534488
:   :   :   <---- t0 OS9$28 <F$SRqMem> {size=a7} #4534488 -> size $100 addr $8100 #4543856
:   :   <---- t0 OS9$80 <I$Attach> {9e35='DD'} #4464064 #4552604
:   {proc=1} t0: OS9$28 <F$SRqMem> {size=100} #4553156
:   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #4553156 -> size $100 addr $8000 #4562592
:   {proc=1} t0: OS9$30 <F$All64> {table=8200} #4562640
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #4562640 -> base $8200 blocknum $2 addr $8280 #4565030
:   {proc=1} t0: OS9$49 <F$LDABX> {} #4565972
:   :   <---- t0 OS9$49 <F$LDABX> {} #4565972 #4566826
:   {proc=1} t0: OS9$10 <F$PrsNam> {path='/DD'} #4566898
:   :   <---- t0 OS9$10 <F$PrsNam> {path='/DD'} #4566898 #4568706
:   {proc=1} t0: OS9$29 <F$SRtMem> {size=100 start=8000} #4663054
:   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=8000} #4663054 #4668000
:   {proc=1} t0: OS9$31 <F$Ret64> {block_num=2 address=8200} #4668466
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=2 address=8200} #4668466 #4669244
:   {proc=1} t0: OS9$81 <I$Detach> {8300} #4669628
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #4670190
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #4670190 #4670930
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #4670950
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #4670950 #4671690
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #4671710
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #4671710 #4672450
:   :   <---- t0 OS9$81 <I$Detach> {8300} #4669628 #4672800
:   {proc=1} t0: OS9$31 <F$Ret64> {block_num=1 address=8200} #4672830
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=1 address=8200} #4672830 #4673608
:   <---- t0 OS9$86 <I$ChgDir> {mode=05, 9e34='/DD'} #4354832 #4673856
{proc=1} t0: OS9$84 <I$Open> {9e37='/Term'} #4674062
:   {proc=1} t0: OS9$30 <F$All64> {table=8200} #4674546
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #4674546 -> base $8200 blocknum $1 addr $8240 #4676900
:   {proc=1} t0: OS9$49 <F$LDABX> {} #4676974
:   :   <---- t0 OS9$49 <F$LDABX> {} #4676974 #4677828
:   {proc=1} t0: OS9$10 <F$PrsNam> {path='/Term'} #4677906
:   :   <---- t0 OS9$10 <F$PrsNam> {path='/Term'} #4677906 #4680164
:   {proc=1} t0: OS9$80 <I$Attach> {9e38='Term'} #4680206
:   :   {proc=1} t0: OS9$34 <F$SLink> {"Term" type f0 name@ 9e38 dat@ 640} #4681422
:   :   :   <---- t0 OS9$34 <F$SLink> {"Term" type f0 name@ 9e38 dat@ 640} #4681422 #4710720
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='VTIO'} #4710852
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='VTIO'} #4710852 -> addr $b8dd entry $be5c #4747624
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='SCF'} #4747710
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='SCF'} #4747710 -> addr $b165 entry $b410 #4786072
:   :   {proc=1} t0: OS9$28 <F$SRqMem> {size=100} #4787856
:   :   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #4787856 -> size $100 addr $8000 #4797240
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=c1 module/file='JoyDrv'} #4805622
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=c1 module/file='JoyDrv'} #4805622 -> addr $c554 entry $c568 #4839968
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=c1 module/file='SndDrv'} #4840124
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=c1 module/file='SndDrv'} #4840124 -> addr $c48e entry $c4a2 #4876370
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=c1 module/file='CoWin'} #4876878
:   :   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$00 <F$Link> {type/lang=c1 module/file='CoWin'} #4876878 #4923758
:   :   {proc=1} t0: OS9$01 <F$Load> {type/lang=c1 filename='CoWin'} #4923798
:   :   :   {proc=1} t0: OS9$4b <F$AllPrc> {} #4924172
:   :   :   :   {proc=1} t0: OS9$28 <F$SRqMem> {size=200} #4924662
:   :   :   :   :   <---- t0 OS9$28 <F$SRqMem> {size=200} #4924662 -> size $200 addr $7e00 #4934472
:   :   :   :   <---- t0 OS9$4b <F$AllPrc> {} #4924172 #4938400
:   :   :   {proc=1} t0: OS9$84 <I$Open> {c1ee='CoWin'} #4938574
:   :   :   :   {proc=1} t0: OS9$30 <F$All64> {table=8200} #4939058
:   :   :   :   :   <---- t0 OS9$30 <F$All64> {table=8200} #4939058 -> base $8200 blocknum $2 addr $8280 #4941448
:   :   :   :   {proc=1} t0: OS9$49 <F$LDABX> {} #4941522
:   :   :   :   :   <---- t0 OS9$49 <F$LDABX> {} #4941522 #4942376
:   :   :   :   {proc=1} t0: OS9$10 <F$PrsNam> {path='DD'} #4942568
:   :   :   :   :   <---- t0 OS9$10 <F$PrsNam> {path='DD'} #4942568 #4944336
:   :   :   :   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #4944378
:   :   :   :   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #4945594
:   :   :   :   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #4945594 #4956832
:   :   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #4956964
:   :   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #4956964 -> addr $e704 entry $e72e #4973148
:   :   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #4973234
:   :   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #4973234 -> addr $9e77 entry $9e88 #5013832
:   :   :   :   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #4944378 #5016022
:   :   :   :   {proc=1} t0: OS9$28 <F$SRqMem> {size=100} #5016518
:   :   :   :   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #5016518 -> size $100 addr $7d00 #5026508
:   :   :   :   {proc=1} t0: OS9$30 <F$All64> {table=8200} #5026556
:   :   :   :   :   <---- t0 OS9$30 <F$All64> {table=8200} #5026556 -> base $8200 blocknum $3 addr $82c0 #5028982
:   :   :   :   {proc=1} t0: OS9$49 <F$LDABX> {} #5029924
:   :   :   :   :   <---- t0 OS9$49 <F$LDABX> {} #5029924 #5030778
:   :   :   :   {proc=1} t0: OS9$10 <F$PrsNam> {path='CoWin'} #5124832
:   :   :   :   :   <---- t0 OS9$10 <F$PrsNam> {path='CoWin'} #5124832 #5127280
:   :   :   :   {proc=1} t0: OS9$38 <F$Move> {srcTask=0 destTask=0 srcPtr=c1ee destPtr=82e0 size=0005} #5127534
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=0 srcPtr=c1ee destPtr=82e0 size=0005} #5127534 #5129282
:   :   :   :   {proc=1} t0: OS9$29 <F$SRtMem> {size=100 start=7d00} #5238314
:   :   :   :   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7d00} #5238314 #5243138
:   :   :   :   {proc=1} t0: OS9$31 <F$Ret64> {block_num=3 address=8200} #5243604
:   :   :   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=3 address=8200} #5243604 #5244382
:   :   :   :   {proc=1} t0: OS9$81 <I$Detach> {8300} #5244628
:   :   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #5245190
:   :   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #5245190 #5245930
:   :   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #5245950
:   :   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #5245950 #5246690
:   :   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #5246710
:   :   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #5246710 #5247450
:   :   :   :   :   <---- t0 OS9$81 <I$Detach> {8300} #5244628 #5247800
:   :   :   :   {proc=1} t0: OS9$31 <F$Ret64> {block_num=2 address=8200} #5247830
:   :   :   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=2 address=8200} #5247830 #5248608
:   :   :   :   <-- ERROR: $d8(E$PNNF   :Path Name Not Found): OS9KERNEL0 OS9$84 <I$Open> {c1ee='CoWin'} #4938574 #5248866
:   :   :   {proc=1} t0: OS9$4c <F$DelPrc> {} #5249028
:   :   :   :   {proc=1} t0: OS9$40 <F$DelTsk> {proc_desc=7e00} #5249480
:   :   :   :   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7e00} #5249480 #5250074
:   :   :   :   {proc=1} t0: OS9$29 <F$SRtMem> {size=200 start=7e00} #5250100
:   :   :   :   :   <---- t0 OS9$29 <F$SRtMem> {size=200 start=7e00} #5250100 #5252932
:   :   :   :   <---- t0 OS9$4c <F$DelPrc> {} #5249028 #5253188
:   :   :   <-- ERROR: $d8(E$PNNF   :Path Name Not Found): OS9KERNEL0 OS9$01 <F$Load> {type/lang=c1 filename='CoWin'} #4923798 #5253540
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=c1 module/file='CoGrf'} #5253860
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=c1 module/file='CoGrf'} #5253860 -> addr $c6d9 entry $c9ec #5287654
:   :   {proc=1} t0: OS9$21 <F$NMLink> {LangType=c1, c7cc='grfdrv'} #5288352
:   :   :   {proc=1} t0: OS9$4e <F$FModul> {"grfdrv" type c1 name@ c7cc dat@ 640} #5288742
:   :   :   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$4e <F$FModul> {"grfdrv" type c1 name@ c7cc dat@ 640} #5288742 #5334166
:   :   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$21 <F$NMLink> {LangType=c1, c7cc='grfdrv'} #5288352 #5334430
:   :   {proc=1} t0: OS9$22 <F$NMLoad> {LangType=c1, c7c4='../CMDS/grfdrv'} #5334692
:   :   :   {proc=1} t0: OS9$4b <F$AllPrc> {} #5335070
:   :   :   :   {proc=1} t0: OS9$28 <F$SRqMem> {size=200} #5335560
:   :   :   :   :   <---- t0 OS9$28 <F$SRqMem> {size=200} #5335560 -> size $200 addr $7e00 #5345034
:   :   :   :   <---- t0 OS9$4b <F$AllPrc> {} #5335070 #5348962
:   :   :   {proc=1} t0: OS9$84 <I$Open> {c7c4='../CMDS/grfdrv'} #5349136
:   :   :   :   {proc=1} t0: OS9$30 <F$All64> {table=8200} #5349620
:   :   :   :   :   <---- t0 OS9$30 <F$All64> {table=8200} #5349620 -> base $8200 blocknum $2 addr $8280 #5352010
:   :   :   :   {proc=1} t0: OS9$49 <F$LDABX> {} #5352084
:   :   :   :   :   <---- t0 OS9$49 <F$LDABX> {} #5352084 #5352938
:   :   :   :   {proc=1} t0: OS9$10 <F$PrsNam> {path='DD'} #5353130
:   :   :   :   :   <---- t0 OS9$10 <F$PrsNam> {path='DD'} #5353130 #5354898
:   :   :   :   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #5354940
:   :   :   :   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #5356156
:   :   :   :   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #5356156 #5367394
:   :   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #5367526
:   :   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #5367526 -> addr $e704 entry $e72e #5383710
:   :   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #5383796
:   :   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #5383796 -> addr $9e77 entry $9e88 #5424394
:   :   :   :   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #5354940 #5426584
:   :   :   :   {proc=1} t0: OS9$28 <F$SRqMem> {size=100} #5427080
:   :   :   :   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #5427080 -> size $100 addr $7d00 #5436932
:   :   :   :   {proc=1} t0: OS9$30 <F$All64> {table=8200} #5436980
:   :   :   :   :   <---- t0 OS9$30 <F$All64> {table=8200} #5436980 -> base $8200 blocknum $3 addr $82c0 #5439406
:   :   :   :   {proc=1} t0: OS9$49 <F$LDABX> {} #5440348
:   :   :   :   :   <---- t0 OS9$49 <F$LDABX> {} #5440348 #5441202
:   :   :   :   {proc=1} t0: OS9$10 <F$PrsNam> {path='../CMDS/grfdrv'} #5535256
:   :   :   :   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path='../CMDS/grfdrv'} #5535256 #5536978
:   :   :   :   {proc=1} t0: OS9$49 <F$LDABX> {} #5537130
:   :   :   :   :   <---- t0 OS9$49 <F$LDABX> {} #5537130 #5537984
:   :   :   :   {proc=1} t0: OS9$49 <F$LDABX> {} #5538146
:   :   :   :   :   <---- t0 OS9$49 <F$LDABX> {} #5538146 #5539000
:   :   :   :   {proc=1} t0: OS9$38 <F$Move> {srcTask=0 destTask=0 srcPtr=c7c4 destPtr=82e0 size=0002} #5539324
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=0 srcPtr=c7c4 destPtr=82e0 size=0002} #5539324 #5540970
:   :   :   :   {proc=1} t0: OS9$10 <F$PrsNam> {path='CMDS/grfdrv'} #5636492
:   :   :   :   :   <---- t0 OS9$10 <F$PrsNam> {path='CMDS/grfdrv'} #5636492 #5639050
:   :   :   :   {proc=1} t0: OS9$38 <F$Move> {srcTask=0 destTask=0 srcPtr=c7c7 destPtr=82e0 size=0004} #5639304
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=0 srcPtr=c7c7 destPtr=82e0 size=0004} #5639304 #5641018
:   :   :   :   {proc=1} t0: OS9$10 <F$PrsNam> {path='grfdrv'} #5740832
:   :   :   :   :   <---- t0 OS9$10 <F$PrsNam> {path='grfdrv'} #5740832 #5743460
:   :   :   :   {proc=1} t0: OS9$38 <F$Move> {srcTask=0 destTask=0 srcPtr=c7cc destPtr=82e0 size=0006} #5743714
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=0 srcPtr=c7cc destPtr=82e0 size=0006} #5743714 #5745496
:   :   :   :   <---- t0 OS9$84 <I$Open> {c7c4='../CMDS/grfdrv'} #5349136 -> path $2 #6074750
:   :   :   {proc=1} t0: OS9$3f <F$AllTsk> {processDesc=7e00} #6074808
:   :   :   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7e00} #6074808 #6075672
:   :   :   {proc=2} t0: OS9$41 <F$SetTsk> {} #6076468
:   :   :   :   <---- t0 OS9$41 <F$SetTsk> {} #6076468 #6077182
:   :   :   {proc=2} t0: OS9$89 <I$Read> {path=2 buf=0000 size=9} #6077234
:   :   :   :   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #6077694
:   :   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #6077694 #6078364
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0000 size=0009} #6129658
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0000 size=0009} #6129658 #6131336
:   :   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=0000 size=9} #6077234 #6132272
:   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0000 destPtr=077e size=0009} #6132384
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0000 destPtr=077e size=0009} #6132384 #6133974
:   :   :   {proc=2} t0: OS9$41 <F$SetTsk> {} #6134790
:   :   :   :   <---- t0 OS9$41 <F$SetTsk> {} #6134790 #6135504
:   :   :   {proc=2} t0: OS9$89 <I$Read> {path=2 buf=0009 size=228d} #6135556
:   :   :   :   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #6136016
:   :   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #6136016 #6136686
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d09 destPtr=0009 size=00f7} #6138498
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d09 destPtr=0009 size=00f7} #6138498 #6144598
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0100 size=0100} #6193210
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0100 size=0100} #6193210 #6199300
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0200 size=0100} #6246724
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0200 size=0100} #6246724 #6252814
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0300 size=0100} #6299856
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0300 size=0100} #6299856 #6305946
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0400 size=0100} #6351418
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0400 size=0100} #6351418 #6357508
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0500 size=0100} #6407690
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0500 size=0100} #6407690 #6413780
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0600 size=0100} #6459252
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0600 size=0100} #6459252 #6465342
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0700 size=0100} #6512384
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0700 size=0100} #6512384 #6518474
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0800 size=0100} #6567086
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0800 size=0100} #6567086 #6573176
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0900 size=0100} #6620218
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0900 size=0100} #6620218 #6626308
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0a00 size=0100} #6672162
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0a00 size=0100} #6672162 #6678252
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0b00 size=0100} #6723724
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0b00 size=0100} #6723724 #6729814
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0c00 size=0100} #6776856
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0c00 size=0100} #6776856 #6782946
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0d00 size=0100} #6829988
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0d00 size=0100} #6829988 #6836078
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0e00 size=0100} #6883120
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0e00 size=0100} #6883120 #6889210
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0f00 size=0100} #6936252
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=0f00 size=0100} #6936252 #6942342
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1000 size=0100} #6989384
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1000 size=0100} #6989384 #6995474
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1100 size=0100} #7042516
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1100 size=0100} #7042516 #7048606
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1200 size=0100} #7096030
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1200 size=0100} #7096030 #7102120
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1300 size=0100} #7149162
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1300 size=0100} #7149162 #7155252
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1400 size=0100} #7202294
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1400 size=0100} #7202294 #7208384
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1500 size=0100} #7255426
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1500 size=0100} #7255426 #7261516
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1600 size=0100} #7308558
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1600 size=0100} #7308558 #7314648
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1700 size=0100} #7361690
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1700 size=0100} #7361690 #7367780
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1800 size=0100} #7414822
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1800 size=0100} #7414822 #7420912
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1900 size=0100} #7468336
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1900 size=0100} #7468336 #7474426
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1a00 size=0100} #7521468
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1a00 size=0100} #7521468 #7527558
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1b00 size=0100} #7573030
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1b00 size=0100} #7573030 #7579120
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1c00 size=0100} #7624592
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1c00 size=0100} #7624592 #7630682
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1d00 size=0100} #7676154
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1d00 size=0100} #7676154 #7682244
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1e00 size=0100} #7729286
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1e00 size=0100} #7729286 #7735376
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1f00 size=0100} #7780848
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=1f00 size=0100} #7780848 #7786938
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=2000 size=0100} #7833980
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=2000 size=0100} #7833980 #7840158
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=2100 size=0100} #7887582
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=2100 size=0100} #7887582 #7893760
:   :   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=2200 size=0096} #7940776
:   :   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7d00 destPtr=2200 size=0096} #7940776 #7945106
:   :   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=0009 size=228d} #6135556 #7946042
:   :   :   {proc=2} t0: OS9$2e <F$VModul> {addr=0000="-" map=[8 1e 27 3b 8 27 3e 8]} #7946180
:   :   :   :   <---- t0 OS9$2e <F$VModul> {addr=0000="-" map=[8 1e 27 3b 8 27 3e 8]} #7946180 #7990002
:   :   :   {proc=2} t0: OS9$89 <I$Read> {path=2 buf=2296 size=9} #7990368
:   :   :   :   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #7990828
:   :   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #7990828 #7991498
:   :   :   :   <-- ERROR: $d3(E$EOF    :End of File): OS9KERNEL0 OS9$89 <I$Read> {path=2 buf=2296 size=9} #7990368 #7993924
:   :   :   {proc=1} t0: OS9$8f <I$Close> {path=2} #7994024
:   :   :   :   {proc=1} t0: OS9$2f <F$Find64> {base=8200 id=2} #7994472
:   :   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #7994472 #7995142
:   :   :   :   {proc=1} t0: OS9$29 <F$SRtMem> {size=100 start=7d00} #7995604
:   :   :   :   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7d00} #7995604 #7999368
:   :   :   :   {proc=1} t0: OS9$31 <F$Ret64> {block_num=3 address=8200} #7999834
:   :   :   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=3 address=8200} #7999834 #8000612
:   :   :   :   {proc=1} t0: OS9$81 <I$Detach> {8300} #8000882
:   :   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #8001444
:   :   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #8001444 #8002184
:   :   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #8002204
:   :   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #8002204 #8002944
:   :   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #8002964
:   :   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #8002964 #8003704
:   :   :   :   :   <---- t0 OS9$81 <I$Detach> {8300} #8000882 #8004054
:   :   :   :   {proc=1} t0: OS9$31 <F$Ret64> {block_num=2 address=8200} #8004084
:   :   :   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=2 address=8200} #8004084 #8004862
:   :   :   :   <---- t0 OS9$8f <I$Close> {path=2} #7994024 #8005110
:   :   :   {proc=1} t0: OS9$4c <F$DelPrc> {} #8005242
:   :   :   :   {proc=1} t0: OS9$40 <F$DelTsk> {proc_desc=7e00} #8005694
:   :   :   :   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7e00} #8005694 #8006364
:   :   :   :   {proc=1} t0: OS9$29 <F$SRtMem> {size=200 start=7e00} #8006390
:   :   :   :   :   <---- t0 OS9$29 <F$SRtMem> {size=200 start=7e00} #8006390 #8009222
:   :   :   :   <---- t0 OS9$4c <F$DelPrc> {} #8005242 #8009478
:   :   :   {proc=1} t0: OS9$48 <F$LDDDXY> {} #8011472
:   :   :   :   <---- t0 OS9$48 <F$LDDDXY> {} #8011472 #8012318
:   :   :   {proc=1} t0: OS9$48 <F$LDDDXY> {} #8012446
:   :   :   :   <---- t0 OS9$48 <F$LDDDXY> {} #8012446 #8013292
:   :   :   <---- t0 OS9$22 <F$NMLoad> {LangType=c1, c7c4='../CMDS/grfdrv'} #5334692 #8013550
:   :   {proc=1} t0: OS9$4e <F$FModul> {"grfdrv" type c1 name@ c7cc dat@ 640} #8015324
:   :   :   <---- t0 OS9$4e <F$FModul> {"grfdrv" type c1 name@ c7cc dat@ 640} #8015324 #8022466
:   :   {proc=1} t0: OS9$48 <F$LDDDXY> {} #8022718
:   :   :   <---- t0 OS9$48 <F$LDDDXY> {} #8022718 #8023564
:   :   {proc=1} t0: OS9$28 <F$SRqMem> {size=2ff} #8023622
:   :   :   <---- t0 OS9$28 <F$SRqMem> {size=2ff} #8023622 -> size $300 addr $7d00 #8033190
:   :   {proc=1} t0: OS9$13 <F$AllBit> {bitmap=1071 first=0 count=1} #8058394
:   :   :   {proc=1} t0: OS9$49 <F$LDABX> {} #8058986
:   :   :   :   <---- t0 OS9$49 <F$LDABX> {} #8058986 #8059752
:   :   :   {proc=1} t0: OS9$4a <F$STABX> {} #8059772
:   :   :   :   <---- t0 OS9$4a <F$STABX> {} #8059772 #8060540
:   :   :   <---- t0 OS9$13 <F$AllBit> {bitmap=1071 first=0 count=1} #8058394 #8060794
:   :   <---- t0 OS9$80 <I$Attach> {9e38='Term'} #4680206 #8062034
:   {proc=1} t0: OS9$10 <F$PrsNam> {path='/Term'} #8063388
:   :   <---- t0 OS9$10 <F$PrsNam> {path='/Term'} #8063388 #8065646
:   {proc=1} t0: OS9$28 <F$SRqMem> {size=100} #8065708
:   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #8065708 -> size $100 addr $7c00 #8075628
:   {proc=1} t0: OS9$80 <I$Attach> {d511='Term'} #8081736
:   :   {proc=1} t0: OS9$34 <F$SLink> {"Term" type f0 name@ d511 dat@ 640} #8082952
:   :   :   <---- t0 OS9$34 <F$SLink> {"Term" type f0 name@ d511 dat@ 640} #8082952 #8115056
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='VTIO'} #8115188
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='VTIO'} #8115188 -> addr $b8dd entry $be5c #8153326
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='SCF'} #8153412
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='SCF'} #8153412 -> addr $b165 entry $b410 #8193140
:   :   <---- t0 OS9$80 <I$Attach> {d511='Term'} #8081736 #8195386
:   <---- t0 OS9$84 <I$Open> {9e37='/Term'} #4674062 -> path $1 #8199346
{proc=1} t0: OS9$82 <I$Dup> {$1} #8199384
:   {proc=1} t0: OS9$2f <F$Find64> {base=8200 id=1} #8199832
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #8199832 #8200502
:   <---- t0 OS9$82 <I$Dup> {$1} #8199384 -> path $1 #8200778
{proc=1} t0: OS9$82 <I$Dup> {$1} #8200800
:   {proc=1} t0: OS9$2f <F$Find64> {base=8200 id=1} #8201248
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #8201248 #8201918
:   <---- t0 OS9$82 <I$Dup> {$1} #8200800 -> path $1 #8202194
{proc=1} t0: OS9$00 <F$Link> {type/lang=c0 module/file='krnp3'} #8202230
:   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$00 <F$Link> {type/lang=c0 module/file='krnp3'} #8202230 #8248804
{proc=1} t0: OS9$03 <F$Fork> {Module/file='SysGo' param="" lang/type=1 pages=0} #8249016
:   {proc=1} t0: OS9$28 <F$SRqMem> {size=200} #8249510
:   :   <---- t0 OS9$28 <F$SRqMem> {size=200} #8249510 -> size $200 addr $7a00 #8259862
:   {proc=1} t0: OS9$82 <I$Dup> {$1} #8264194
:   :   {proc=1} t0: OS9$2f <F$Find64> {base=8200 id=1} #8264642
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #8264642 #8265312
:   :   <---- t0 OS9$82 <I$Dup> {$1} #8264194 -> path $1 #8265588
:   {proc=1} t0: OS9$82 <I$Dup> {$1} #8265652
:   :   {proc=1} t0: OS9$2f <F$Find64> {base=8200 id=1} #8266100
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #8266100 #8266770
:   :   <---- t0 OS9$82 <I$Dup> {$1} #8265652 -> path $1 #8267046
:   {proc=1} t0: OS9$82 <I$Dup> {$1} #8267110
:   :   {proc=1} t0: OS9$2f <F$Find64> {base=8200 id=1} #8267558
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #8267558 #8268228
:   :   <---- t0 OS9$82 <I$Dup> {$1} #8267110 -> path $1 #8268504
:   {proc=2} t0: OS9$34 <F$SLink> {"SysGo" type 1 name@ 9e2f dat@ 640} #8268694
:   :   <---- t0 OS9$34 <F$SLink> {"SysGo" type 1 name@ 9e2f dat@ 640} #8268694 #8290974
:   {proc=2} t0: OS9$48 <F$LDDDXY> {} #8291140
:   :   <---- t0 OS9$48 <F$LDDDXY> {} #8291140 #8292226
:   {proc=2} t0: OS9$07 <F$Mem> {desired_size=fc} #8292254
:   :   {proc=2} t0: OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7a00} #8292772
:   :   :   <---- t0 OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7a00} #8292772 #8294362
:   :   <---- t0 OS9$07 <F$Mem> {desired_size=fc} #8292254 #8294668
:   {proc=1} t0: OS9$3f <F$AllTsk> {processDesc=7a00} #8294992
:   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7a00} #8294992 #8295856
:   {proc=1} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=9e00 destPtr=0100 size=0000} #8295966
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=9e00 destPtr=0100 size=0000} #8295966 #8296564
:   {proc=1} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7bf4 destPtr=00f4 size=000c} #8296630
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7bf4 destPtr=00f4 size=000c} #8296630 #8298410
:   {proc=1} t0: OS9$40 <F$DelTsk> {proc_desc=7a00} #8298440
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7a00} #8298440 #8299110
:   {proc=1} t0: OS9$2c <F$AProc> {proc=7a00} #8299216
:   :   <---- t0 OS9$2c <F$AProc> {proc=7a00} #8299216 #8299940
:   <---- t0 OS9$03 <F$Fork> {Module/file='SysGo' param="" lang/type=1 pages=0} #8249016 #8300180
{proc=1} t0: OS9$2d <F$NProc> {} #8300198
{proc=2"SysGo"} t1: OS9$09 <F$Icpt> {routine=da8b storage=0000} #8301674
:   <---- t1 OS9$09 <F$Icpt> {routine=da8b storage=0000} #8301674 #8303952
{proc=2"SysGo"} t1: OS9$0c <F$ID> {} #8303964
:   <---- t1 OS9$0c <F$ID> {} #8303964 #8306236
{proc=2"SysGo"} t1: OS9$0d <F$SPrior> {pid=02 priority=80} #8306252
:   <---- t1 OS9$0d <F$SPrior> {pid=02 priority=80} #8306252 #8308656
{proc=2"SysGo"} t1: OS9$00 <F$Link> {type/lang=00 module/file='Init'} #8308696
:   <---- t1 OS9$00 <F$Link> {type/lang=00 module/file='Init'} #8308696 -> addr $9e00 entry $ade0 #8352448
{proc=2"SysGo"} t1: OS9$8a <I$Write> {NitrOS-9/6809 Level 2 V3.3.0} #8353444
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #8354712
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #8354712 #8355382
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=9e40 destPtr=7ce4 size=001c} #8356526
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=9e40 destPtr=7ce4 size=001c} #8356526 #8358622
{proc=2"d"} t1: OS9$53 <F$AlHRAM> {} #8361902
:   <---- t1 OS9$53 <F$AlHRAM> {} #8361902 #8362782
<---- t1 OS9$8a <I$Write> {NitrOS-9/6809 Level 2 V3.3.0} #8353444 #8411218
{proc=2"SysGo"} t1: OS9$8c <I$WritLn> {} #8411284
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #8412552
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #8412552 #8413222
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d8c3 destPtr=7cff size=0001} #8414358
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d8c3 destPtr=7cff size=0001} #8414358 #8415970
:   <---- t1 OS9$8c <I$WritLn> {} #8411284 #8423858
{proc=2"SysGo"} t1: OS9$8a <I$Write> {Tandy Color Computer 3} #8424676
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #8425944
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #8425944 #8426614
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=9e5d destPtr=7cea size=0016} #8427758
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=9e5d destPtr=7cea size=0016} #8427758 #8429808
:   <---- t1 OS9$8a <I$Write> {Tandy Color Computer 3} #8424676 #8440428
{proc=2"SysGo"} t1: OS9$8c <I$WritLn> {} #8440494
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #8441762
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #8441762 #8442432
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d8c3 destPtr=7cff size=0001} #8443568
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d8c3 destPtr=7cff size=0001} #8443568 #8445180
:   <---- t1 OS9$8c <I$WritLn> {} #8440494 #8453068
{proc=2"SysGo"} t1: OS9$8a <I$Write> {(C) 2014 The NitrOS-9 Project
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #8454410
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #8454410 #8455080
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d8a6 destPtr=7ce0 size=0020} #8456230
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d8a6 destPtr=7ce0 size=0020} #8456230 #8458304
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d8c6 destPtr=7ce0 size=0020} #8478674
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d8c6 destPtr=7ce0 size=0020} #8478674 #8480748
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d8e6 destPtr=7ce0 size=0020} #8500408
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d8e6 destPtr=7ce0 size=0020} #8500408 #8502482
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d906 destPtr=7ce0 size=0020} #8522078
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d906 destPtr=7ce0 size=0020} #8522078 #8524152
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d926 destPtr=7cf4 size=000c} #8543662
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d926 destPtr=7cf4 size=000c} #8543662 #8545530
:   <---- t1 OS9$8a <I$Write> {(C) 2014 The NitrOS-9 Project
{proc=2"SysGo"} t1: OS9$16 <F$STime> {y85 m12 d31 h23 m59 s59} #8561324
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d95d destPtr=0028 size=0006} #8562440
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d95d destPtr=0028 size=0006} #8562440 #8564134
:   {proc=1} t0: OS9$00 <F$Link> {type/lang=c1 module/file='Clock'} #8564204
:   :   <---- t0 OS9$00 <F$Link> {type/lang=c1 module/file='Clock'} #8564204 -> addr $d626 entry $d7c2 #8589886
:   {proc=1} t0: OS9$00 <F$Link> {type/lang=21 module/file='Clock2'} #8589992
:   :   <---- t0 OS9$00 <F$Link> {type/lang=21 module/file='Clock2'} #8589992 -> addr $d82d entry $d841 #8614326
:   {proc=2} t0: OS9$32 <F$SSvc> {table=d639} #8614610
:   :   <---- t0 OS9$32 <F$SSvc> {table=d639} #8614610 #8615766
:   <---- t1 OS9$16 <F$STime> {y85 m12 d31 h23 m59 s59} #8561324 #8618608
{proc=2"SysGo"} t1: OS9$86 <I$ChgDir> {mode=04, d93a='CMDS'} #8618642
:   {proc=2} t0: OS9$30 <F$All64> {table=8200} #8619862
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #8619862 -> base $8200 blocknum $2 addr $8280 #8622252
:   {proc=2} t0: OS9$49 <F$LDABX> {} #8622326
:   :   <---- t0 OS9$49 <F$LDABX> {} #8622326 #8623180
:   {proc=1} t0: OS9$10 <F$PrsNam> {path='DD'} #8623372
:   :   <---- t0 OS9$10 <F$PrsNam> {path='DD'} #8623372 #8625140
:   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #8625182
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #8626398
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #8626398 #8639042
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #8639174
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #8639174 -> addr $e704 entry $e72e #8661100
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #8661186
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #8661186 -> addr $9e77 entry $9e88 #8704824
:   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #8625182 #8708704
:   {proc=2} t0: OS9$28 <F$SRqMem> {size=100} #8709256
:   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #8709256 -> size $100 addr $7900 #8719096
:   {proc=2} t0: OS9$30 <F$All64> {table=8200} #8719144
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #8719144 -> base $8200 blocknum $3 addr $82c0 #8721570
:   {proc=2} t0: OS9$49 <F$LDABX> {} #8722508
:   :   <---- t0 OS9$49 <F$LDABX> {} #8722508 #8723362
:   {proc=2} t0: OS9$10 <F$PrsNam> {path='CMDS'} #8827580
:   :   <---- t0 OS9$10 <F$PrsNam> {path='CMDS'} #8827580 #8830138
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d93a destPtr=82e0 size=0004} #8830392
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d93a destPtr=82e0 size=0004} #8830392 #8832106
:   {proc=2} t0: OS9$10 <F$PrsNam> {path=''} #8931416
:   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path=''} #8931416 #8932898
:   {proc=2} t0: OS9$29 <F$SRtMem> {size=100 start=7900} #8933828
:   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7900} #8933828 #8938540
:   {proc=2} t0: OS9$31 <F$Ret64> {block_num=3 address=8200} #8939006
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=3 address=8200} #8939006 #8939784
:   {proc=2} t0: OS9$81 <I$Detach> {8300} #8940124
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #8940686
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #8940686 #8941426
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #8941446
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #8941446 #8942186
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #8942206
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #8942206 #8942946
:   :   <---- t0 OS9$81 <I$Detach> {8300} #8940124 #8943296
:   {proc=2} t0: OS9$31 <F$Ret64> {block_num=2 address=8200} #8943326
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=2 address=8200} #8943326 #8945738
:   <---- t1 OS9$86 <I$ChgDir> {mode=04, d93a='CMDS'} #8618642 #8947230
{proc=2"SysGo"} t1: OS9$86 <I$ChgDir> {mode=01, d932='/DD'} #8947264
:   {proc=2} t0: OS9$30 <F$All64> {table=8200} #8948484
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #8948484 -> base $8200 blocknum $2 addr $8280 #8950874
:   {proc=2} t0: OS9$49 <F$LDABX> {} #8950948
:   :   <---- t0 OS9$49 <F$LDABX> {} #8950948 #8951802
:   {proc=2} t0: OS9$10 <F$PrsNam> {path='/DD'} #8951880
:   :   <---- t0 OS9$10 <F$PrsNam> {path='/DD'} #8951880 #8954058
:   {proc=2} t0: OS9$80 <I$Attach> {d933='DD'} #8954100
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ d933 dat@ 7a40} #8955316
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ d933 dat@ 7a40} #8955316 #8967970
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #8968102
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #8968102 -> addr $e704 entry $e72e #8987326
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #8987412
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #8987412 -> addr $9e77 entry $9e88 #9031050
:   :   <---- t0 OS9$80 <I$Attach> {d933='DD'} #8954100 #9033296
:   {proc=2} t0: OS9$28 <F$SRqMem> {size=100} #9035482
:   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #9035482 -> size $100 addr $7900 #9045606
:   {proc=2} t0: OS9$30 <F$All64> {table=8200} #9045654
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #9045654 -> base $8200 blocknum $3 addr $82c0 #9048080
:   {proc=2} t0: OS9$49 <F$LDABX> {} #9049018
:   :   <---- t0 OS9$49 <F$LDABX> {} #9049018 #9049872
:   {proc=2} t0: OS9$10 <F$PrsNam> {path='/DD'} #9049944
:   :   <---- t0 OS9$10 <F$PrsNam> {path='/DD'} #9049944 #9052122
:   {proc=2} t0: OS9$10 <F$PrsNam> {path=''} #9148802
:   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path=''} #9148802 #9150284
:   {proc=2} t0: OS9$29 <F$SRtMem> {size=100 start=7900} #9151214
:   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7900} #9151214 #9156500
:   {proc=2} t0: OS9$31 <F$Ret64> {block_num=3 address=8200} #9156966
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=3 address=8200} #9156966 #9157744
:   {proc=2} t0: OS9$81 <I$Detach> {8300} #9158084
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #9158646
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #9158646 #9159386
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #9159406
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #9159406 #9160146
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #9160166
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #9160166 #9160906
:   :   <---- t0 OS9$81 <I$Detach> {8300} #9158084 #9161256
:   {proc=2} t0: OS9$31 <F$Ret64> {block_num=2 address=8200} #9161286
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=2 address=8200} #9161286 #9162064
:   <---- t1 OS9$86 <I$ChgDir> {mode=01, d932='/DD'} #8947264 #9163556
{proc=2"SysGo"} t1: OS9$86 <I$ChgDir> {mode=04, d936='/DD/CMDS'} #9163596
:   {proc=2} t0: OS9$30 <F$All64> {table=8200} #9164816
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #9164816 -> base $8200 blocknum $2 addr $8280 #9167206
:   {proc=2} t0: OS9$49 <F$LDABX> {} #9167280
:   :   <---- t0 OS9$49 <F$LDABX> {} #9167280 #9168134
:   {proc=2} t0: OS9$10 <F$PrsNam> {path='/DD/CMDS'} #9168212
:   :   <---- t0 OS9$10 <F$PrsNam> {path='/DD/CMDS'} #9168212 #9170390
:   {proc=2} t0: OS9$80 <I$Attach> {d937='DD/CMDS'} #9170432
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD/CMDS" type f0 name@ d937 dat@ 7a40} #9171648
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD/CMDS" type f0 name@ d937 dat@ 7a40} #9171648 #9185936
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #9186068
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #9186068 -> addr $e704 entry $e72e #9203658
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #9203744
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #9203744 -> addr $9e77 entry $9e88 #9249016
:   :   <---- t0 OS9$80 <I$Attach> {d937='DD/CMDS'} #9170432 #9251262
:   {proc=2} t0: OS9$28 <F$SRqMem> {size=100} #9251814
:   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #9251814 -> size $100 addr $7900 #9261938
:   {proc=2} t0: OS9$30 <F$All64> {table=8200} #9261986
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #9261986 -> base $8200 blocknum $3 addr $82c0 #9264412
:   {proc=2} t0: OS9$49 <F$LDABX> {} #9265350
:   :   <---- t0 OS9$49 <F$LDABX> {} #9265350 #9266204
:   {proc=2} t0: OS9$10 <F$PrsNam> {path='/DD/CMDS'} #9266276
:   :   <---- t0 OS9$10 <F$PrsNam> {path='/DD/CMDS'} #9266276 #9268454
:   {proc=2} t0: OS9$10 <F$PrsNam> {path='CMDS'} #9369280
:   :   <---- t0 OS9$10 <F$PrsNam> {path='CMDS'} #9369280 #9371838
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d93a destPtr=82e0 size=0004} #9372092
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d93a destPtr=82e0 size=0004} #9372092 #9373806
:   {proc=2} t0: OS9$10 <F$PrsNam> {path=''} #9473740
:   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path=''} #9473740 #9475222
:   {proc=2} t0: OS9$29 <F$SRtMem> {size=100 start=7900} #9476152
:   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7900} #9476152 #9479804
:   {proc=2} t0: OS9$31 <F$Ret64> {block_num=3 address=8200} #9480270
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=3 address=8200} #9480270 #9482682
:   {proc=2} t0: OS9$81 <I$Detach> {8300} #9483022
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #9483584
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #9483584 #9484324
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #9484344
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #9484344 #9485084
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #9485104
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #9485104 #9485844
:   :   <---- t0 OS9$81 <I$Detach> {8300} #9483022 #9486194
:   {proc=2} t0: OS9$31 <F$Ret64> {block_num=2 address=8200} #9486224
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=2 address=8200} #9486224 #9487002
:   <---- t1 OS9$86 <I$ChgDir> {mode=04, d936='/DD/CMDS'} #9163596 #9488494
{proc=2"SysGo"} t1: OS9$0c <F$ID> {} #9488524
:   <---- t1 OS9$0c <F$ID> {} #9488524 #9490796
{proc=2"SysGo"} t1: OS9$18 <F$GPrDsc> {} #9490828
:   {proc=2} t0: OS9$37 <F$GProcP> {id=02} #9491918
:   :   <---- t0 OS9$37 <F$GProcP> {id=02} #9491918 #9492626
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7a00 destPtr=0000 size=0200} #9492680
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7a00 destPtr=0000 size=0200} #9492680 #9503786
:   <---- t1 OS9$18 <F$GPrDsc> {} #9490828 #9504954
{proc=2"SysGo"} t1: OS9$4f <F$MapBlk> {beginningBlock=0 numBlocks=1} #9504996
:   {proc=2} t0: OS9$3e <F$FreeHB> {} #9506184
:   :   <---- t0 OS9$3e <F$FreeHB> {} #9506184 #9507302
:   {proc=2} t0: OS9$3c <F$SetImg> {} #9507392
:   :   <---- t0 OS9$3c <F$SetImg> {} #9507392 #9508160
:   <---- t1 OS9$4f <F$MapBlk> {beginningBlock=0 numBlocks=1} #9504996 #9509802
{proc=2"SysGo"} t1: OS9$8d <I$GetStt> {path=1 27==SS.KySns  : Getstat/SetStat for COCO keyboard} #9510374
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #9513336
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #9513336 #9514006
:   <---- t1 OS9$8d <I$GetStt> {path=1 27==SS.KySns  : Getstat/SetStat for COCO keyboard} #9510374 #9515938
{proc=2"SysGo"} t1: OS9$03 <F$Fork> {Module/file='Shell' param="startup{32}-p{13}" lang/type=1 pages=0} #9516016
:   {proc=2} t0: OS9$28 <F$SRqMem> {size=200} #9517254
:   :   <---- t0 OS9$28 <F$SRqMem> {size=200} #9517254 -> size $200 addr $7800 #9527472
:   {proc=2} t0: OS9$82 <I$Dup> {$1} #9531804
:   :   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #9532252
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #9532252 #9532922
:   :   <---- t0 OS9$82 <I$Dup> {$1} #9531804 -> path $1 #9533198
:   {proc=2} t0: OS9$82 <I$Dup> {$1} #9533262
:   :   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #9533710
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #9533710 #9534380
:   :   <---- t0 OS9$82 <I$Dup> {$1} #9533262 -> path $1 #9534656
:   {proc=2} t0: OS9$82 <I$Dup> {$1} #9534720
:   :   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #9535168
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #9535168 #9535838
:   :   <---- t0 OS9$82 <I$Dup> {$1} #9534720 -> path $1 #9536114
:   {proc=3} t0: OS9$34 <F$SLink> {"Shell" type 1 name@ d93f dat@ 7a40} #9536304
:   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$34 <F$SLink> {"Shell" type 1 name@ d93f dat@ 7a40} #9536304 #9587664
:   {proc=2} t0: OS9$01 <F$Load> {type/lang=7a filename='Shell'} #9587714
:   :   {proc=2} t0: OS9$4b <F$AllPrc> {} #9588088
:   :   :   {proc=2} t0: OS9$28 <F$SRqMem> {size=200} #9588614
:   :   :   :   <---- t0 OS9$28 <F$SRqMem> {size=200} #9588614 -> size $200 addr $7600 #9599238
:   :   :   <---- t0 OS9$4b <F$AllPrc> {} #9588088 #9604800
:   :   {proc=2} t0: OS9$84 <I$Open> {d93f='Shell'} #9604974
:   :   :   {proc=2} t0: OS9$30 <F$All64> {table=8200} #9605458
:   :   :   :   <---- t0 OS9$30 <F$All64> {table=8200} #9605458 -> base $8200 blocknum $2 addr $8280 #9607848
:   :   :   {proc=2} t0: OS9$49 <F$LDABX> {} #9607922
:   :   :   :   <---- t0 OS9$49 <F$LDABX> {} #9607922 #9608776
:   :   :   {proc=1} t0: OS9$10 <F$PrsNam> {path='DD'} #9608968
:   :   :   :   <---- t0 OS9$10 <F$PrsNam> {path='DD'} #9608968 #9610736
:   :   :   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #9610778
:   :   :   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #9611994
:   :   :   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #9611994 #9624638
:   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #9624770
:   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #9624770 -> addr $e704 entry $e72e #9643994
:   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #9644080
:   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #9644080 -> addr $9e77 entry $9e88 #9687718
:   :   :   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #9610778 #9691598
:   :   :   {proc=2} t0: OS9$28 <F$SRqMem> {size=100} #9692094
:   :   :   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #9692094 -> size $100 addr $7500 #9702206
:   :   :   {proc=2} t0: OS9$30 <F$All64> {table=8200} #9702254
:   :   :   :   <---- t0 OS9$30 <F$All64> {table=8200} #9702254 -> base $8200 blocknum $3 addr $82c0 #9704680
:   :   :   {proc=2} t0: OS9$49 <F$LDABX> {} #9705618
:   :   :   :   <---- t0 OS9$49 <F$LDABX> {} #9705618 #9706472
:   :   :   {proc=2} t0: OS9$10 <F$PrsNam> {path='Shell'} #9811700
:   :   :   :   <---- t0 OS9$10 <F$PrsNam> {path='Shell'} #9811700 #9814408
:   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d93f destPtr=82e0 size=0005} #9814662
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d93f destPtr=82e0 size=0005} #9814662 #9816410
:   :   :   {proc=2} t0: OS9$10 <F$PrsNam> {path=''} #10453308
:   :   :   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path=''} #10453308 #10454790
:   :   :   <---- t0 OS9$84 <I$Open> {d93f='Shell'} #9604974 -> path $2 #10455994
:   :   {proc=2} t0: OS9$3f <F$AllTsk> {processDesc=7600} #10456052
:   :   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7600} #10456052 #10456946
:   :   {proc=4} t0: OS9$41 <F$SetTsk> {} #10457938
:   :   :   <---- t0 OS9$41 <F$SetTsk> {} #10457938 #10458652
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=0000 size=9} #10458704
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #10459164
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #10459164 #10459834
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0000 size=0009} #10509618
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0000 size=0009} #10509618 #10511296
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=0000 size=9} #10458704 #10512232
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=0000 destPtr=7bb6 size=0009} #10512344
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=0000 destPtr=7bb6 size=0009} #10512344 #10514022
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=0009 size=1b4e} #10514216
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #10514676
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #10514676 #10515346
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7509 destPtr=0009 size=00f7} #10517108
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7509 destPtr=0009 size=00f7} #10517108 #10523208
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0100 size=0100} #10575020
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0100 size=0100} #10575020 #10581110
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0200 size=0100} #10629782
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0200 size=0100} #10629782 #10635872
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0300 size=0100} #10684544
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0300 size=0100} #10684544 #10690634
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0400 size=0100} #10739306
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0400 size=0100} #10739306 #10745396
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0500 size=0100} #10795744
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0500 size=0100} #10795744 #10801834
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0600 size=0100} #10850506
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0600 size=0100} #10850506 #10858230
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0700 size=0100} #10907284
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0700 size=0100} #10907284 #10915008
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0800 size=0100} #10962110
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0800 size=0100} #10962110 #10968200
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0900 size=0100} #11016872
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0900 size=0100} #11016872 #11022962
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0a00 size=0100} #11071634
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0a00 size=0100} #11071634 #11077724
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0b00 size=0100} #11126396
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0b00 size=0100} #11126396 #11132486
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0c00 size=0100} #11181158
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0c00 size=0100} #11181158 #11188882
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0d00 size=0100} #11237554
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0d00 size=0100} #11237554 #11245278
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0e00 size=0100} #11293950
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0e00 size=0100} #11293950 #11300040
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0f00 size=0100} #11350728
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0f00 size=0100} #11350728 #11356818
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1000 size=0100} #11405490
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1000 size=0100} #11405490 #11411580
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1100 size=0100} #11460252
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1100 size=0100} #11460252 #11466342
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1200 size=0100} #11515014
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1200 size=0100} #11515014 #11521104
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1300 size=0100} #11571410
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1300 size=0100} #11571410 #11577500
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1400 size=0100} #11626172
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1400 size=0100} #11626172 #11633896
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1500 size=0100} #11680998
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1500 size=0100} #11680998 #11687088
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1600 size=0100} #11737712
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1600 size=0100} #11737712 #11743802
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1700 size=0100} #11795614
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1700 size=0100} #11795614 #11801704
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1800 size=0100} #11850376
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1800 size=0100} #11850376 #11856466
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1900 size=0100} #11903568
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1900 size=0100} #11903568 #11909658
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1a00 size=0100} #11956760
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1a00 size=0100} #11956760 #11964484
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1b00 size=0057} #12014700
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1b00 size=0057} #12014700 #12019242
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=0009 size=1b4e} #10514216 #12020178
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=0000="-" map=[2e 1 26 2e 26 26 26 26]} #12020316
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=0000="-" map=[2e 1 26 2e 26 26 26 26]} #12020316 #12067806
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1b57 size=9} #12068168
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #12068628
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12068628 #12069298
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7557 destPtr=1b57 size=0009} #12071072
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7557 destPtr=1b57 size=0009} #12071072 #12072750
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1b57 size=9} #12068168 #12073686
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1b57 destPtr=7bb6 size=0009} #12073798
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1b57 destPtr=7bb6 size=0009} #12073798 #12075476
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1b60 size=e8} #12075670
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #12076130
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12076130 #12078434
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7560 destPtr=1b60 size=00a0} #12080196
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7560 destPtr=1b60 size=00a0} #12080196 #12084462
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1c00 size=0048} #12134678
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1c00 size=0048} #12134678 #12138868
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1b60 size=e8} #12075670 #12139804
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=1b57="-" map=[2e 1 26 2e 26 26 26 26]} #12139942
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=1b57="-" map=[2e 1 26 2e 26 26 26 26]} #12139942 #12187826
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1c48 size=9} #12188092
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #12188552
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12188552 #12189222
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7548 destPtr=1c48 size=0009} #12190996
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7548 destPtr=1c48 size=0009} #12190996 #12192674
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1c48 size=9} #12188092 #12193610
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1c48 destPtr=7bb6 size=0009} #12193722
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1c48 destPtr=7bb6 size=0009} #12193722 #12195400
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1c51 size=4a} #12195594
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #12197688
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12197688 #12198358
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7551 destPtr=1c51 size=004a} #12200132
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7551 destPtr=1c51 size=004a} #12200132 #12202756
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1c51 size=4a} #12195594 #12203692
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=1c48="-" map=[2e 1 26 2e 26 26 26 26]} #12203830
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=1c48="-" map=[2e 1 26 2e 26 26 26 26]} #12203830 #12253856
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1c9b size=9} #12254122
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #12254582
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12254582 #12255252
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=759b destPtr=1c9b size=0009} #12258660
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=759b destPtr=1c9b size=0009} #12258660 #12260338
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1c9b size=9} #12254122 #12261274
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1c9b destPtr=7bb6 size=0009} #12261386
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1c9b destPtr=7bb6 size=0009} #12261386 #12263064
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1ca4 size=19} #12263258
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #12263718
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12263718 #12264388
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75a4 destPtr=1ca4 size=0019} #12266162
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75a4 destPtr=1ca4 size=0019} #12266162 #12268068
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1ca4 size=19} #12263258 #12269004
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=1c9b="-" map=[2e 1 26 2e 26 26 26 26]} #12269142
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=1c9b="-" map=[2e 1 26 2e 26 26 26 26]} #12269142 #12320692
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1cbd size=9} #12320958
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #12321418
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12321418 #12322088
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75bd destPtr=1cbd size=0009} #12323862
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75bd destPtr=1cbd size=0009} #12323862 #12325540
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1cbd size=9} #12320958 #12326476
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1cbd destPtr=7bb6 size=0009} #12326588
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1cbd destPtr=7bb6 size=0009} #12326588 #12328266
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1cc6 size=3e} #12328460
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #12328920
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12328920 #12329590
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75c6 destPtr=1cc6 size=003a} #12331352
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75c6 destPtr=1cc6 size=003a} #12331352 #12333748
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1d00 size=0004} #12385534
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1d00 size=0004} #12385534 #12387160
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1cc6 size=3e} #12328460 #12388096
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=1cbd="-" map=[2e 1 26 2e 26 26 26 26]} #12388234
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=1cbd="-" map=[2e 1 26 2e 26 26 26 26]} #12388234 #12442658
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d04 size=9} #12442924
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #12443384
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12443384 #12444054
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7504 destPtr=1d04 size=0009} #12445828
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7504 destPtr=1d04 size=0009} #12445828 #12447506
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d04 size=9} #12442924 #12448442
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1d04 destPtr=7bb6 size=0009} #12448554
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1d04 destPtr=7bb6 size=0009} #12448554 #12450232
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d0d size=23} #12450426
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #12450886
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12450886 #12451556
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=750d destPtr=1d0d size=0023} #12453330
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=750d destPtr=1d0d size=0023} #12453330 #12455418
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d0d size=23} #12450426 #12456354
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=1d04="-" map=[2e 1 26 2e 26 26 26 26]} #12456492
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=1d04="-" map=[2e 1 26 2e 26 26 26 26]} #12456492 #12511720
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d30 size=9} #12511986
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #12512446
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12512446 #12513116
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7530 destPtr=1d30 size=0009} #12514890
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7530 destPtr=1d30 size=0009} #12514890 #12516568
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d30 size=9} #12511986 #12517504
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1d30 destPtr=7bb6 size=0009} #12517616
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1d30 destPtr=7bb6 size=0009} #12517616 #12519294
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d39 size=1b} #12519488
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #12519948
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12519948 #12520618
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7539 destPtr=1d39 size=001b} #12522392
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7539 destPtr=1d39 size=001b} #12522392 #12526000
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d39 size=1b} #12519488 #12526936
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=1d30="-" map=[2e 1 26 2e 26 26 26 26]} #12527074
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=1d30="-" map=[2e 1 26 2e 26 26 26 26]} #12527074 #12582300
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d54 size=9} #12582566
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #12583026
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12583026 #12583696
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7554 destPtr=1d54 size=0009} #12587104
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7554 destPtr=1d54 size=0009} #12587104 #12588782
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d54 size=9} #12582566 #12589718
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1d54 destPtr=7bb6 size=0009} #12589830
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1d54 destPtr=7bb6 size=0009} #12589830 #12591508
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d5d size=5e} #12591702
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #12592162
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12592162 #12592832
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=755d destPtr=1d5d size=005e} #12594606
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=755d destPtr=1d5d size=005e} #12594606 #12597594
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d5d size=5e} #12591702 #12598530
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=1d54="-" map=[2e 1 26 2e 26 26 26 26]} #12598668
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=1d54="-" map=[2e 1 26 2e 26 26 26 26]} #12598668 #12656746
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1dbb size=9} #12657012
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #12657472
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12657472 #12658142
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75bb destPtr=1dbb size=0009} #12659916
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75bb destPtr=1dbb size=0009} #12659916 #12661594
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1dbb size=9} #12657012 #12662530
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1dbb destPtr=7bb6 size=0009} #12662642
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1dbb destPtr=7bb6 size=0009} #12662642 #12664320
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1dc4 size=1e} #12664514
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #12664974
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12664974 #12665644
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75c4 destPtr=1dc4 size=001e} #12667418
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75c4 destPtr=1dc4 size=001e} #12667418 #12669494
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1dc4 size=1e} #12664514 #12670430
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=1dbb="-" map=[2e 1 26 2e 26 26 26 26]} #12670568
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=1dbb="-" map=[2e 1 26 2e 26 26 26 26]} #12670568 #12728688
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1de2 size=9} #12728954
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #12729414
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12729414 #12730084
:   :   :   <-- ERROR: $d3(E$EOF    :End of File): OS9KERNEL0 OS9$89 <I$Read> {path=2 buf=1de2 size=9} #12728954 #12732510
:   :   {proc=2} t0: OS9$8f <I$Close> {path=2} #12732610
:   :   :   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #12734692
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #12734692 #12735362
:   :   :   {proc=2} t0: OS9$29 <F$SRtMem> {size=100 start=7500} #12735824
:   :   :   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7500} #12735824 #12740424
:   :   :   {proc=2} t0: OS9$31 <F$Ret64> {block_num=3 address=8200} #12740890
:   :   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=3 address=8200} #12740890 #12741668
:   :   :   {proc=2} t0: OS9$81 <I$Detach> {8300} #12741938
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #12742500
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #12742500 #12743240
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #12743260
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #12743260 #12744000
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #12744020
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #12744020 #12744760
:   :   :   :   <---- t0 OS9$81 <I$Detach> {8300} #12741938 #12745110
:   :   :   {proc=2} t0: OS9$31 <F$Ret64> {block_num=2 address=8200} #12745140
:   :   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=2 address=8200} #12745140 #12745918
:   :   :   <---- t0 OS9$8f <I$Close> {path=2} #12732610 #12746166
:   :   {proc=2} t0: OS9$4c <F$DelPrc> {} #12746298
:   :   :   {proc=2} t0: OS9$40 <F$DelTsk> {proc_desc=7600} #12746750
:   :   :   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7600} #12746750 #12747420
:   :   :   {proc=2} t0: OS9$29 <F$SRtMem> {size=200 start=7600} #12747446
:   :   :   :   <---- t0 OS9$29 <F$SRtMem> {size=200 start=7600} #12747446 #12750016
:   :   :   <---- t0 OS9$4c <F$DelPrc> {} #12746298 #12750272
:   :   {proc=3} t0: OS9$48 <F$LDDDXY> {} #12752432
:   :   :   <---- t0 OS9$48 <F$LDDDXY> {} #12752432 #12753278
:   :   {proc=3} t0: OS9$4d <F$ELink> {} #12753320
:   :   :   <---- t0 OS9$4d <F$ELink> {} #12753320 #12756098
:   :   <---- t0 OS9$01 <F$Load> {type/lang=7a filename='Shell'} #9587714 #12756430
:   {proc=3} t0: OS9$48 <F$LDDDXY> {} #12756596
:   :   <---- t0 OS9$48 <F$LDDDXY> {} #12756596 #12757722
:   {proc=3} t0: OS9$07 <F$Mem> {desired_size=1f00} #12757750
:   :   {proc=3} t0: OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7800} #12758268
:   :   :   <---- t0 OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7800} #12758268 #12759594
:   :   <---- t0 OS9$07 <F$Mem> {desired_size=1f00} #12757750 #12759900
:   {proc=2} t0: OS9$3f <F$AllTsk> {processDesc=7800} #12760224
:   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7800} #12760224 #12761118
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=3 srcPtr=d94d destPtr=1ef5 size=000b} #12761228
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=3 srcPtr=d94d destPtr=1ef5 size=000b} #12761228 #12764608
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=79f4 destPtr=1ee9 size=000c} #12764674
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=79f4 destPtr=1ee9 size=000c} #12764674 #12766454
:   {proc=2} t0: OS9$40 <F$DelTsk> {proc_desc=7800} #12766484
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7800} #12766484 #12767154
:   {proc=2} t0: OS9$2c <F$AProc> {proc=7800} #12767260
:   :   <---- t0 OS9$2c <F$AProc> {proc=7800} #12767260 #12767984
:   <---- t1 OS9$03 <F$Fork> {Module/file='Shell' param="startup{32}-p{13}" lang/type=1 pages=0} #9516016 #12769468
{proc=2"SysGo"} t1: OS9$04 <F$Wait> {} #12769486
:   {proc=2} t0: OS9$40 <F$DelTsk> {proc_desc=7a00} #12770920
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7a00} #12770920 #12771590
:   {proc=2} t0: OS9$2d <F$NProc> {} #12771644
{proc=3"Shell"} t1: OS9$09 <F$Icpt> {routine=e06b storage=0000} #12780034
:   <---- t1 OS9$09 <F$Icpt> {routine=e06b storage=0000} #12780034 #12782312
{proc=3"Shell"} t1: OS9$0c <F$ID> {} #12782364
:   <---- t1 OS9$0c <F$ID> {} #12782364 #12784636
{proc=3"Shell"} t1: OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #12787256
:   {proc=3} t0: OS9$2f <F$Find64> {base=8200 id=1} #12788584
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #12788584 #12789254
:   {proc=3} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=d511 destPtr=0214 size=0020} #12790528
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=d511 destPtr=0214 size=0020} #12790528 #12792514
:   <---- t1 OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #12787256 #12795290
{proc=3"Shell"} t1: OS9$84 <I$Open> {f7f4='.'} #12796032
:   {proc=3} t0: OS9$30 <F$All64> {table=8200} #12797464
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #12797464 -> base $8200 blocknum $2 addr $8280 #12799854
:   {proc=3} t0: OS9$49 <F$LDABX> {} #12799928
:   :   <---- t0 OS9$49 <F$LDABX> {} #12799928 #12800782
:   {proc=1} t0: OS9$10 <F$PrsNam> {path='DD'} #12800968
:   :   <---- t0 OS9$10 <F$PrsNam> {path='DD'} #12800968 #12802736
:   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #12802778
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #12803994
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #12803994 #12831734
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #12831866
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #12831866 -> addr $e704 entry $e72e #12863744
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #12863830
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #12863830 -> addr $9e77 entry $9e88 #12921798
:   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #12802778 #12924044
:   {proc=3} t0: OS9$28 <F$SRqMem> {size=100} #12924540
:   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #12924540 -> size $100 addr $7700 #12934998
:   {proc=3} t0: OS9$30 <F$All64> {table=8200} #12935046
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #12935046 -> base $8200 blocknum $3 addr $82c0 #12937472
:   {proc=3} t0: OS9$49 <F$LDABX> {} #12938410
:   :   <---- t0 OS9$49 <F$LDABX> {} #12938410 #12939264
:   {proc=3} t0: OS9$10 <F$PrsNam> {path='a'} #13035380
:   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path='a'} #13035380 #13037142
:   {proc=3} t0: OS9$49 <F$LDABX> {} #13037294
:   :   <---- t0 OS9$49 <F$LDABX> {} #13037294 #13038148
:   {proc=3} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=f7f4 destPtr=82e0 size=0001} #13038472
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=f7f4 destPtr=82e0 size=0001} #13038472 #13040084
:   {proc=3} t0: OS9$10 <F$PrsNam> {path='"'} #13139556
:   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path='"'} #13139556 #13141078
:   <---- t1 OS9$84 <I$Open> {f7f4='.'} #12796032 -> path $3 #13143592
{proc=3"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #13143686
:   {proc=3} t0: OS9$2f <F$Find64> {base=8200 id=2} #13144948
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #13144948 #13145618
:   {proc=3} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7700 destPtr=03f8 size=0020} #13193832
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7700 destPtr=03f8 size=0020} #13193832 #13195818
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #13143686 #13197682
{proc=3"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #13197802
:   {proc=3} t0: OS9$2f <F$Find64> {base=8200 id=2} #13199064
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #13199064 #13199734
:   {proc=3} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7720 destPtr=03f8 size=0020} #13201508
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7720 destPtr=03f8 size=0020} #13201508 #13203494
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #13197802 #13205358
{proc=3"Shell"} t1: OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #13205626
:   {proc=3} t0: OS9$2f <F$Find64> {base=8200 id=2} #13206954
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #13206954 #13207624
:   {proc=3} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=ec52 destPtr=00b5 size=0020} #13208874
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=ec52 destPtr=00b5 size=0020} #13208874 #13212494
:   <---- t1 OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #13205626 #13213688
{proc=3"Shell"} t1: OS9$10 <F$PrsNam> {path='DD'} #13213714
:   <---- t1 OS9$10 <F$PrsNam> {path='DD'} #13213714 #13216856
{proc=3"Shell"} t1: OS9$8f <I$Close> {path=3} #13217092
:   {proc=3} t0: OS9$2f <F$Find64> {base=8200 id=2} #13218404
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #13218404 #13219074
:   {proc=3} t0: OS9$29 <F$SRtMem> {size=100 start=7700} #13219536
:   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7700} #13219536 #13223132
:   {proc=3} t0: OS9$31 <F$Ret64> {block_num=3 address=8200} #13223598
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=3 address=8200} #13223598 #13224376
:   {proc=3} t0: OS9$81 <I$Detach> {8300} #13224646
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #13225208
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #13225208 #13225948
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #13225968
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #13225968 #13226708
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #13226728
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #13226728 #13227468
:   :   <---- t0 OS9$81 <I$Detach> {8300} #13224646 #13227818
:   {proc=3} t0: OS9$31 <F$Ret64> {block_num=2 address=8200} #13227848
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=2 address=8200} #13227848 #13228626
:   <---- t1 OS9$8f <I$Close> {path=3} #13217092 #13229802
{proc=3"Shell"} t1: OS9$1c <F$SUser> {} #13232352
:   <---- t1 OS9$1c <F$SUser> {} #13232352 #13234606
{proc=3"Shell"} t1: OS9$10 <F$PrsNam> {path='startup'} #13245986
:   <---- t1 OS9$10 <F$PrsNam> {path='startup'} #13245986 #13250558
{proc=3"Shell"} t1: OS9$10 <F$PrsNam> {path=''} #13250618
:   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL1 OS9$10 <F$PrsNam> {path=''} #13250618 #13253694
{proc=3"Shell"} t1: OS9$21 <F$NMLink> {LangType=00, 0e6d='startup'} #13257724
:   {proc=3} t0: OS9$4e <F$FModul> {"startup" type 0 name@ e6d dat@ 7840} #13258840
:   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$4e <F$FModul> {"startup" type 0 name@ e6d dat@ 7840} #13258840 #13314140
:   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL1 OS9$21 <F$NMLink> {LangType=00, 0e6d='startup'} #13257724 #13315332
{proc=3"Shell"} t1: OS9$84 <I$Open> {0e6d='startup'} #13315360
:   {proc=3} t0: OS9$30 <F$All64> {table=8200} #13316792
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #13316792 -> base $8200 blocknum $2 addr $8280 #13319182
:   {proc=3} t0: OS9$49 <F$LDABX> {} #13319256
:   :   <---- t0 OS9$49 <F$LDABX> {} #13319256 #13320022
:   {proc=1} t0: OS9$10 <F$PrsNam> {path='DD'} #13320214
:   :   <---- t0 OS9$10 <F$PrsNam> {path='DD'} #13320214 #13321982
:   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #13322024
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #13323240
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #13323240 #13350980
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #13351112
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #13351112 -> addr $e704 entry $e72e #13382990
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #13383076
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #13383076 -> addr $9e77 entry $9e88 #13441002
:   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #13322024 #13443248
:   {proc=3} t0: OS9$28 <F$SRqMem> {size=100} #13443744
:   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #13443744 -> size $100 addr $7700 #13455836
:   {proc=3} t0: OS9$30 <F$All64> {table=8200} #13455884
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #13455884 -> base $8200 blocknum $3 addr $82c0 #13458310
:   {proc=3} t0: OS9$49 <F$LDABX> {} #13459248
:   :   <---- t0 OS9$49 <F$LDABX> {} #13459248 #13460014
:   {proc=3} t0: OS9$10 <F$PrsNam> {path=''} #13558898
:   :   <---- t0 OS9$10 <F$PrsNam> {path=''} #13558898 #13561816
:   {proc=3} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0e6d destPtr=82e0 size=0007} #13562070
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0e6d destPtr=82e0 size=0007} #13562070 #13563798
:   {proc=3} t0: OS9$29 <F$SRtMem> {size=100 start=7700} #14380798
:   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7700} #14380798 #14384394
:   {proc=3} t0: OS9$31 <F$Ret64> {block_num=3 address=8200} #14384860
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=3 address=8200} #14384860 #14385638
:   {proc=3} t0: OS9$81 <I$Detach> {8300} #14385884
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #14386446
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #14386446 #14387186
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #14387206
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #14387206 #14387946
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #14387966
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #14387966 #14388706
:   :   <---- t0 OS9$81 <I$Detach> {8300} #14385884 #14389056
:   {proc=3} t0: OS9$31 <F$Ret64> {block_num=2 address=8200} #14389086
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=2 address=8200} #14389086 #14389864
:   <-- ERROR: $d8(E$PNNF   :Path Name Not Found): OS9KERNEL1 OS9$84 <I$Open> {0e6d='startup'} #13315360 #14391782
{proc=3"Shell"} t1: OS9$82 <I$Dup> {$0} #14392028
:   {proc=3} t0: OS9$2f <F$Find64> {base=8200 id=1} #14393430
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #14393430 #14394100
:   <---- t1 OS9$82 <I$Dup> {$0} #14392028 -> path $3 #14395356
{proc=3"Shell"} t1: OS9$8f <I$Close> {path=0} #14395406
:   {proc=3} t0: OS9$2f <F$Find64> {base=8200 id=1} #14396718
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #14396718 #14397388
:   <---- t1 OS9$8f <I$Close> {path=0} #14395406 #14399308
{proc=3"Shell"} t1: OS9$84 <I$Open> {0e6d='startup'} #14399528
:   {proc=3} t0: OS9$30 <F$All64> {table=8200} #14400858
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #14400858 -> base $8200 blocknum $2 addr $8280 #14403248
:   {proc=3} t0: OS9$49 <F$LDABX> {} #14403322
:   :   <---- t0 OS9$49 <F$LDABX> {} #14403322 #14405670
:   {proc=1} t0: OS9$10 <F$PrsNam> {path='DD'} #14405856
:   :   <---- t0 OS9$10 <F$PrsNam> {path='DD'} #14405856 #14407624
:   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #14407666
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #14408882
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #14408882 #14436622
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #14436754
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #14436754 -> addr $e704 entry $e72e #14468632
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #14468718
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #14468718 -> addr $9e77 entry $9e88 #14526644
:   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #14407666 #14528890
:   {proc=3} t0: OS9$28 <F$SRqMem> {size=100} #14529386
:   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #14529386 -> size $100 addr $7700 #14539844
:   {proc=3} t0: OS9$30 <F$All64> {table=8200} #14539892
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #14539892 -> base $8200 blocknum $3 addr $82c0 #14542318
:   {proc=3} t0: OS9$49 <F$LDABX> {} #14543256
:   :   <---- t0 OS9$49 <F$LDABX> {} #14543256 #14544022
:   {proc=3} t0: OS9$10 <F$PrsNam> {path=''} #14644530
:   :   <---- t0 OS9$10 <F$PrsNam> {path=''} #14644530 #14647448
:   {proc=3} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0e6d destPtr=82e0 size=0007} #14647702
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0e6d destPtr=82e0 size=0007} #14647702 #14649430
:   {proc=3} t0: OS9$10 <F$PrsNam> {path=''} #14804004
:   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path=''} #14804004 #14805426
:   <---- t1 OS9$84 <I$Open> {0e6d='startup'} #14399528 -> path $0 #14807940
{proc=3"Shell"} t1: OS9$21 <F$NMLink> {LangType=11, e00d='Shell'} #14808826
:   {proc=3} t0: OS9$4e <F$FModul> {"Shell" type 11 name@ e00d dat@ 7840} #14809942
:   :   <---- t0 OS9$4e <F$FModul> {"Shell" type 11 name@ e00d dat@ 7840} #14809942 #14829876
:   {proc=3} t0: OS9$48 <F$LDDDXY> {} #14830026
:   :   <---- t0 OS9$48 <F$LDDDXY> {} #14830026 #14830872
:   <---- t1 OS9$21 <F$NMLink> {LangType=11, e00d='Shell'} #14808826 #14832058
{proc=3"Shell"} t1: OS9$03 <F$Fork> {Module/file='Shell' param="-P{32}X{32}PATH=;-p{13}" lang/type=11 pages=1f} #14832254
:   {proc=3} t0: OS9$28 <F$SRqMem> {size=200} #14833510
:   :   <---- t0 OS9$28 <F$SRqMem> {size=200} #14833510 -> size $200 addr $7500 #14844400
:   {proc=3} t0: OS9$82 <I$Dup> {$2} #14848732
:   :   {proc=3} t0: OS9$2f <F$Find64> {base=8200 id=2} #14849180
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #14849180 #14849850
:   :   <---- t0 OS9$82 <I$Dup> {$2} #14848732 -> path $2 #14850126
:   {proc=3} t0: OS9$82 <I$Dup> {$1} #14850190
:   :   {proc=3} t0: OS9$2f <F$Find64> {base=8200 id=1} #14850638
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #14850638 #14851308
:   :   <---- t0 OS9$82 <I$Dup> {$1} #14850190 -> path $1 #14853218
:   {proc=3} t0: OS9$82 <I$Dup> {$1} #14853282
:   :   {proc=3} t0: OS9$2f <F$Find64> {base=8200 id=1} #14853730
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #14853730 #14854400
:   :   <---- t0 OS9$82 <I$Dup> {$1} #14853282 -> path $1 #14854676
:   {proc=4} t0: OS9$34 <F$SLink> {"Shell" type 11 name@ e00d dat@ 7840} #14854866
:   :   <---- t0 OS9$34 <F$SLink> {"Shell" type 11 name@ e00d dat@ 7840} #14854866 #14875356
:   {proc=4} t0: OS9$48 <F$LDDDXY> {} #14875522
:   :   <---- t0 OS9$48 <F$LDDDXY> {} #14875522 #14876648
:   {proc=4} t0: OS9$07 <F$Mem> {desired_size=1f00} #14876676
:   :   {proc=4} t0: OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7500} #14877194
:   :   :   <---- t0 OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7500} #14877194 #14878982
:   :   <---- t0 OS9$07 <F$Mem> {desired_size=1f00} #14876676 #14879288
:   {proc=3} t0: OS9$3f <F$AllTsk> {processDesc=7500} #14879612
:   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7500} #14879612 #14880506
:   {proc=3} t0: OS9$38 <F$Move> {srcTask=2 destTask=3 srcPtr=166d destPtr=1ef2 size=000e} #14880616
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=3 srcPtr=166d destPtr=1ef2 size=000e} #14880616 #14884010
:   {proc=3} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=76f4 destPtr=1ee6 size=000c} #14884076
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=76f4 destPtr=1ee6 size=000c} #14884076 #14885856
:   {proc=3} t0: OS9$40 <F$DelTsk> {proc_desc=7500} #14885886
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7500} #14885886 #14886556
:   {proc=3} t0: OS9$2c <F$AProc> {proc=7500} #14886662
:   :   <---- t0 OS9$2c <F$AProc> {proc=7500} #14886662 #14887386
:   <---- t1 OS9$03 <F$Fork> {Module/file='Shell' param="-P{32}X{32}PATH=;-p{13}" lang/type=11 pages=1f} #14832254 #14888870
{proc=3"Shell"} t1: OS9$1d <F$UnLoad> {} #14888944
:   {proc=3} t0: OS9$4e <F$FModul> {"Shell" type 11 name@ e00d dat@ 7840} #14890060
:   :   <---- t0 OS9$4e <F$FModul> {"Shell" type 11 name@ e00d dat@ 7840} #14890060 #14908360
:   <---- t1 OS9$1d <F$UnLoad> {} #14888944 #14909610
{proc=3"Shell"} t1: OS9$8f <I$Close> {path=0} #14909852
:   {proc=3} t0: OS9$2f <F$Find64> {base=8200 id=2} #14912798
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #14912798 #14913468
:   <---- t1 OS9$8f <I$Close> {path=0} #14909852 #14915110
{proc=3"Shell"} t1: OS9$82 <I$Dup> {$3} #14915132
:   {proc=3} t0: OS9$2f <F$Find64> {base=8200 id=1} #14916432
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #14916432 #14917102
:   <---- t1 OS9$82 <I$Dup> {$3} #14915132 -> path $0 #14918358
{proc=3"Shell"} t1: OS9$8f <I$Close> {path=3} #14918408
:   {proc=3} t0: OS9$2f <F$Find64> {base=8200 id=1} #14919720
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #14919720 #14920390
:   <---- t1 OS9$8f <I$Close> {path=3} #14918408 #14922310
{proc=3"Shell"} t1: OS9$04 <F$Wait> {} #14922642
:   {proc=3} t0: OS9$40 <F$DelTsk> {proc_desc=7800} #14924076
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7800} #14924076 #14924746
:   {proc=3} t0: OS9$2d <F$NProc> {} #14924800
{proc=4"Shell"} t1: OS9$09 <F$Icpt> {routine=e06b storage=0000} #14933190
:   <---- t1 OS9$09 <F$Icpt> {routine=e06b storage=0000} #14933190 #14935468
{proc=4"Shell"} t1: OS9$0c <F$ID> {} #14935520
:   <---- t1 OS9$0c <F$ID> {} #14935520 #14937792
{proc=4"Shell"} t1: OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #14940412
:   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=1} #14943322
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #14943322 #14943992
:   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=d511 destPtr=0214 size=0020} #14945266
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=d511 destPtr=0214 size=0020} #14945266 #14947252
:   <---- t1 OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #14940412 #14948446
{proc=4"Shell"} t1: OS9$84 <I$Open> {f7f4='.'} #14949188
:   {proc=4} t0: OS9$30 <F$All64> {table=8200} #14950620
:   :   {proc=4} t0: OS9$28 <F$SRqMem> {size=100} #14956782
:   :   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #14956782 -> size $100 addr $7400 #14966962
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #14950620 -> base $8200 blocknum $4 addr $7400 #14977212
:   {proc=4} t0: OS9$49 <F$LDABX> {} #14977286
:   :   <---- t0 OS9$49 <F$LDABX> {} #14977286 #14978140
:   {proc=1} t0: OS9$10 <F$PrsNam> {path='DD'} #14978326
:   :   <---- t0 OS9$10 <F$PrsNam> {path='DD'} #14978326 #14980094
:   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #14980136
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #14981352
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #14981352 #15009092
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #15009224
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #15009224 -> addr $e704 entry $e72e #15041144
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #15041230
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #15041230 -> addr $9e77 entry $9e88 #15099156
:   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #14980136 #15101402
:   {proc=4} t0: OS9$28 <F$SRqMem> {size=100} #15101898
:   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #15101898 -> size $100 addr $7300 #15112146
:   {proc=4} t0: OS9$30 <F$All64> {table=8200} #15112194
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #15112194 -> base $8200 blocknum $5 addr $7440 #15114796
:   {proc=4} t0: OS9$49 <F$LDABX> {} #15115734
:   :   <---- t0 OS9$49 <F$LDABX> {} #15115734 #15116588
:   {proc=4} t0: OS9$10 <F$PrsNam> {path='a'} #15220222
:   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path='a'} #15220222 #15221984
:   {proc=4} t0: OS9$49 <F$LDABX> {} #15222136
:   :   <---- t0 OS9$49 <F$LDABX> {} #15222136 #15222990
:   {proc=4} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=f7f4 destPtr=7460 size=0001} #15223314
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=f7f4 destPtr=7460 size=0001} #15223314 #15224926
:   {proc=4} t0: OS9$10 <F$PrsNam> {path='"'} #15330728
:   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path='"'} #15330728 #15332250
:   <---- t1 OS9$84 <I$Open> {f7f4='.'} #14949188 -> path $3 #15334764
{proc=4"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #15334858
:   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=4} #15336120
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #15336120 #15336790
:   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7300 destPtr=03f8 size=0020} #15386574
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7300 destPtr=03f8 size=0020} #15386574 #15390194
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #15334858 #15392058
{proc=4"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #15392178
:   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=4} #15393440
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #15393440 #15394110
:   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7320 destPtr=03f8 size=0020} #15395884
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7320 destPtr=03f8 size=0020} #15395884 #15397870
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #15392178 #15399734
{proc=4"Shell"} t1: OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #15400002
:   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=4} #15401330
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #15401330 #15402000
:   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=ec52 destPtr=00b5 size=0020} #15403250
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=ec52 destPtr=00b5 size=0020} #15403250 #15405236
:   <---- t1 OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #15400002 #15406430
{proc=4"Shell"} t1: OS9$10 <F$PrsNam> {path='DD'} #15406456
:   <---- t1 OS9$10 <F$PrsNam> {path='DD'} #15406456 #15409598
{proc=4"Shell"} t1: OS9$8f <I$Close> {path=3} #15409834
:   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=4} #15411146
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #15411146 #15411816
:   {proc=4} t0: OS9$29 <F$SRtMem> {size=100 start=7300} #15412278
:   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7300} #15412278 #15416822
:   {proc=4} t0: OS9$31 <F$Ret64> {block_num=5 address=8200} #15417288
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=5 address=8200} #15417288 #15418066
:   {proc=4} t0: OS9$81 <I$Detach> {8300} #15419970
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #15420532
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #15420532 #15421272
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #15421292
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #15421292 #15422032
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #15422052
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #15422052 #15422792
:   :   <---- t0 OS9$81 <I$Detach> {8300} #15419970 #15423142
:   {proc=4} t0: OS9$31 <F$Ret64> {block_num=4 address=8200} #15423172
:   :   {proc=4} t0: OS9$29 <F$SRtMem> {size=100 start=7400} #15423808
:   :   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7400} #15423808 #15426260
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=4 address=8200} #15423172 #15426558
:   <---- t1 OS9$8f <I$Close> {path=3} #15409834 #15427734
{proc=4"Shell"} t1: OS9$1c <F$SUser> {} #15430386
:   <---- t1 OS9$1c <F$SUser> {} #15430386 #15432640
{proc=4"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #15461340
:   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #15462586
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #15462586 #15463256
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #15461340 #15465484
{proc=4"Shell"} t1: OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #15465526
:   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #15466854
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #15466854 #15467524
:   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=82a0 destPtr=0124 size=0020} #15468032
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=82a0 destPtr=0124 size=0020} #15468032 #15470018
:   <---- t1 OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #15465526 #15471212
{proc=4"Shell"} t1: OS9$1c <F$SUser> {} #15471490
:   <---- t1 OS9$1c <F$SUser> {} #15471490 #15473744
{proc=4"Shell"} t1: OS9$8b <I$ReadLn> {} #15473886
:   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #15475148
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #15475148 #15475818
:   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7700 destPtr=0124 size=0001} #15525682
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7700 destPtr=0124 size=0001} #15525682 #15527206
:   {proc=4} t0: OS9$49 <F$LDABX> {} #15527502
:   :   <---- t0 OS9$49 <F$LDABX> {} #15527502 #15528268
:   <---- t1 OS9$8b <I$ReadLn> {} #15473886 #15530026
{proc=4"Shell"} t1: OS9$8b <I$ReadLn> {} #15530422
:   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #15531684
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #15531684 #15532354
:   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7701 destPtr=0124 size=0008} #15534462
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7701 destPtr=0124 size=0008} #15534462 #15536106
:   {proc=4} t0: OS9$49 <F$LDABX> {} #15536402
:   :   <---- t0 OS9$49 <F$LDABX> {} #15536402 #15537168
:   <---- t1 OS9$8b <I$ReadLn> {} #15530422 #15540560
{proc=4"Shell"} t1: OS9$1c <F$SUser> {} #15544186
:   <---- t1 OS9$1c <F$SUser> {} #15544186 #15546440
{proc=4"Shell"} t1: OS9$10 <F$PrsNam> {path='dir'} #15555294
:   <---- t1 OS9$10 <F$PrsNam> {path='dir'} #15555294 #15558906
{proc=4"Shell"} t1: OS9$10 <F$PrsNam> {path=''} #15558966
:   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL1 OS9$10 <F$PrsNam> {path=''} #15558966 #15562042
{proc=4"Shell"} t1: OS9$21 <F$NMLink> {LangType=00, 0e6d='dir'} #15566544
:   {proc=4} t0: OS9$4e <F$FModul> {"dir" type 0 name@ e6d dat@ 7540} #15569294
:   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$4e <F$FModul> {"dir" type 0 name@ e6d dat@ 7540} #15569294 #15621214
:   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL1 OS9$21 <F$NMLink> {LangType=00, 0e6d='dir'} #15566544 #15622406
{proc=4"Shell"} t1: OS9$84 <I$Open> {0e6d='dir'} #15622434
:   {proc=4} t0: OS9$30 <F$All64> {table=8200} #15623866
:   :   {proc=4} t0: OS9$28 <F$SRqMem> {size=100} #15631662
:   :   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #15631662 -> size $100 addr $7400 #15642324
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #15623866 -> base $8200 blocknum $4 addr $7400 #15650940
:   {proc=4} t0: OS9$49 <F$LDABX> {} #15651014
:   :   <---- t0 OS9$49 <F$LDABX> {} #15651014 #15651780
:   {proc=1} t0: OS9$10 <F$PrsNam> {path='DD'} #15651972
:   :   <---- t0 OS9$10 <F$PrsNam> {path='DD'} #15651972 #15653740
:   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #15653782
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #15654998
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #15654998 #15682738
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #15682870
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #15682870 -> addr $e704 entry $e72e #15714748
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #15714834
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #15714834 -> addr $9e77 entry $9e88 #15772760
:   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #15653782 #15775006
:   {proc=4} t0: OS9$28 <F$SRqMem> {size=100} #15775502
:   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #15775502 -> size $100 addr $7300 #15788136
:   {proc=4} t0: OS9$30 <F$All64> {table=8200} #15788184
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #15788184 -> base $8200 blocknum $5 addr $7440 #15790786
:   {proc=4} t0: OS9$49 <F$LDABX> {} #15791724
:   :   <---- t0 OS9$49 <F$LDABX> {} #15791724 #15792490
:   {proc=4} t0: OS9$10 <F$PrsNam> {path=''} #15885526
:   :   <---- t0 OS9$10 <F$PrsNam> {path=''} #15885526 #15887484
:   {proc=4} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0e6d destPtr=7460 size=0003} #15887738
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0e6d destPtr=7460 size=0003} #15887738 #15889330
:   {proc=4} t0: OS9$10 <F$PrsNam> {path=''} #16106296
:   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path=''} #16106296 #16107718
:   <---- t1 OS9$84 <I$Open> {0e6d='dir'} #15622434 -> path $3 #16110232
{proc=4"Shell"} t1: OS9$89 <I$Read> {path=3 buf=00d6 size=4d} #16110274
:   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=4} #16111536
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #16111536 #16112206
:   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7300 destPtr=00d6 size=004d} #16160420
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7300 destPtr=00d6 size=004d} #16160420 #16163146
:   <---- t1 OS9$89 <I$Read> {path=3 buf=00d6 size=4d} #16110274 #16166644
{proc=4"Shell"} t1: OS9$10 <F$PrsNam> {path='Dir'} #16166812
:   <---- t1 OS9$10 <F$PrsNam> {path='Dir'} #16166812 #16170164
{proc=4"Shell"} t1: OS9$8f <I$Close> {path=3} #16170222
:   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=4} #16171534
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #16171534 #16172204
:   {proc=4} t0: OS9$29 <F$SRtMem> {size=100 start=7300} #16172666
:   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7300} #16172666 #16177210
:   {proc=4} t0: OS9$31 <F$Ret64> {block_num=5 address=8200} #16177728
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=5 address=8200} #16177728 #16178506
:   {proc=4} t0: OS9$81 <I$Detach> {8300} #16178776
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #16179338
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #16179338 #16180078
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #16180098
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #16180098 #16180838
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #16180858
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #16180858 #16181598
:   :   <---- t0 OS9$81 <I$Detach> {8300} #16178776 #16181948
:   {proc=4} t0: OS9$31 <F$Ret64> {block_num=4 address=8200} #16181978
:   :   {proc=4} t0: OS9$29 <F$SRtMem> {size=100 start=7400} #16182614
:   :   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7400} #16182614 #16185066
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=4 address=8200} #16181978 #16185364
:   <---- t1 OS9$8f <I$Close> {path=3} #16170222 #16186920
{proc=4"Shell"} t1: OS9$21 <F$NMLink> {LangType=11, 0e6d='dir'} #16187084
:   {proc=4} t0: OS9$4e <F$FModul> {"dir" type 11 name@ e6d dat@ 7540} #16188200
:   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$4e <F$FModul> {"dir" type 11 name@ e6d dat@ 7540} #16188200 #16241702
:   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL1 OS9$21 <F$NMLink> {LangType=11, 0e6d='dir'} #16187084 #16242894
{proc=4"Shell"} t1: OS9$22 <F$NMLoad> {LangType=11, 0e6d='dir'} #16242912
:   {proc=4} t0: OS9$4b <F$AllPrc> {} #16244016
:   :   {proc=4} t0: OS9$28 <F$SRqMem> {size=200} #16244560
:   :   :   <---- t0 OS9$28 <F$SRqMem> {size=200} #16244560 -> size $200 addr $7300 #16256950
:   :   <---- t0 OS9$4b <F$AllPrc> {} #16244016 #16260878
:   {proc=4} t0: OS9$84 <I$Open> {0e6d=''} #16261052
:   :   {proc=4} t0: OS9$30 <F$All64> {table=8200} #16261536
:   :   :   {proc=4} t0: OS9$28 <F$SRqMem> {size=100} #16267698
:   :   :   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #16267698 -> size $100 addr $7200 #16278766
:   :   :   <---- t0 OS9$30 <F$All64> {table=8200} #16261536 -> base $8200 blocknum $4 addr $7200 #16289016
:   :   {proc=4} t0: OS9$49 <F$LDABX> {} #16289090
:   :   :   <---- t0 OS9$49 <F$LDABX> {} #16289090 #16289856
:   :   {proc=1} t0: OS9$10 <F$PrsNam> {path='DD'} #16290048
:   :   :   <---- t0 OS9$10 <F$PrsNam> {path='DD'} #16290048 #16291816
:   :   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #16291858
:   :   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #16293074
:   :   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #16293074 #16320814
:   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #16320946
:   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #16320946 -> addr $e704 entry $e72e #16352824
:   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #16352910
:   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #16352910 -> addr $9e77 entry $9e88 #16410836
:   :   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #16291858 #16413082
:   :   {proc=4} t0: OS9$28 <F$SRqMem> {size=100} #16413578
:   :   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #16413578 -> size $100 addr $7100 #16423962
:   :   {proc=4} t0: OS9$30 <F$All64> {table=8200} #16424010
:   :   :   <---- t0 OS9$30 <F$All64> {table=8200} #16424010 -> base $8200 blocknum $5 addr $7240 #16426612
:   :   {proc=4} t0: OS9$49 <F$LDABX> {} #16427550
:   :   :   <---- t0 OS9$49 <F$LDABX> {} #16427550 #16428316
:   :   {proc=4} t0: OS9$10 <F$PrsNam> {path=''} #16521352
:   :   :   <---- t0 OS9$10 <F$PrsNam> {path=''} #16521352 #16524944
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0e6d destPtr=7260 size=0003} #16525198
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0e6d destPtr=7260 size=0003} #16525198 #16526790
:   :   {proc=4} t0: OS9$10 <F$PrsNam> {path=''} #16740680
:   :   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path=''} #16740680 #16742102
:   :   <---- t0 OS9$84 <I$Open> {0e6d=''} #16261052 -> path $4 #16743306
:   {proc=4} t0: OS9$3f <F$AllTsk> {processDesc=7300} #16743364
:   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7300} #16743364 #16744258
:   {proc=5} t0: OS9$41 <F$SetTsk> {} #16745400
:   :   <---- t0 OS9$41 <F$SetTsk> {} #16745400 #16746114
:   {proc=5} t0: OS9$89 <I$Read> {path=4 buf=0000 size=9} #16746166
:   :   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #16746626
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #16746626 #16747296
:   :   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7100 destPtr=0000 size=0009} #16793940
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7100 destPtr=0000 size=0009} #16793940 #16795618
:   :   <---- t0 OS9$89 <I$Read> {path=4 buf=0000 size=9} #16746166 #16796554
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=0000 destPtr=76d9 size=0009} #16796666
:   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=0000 destPtr=76d9 size=0009} #16796666 #16798344
:   {proc=5} t0: OS9$89 <I$Read> {path=4 buf=0009 size=398} #16798538
:   :   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #16798998
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #16798998 #16799668
:   :   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7109 destPtr=0009 size=00f7} #16801430
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7109 destPtr=0009 size=00f7} #16801430 #16807530
:   :   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7100 destPtr=0100 size=0100} #16853062
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7100 destPtr=0100 size=0100} #16853062 #16859152
:   :   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7100 destPtr=0200 size=0100} #16906636
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7100 destPtr=0200 size=0100} #16906636 #16914360
:   :   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7100 destPtr=0300 size=00a1} #16959866
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7100 destPtr=0300 size=00a1} #16959866 #16964166
:   :   <---- t0 OS9$89 <I$Read> {path=4 buf=0009 size=398} #16798538 #16965102
:   {proc=5} t0: OS9$2e <F$VModul> {addr=0000="-" map=[20 5 20 1 d 34 17 3c]} #16965240
:   :   <---- t0 OS9$2e <F$VModul> {addr=0000="-" map=[20 5 20 1 d 34 17 3c]} #16965240 #17024656
:   {proc=5} t0: OS9$89 <I$Read> {path=4 buf=03a1 size=9} #17025018
:   :   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #17025478
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #17025478 #17026148
:   :   <-- ERROR: $d3(E$EOF    :End of File): OS9KERNEL0 OS9$89 <I$Read> {path=4 buf=03a1 size=9} #17025018 #17028574
:   {proc=4} t0: OS9$8f <I$Close> {path=4} #17028674
:   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=4} #17029122
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #17029122 #17031426
:   :   {proc=4} t0: OS9$29 <F$SRtMem> {size=100 start=7100} #17031888
:   :   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7100} #17031888 #17036376
:   :   {proc=4} t0: OS9$31 <F$Ret64> {block_num=5 address=8200} #17036894
:   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=5 address=8200} #17036894 #17037672
:   :   {proc=4} t0: OS9$81 <I$Detach> {8300} #17037942
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #17038504
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #17038504 #17039244
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #17039264
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #17039264 #17040004
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #17040024
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #17040024 #17040764
:   :   :   <---- t0 OS9$81 <I$Detach> {8300} #17037942 #17041114
:   :   {proc=4} t0: OS9$31 <F$Ret64> {block_num=4 address=8200} #17041144
:   :   :   {proc=4} t0: OS9$29 <F$SRtMem> {size=100 start=7200} #17041780
:   :   :   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7200} #17041780 #17044176
:   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=4 address=8200} #17041144 #17044474
:   :   <---- t0 OS9$8f <I$Close> {path=4} #17028674 #17044722
:   {proc=4} t0: OS9$4c <F$DelPrc> {} #17044854
:   :   {proc=4} t0: OS9$40 <F$DelTsk> {proc_desc=7300} #17045306
:   :   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7300} #17045306 #17045976
:   :   {proc=4} t0: OS9$29 <F$SRtMem> {size=200 start=7300} #17046002
:   :   :   <---- t0 OS9$29 <F$SRtMem> {size=200 start=7300} #17046002 #17048488
:   :   <---- t0 OS9$4c <F$DelPrc> {} #17044854 #17048744
:   {proc=4} t0: OS9$48 <F$LDDDXY> {} #17051390
:   :   <---- t0 OS9$48 <F$LDDDXY> {} #17051390 #17052236
:   {proc=4} t0: OS9$48 <F$LDDDXY> {} #17052364
:   :   <---- t0 OS9$48 <F$LDDDXY> {} #17052364 #17053210
:   <---- t1 OS9$22 <F$NMLoad> {LangType=11, 0e6d='dir'} #16242912 #17054712
{proc=4"Shell"} t1: OS9$03 <F$Fork> {Module/file='dir' param="/dd{13}" lang/type=11 pages=1f} #17054920
:   {proc=4} t0: OS9$28 <F$SRqMem> {size=200} #17056194
:   :   <---- t0 OS9$28 <F$SRqMem> {size=200} #17056194 -> size $200 addr $7300 #17068584
:   {proc=4} t0: OS9$82 <I$Dup> {$2} #17072916
:   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #17073364
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #17073364 #17074034
:   :   <---- t0 OS9$82 <I$Dup> {$2} #17072916 -> path $2 #17074310
:   {proc=4} t0: OS9$82 <I$Dup> {$1} #17074374
:   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=1} #17074822
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #17074822 #17075492
:   :   <---- t0 OS9$82 <I$Dup> {$1} #17074374 -> path $1 #17075768
:   {proc=4} t0: OS9$82 <I$Dup> {$1} #17075832
:   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=1} #17076280
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #17076280 #17076950
:   :   <---- t0 OS9$82 <I$Dup> {$1} #17075832 -> path $1 #17077226
:   {proc=5} t0: OS9$34 <F$SLink> {"dir" type 11 name@ e6d dat@ 7540} #17077416
:   :   <---- t0 OS9$34 <F$SLink> {"dir" type 11 name@ e6d dat@ 7540} #17077416 #17084356
:   {proc=5} t0: OS9$48 <F$LDDDXY> {} #17084522
:   :   <---- t0 OS9$48 <F$LDDDXY> {} #17084522 #17085648
:   {proc=5} t0: OS9$07 <F$Mem> {desired_size=1f00} #17085690
:   :   {proc=5} t0: OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7300} #17086208
:   :   :   <---- t0 OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7300} #17086208 #17087930
:   :   <---- t0 OS9$07 <F$Mem> {desired_size=1f00} #17085690 #17088236
:   {proc=4} t0: OS9$3f <F$AllTsk> {processDesc=7300} #17088560
:   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7300} #17088560 #17091088
:   {proc=4} t0: OS9$38 <F$Move> {srcTask=2 destTask=3 srcPtr=0e71 destPtr=1efc size=0004} #17091198
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=3 srcPtr=0e71 destPtr=1efc size=0004} #17091198 #17092736
:   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=74f4 destPtr=1ef0 size=000c} #17092802
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=74f4 destPtr=1ef0 size=000c} #17092802 #17094582
:   {proc=4} t0: OS9$40 <F$DelTsk> {proc_desc=7300} #17094612
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7300} #17094612 #17095282
:   {proc=4} t0: OS9$2c <F$AProc> {proc=7300} #17095388
:   :   <---- t0 OS9$2c <F$AProc> {proc=7300} #17095388 #17096112
:   <---- t1 OS9$03 <F$Fork> {Module/file='dir' param="/dd{13}" lang/type=11 pages=1f} #17054920 #17097596
{proc=4"Shell"} t1: OS9$1d <F$UnLoad> {} #17097660
:   {proc=4} t0: OS9$4e <F$FModul> {"Dir" type 11 name@ e3 dat@ 7540} #17098776
:   :   <---- t0 OS9$4e <F$FModul> {"Dir" type 11 name@ e3 dat@ 7540} #17098776 #17103266
:   <---- t1 OS9$1d <F$UnLoad> {} #17097660 #17104516
{proc=4"Shell"} t1: OS9$04 <F$Wait> {} #17105078
:   {proc=4} t0: OS9$40 <F$DelTsk> {proc_desc=7500} #17106512
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7500} #17106512 #17107182
:   {proc=4} t0: OS9$2d <F$NProc> {} #17107236
{proc=5"Dir"} t1: OS9$8d <I$GetStt> {path=1 26==SS.ScSiz  : Return screen size for COCO} #17108792
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=1} #17110120
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #17110120 #17110790
:   <---- t1 OS9$8d <I$GetStt> {path=1 26==SS.ScSiz  : Return screen size for COCO} #17108792 #17113298
{proc=5"Dir"} t1: OS9$84 <I$Open> {1efc='/dd'} #17113842
:   {proc=5} t0: OS9$30 <F$All64> {table=8200} #17115274
:   :   {proc=5} t0: OS9$28 <F$SRqMem> {size=100} #17123018
:   :   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #17123018 -> size $100 addr $7200 #17134284
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #17115274 -> base $8200 blocknum $4 addr $7200 #17142900
:   {proc=5} t0: OS9$49 <F$LDABX> {} #17142974
:   :   <---- t0 OS9$49 <F$LDABX> {} #17142974 #17143740
:   {proc=5} t0: OS9$10 <F$PrsNam> {path=''} #17143818
:   :   <---- t0 OS9$10 <F$PrsNam> {path=''} #17143818 #17145696
:   {proc=5} t0: OS9$80 <I$Attach> {1efd=''} #17145738
:   :   {proc=1} t0: OS9$34 <F$SLink> {"dd" type f0 name@ 1efd dat@ 7340} #17146954
:   :   :   <---- t0 OS9$34 <F$SLink> {"dd" type f0 name@ 1efd dat@ 7340} #17146954 #17172016
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #17172148
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #17172148 -> addr $e704 entry $e72e #17205432
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #17205518
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #17205518 -> addr $9e77 entry $9e88 #17264850
:   :   <---- t0 OS9$80 <I$Attach> {1efd=''} #17145738 #17267096
:   {proc=5} t0: OS9$28 <F$SRqMem> {size=100} #17267592
:   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #17267592 -> size $100 addr $7100 #17279610
:   {proc=5} t0: OS9$30 <F$All64> {table=8200} #17279658
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #17279658 -> base $8200 blocknum $5 addr $7240 #17282260
:   {proc=5} t0: OS9$49 <F$LDABX> {} #17283198
:   :   <---- t0 OS9$49 <F$LDABX> {} #17283198 #17283964
:   {proc=5} t0: OS9$10 <F$PrsNam> {path=''} #17284036
:   :   <---- t0 OS9$10 <F$PrsNam> {path=''} #17284036 #17285914
:   {proc=5} t0: OS9$10 <F$PrsNam> {path=''} #17385784
:   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path=''} #17385784 #17387026
:   <---- t1 OS9$84 <I$Open> {1efc='/dd'} #17113842 -> path $3 #17391174
{proc=5"Dir"} t1: OS9$86 <I$ChgDir> {mode=81, 1efc='/dd'} #17391222
:   {proc=5} t0: OS9$30 <F$All64> {table=8200} #17392442
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #17392442 -> base $8200 blocknum $6 addr $7280 #17395080
:   {proc=5} t0: OS9$49 <F$LDABX> {} #17395154
:   :   <---- t0 OS9$49 <F$LDABX> {} #17395154 #17395920
:   {proc=5} t0: OS9$10 <F$PrsNam> {path=''} #17395998
:   :   <---- t0 OS9$10 <F$PrsNam> {path=''} #17395998 #17397876
:   {proc=5} t0: OS9$80 <I$Attach> {1efd=''} #17397918
:   :   {proc=1} t0: OS9$34 <F$SLink> {"dd" type f0 name@ 1efd dat@ 7340} #17399134
:   :   :   <---- t0 OS9$34 <F$SLink> {"dd" type f0 name@ 1efd dat@ 7340} #17399134 #17424154
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #17424286
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #17424286 -> addr $e704 entry $e72e #17457570
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #17457656
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #17457656 -> addr $9e77 entry $9e88 #17516988
:   :   <---- t0 OS9$80 <I$Attach> {1efd=''} #17397918 #17519234
:   {proc=5} t0: OS9$28 <F$SRqMem> {size=100} #17519786
:   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #17519786 -> size $100 addr $7000 #17530238
:   {proc=5} t0: OS9$30 <F$All64> {table=8200} #17530286
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #17530286 -> base $8200 blocknum $7 addr $72c0 #17532960
:   {proc=5} t0: OS9$49 <F$LDABX> {} #17533898
:   :   <---- t0 OS9$49 <F$LDABX> {} #17533898 #17534664
:   {proc=5} t0: OS9$10 <F$PrsNam> {path=''} #17534736
:   :   <---- t0 OS9$10 <F$PrsNam> {path=''} #17534736 #17538248
:   {proc=5} t0: OS9$10 <F$PrsNam> {path=''} #17635062
:   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path=''} #17635062 #17636304
:   {proc=5} t0: OS9$29 <F$SRtMem> {size=100 start=7000} #17637282
:   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7000} #17637282 #17641742
:   {proc=5} t0: OS9$31 <F$Ret64> {block_num=7 address=8200} #17642284
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=7 address=8200} #17642284 #17643062
:   {proc=5} t0: OS9$81 <I$Detach> {8300} #17643402
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #17643964
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #17643964 #17644704
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #17644724
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #17644724 #17645464
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #17645484
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #17645484 #17646224
:   :   <---- t0 OS9$81 <I$Detach> {8300} #17643402 #17646574
:   {proc=5} t0: OS9$31 <F$Ret64> {block_num=6 address=8200} #17646604
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=6 address=8200} #17646604 #17647382
:   <---- t1 OS9$86 <I$ChgDir> {mode=81, 1efc='/dd'} #17391222 #17648874
{proc=5"Dir"} t1: OS9$15 <F$Time> {buf=d} #17651396
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=0028 destPtr=000d size=0006} #17652512
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=0028 destPtr=000d size=0006} #17652512 #17654118
:   <---- t1 OS9$15 <F$Time> {buf=d} #17651396 #17655286
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #17658852
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=1} #17660120
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #17660120 #17660790
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0040 destPtr=7ce0 size=0020} #17661932
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0040 destPtr=7ce0 size=0020} #17661932 #17663918
:   <---- t1 OS9$8c <I$WritLn> {} #17658852 #17688664
{proc=5"Dir"} t1: OS9$88 <I$Seek> {path=3 pos=00000040} #17688754
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #17690000
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #17690000 #17690670
:   <---- t1 OS9$88 <I$Seek> {path=3 pos=00000040} #17688754 #17692280
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17692342
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #17693604
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #17693604 #17694274
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7140 destPtr=0013 size=0020} #17748708
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7140 destPtr=0013 size=0020} #17748708 #17750694
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17692342 #17752874
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17754374
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #17755636
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #17755636 #17756306
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7160 destPtr=0013 size=0020} #17758080
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7160 destPtr=0013 size=0020} #17758080 #17760066
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17754374 #17761930
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17763514
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #17764776
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #17764776 #17765446
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7180 destPtr=0013 size=0020} #17767220
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7180 destPtr=0013 size=0020} #17767220 #17769206
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17763514 #17771070
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17772582
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #17773844
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #17773844 #17774514
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=71a0 destPtr=0013 size=0020} #17777922
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=71a0 destPtr=0013 size=0020} #17777922 #17779908
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17772582 #17781772
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #17782536
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=1} #17783804
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #17783804 #17784474
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0040 destPtr=7ce0 size=0020} #17785616
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0040 destPtr=7ce0 size=0020} #17785616 #17787602
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0060 destPtr=7ce0 size=0020} #17799990
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0060 destPtr=7ce0 size=0020} #17799990 #17801976
:   <---- t1 OS9$8c <I$WritLn> {} #17782536 #17814172
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17814242
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #17815504
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #17815504 #17816174
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=71c0 destPtr=0013 size=0020} #17817948
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=71c0 destPtr=0013 size=0020} #17817948 #17819934
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17814242 #17821798
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17823250
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #17824512
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #17824512 #17825182
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=71e0 destPtr=0013 size=0020} #17826944
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=71e0 destPtr=0013 size=0020} #17826944 #17828930
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17823250 #17830946
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17832416
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #17833678
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #17833678 #17834348
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7100 destPtr=0013 size=0020} #17889772
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7100 destPtr=0013 size=0020} #17889772 #17891758
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17832416 #17893938
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17897902
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #17899164
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #17899164 #17899834
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7120 destPtr=0013 size=0020} #17901608
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7120 destPtr=0013 size=0020} #17901608 #17903594
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17897902 #17905458
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #17906466
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=1} #17907734
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #17907734 #17908404
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0040 destPtr=7ce0 size=0020} #17909546
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0040 destPtr=7ce0 size=0020} #17909546 #17911532
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0060 destPtr=7ce0 size=0020} #17923744
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0060 destPtr=7ce0 size=0020} #17923744 #17927332
:   <---- t1 OS9$8c <I$WritLn> {} #17906466 #17938386
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17938456
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #17939718
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #17939718 #17940388
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7140 destPtr=0013 size=0020} #17942162
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7140 destPtr=0013 size=0020} #17942162 #17944148
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17938456 #17946012
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #17947440
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #17948702
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #17948702 #17949372
:   <-- ERROR: $d3(E$EOF    :End of File): OS9KERNEL1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #17947440 #17952726
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #17952932
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=1} #17955854
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #17955854 #17956524
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0040 destPtr=7ce0 size=0020} #17957666
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0040 destPtr=7ce0 size=0020} #17957666 #17959652
:   <---- t1 OS9$8c <I$WritLn> {} #17952932 #17973020
{proc=5"Dir"} t1: OS9$06 <F$Exit> {status=0} #17973058
:   {proc=5} t0: OS9$8f <I$Close> {path=2} #17974214
:   :   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=2} #17974662
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #17974662 #17975332
:   :   <---- t0 OS9$8f <I$Close> {path=2} #17974214 #17976046
:   {proc=5} t0: OS9$8f <I$Close> {path=1} #17976124
:   :   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=1} #17976572
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #17976572 #17977242
:   :   <---- t0 OS9$8f <I$Close> {path=1} #17976124 #17978434
:   {proc=5} t0: OS9$8f <I$Close> {path=1} #17978512
:   :   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=1} #17978960
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #17978960 #17979630
:   :   {proc=5} t0: OS9$37 <F$GProcP> {id=04} #17980828
:   :   :   <---- t0 OS9$37 <F$GProcP> {id=04} #17980828 #17981536
:   :   <---- t0 OS9$8f <I$Close> {path=1} #17978512 #17982274
:   {proc=5} t0: OS9$8f <I$Close> {path=4} #17982352
:   :   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #17982800
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #17982800 #17983470
:   :   {proc=5} t0: OS9$29 <F$SRtMem> {size=100 start=7100} #17983932
:   :   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7100} #17983932 #17987954
:   :   {proc=5} t0: OS9$31 <F$Ret64> {block_num=5 address=8200} #17988420
:   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=5 address=8200} #17988420 #17989198
:   :   {proc=5} t0: OS9$81 <I$Detach> {8300} #17989468
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #17990030
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #17990030 #17990770
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #17990790
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #17990790 #17991530
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #17991550
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #17991550 #17992290
:   :   :   <---- t0 OS9$81 <I$Detach> {8300} #17989468 #17992640
:   :   {proc=5} t0: OS9$31 <F$Ret64> {block_num=4 address=8200} #17992670
:   :   :   {proc=5} t0: OS9$29 <F$SRtMem> {size=100 start=7200} #17993306
:   :   :   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7200} #17993306 #17995702
:   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=4 address=8200} #17992670 #17996000
:   :   <---- t0 OS9$8f <I$Close> {path=4} #17982352 #17996248
:   {proc=5} t0: OS9$3b <F$DelImg> {beginBlock=0 numBlocks=1 processDesc=7300} #17996662
:   :   <---- t0 OS9$3b <F$DelImg> {beginBlock=0 numBlocks=1 processDesc=7300} #17996662 #17997440
:   {proc=5} t0: OS9$02 <F$UnLink> {u=e000 magic=ae48 module='''} #17997498
:   :   {proc=5} t0: OS9$48 <F$LDDDXY> {} #18000578
:   :   :   <---- t0 OS9$48 <F$LDDDXY> {} #18000578 #18001704
:   :   <---- t0 OS9$02 <F$UnLink> {u=e000 magic=ae48 module='''} #17997498 #18007114
:   {proc=5} t0: OS9$40 <F$DelTsk> {proc_desc=7300} #18007154
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7300} #18007154 #18007824
:   {proc=5} t0: OS9$40 <F$DelTsk> {proc_desc=7300} #18008276
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7300} #18008276 #18008870
:   {proc=5} t0: OS9$29 <F$SRtMem> {size=200 start=7300} #18008896
:   :   <---- t0 OS9$29 <F$SRtMem> {size=200 start=7300} #18008896 #18011382
:   {proc=5} t0: OS9$2c <F$AProc> {proc=7500} #18011420
:   :   <---- t0 OS9$2c <F$AProc> {proc=7500} #18011420 #18012144
:   {proc=5} t0: OS9$2d <F$NProc> {} #18012156
:   :   <---- t1 OS9$04 <F$Wait> {} #17105078 #18014642
{proc=4"Shell"} t1: OS9$8b <I$ReadLn> {} #18020766
:   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #18022028
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #18022028 #18022698
:   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7709 destPtr=0124 size=0008} #18024806
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7709 destPtr=0124 size=0008} #18024806 #18026450
:   {proc=4} t0: OS9$49 <F$LDABX> {} #18026746
:   :   <---- t0 OS9$49 <F$LDABX> {} #18026746 #18027512
:   <---- t1 OS9$8b <I$ReadLn> {} #18020766 #18029270
{proc=4"Shell"} t1: OS9$1c <F$SUser> {} #18032896
:   <---- t1 OS9$1c <F$SUser> {} #18032896 #18035150
{proc=4"Shell"} t1: OS9$10 <F$PrsNam> {path='dir'} #18046456
:   <---- t1 OS9$10 <F$PrsNam> {path='dir'} #18046456 #18050068
{proc=4"Shell"} t1: OS9$10 <F$PrsNam> {path=''} #18050128
:   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL1 OS9$10 <F$PrsNam> {path=''} #18050128 #18053204
{proc=4"Shell"} t1: OS9$21 <F$NMLink> {LangType=00, 0e6d='dir'} #18057706
:   {proc=4} t0: OS9$4e <F$FModul> {"dir" type 0 name@ e6d dat@ 7540} #18058822
:   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$4e <F$FModul> {"dir" type 0 name@ e6d dat@ 7540} #18058822 #18112484
:   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL1 OS9$21 <F$NMLink> {LangType=00, 0e6d='dir'} #18057706 #18113676
{proc=4"Shell"} t1: OS9$84 <I$Open> {0e6d='dir'} #18113704
:   {proc=4} t0: OS9$30 <F$All64> {table=8200} #18115136
:   :   {proc=4} t0: OS9$28 <F$SRqMem> {size=100} #18121298
:   :   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #18121298 -> size $100 addr $7400 #18131960
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #18115136 -> base $8200 blocknum $4 addr $7400 #18142210
:   {proc=4} t0: OS9$49 <F$LDABX> {} #18142284
:   :   <---- t0 OS9$49 <F$LDABX> {} #18142284 #18143050
:   {proc=1} t0: OS9$10 <F$PrsNam> {path='DD'} #18143242
:   :   <---- t0 OS9$10 <F$PrsNam> {path='DD'} #18143242 #18145010
:   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #18145052
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #18146268
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #18146268 #18174168
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #18174300
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #18174300 -> addr $e704 entry $e72e #18206338
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #18206424
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #18206424 -> addr $9e77 entry $9e88 #18264510
:   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #18145052 #18266756
:   {proc=4} t0: OS9$28 <F$SRqMem> {size=100} #18267252
:   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #18267252 -> size $100 addr $7300 #18278252
:   {proc=4} t0: OS9$30 <F$All64> {table=8200} #18278300
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #18278300 -> base $8200 blocknum $5 addr $7440 #18280902
:   {proc=4} t0: OS9$49 <F$LDABX> {} #18281840
:   :   <---- t0 OS9$49 <F$LDABX> {} #18281840 #18284240
:   {proc=4} t0: OS9$10 <F$PrsNam> {path=''} #18381986
:   :   <---- t0 OS9$10 <F$PrsNam> {path=''} #18381986 #18383944
:   {proc=4} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0e6d destPtr=7460 size=0003} #18384198
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0e6d destPtr=7460 size=0003} #18384198 #18385790
:   {proc=4} t0: OS9$10 <F$PrsNam> {path=''} #18607402
:   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path=''} #18607402 #18608824
:   <---- t1 OS9$84 <I$Open> {0e6d='dir'} #18113704 -> path $3 #18611338
{proc=4"Shell"} t1: OS9$89 <I$Read> {path=3 buf=00d6 size=4d} #18613832
:   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=4} #18615094
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #18615094 #18615764
:   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7300 destPtr=00d6 size=004d} #18665496
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7300 destPtr=00d6 size=004d} #18665496 #18668222
:   <---- t1 OS9$89 <I$Read> {path=3 buf=00d6 size=4d} #18613832 #18670086
{proc=4"Shell"} t1: OS9$10 <F$PrsNam> {path='Dir'} #18672706
:   <---- t1 OS9$10 <F$PrsNam> {path='Dir'} #18672706 #18676438
{proc=4"Shell"} t1: OS9$8f <I$Close> {path=3} #18676496
:   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=4} #18677808
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #18677808 #18678478
:   {proc=4} t0: OS9$29 <F$SRtMem> {size=100 start=7300} #18678940
:   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7300} #18678940 #18683484
:   {proc=4} t0: OS9$31 <F$Ret64> {block_num=5 address=8200} #18684002
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=5 address=8200} #18684002 #18684780
:   {proc=4} t0: OS9$81 <I$Detach> {8300} #18685050
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #18685612
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #18685612 #18686352
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #18686372
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #18686372 #18687112
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #18687132
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #18687132 #18687872
:   :   <---- t0 OS9$81 <I$Detach> {8300} #18685050 #18688222
:   {proc=4} t0: OS9$31 <F$Ret64> {block_num=4 address=8200} #18688252
:   :   {proc=4} t0: OS9$29 <F$SRtMem> {size=100 start=7400} #18688888
:   :   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7400} #18688888 #18691340
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=4 address=8200} #18688252 #18691638
:   <---- t1 OS9$8f <I$Close> {path=3} #18676496 #18692814
{proc=4"Shell"} t1: OS9$21 <F$NMLink> {LangType=11, 0e6d='dir'} #18692978
:   {proc=4} t0: OS9$4e <F$FModul> {"dir" type 11 name@ e6d dat@ 7540} #18694094
:   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$4e <F$FModul> {"dir" type 11 name@ e6d dat@ 7540} #18694094 #18747756
:   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL1 OS9$21 <F$NMLink> {LangType=11, 0e6d='dir'} #18692978 #18748948
{proc=4"Shell"} t1: OS9$22 <F$NMLoad> {LangType=11, 0e6d='dir'} #18748966
:   {proc=4} t0: OS9$4b <F$AllPrc> {} #18750070
:   :   {proc=4} t0: OS9$28 <F$SRqMem> {size=200} #18750614
:   :   :   <---- t0 OS9$28 <F$SRqMem> {size=200} #18750614 -> size $200 addr $7300 #18763004
:   :   <---- t0 OS9$4b <F$AllPrc> {} #18750070 #18766932
:   {proc=4} t0: OS9$84 <I$Open> {0e6d=''} #18767106
:   :   {proc=4} t0: OS9$30 <F$All64> {table=8200} #18767590
:   :   :   {proc=4} t0: OS9$28 <F$SRqMem> {size=100} #18773752
:   :   :   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #18773752 -> size $100 addr $7200 #18784820
:   :   :   <---- t0 OS9$30 <F$All64> {table=8200} #18767590 -> base $8200 blocknum $4 addr $7200 #18795070
:   :   {proc=4} t0: OS9$49 <F$LDABX> {} #18795144
:   :   :   <---- t0 OS9$49 <F$LDABX> {} #18795144 #18795910
:   :   {proc=1} t0: OS9$10 <F$PrsNam> {path='DD'} #18796102
:   :   :   <---- t0 OS9$10 <F$PrsNam> {path='DD'} #18796102 #18797870
:   :   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #18797912
:   :   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #18799128
:   :   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #18799128 #18827028
:   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #18827160
:   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #18827160 -> addr $e704 entry $e72e #18859198
:   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #18859284
:   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #18859284 -> addr $9e77 entry $9e88 #18917370
:   :   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #18797912 #18919616
:   :   {proc=4} t0: OS9$28 <F$SRqMem> {size=100} #18920112
:   :   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #18920112 -> size $100 addr $7100 #18930496
:   :   {proc=4} t0: OS9$30 <F$All64> {table=8200} #18930544
:   :   :   <---- t0 OS9$30 <F$All64> {table=8200} #18930544 -> base $8200 blocknum $5 addr $7240 #18933146
:   :   {proc=4} t0: OS9$49 <F$LDABX> {} #18934084
:   :   :   <---- t0 OS9$49 <F$LDABX> {} #18934084 #18934850
:   :   {proc=4} t0: OS9$10 <F$PrsNam> {path=''} #19032596
:   :   :   <---- t0 OS9$10 <F$PrsNam> {path=''} #19032596 #19034554
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0e6d destPtr=7260 size=0003} #19034808
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0e6d destPtr=7260 size=0003} #19034808 #19036400
:   :   {proc=4} t0: OS9$10 <F$PrsNam> {path=''} #19262956
:   :   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path=''} #19262956 #19264378
:   :   <---- t0 OS9$84 <I$Open> {0e6d=''} #18767106 -> path $4 #19265582
:   {proc=4} t0: OS9$3f <F$AllTsk> {processDesc=7300} #19265640
:   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7300} #19265640 #19266534
:   {proc=5} t0: OS9$41 <F$SetTsk> {} #19269310
:   :   <---- t0 OS9$41 <F$SetTsk> {} #19269310 #19270024
:   {proc=5} t0: OS9$89 <I$Read> {path=4 buf=0000 size=9} #19270076
:   :   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #19270536
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #19270536 #19271206
:   :   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7100 destPtr=0000 size=0009} #19320990
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7100 destPtr=0000 size=0009} #19320990 #19322668
:   :   <---- t0 OS9$89 <I$Read> {path=4 buf=0000 size=9} #19270076 #19323604
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=0000 destPtr=76d9 size=0009} #19323716
:   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=0000 destPtr=76d9 size=0009} #19323716 #19325394
:   {proc=5} t0: OS9$89 <I$Read> {path=4 buf=0009 size=398} #19325588
:   :   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #19326048
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #19326048 #19328352
:   :   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7109 destPtr=0009 size=00f7} #19330114
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7109 destPtr=0009 size=00f7} #19330114 #19336214
:   :   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7100 destPtr=0100 size=0100} #19384886
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7100 destPtr=0100 size=0100} #19384886 #19392610
:   :   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7100 destPtr=0200 size=0100} #19441282
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7100 destPtr=0200 size=0100} #19441282 #19449006
:   :   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7100 destPtr=0300 size=00a1} #19496464
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7100 destPtr=0300 size=00a1} #19496464 #19500764
:   :   <---- t0 OS9$89 <I$Read> {path=4 buf=0009 size=398} #19325588 #19501700
:   {proc=5} t0: OS9$2e <F$VModul> {addr=0000="-" map=[20 5 20 1 d 34 17 3c]} #19501838
:   :   <---- t0 OS9$2e <F$VModul> {addr=0000="-" map=[20 5 20 1 d 34 17 3c]} #19501838 #19561316
:   {proc=5} t0: OS9$89 <I$Read> {path=4 buf=03a1 size=9} #19561678
:   :   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #19562138
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #19562138 #19562808
:   :   <-- ERROR: $d3(E$EOF    :End of File): OS9KERNEL0 OS9$89 <I$Read> {path=4 buf=03a1 size=9} #19561678 #19565234
:   {proc=4} t0: OS9$8f <I$Close> {path=4} #19566968
:   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=4} #19567416
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #19567416 #19568086
:   :   {proc=4} t0: OS9$29 <F$SRtMem> {size=100 start=7100} #19568548
:   :   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7100} #19568548 #19573036
:   :   {proc=4} t0: OS9$31 <F$Ret64> {block_num=5 address=8200} #19573554
:   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=5 address=8200} #19573554 #19574332
:   :   {proc=4} t0: OS9$81 <I$Detach> {8300} #19574602
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #19575164
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #19575164 #19575904
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #19575924
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #19575924 #19576664
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #19576684
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #19576684 #19577424
:   :   :   <---- t0 OS9$81 <I$Detach> {8300} #19574602 #19577774
:   :   {proc=4} t0: OS9$31 <F$Ret64> {block_num=4 address=8200} #19577804
:   :   :   {proc=4} t0: OS9$29 <F$SRtMem> {size=100 start=7200} #19578440
:   :   :   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7200} #19578440 #19580836
:   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=4 address=8200} #19577804 #19581134
:   :   <---- t0 OS9$8f <I$Close> {path=4} #19566968 #19581382
:   {proc=4} t0: OS9$4c <F$DelPrc> {} #19581514
:   :   {proc=4} t0: OS9$40 <F$DelTsk> {proc_desc=7300} #19581966
:   :   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7300} #19581966 #19582636
:   :   {proc=4} t0: OS9$29 <F$SRtMem> {size=200 start=7300} #19582662
:   :   :   <---- t0 OS9$29 <F$SRtMem> {size=200 start=7300} #19582662 #19585148
:   :   <---- t0 OS9$4c <F$DelPrc> {} #19581514 #19585404
:   {proc=4} t0: OS9$48 <F$LDDDXY> {} #19588050
:   :   <---- t0 OS9$48 <F$LDDDXY> {} #19588050 #19588896
:   {proc=4} t0: OS9$48 <F$LDDDXY> {} #19589024
:   :   <---- t0 OS9$48 <F$LDDDXY> {} #19589024 #19589870
:   <---- t1 OS9$22 <F$NMLoad> {LangType=11, 0e6d='dir'} #18748966 #19591372
{proc=4"Shell"} t1: OS9$03 <F$Fork> {Module/file='dir' param="/b1{13}" lang/type=11 pages=1f} #19591580
:   {proc=4} t0: OS9$28 <F$SRqMem> {size=200} #19592854
:   :   <---- t0 OS9$28 <F$SRqMem> {size=200} #19592854 -> size $200 addr $7300 #19605244
:   {proc=4} t0: OS9$82 <I$Dup> {$2} #19609576
:   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #19610024
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #19610024 #19610694
:   :   <---- t0 OS9$82 <I$Dup> {$2} #19609576 -> path $2 #19610970
:   {proc=4} t0: OS9$82 <I$Dup> {$1} #19611034
:   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=1} #19611482
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #19611482 #19612152
:   :   <---- t0 OS9$82 <I$Dup> {$1} #19611034 -> path $1 #19612428
:   {proc=4} t0: OS9$82 <I$Dup> {$1} #19612492
:   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=1} #19612940
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #19612940 #19613610
:   :   <---- t0 OS9$82 <I$Dup> {$1} #19612492 -> path $1 #19613886
:   {proc=5} t0: OS9$34 <F$SLink> {"dir" type 11 name@ e6d dat@ 7540} #19614076
:   :   <---- t0 OS9$34 <F$SLink> {"dir" type 11 name@ e6d dat@ 7540} #19614076 #19621016
:   {proc=5} t0: OS9$48 <F$LDDDXY> {} #19621182
:   :   <---- t0 OS9$48 <F$LDDDXY> {} #19621182 #19622308
:   {proc=5} t0: OS9$07 <F$Mem> {desired_size=1f00} #19622350
:   :   {proc=5} t0: OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7300} #19622868
:   :   :   <---- t0 OS9$3a <F$AllImg> {beginBlock=0 numBlocks=1 processDesc=7300} #19622868 #19624590
:   :   <---- t0 OS9$07 <F$Mem> {desired_size=1f00} #19622350 #19624896
:   {proc=4} t0: OS9$3f <F$AllTsk> {processDesc=7300} #19626854
:   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7300} #19626854 #19627748
:   {proc=4} t0: OS9$38 <F$Move> {srcTask=2 destTask=3 srcPtr=0e71 destPtr=1efc size=0004} #19627858
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=3 srcPtr=0e71 destPtr=1efc size=0004} #19627858 #19629396
:   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=74f4 destPtr=1ef0 size=000c} #19629462
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=74f4 destPtr=1ef0 size=000c} #19629462 #19631242
:   {proc=4} t0: OS9$40 <F$DelTsk> {proc_desc=7300} #19631272
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7300} #19631272 #19631942
:   {proc=4} t0: OS9$2c <F$AProc> {proc=7300} #19632048
:   :   <---- t0 OS9$2c <F$AProc> {proc=7300} #19632048 #19632772
:   <---- t1 OS9$03 <F$Fork> {Module/file='dir' param="/b1{13}" lang/type=11 pages=1f} #19591580 #19634256
{proc=4"Shell"} t1: OS9$1d <F$UnLoad> {} #19634320
:   {proc=4} t0: OS9$4e <F$FModul> {"Dir" type 11 name@ e3 dat@ 7540} #19635436
:   :   <---- t0 OS9$4e <F$FModul> {"Dir" type 11 name@ e3 dat@ 7540} #19635436 #19639926
:   <---- t1 OS9$1d <F$UnLoad> {} #19634320 #19641176
{proc=4"Shell"} t1: OS9$04 <F$Wait> {} #19641738
:   {proc=4} t0: OS9$40 <F$DelTsk> {proc_desc=7500} #19643172
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7500} #19643172 #19643842
:   {proc=4} t0: OS9$2d <F$NProc> {} #19643896
{proc=5"Dir"} t1: OS9$8d <I$GetStt> {path=1 26==SS.ScSiz  : Return screen size for COCO} #19645452
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=1} #19646780
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #19646780 #19647450
:   <---- t1 OS9$8d <I$GetStt> {path=1 26==SS.ScSiz  : Return screen size for COCO} #19645452 #19649958
{proc=5"Dir"} t1: OS9$84 <I$Open> {1efc='/b1'} #19650502
:   {proc=5} t0: OS9$30 <F$All64> {table=8200} #19651934
:   :   {proc=5} t0: OS9$28 <F$SRqMem> {size=100} #19659678
:   :   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #19659678 -> size $100 addr $7200 #19670944
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #19651934 -> base $8200 blocknum $4 addr $7200 #19679560
:   {proc=5} t0: OS9$49 <F$LDABX> {} #19679634
:   :   <---- t0 OS9$49 <F$LDABX> {} #19679634 #19680400
:   {proc=5} t0: OS9$10 <F$PrsNam> {path=''} #19680478
:   :   <---- t0 OS9$10 <F$PrsNam> {path=''} #19680478 #19682406
:   {proc=5} t0: OS9$80 <I$Attach> {1efd=''} #19682448
:   :   {proc=1} t0: OS9$34 <F$SLink> {"b1" type f0 name@ 1efd dat@ 7340} #19683664
:   :   :   <---- t0 OS9$34 <F$SLink> {"b1" type f0 name@ 1efd dat@ 7340} #19683664 #19706984
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #19707116
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #19707116 -> addr $e704 entry $e72e #19740400
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #19740486
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #19740486 -> addr $9e77 entry $9e88 #19799818
:   :   <---- t0 OS9$80 <I$Attach> {1efd=''} #19682448 #19802538
:   {proc=5} t0: OS9$28 <F$SRqMem> {size=100} #19803034
:   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #19803034 -> size $100 addr $7100 #19815052
:   {proc=5} t0: OS9$30 <F$All64> {table=8200} #19815100
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #19815100 -> base $8200 blocknum $5 addr $7240 #19817702
:   {proc=5} t0: OS9$49 <F$LDABX> {} #19818640
:   :   <---- t0 OS9$49 <F$LDABX> {} #19818640 #19819406
:   {proc=5} t0: OS9$10 <F$PrsNam> {path=''} #19819478
:   :   <---- t0 OS9$10 <F$PrsNam> {path=''} #19819478 #19821406
:   {proc=5} t0: OS9$10 <F$PrsNam> {path=''} #19919666
:   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path=''} #19919666 #19920908
:   <---- t1 OS9$84 <I$Open> {1efc='/b1'} #19650502 -> path $3 #19923422
{proc=5"Dir"} t1: OS9$86 <I$ChgDir> {mode=81, 1efc='/b1'} #19925922
:   {proc=5} t0: OS9$30 <F$All64> {table=8200} #19927142
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #19927142 -> base $8200 blocknum $6 addr $7280 #19929780
:   {proc=5} t0: OS9$49 <F$LDABX> {} #19929854
:   :   <---- t0 OS9$49 <F$LDABX> {} #19929854 #19930620
:   {proc=5} t0: OS9$10 <F$PrsNam> {path=''} #19930698
:   :   <---- t0 OS9$10 <F$PrsNam> {path=''} #19930698 #19932626
:   {proc=5} t0: OS9$80 <I$Attach> {1efd=''} #19932668
:   :   {proc=1} t0: OS9$34 <F$SLink> {"b1" type f0 name@ 1efd dat@ 7340} #19933884
:   :   :   <---- t0 OS9$34 <F$SLink> {"b1" type f0 name@ 1efd dat@ 7340} #19933884 #19957152
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #19957284
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #19957284 -> addr $e704 entry $e72e #19990568
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #19990654
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #19990654 -> addr $9e77 entry $9e88 #20049986
:   :   <---- t0 OS9$80 <I$Attach> {1efd=''} #19932668 #20052324
:   {proc=5} t0: OS9$28 <F$SRqMem> {size=100} #20052876
:   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #20052876 -> size $100 addr $7000 #20063328
:   {proc=5} t0: OS9$30 <F$All64> {table=8200} #20063376
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #20063376 -> base $8200 blocknum $7 addr $72c0 #20066050
:   {proc=5} t0: OS9$49 <F$LDABX> {} #20066988
:   :   <---- t0 OS9$49 <F$LDABX> {} #20066988 #20067754
:   {proc=5} t0: OS9$10 <F$PrsNam> {path=''} #20067826
:   :   <---- t0 OS9$10 <F$PrsNam> {path=''} #20067826 #20069754
:   {proc=5} t0: OS9$10 <F$PrsNam> {path=''} #20166578
:   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path=''} #20166578 #20167820
:   {proc=5} t0: OS9$29 <F$SRtMem> {size=100 start=7000} #20168798
:   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7000} #20168798 #20173258
:   {proc=5} t0: OS9$31 <F$Ret64> {block_num=7 address=8200} #20173748
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=7 address=8200} #20173748 #20174526
:   {proc=5} t0: OS9$81 <I$Detach> {831a} #20174866
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #20175428
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #20175428 #20176168
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #20176188
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #20176188 #20176928
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec61 magic=87cd module='B1'} #20176948
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec61 magic=87cd module='B1'} #20176948 #20177688
:   :   <---- t0 OS9$81 <I$Detach> {831a} #20174866 #20178038
:   {proc=5} t0: OS9$31 <F$Ret64> {block_num=6 address=8200} #20178068
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=6 address=8200} #20178068 #20178846
:   <---- t1 OS9$86 <I$ChgDir> {mode=81, 1efc='/b1'} #19925922 #20180338
{proc=5"Dir"} t1: OS9$15 <F$Time> {buf=d} #20182860
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=0028 destPtr=000d size=0006} #20183976
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=0028 destPtr=000d size=0006} #20183976 #20185582
:   <---- t1 OS9$15 <F$Time> {buf=d} #20182860 #20186750
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #20187864
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=1} #20189132
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #20189132 #20189802
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0040 destPtr=7ce0 size=0020} #20190944
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0040 destPtr=7ce0 size=0020} #20190944 #20194564
:   <---- t1 OS9$8c <I$WritLn> {} #20187864 #20217724
{proc=5"Dir"} t1: OS9$88 <I$Seek> {path=3 pos=00000040} #20217814
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #20219060
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #20219060 #20219730
:   <---- t1 OS9$88 <I$Seek> {path=3 pos=00000040} #20217814 #20221340
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20221402
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #20224318
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #20224318 #20224988
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7140 destPtr=0013 size=0020} #20280308
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7140 destPtr=0013 size=0020} #20280308 #20283928
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20221402 #20286108
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20287608
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #20288870
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #20288870 #20289540
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7160 destPtr=0013 size=0020} #20291314
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7160 destPtr=0013 size=0020} #20291314 #20293300
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20287608 #20295164
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20296748
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #20298010
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #20298010 #20298680
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7180 destPtr=0013 size=0020} #20300454
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7180 destPtr=0013 size=0020} #20300454 #20302440
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20296748 #20304304
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20305816
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #20307078
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #20307078 #20307748
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=71a0 destPtr=0013 size=0020} #20309522
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=71a0 destPtr=0013 size=0020} #20309522 #20313142
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20305816 #20315006
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #20315770
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=1} #20317038
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #20317038 #20317708
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0040 destPtr=7ce0 size=0020} #20318850
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0040 destPtr=7ce0 size=0020} #20318850 #20320836
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0060 destPtr=7ce0 size=0020} #20333224
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0060 destPtr=7ce0 size=0020} #20333224 #20335210
:   <---- t1 OS9$8c <I$WritLn> {} #20315770 #20347466
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20347536
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #20348798
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #20348798 #20349468
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=71c0 destPtr=0013 size=0020} #20351242
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=71c0 destPtr=0013 size=0020} #20351242 #20353228
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20347536 #20355092
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20356544
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #20357806
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #20357806 #20358476
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=71e0 destPtr=0013 size=0020} #20360238
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=71e0 destPtr=0013 size=0020} #20360238 #20362224
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20356544 #20364240
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20365710
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #20366972
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #20366972 #20367642
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7100 destPtr=0013 size=0020} #20420586
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7100 destPtr=0013 size=0020} #20420586 #20422572
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20365710 #20424436
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20425948
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #20427210
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #20427210 #20427880
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7120 destPtr=0013 size=0020} #20429654
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7120 destPtr=0013 size=0020} #20429654 #20435690
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20425948 #20437870
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #20438878
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=1} #20440146
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #20440146 #20440816
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0040 destPtr=7ce0 size=0020} #20441958
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0040 destPtr=7ce0 size=0020} #20441958 #20443944
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0060 destPtr=7ce0 size=0020} #20456156
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0060 destPtr=7ce0 size=0020} #20456156 #20458142
:   <---- t1 OS9$8c <I$WritLn> {} #20438878 #20470910
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20470980
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #20472242
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #20472242 #20472912
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7140 destPtr=0013 size=0020} #20474686
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7140 destPtr=0013 size=0020} #20474686 #20476672
:   <---- t1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20470980 #20478536
{proc=5"Dir"} t1: OS9$89 <I$Read> {path=3 buf=0013 size=20} #20479964
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #20481226
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #20481226 #20481896
:   <-- ERROR: $d3(E$EOF    :End of File): OS9KERNEL1 OS9$89 <I$Read> {path=3 buf=0013 size=20} #20479964 #20485250
{proc=5"Dir"} t1: OS9$8c <I$WritLn> {} #20485456
:   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=1} #20486724
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #20486724 #20487394
:   {proc=5} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0040 destPtr=7ce0 size=0020} #20488536
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0040 destPtr=7ce0 size=0020} #20488536 #20492176
:   <---- t1 OS9$8c <I$WritLn> {} #20485456 #20505544
{proc=5"Dir"} t1: OS9$06 <F$Exit> {status=0} #20505582
:   {proc=5} t0: OS9$8f <I$Close> {path=2} #20506738
:   :   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=2} #20507186
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #20507186 #20507856
:   :   <---- t0 OS9$8f <I$Close> {path=2} #20506738 #20508570
:   {proc=5} t0: OS9$8f <I$Close> {path=1} #20508648
:   :   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=1} #20509096
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #20509096 #20509766
:   :   <---- t0 OS9$8f <I$Close> {path=1} #20508648 #20510958
:   {proc=5} t0: OS9$8f <I$Close> {path=1} #20511036
:   :   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=1} #20511484
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #20511484 #20512154
:   :   {proc=5} t0: OS9$37 <F$GProcP> {id=04} #20513352
:   :   :   <---- t0 OS9$37 <F$GProcP> {id=04} #20513352 #20514060
:   :   <---- t0 OS9$8f <I$Close> {path=1} #20511036 #20514798
:   {proc=5} t0: OS9$8f <I$Close> {path=4} #20514876
:   :   {proc=5} t0: OS9$2f <F$Find64> {base=8200 id=4} #20515324
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=4} #20515324 #20515994
:   :   {proc=5} t0: OS9$29 <F$SRtMem> {size=100 start=7100} #20516456
:   :   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7100} #20516456 #20518824
:   :   {proc=5} t0: OS9$31 <F$Ret64> {block_num=5 address=8200} #20519290
:   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=5 address=8200} #20519290 #20521722
:   :   {proc=5} t0: OS9$81 <I$Detach> {831a} #20521992
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #20522554
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #20522554 #20523294
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #20523314
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #20523314 #20524054
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec61 magic=87cd module='B1'} #20524074
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec61 magic=87cd module='B1'} #20524074 #20524814
:   :   :   <---- t0 OS9$81 <I$Detach> {831a} #20521992 #20525164
:   :   {proc=5} t0: OS9$31 <F$Ret64> {block_num=4 address=8200} #20525194
:   :   :   {proc=5} t0: OS9$29 <F$SRtMem> {size=100 start=7200} #20525830
:   :   :   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7200} #20525830 #20528226
:   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=4 address=8200} #20525194 #20528524
:   :   <---- t0 OS9$8f <I$Close> {path=4} #20514876 #20528772
:   {proc=5} t0: OS9$3b <F$DelImg> {beginBlock=0 numBlocks=1 processDesc=7300} #20529186
:   :   <---- t0 OS9$3b <F$DelImg> {beginBlock=0 numBlocks=1 processDesc=7300} #20529186 #20529964
:   {proc=5} t0: OS9$02 <F$UnLink> {u=e000 magic=ae48 module='''} #20530022
:   :   {proc=5} t0: OS9$48 <F$LDDDXY> {} #20533102
:   :   :   <---- t0 OS9$48 <F$LDDDXY> {} #20533102 #20534228
:   :   <---- t0 OS9$02 <F$UnLink> {u=e000 magic=ae48 module='''} #20530022 #20539638
:   {proc=5} t0: OS9$40 <F$DelTsk> {proc_desc=7300} #20539678
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7300} #20539678 #20540348
:   {proc=5} t0: OS9$40 <F$DelTsk> {proc_desc=7300} #20540800
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7300} #20540800 #20541394
:   {proc=5} t0: OS9$29 <F$SRtMem> {size=200 start=7300} #20541420
:   :   <---- t0 OS9$29 <F$SRtMem> {size=200 start=7300} #20541420 #20543906
:   {proc=5} t0: OS9$2c <F$AProc> {proc=7500} #20543944
:   :   <---- t0 OS9$2c <F$AProc> {proc=7500} #20543944 #20544668
:   {proc=5} t0: OS9$2d <F$NProc> {} #20544680
:   :   <---- t1 OS9$04 <F$Wait> {} #19641738 #20547166
{proc=4"Shell"} t1: OS9$8b <I$ReadLn> {} #20548542
:   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #20553802
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #20553802 #20554472
:   <-- ERROR: $d3(E$EOF    :End of File): OS9KERNEL1 OS9$8b <I$ReadLn> {} #20548542 #20558142
{proc=4"Shell"} t1: OS9$1c <F$SUser> {} #20558528
:   <---- t1 OS9$1c <F$SUser> {} #20558528 #20560782
{proc=4"Shell"} t1: OS9$06 <F$Exit> {status=0} #20560826
:   {proc=4} t0: OS9$8f <I$Close> {path=2} #20561982
:   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #20562430
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #20562430 #20563100
:   :   {proc=4} t0: OS9$29 <F$SRtMem> {size=100 start=7700} #20563562
:   :   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7700} #20563562 #20566014
:   :   {proc=4} t0: OS9$31 <F$Ret64> {block_num=3 address=8200} #20566480
:   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=3 address=8200} #20566480 #20567258
:   :   {proc=4} t0: OS9$81 <I$Detach> {8300} #20567528
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #20568090
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #20568090 #20568830
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #20568850
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #20568850 #20569590
:   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #20569610
:   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #20569610 #20570350
:   :   :   <---- t0 OS9$81 <I$Detach> {8300} #20567528 #20570700
:   :   {proc=4} t0: OS9$31 <F$Ret64> {block_num=2 address=8200} #20570730
:   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=2 address=8200} #20570730 #20571508
:   :   <---- t0 OS9$8f <I$Close> {path=2} #20561982 #20571756
:   {proc=4} t0: OS9$8f <I$Close> {path=1} #20571834
:   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=1} #20572282
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #20572282 #20572952
:   :   <---- t0 OS9$8f <I$Close> {path=1} #20571834 #20574144
:   {proc=4} t0: OS9$8f <I$Close> {path=1} #20574222
:   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=1} #20574670
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #20574670 #20575340
:   :   {proc=4} t0: OS9$37 <F$GProcP> {id=03} #20576538
:   :   :   <---- t0 OS9$37 <F$GProcP> {id=03} #20576538 #20577246
:   :   <---- t0 OS9$8f <I$Close> {path=1} #20574222 #20577954
:   {proc=4} t0: OS9$3b <F$DelImg> {beginBlock=0 numBlocks=1 processDesc=7500} #20578396
:   :   <---- t0 OS9$3b <F$DelImg> {beginBlock=0 numBlocks=1 processDesc=7500} #20578396 #20579174
:   {proc=4} t0: OS9$02 <F$UnLink> {u=e000 magic=ae48 module='''} #20579232
:   :   <---- t0 OS9$02 <F$UnLink> {u=e000 magic=ae48 module='''} #20579232 #20583756
:   {proc=4} t0: OS9$40 <F$DelTsk> {proc_desc=7500} #20583796
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7500} #20583796 #20584466
:   {proc=4} t0: OS9$40 <F$DelTsk> {proc_desc=7500} #20584918
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7500} #20584918 #20585512
:   {proc=4} t0: OS9$29 <F$SRtMem> {size=200 start=7500} #20585538
:   :   <---- t0 OS9$29 <F$SRtMem> {size=200 start=7500} #20585538 #20588108
:   {proc=4} t0: OS9$2c <F$AProc> {proc=7800} #20588146
:   :   <---- t0 OS9$2c <F$AProc> {proc=7800} #20588146 #20588870
:   {proc=4} t0: OS9$2d <F$NProc> {} #20588882
:   :   <---- t1 OS9$04 <F$Wait> {} #14922642 #20591368
{proc=3"Shell"} t1: OS9$1c <F$SUser> {} #20592658
:   <---- t1 OS9$1c <F$SUser> {} #20592658 #20594912
{proc=3"Shell"} t1: OS9$06 <F$Exit> {status=0} #20594956
:   {proc=3} t0: OS9$8f <I$Close> {path=1} #20596112
:   :   {proc=3} t0: OS9$2f <F$Find64> {base=8200 id=1} #20596560
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #20596560 #20597230
:   :   <---- t0 OS9$8f <I$Close> {path=1} #20596112 #20598362
:   {proc=3} t0: OS9$8f <I$Close> {path=1} #20598440
:   :   {proc=3} t0: OS9$2f <F$Find64> {base=8200 id=1} #20598888
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #20598888 #20599558
:   :   <---- t0 OS9$8f <I$Close> {path=1} #20598440 #20600750
:   {proc=3} t0: OS9$8f <I$Close> {path=1} #20600828
:   :   {proc=3} t0: OS9$2f <F$Find64> {base=8200 id=1} #20601276
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #20601276 #20601946
:   :   {proc=3} t0: OS9$37 <F$GProcP> {id=02} #20603144
:   :   :   <---- t0 OS9$37 <F$GProcP> {id=02} #20603144 #20603852
:   :   <---- t0 OS9$8f <I$Close> {path=1} #20600828 #20604560
:   {proc=3} t0: OS9$3b <F$DelImg> {beginBlock=0 numBlocks=1 processDesc=7800} #20605002
:   :   <---- t0 OS9$3b <F$DelImg> {beginBlock=0 numBlocks=1 processDesc=7800} #20605002 #20605780
:   {proc=3} t0: OS9$02 <F$UnLink> {u=e000 magic=ae48 module='''} #20605838
:   :   {proc=3} t0: OS9$48 <F$LDDDXY> {} #20608352
:   :   :   <---- t0 OS9$48 <F$LDDDXY> {} #20608352 #20611060
:   :   <---- t0 OS9$02 <F$UnLink> {u=e000 magic=ae48 module='''} #20605838 #20616822
:   {proc=3} t0: OS9$40 <F$DelTsk> {proc_desc=7800} #20616862
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7800} #20616862 #20617532
:   {proc=3} t0: OS9$40 <F$DelTsk> {proc_desc=7800} #20617984
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7800} #20617984 #20618578
:   {proc=3} t0: OS9$29 <F$SRtMem> {size=200 start=7800} #20618604
:   :   <---- t0 OS9$29 <F$SRtMem> {size=200 start=7800} #20618604 #20621230
:   {proc=3} t0: OS9$2c <F$AProc> {proc=7a00} #20621268
:   :   <---- t0 OS9$2c <F$AProc> {proc=7a00} #20621268 #20621992
:   {proc=3} t0: OS9$2d <F$NProc> {} #20622004
:   :   <---- t1 OS9$04 <F$Wait> {} #12769486 #20624490
{proc=2"SysGo"} t1: OS9$03 <F$Fork> {Module/file='AutoEx' param="{13}" lang/type=1 pages=0} #20624552
:   {proc=2} t0: OS9$28 <F$SRqMem> {size=200} #20625790
:   :   <---- t0 OS9$28 <F$SRqMem> {size=200} #20625790 -> size $200 addr $7800 #20636008
:   {proc=2} t0: OS9$82 <I$Dup> {$1} #20641922
:   :   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #20642370
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #20642370 #20643040
:   :   <---- t0 OS9$82 <I$Dup> {$1} #20641922 -> path $1 #20643316
:   {proc=2} t0: OS9$82 <I$Dup> {$1} #20643380
:   :   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #20643828
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #20643828 #20644498
:   :   <---- t0 OS9$82 <I$Dup> {$1} #20643380 -> path $1 #20644774
:   {proc=2} t0: OS9$82 <I$Dup> {$1} #20644838
:   :   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #20645286
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #20645286 #20645956
:   :   <---- t0 OS9$82 <I$Dup> {$1} #20644838 -> path $1 #20646232
:   {proc=3} t0: OS9$34 <F$SLink> {"AutoEx" type 1 name@ d945 dat@ 7a40} #20646422
:   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$34 <F$SLink> {"AutoEx" type 1 name@ d945 dat@ 7a40} #20646422 #20696668
:   {proc=2} t0: OS9$01 <F$Load> {type/lang=7a filename='AutoEx'} #20696718
:   :   {proc=2} t0: OS9$4b <F$AllPrc> {} #20697092
:   :   :   {proc=2} t0: OS9$28 <F$SRqMem> {size=200} #20697618
:   :   :   :   <---- t0 OS9$28 <F$SRqMem> {size=200} #20697618 -> size $200 addr $7600 #20709876
:   :   :   <---- t0 OS9$4b <F$AllPrc> {} #20697092 #20713804
:   :   {proc=2} t0: OS9$84 <I$Open> {d945='AutoEx'} #20713978
:   :   :   {proc=2} t0: OS9$30 <F$All64> {table=8200} #20714462
:   :   :   :   <---- t0 OS9$30 <F$All64> {table=8200} #20714462 -> base $8200 blocknum $2 addr $8280 #20716852
:   :   :   {proc=2} t0: OS9$49 <F$LDABX> {} #20716926
:   :   :   :   <---- t0 OS9$49 <F$LDABX> {} #20716926 #20717780
:   :   :   {proc=1} t0: OS9$10 <F$PrsNam> {path='DD'} #20717972
:   :   :   :   <---- t0 OS9$10 <F$PrsNam> {path='DD'} #20717972 #20719740
:   :   :   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #20719782
:   :   :   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #20720998
:   :   :   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #20720998 #20736768
:   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #20736900
:   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #20736900 -> addr $e704 entry $e72e #20755982
:   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #20756068
:   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #20756068 -> addr $9e77 entry $9e88 #20802832
:   :   :   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #20719782 #20805170
:   :   :   {proc=2} t0: OS9$28 <F$SRqMem> {size=100} #20805666
:   :   :   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #20805666 -> size $100 addr $7500 #20815778
:   :   :   {proc=2} t0: OS9$30 <F$All64> {table=8200} #20815826
:   :   :   :   <---- t0 OS9$30 <F$All64> {table=8200} #20815826 -> base $8200 blocknum $3 addr $82c0 #20818252
:   :   :   {proc=2} t0: OS9$49 <F$LDABX> {} #20820824
:   :   :   :   <---- t0 OS9$49 <F$LDABX> {} #20820824 #20821678
:   :   :   {proc=2} t0: OS9$10 <F$PrsNam> {path='AutoEx'} #20917804
:   :   :   :   <---- t0 OS9$10 <F$PrsNam> {path='AutoEx'} #20917804 #20920782
:   :   :   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d945 destPtr=82e0 size=0006} #20921036
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d945 destPtr=82e0 size=0006} #20921036 #20922818
:   :   :   {proc=2} t0: OS9$29 <F$SRtMem> {size=100 start=7500} #21742800
:   :   :   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7500} #21742800 #21749034
:   :   :   {proc=2} t0: OS9$31 <F$Ret64> {block_num=3 address=8200} #21749500
:   :   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=3 address=8200} #21749500 #21750278
:   :   :   {proc=2} t0: OS9$81 <I$Detach> {8300} #21750524
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #21751086
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #21751086 #21751826
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #21751846
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #21751846 #21752586
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #21752606
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #21752606 #21753346
:   :   :   :   <---- t0 OS9$81 <I$Detach> {8300} #21750524 #21753696
:   :   :   {proc=2} t0: OS9$31 <F$Ret64> {block_num=2 address=8200} #21753726
:   :   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=2 address=8200} #21753726 #21754504
:   :   :   <-- ERROR: $d8(E$PNNF   :Path Name Not Found): OS9KERNEL0 OS9$84 <I$Open> {d945='AutoEx'} #20713978 #21754762
:   :   {proc=2} t0: OS9$4c <F$DelPrc> {} #21754924
:   :   :   {proc=2} t0: OS9$40 <F$DelTsk> {proc_desc=7600} #21755376
:   :   :   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7600} #21755376 #21755970
:   :   :   {proc=2} t0: OS9$29 <F$SRtMem> {size=200 start=7600} #21755996
:   :   :   :   <---- t0 OS9$29 <F$SRtMem> {size=200 start=7600} #21755996 #21758566
:   :   :   <---- t0 OS9$4c <F$DelPrc> {} #21754924 #21758822
:   :   <-- ERROR: $d8(E$PNNF   :Path Name Not Found): OS9KERNEL0 OS9$01 <F$Load> {type/lang=7a filename='AutoEx'} #20696718 #21759174
:   {proc=2} t0: OS9$8f <I$Close> {path=1} #21759350
:   :   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #21759798
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #21759798 #21760468
:   :   <---- t0 OS9$8f <I$Close> {path=1} #21759350 #21761540
:   {proc=2} t0: OS9$8f <I$Close> {path=1} #21761618
:   :   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #21762066
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #21762066 #21762736
:   :   <---- t0 OS9$8f <I$Close> {path=1} #21761618 #21763808
:   {proc=2} t0: OS9$8f <I$Close> {path=1} #21763886
:   :   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #21764334
:   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #21764334 #21765004
:   :   <---- t0 OS9$8f <I$Close> {path=1} #21763886 #21766076
:   {proc=3} t0: OS9$02 <F$UnLink> {u=0000 magic=0000 module=''} #21766540
:   :   <---- t0 OS9$02 <F$UnLink> {u=0000 magic=0000 module=''} #21766540 #21767208
:   {proc=2} t0: OS9$40 <F$DelTsk> {proc_desc=7800} #21767248
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7800} #21767248 #21767842
:   {proc=2} t0: OS9$40 <F$DelTsk> {proc_desc=7800} #21767992
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7800} #21767992 #21768586
:   {proc=2} t0: OS9$29 <F$SRtMem> {size=200 start=7800} #21768612
:   :   <---- t0 OS9$29 <F$SRtMem> {size=200 start=7800} #21768612 #21771238
:   <-- ERROR: $d8(E$PNNF   :Path Name Not Found): OS9KERNEL1 OS9$03 <F$Fork> {Module/file='AutoEx' param="{13}" lang/type=1 pages=0} #20624552 #21772772
{proc=2"SysGo"} t1: OS9$05 <F$Chain> {Module/file='Shell' param="i=/1{13}" lang/type=1 pages=0} #21775490
:   {proc=2} t0: OS9$28 <F$SRqMem> {size=200} #21776728
:   :   <---- t0 OS9$28 <F$SRqMem> {size=200} #21776728 -> size $200 addr $7800 #21786946
:   {proc=2} t0: OS9$02 <F$UnLink> {u=d893 magic=87cd module='SysGo'} #21795370
:   :   <---- t0 OS9$02 <F$UnLink> {u=d893 magic=87cd module='SysGo'} #21795370 #21796110
:   {proc=2} t0: OS9$34 <F$SLink> {"Shell" type 1 name@ d93f dat@ 7840} #21796552
:   :   <-- ERROR: $dd(E$MNF    :Module Not Found): OS9KERNEL0 OS9$34 <F$SLink> {"Shell" type 1 name@ d93f dat@ 7840} #21796552 #21849352
:   {proc=3} t0: OS9$01 <F$Load> {type/lang=78 filename='Shell'} #21849402
:   :   {proc=3} t0: OS9$4b <F$AllPrc> {} #21849776
:   :   :   {proc=3} t0: OS9$28 <F$SRqMem> {size=200} #21850302
:   :   :   :   <---- t0 OS9$28 <F$SRqMem> {size=200} #21850302 -> size $200 addr $7600 #21860926
:   :   :   <---- t0 OS9$4b <F$AllPrc> {} #21849776 #21866488
:   :   {proc=3} t0: OS9$84 <I$Open> {d93f='Shell'} #21866662
:   :   :   {proc=3} t0: OS9$30 <F$All64> {table=8200} #21867146
:   :   :   :   <---- t0 OS9$30 <F$All64> {table=8200} #21867146 -> base $8200 blocknum $2 addr $8280 #21869536
:   :   :   {proc=3} t0: OS9$49 <F$LDABX> {} #21869610
:   :   :   :   <---- t0 OS9$49 <F$LDABX> {} #21869610 #21870464
:   :   :   {proc=1} t0: OS9$10 <F$PrsNam> {path='DD'} #21870656
:   :   :   :   <---- t0 OS9$10 <F$PrsNam> {path='DD'} #21870656 #21872424
:   :   :   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #21872466
:   :   :   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #21873682
:   :   :   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #21873682 #21887818
:   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #21887950
:   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #21887950 -> addr $e704 entry $e72e #21908666
:   :   :   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #21908752
:   :   :   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #21908752 -> addr $9e77 entry $9e88 #21955516
:   :   :   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #21872466 #21957854
:   :   :   {proc=3} t0: OS9$28 <F$SRqMem> {size=100} #21958350
:   :   :   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #21958350 -> size $100 addr $7500 #21968462
:   :   :   {proc=3} t0: OS9$30 <F$All64> {table=8200} #21968510
:   :   :   :   <---- t0 OS9$30 <F$All64> {table=8200} #21968510 -> base $8200 blocknum $3 addr $82c0 #21970936
:   :   :   {proc=3} t0: OS9$49 <F$LDABX> {} #21971874
:   :   :   :   <---- t0 OS9$49 <F$LDABX> {} #21971874 #21972728
:   :   :   {proc=3} t0: OS9$10 <F$PrsNam> {path='Shell'} #22068854
:   :   :   :   <---- t0 OS9$10 <F$PrsNam> {path='Shell'} #22068854 #22073196
:   :   :   {proc=3} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d93f destPtr=82e0 size=0005} #22073450
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=d93f destPtr=82e0 size=0005} #22073450 #22075198
:   :   :   {proc=3} t0: OS9$10 <F$PrsNam> {path=''} #22710530
:   :   :   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path=''} #22710530 #22712012
:   :   :   <---- t0 OS9$84 <I$Open> {d93f='Shell'} #21866662 -> path $2 #22713216
:   :   {proc=3} t0: OS9$3f <F$AllTsk> {processDesc=7600} #22713274
:   :   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7600} #22713274 #22714168
:   :   {proc=4} t0: OS9$41 <F$SetTsk> {} #22715160
:   :   :   <---- t0 OS9$41 <F$SetTsk> {} #22715160 #22715874
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=0000 size=9} #22715926
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #22716386
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #22716386 #22717056
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0000 size=0009} #22765270
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0000 size=0009} #22765270 #22766948
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=0000 size=9} #22715926 #22767884
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=0000 destPtr=7bb4 size=0009} #22767996
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=0000 destPtr=7bb4 size=0009} #22767996 #22769674
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=0009 size=1b4e} #22769868
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #22770328
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #22770328 #22770998
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7509 destPtr=0009 size=00f7} #22772760
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7509 destPtr=0009 size=00f7} #22772760 #22778860
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0100 size=0100} #22825962
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0100 size=0100} #22825962 #22832052
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0200 size=0100} #22883864
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0200 size=0100} #22883864 #22889954
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0300 size=0100} #22939008
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0300 size=0100} #22939008 #22945098
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0400 size=0100} #22993770
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0400 size=0100} #22993770 #23001494
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0500 size=0100} #23050166
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0500 size=0100} #23050166 #23057890
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0600 size=0100} #23106562
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0600 size=0100} #23106562 #23112652
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0700 size=0100} #23161324
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0700 size=0100} #23161324 #23167414
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0800 size=0100} #23216086
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0800 size=0100} #23216086 #23222176
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0900 size=0100} #23270848
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0900 size=0100} #23270848 #23276938
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0a00 size=0100} #23326234
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0a00 size=0100} #23326234 #23332324
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0b00 size=0100} #23381378
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0b00 size=0100} #23381378 #23389144
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0c00 size=0100} #23437816
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0c00 size=0100} #23437816 #23445540
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0d00 size=0100} #23494212
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0d00 size=0100} #23494212 #23500302
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0e00 size=0100} #23548974
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0e00 size=0100} #23548974 #23555064
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0f00 size=0100} #23603736
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=0f00 size=0100} #23603736 #23609826
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1000 size=0100} #23658498
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1000 size=0100} #23658498 #23664588
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1100 size=0100} #23711690
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1100 size=0100} #23711690 #23719414
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1200 size=0100} #23768086
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1200 size=0100} #23768086 #23775810
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1300 size=0100} #23824864
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1300 size=0100} #23824864 #23830954
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1400 size=0100} #23881260
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1400 size=0100} #23881260 #23887350
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1500 size=0100} #23936022
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1500 size=0100} #23936022 #23942112
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1600 size=0100} #23990784
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1600 size=0100} #23990784 #23996874
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1700 size=0100} #24043030
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1700 size=0100} #24043030 #24049120
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1800 size=0100} #24097792
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1800 size=0100} #24097792 #24105516
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1900 size=0100} #24154188
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1900 size=0100} #24154188 #24161912
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1a00 size=0100} #24210584
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1a00 size=0100} #24210584 #24216674
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1b00 size=0057} #24265702
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1b00 size=0057} #24265702 #24268610
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=0009 size=1b4e} #22769868 #24269546
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=0000="-" map=[2e 1 26 2e 26 26 26 26]} #24269684
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=0000="-" map=[2e 1 26 2e 26 26 26 26]} #24269684 #24320850
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1b57 size=9} #24321212
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #24321672
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24321672 #24322342
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7557 destPtr=1b57 size=0009} #24324116
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7557 destPtr=1b57 size=0009} #24324116 #24325794
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1b57 size=9} #24321212 #24326730
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1b57 destPtr=7bb4 size=0009} #24326842
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1b57 destPtr=7bb4 size=0009} #24326842 #24328520
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1b60 size=e8} #24328714
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #24329174
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24329174 #24329844
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7560 destPtr=1b60 size=00a0} #24331606
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7560 destPtr=1b60 size=00a0} #24331606 #24335872
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1c00 size=0048} #24386088
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1c00 size=0048} #24386088 #24388644
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1b60 size=e8} #24328714 #24389580
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=1b57="-" map=[2e 1 26 2e 26 26 26 26]} #24389718
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=1b57="-" map=[2e 1 26 2e 26 26 26 26]} #24389718 #24441056
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1c48 size=9} #24441322
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #24441782
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24441782 #24442452
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7548 destPtr=1c48 size=0009} #24444226
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7548 destPtr=1c48 size=0009} #24444226 #24445904
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1c48 size=9} #24441322 #24446840
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1c48 destPtr=7bb4 size=0009} #24446952
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1c48 destPtr=7bb4 size=0009} #24446952 #24448630
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1c51 size=4a} #24448824
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #24449284
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24449284 #24449954
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7551 destPtr=1c51 size=004a} #24451728
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7551 destPtr=1c51 size=004a} #24451728 #24454352
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1c51 size=4a} #24448824 #24455288
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=1c48="-" map=[2e 1 26 2e 26 26 26 26]} #24455426
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=1c48="-" map=[2e 1 26 2e 26 26 26 26]} #24455426 #24508686
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1c9b size=9} #24508952
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #24509412
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24509412 #24510082
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=759b destPtr=1c9b size=0009} #24511856
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=759b destPtr=1c9b size=0009} #24511856 #24513534
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1c9b size=9} #24508952 #24514470
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1c9b destPtr=7bb4 size=0009} #24514582
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1c9b destPtr=7bb4 size=0009} #24514582 #24516260
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1ca4 size=19} #24516454
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #24516914
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24516914 #24517584
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75a4 destPtr=1ca4 size=0019} #24520992
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75a4 destPtr=1ca4 size=0019} #24520992 #24522898
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1ca4 size=19} #24516454 #24523834
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=1c9b="-" map=[2e 1 26 2e 26 26 26 26]} #24523972
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=1c9b="-" map=[2e 1 26 2e 26 26 26 26]} #24523972 #24575268
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1cbd size=9} #24575534
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #24575994
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24575994 #24576664
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75bd destPtr=1cbd size=0009} #24580072
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75bd destPtr=1cbd size=0009} #24580072 #24581750
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1cbd size=9} #24575534 #24582686
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1cbd destPtr=7bb4 size=0009} #24582798
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1cbd destPtr=7bb4 size=0009} #24582798 #24584476
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1cc6 size=3e} #24584670
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #24585130
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24585130 #24585800
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75c6 destPtr=1cc6 size=003a} #24587562
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75c6 destPtr=1cc6 size=003a} #24587562 #24589958
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1d00 size=0004} #24639228
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7500 destPtr=1d00 size=0004} #24639228 #24640854
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1cc6 size=3e} #24584670 #24641790
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=1cbd="-" map=[2e 1 26 2e 26 26 26 26]} #24641928
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=1cbd="-" map=[2e 1 26 2e 26 26 26 26]} #24641928 #24695878
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d04 size=9} #24696144
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #24696604
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24696604 #24698908
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7504 destPtr=1d04 size=0009} #24700682
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7504 destPtr=1d04 size=0009} #24700682 #24702360
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d04 size=9} #24696144 #24703296
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1d04 destPtr=7bb4 size=0009} #24703408
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1d04 destPtr=7bb4 size=0009} #24703408 #24705086
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d0d size=23} #24705280
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #24705740
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24705740 #24706410
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=750d destPtr=1d0d size=0023} #24708184
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=750d destPtr=1d0d size=0023} #24708184 #24710272
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d0d size=23} #24705280 #24711208
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=1d04="-" map=[2e 1 26 2e 26 26 26 26]} #24711346
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=1d04="-" map=[2e 1 26 2e 26 26 26 26]} #24711346 #24767514
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d30 size=9} #24767780
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #24768240
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24768240 #24768910
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7530 destPtr=1d30 size=0009} #24770684
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7530 destPtr=1d30 size=0009} #24770684 #24772362
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d30 size=9} #24767780 #24773298
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1d30 destPtr=7bb4 size=0009} #24773410
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1d30 destPtr=7bb4 size=0009} #24773410 #24775088
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d39 size=1b} #24775282
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #24775742
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24775742 #24776412
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7539 destPtr=1d39 size=001b} #24778186
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7539 destPtr=1d39 size=001b} #24778186 #24780160
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d39 size=1b} #24775282 #24781096
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=1d30="-" map=[2e 1 26 2e 26 26 26 26]} #24781234
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=1d30="-" map=[2e 1 26 2e 26 26 26 26]} #24781234 #24838814
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d54 size=9} #24839080
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #24839540
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24839540 #24840210
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7554 destPtr=1d54 size=0009} #24841984
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7554 destPtr=1d54 size=0009} #24841984 #24843662
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d54 size=9} #24839080 #24844598
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1d54 destPtr=7bb4 size=0009} #24844710
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1d54 destPtr=7bb4 size=0009} #24844710 #24848022
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1d5d size=5e} #24848216
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #24848676
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24848676 #24849346
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=755d destPtr=1d5d size=005e} #24851120
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=755d destPtr=1d5d size=005e} #24851120 #24854108
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1d5d size=5e} #24848216 #24855044
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=1d54="-" map=[2e 1 26 2e 26 26 26 26]} #24855182
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=1d54="-" map=[2e 1 26 2e 26 26 26 26]} #24855182 #24913760
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1dbb size=9} #24914026
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #24914486
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24914486 #24915156
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75bb destPtr=1dbb size=0009} #24916930
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75bb destPtr=1dbb size=0009} #24916930 #24918608
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1dbb size=9} #24914026 #24919544
:   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1dbb destPtr=7bb4 size=0009} #24919656
:   :   :   <---- t0 OS9$38 <F$Move> {srcTask=3 destTask=0 srcPtr=1dbb destPtr=7bb4 size=0009} #24919656 #24921334
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1dc4 size=1e} #24921528
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #24921988
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24921988 #24922658
:   :   :   {proc=4} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75c4 destPtr=1dc4 size=001e} #24924432
:   :   :   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=75c4 destPtr=1dc4 size=001e} #24924432 #24926508
:   :   :   <---- t0 OS9$89 <I$Read> {path=2 buf=1dc4 size=1e} #24921528 #24927444
:   :   {proc=4} t0: OS9$2e <F$VModul> {addr=1dbb="-" map=[2e 1 26 2e 26 26 26 26]} #24927582
:   :   :   <---- t0 OS9$2e <F$VModul> {addr=1dbb="-" map=[2e 1 26 2e 26 26 26 26]} #24927582 #24985982
:   :   {proc=4} t0: OS9$89 <I$Read> {path=2 buf=1de2 size=9} #24986248
:   :   :   {proc=4} t0: OS9$2f <F$Find64> {base=8200 id=2} #24986708
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24986708 #24987378
:   :   :   <-- ERROR: $d3(E$EOF    :End of File): OS9KERNEL0 OS9$89 <I$Read> {path=2 buf=1de2 size=9} #24986248 #24989804
:   :   {proc=3} t0: OS9$8f <I$Close> {path=2} #24989904
:   :   :   {proc=3} t0: OS9$2f <F$Find64> {base=8200 id=2} #24990352
:   :   :   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #24990352 #24991022
:   :   :   {proc=3} t0: OS9$29 <F$SRtMem> {size=100 start=7500} #24991484
:   :   :   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7500} #24991484 #24997718
:   :   :   {proc=3} t0: OS9$31 <F$Ret64> {block_num=3 address=8200} #24998184
:   :   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=3 address=8200} #24998184 #24998962
:   :   :   {proc=3} t0: OS9$81 <I$Detach> {8300} #24999232
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #24999794
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #24999794 #25000534
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #25000554
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #25000554 #25001294
:   :   :   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #25001314
:   :   :   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #25001314 #25002054
:   :   :   :   <---- t0 OS9$81 <I$Detach> {8300} #24999232 #25002404
:   :   :   {proc=3} t0: OS9$31 <F$Ret64> {block_num=2 address=8200} #25002434
:   :   :   :   <---- t0 OS9$31 <F$Ret64> {block_num=2 address=8200} #25002434 #25003212
:   :   :   <---- t0 OS9$8f <I$Close> {path=2} #24989904 #25003460
:   :   {proc=3} t0: OS9$4c <F$DelPrc> {} #25003592
:   :   :   {proc=3} t0: OS9$40 <F$DelTsk> {proc_desc=7600} #25004044
:   :   :   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7600} #25004044 #25004714
:   :   :   {proc=3} t0: OS9$29 <F$SRtMem> {size=200 start=7600} #25004740
:   :   :   :   <---- t0 OS9$29 <F$SRtMem> {size=200 start=7600} #25004740 #25007310
:   :   :   <---- t0 OS9$4c <F$DelPrc> {} #25003592 #25007566
:   :   {proc=2} t0: OS9$48 <F$LDDDXY> {} #25010292
:   :   :   <---- t0 OS9$48 <F$LDDDXY> {} #25010292 #25011138
:   :   {proc=2} t0: OS9$4d <F$ELink> {} #25011180
:   :   :   <---- t0 OS9$4d <F$ELink> {} #25011180 #25013958
:   :   <---- t0 OS9$01 <F$Load> {type/lang=78 filename='Shell'} #21849402 #25014290
:   {proc=2} t0: OS9$48 <F$LDDDXY> {} #25014456
:   :   <---- t0 OS9$48 <F$LDDDXY> {} #25014456 #25015582
:   {proc=2} t0: OS9$07 <F$Mem> {desired_size=1f00} #25015610
:   :   <---- t0 OS9$07 <F$Mem> {desired_size=1f00} #25015610 #25016410
:   {proc=3} t0: OS9$3f <F$AllTsk> {processDesc=7a00} #25016740
:   :   <---- t0 OS9$3f <F$AllTsk> {processDesc=7a00} #25016740 #25017634
:   {proc=3} t0: OS9$38 <F$Move> {srcTask=2 destTask=3 srcPtr=0000 destPtr=1efb size=0005} #25017776
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=3 srcPtr=0000 destPtr=1efb size=0005} #25017776 #25019348
:   {proc=3} t0: OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7bf4 destPtr=1eef size=000c} #25019414
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=3 srcPtr=7bf4 destPtr=1eef size=000c} #25019414 #25021194
:   {proc=3} t0: OS9$40 <F$DelTsk> {proc_desc=7800} #25021352
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7800} #25021352 #25022022
:   {proc=3} t0: OS9$29 <F$SRtMem> {size=200 start=7800} #25022048
:   :   <---- t0 OS9$29 <F$SRtMem> {size=200 start=7800} #25022048 #25024674
:   {proc=3} t0: OS9$40 <F$DelTsk> {proc_desc=7a00} #25024712
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7a00} #25024712 #25027016
:   {proc=1} t0: OS9$2c <F$AProc> {proc=7a00} #25027078
:   :   <---- t0 OS9$2c <F$AProc> {proc=7a00} #25027078 #25027802
:   {proc=1} t0: OS9$2d <F$NProc> {} #25027814
{proc=2"Shell"} t1: OS9$09 <F$Icpt> {routine=e06b storage=0000} #25036204
:   <---- t1 OS9$09 <F$Icpt> {routine=e06b storage=0000} #25036204 #25038862
{proc=2"Shell"} t1: OS9$0c <F$ID> {} #25038914
:   <---- t1 OS9$0c <F$ID> {} #25038914 #25041186
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25043806
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #25045134
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #25045134 #25045804
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=d511 destPtr=0214 size=0020} #25047078
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=d511 destPtr=0214 size=0020} #25047078 #25049064
:   <---- t1 OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25043806 #25050258
{proc=2"Shell"} t1: OS9$84 <I$Open> {f7f4='.'} #25051000
:   {proc=2} t0: OS9$30 <F$All64> {table=8200} #25052432
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #25052432 -> base $8200 blocknum $2 addr $8280 #25054822
:   {proc=2} t0: OS9$49 <F$LDABX> {} #25056478
:   :   <---- t0 OS9$49 <F$LDABX> {} #25056478 #25057332
:   {proc=1} t0: OS9$10 <F$PrsNam> {path='DD'} #25057518
:   :   <---- t0 OS9$10 <F$PrsNam> {path='DD'} #25057518 #25059286
:   {proc=1} t0: OS9$80 <I$Attach> {ec52='DD'} #25059328
:   :   {proc=1} t0: OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #25060544
:   :   :   <---- t0 OS9$34 <F$SLink> {"DD" type f0 name@ ec52 dat@ 640} #25060544 #25088444
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #25088576
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='RBLemma'} #25088576 -> addr $e704 entry $e72e #25120614
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #25120700
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='RBF'} #25120700 -> addr $9e77 entry $9e88 #25178786
:   :   <---- t0 OS9$80 <I$Attach> {ec52='DD'} #25059328 #25181124
:   {proc=2} t0: OS9$28 <F$SRqMem> {size=100} #25181620
:   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #25181620 -> size $100 addr $7900 #25191744
:   {proc=2} t0: OS9$30 <F$All64> {table=8200} #25191792
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #25191792 -> base $8200 blocknum $3 addr $82c0 #25194218
:   {proc=2} t0: OS9$49 <F$LDABX> {} #25195156
:   :   <---- t0 OS9$49 <F$LDABX> {} #25195156 #25196010
:   {proc=2} t0: OS9$10 <F$PrsNam> {path='a'} #25301164
:   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path='a'} #25301164 #25302926
:   {proc=2} t0: OS9$49 <F$LDABX> {} #25303078
:   :   <---- t0 OS9$49 <F$LDABX> {} #25303078 #25303932
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=f7f4 destPtr=82e0 size=0001} #25304256
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=f7f4 destPtr=82e0 size=0001} #25304256 #25305868
:   {proc=2} t0: OS9$10 <F$PrsNam> {path='"'} #25406910
:   :   <-- ERROR: $eb(E$BNam   :Bad Name): OS9KERNEL0 OS9$10 <F$PrsNam> {path='"'} #25406910 #25408432
:   <---- t1 OS9$84 <I$Open> {f7f4='.'} #25051000 -> path $3 #25410946
{proc=2"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #25411040
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #25412302
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #25412302 #25414606
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7900 destPtr=03f8 size=0020} #25464814
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7900 destPtr=03f8 size=0020} #25464814 #25466800
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #25411040 #25468664
{proc=2"Shell"} t1: OS9$89 <I$Read> {path=3 buf=03f8 size=20} #25468784
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #25470046
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #25470046 #25470716
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7920 destPtr=03f8 size=0020} #25474124
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=7920 destPtr=03f8 size=0020} #25474124 #25476110
:   <---- t1 OS9$89 <I$Read> {path=3 buf=03f8 size=20} #25468784 #25477974
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25478242
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #25479570
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #25479570 #25480240
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=ec52 destPtr=00b5 size=0020} #25481490
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=ec52 destPtr=00b5 size=0020} #25481490 #25483476
:   <---- t1 OS9$8d <I$GetStt> {path=3 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25478242 #25484670
{proc=2"Shell"} t1: OS9$10 <F$PrsNam> {path='DD'} #25484696
:   <---- t1 OS9$10 <F$PrsNam> {path='DD'} #25484696 #25487838
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=3} #25488074
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #25489386
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #25489386 #25490056
:   {proc=2} t0: OS9$29 <F$SRtMem> {size=100 start=7900} #25490518
:   :   <---- t0 OS9$29 <F$SRtMem> {size=100 start=7900} #25490518 #25494170
:   {proc=2} t0: OS9$31 <F$Ret64> {block_num=3 address=8200} #25494636
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=3 address=8200} #25494636 #25495414
:   {proc=2} t0: OS9$81 <I$Detach> {8300} #25495684
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #25496246
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=9e77 magic=87cd module='RBF'} #25496246 #25496986
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #25497006
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=e704 magic=87cd module='RBLemma'} #25497006 #25497746
:   :   {proc=1} t0: OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #25497766
:   :   :   <---- t0 OS9$02 <F$UnLink> {u=ec3c magic=87cd module='DD'} #25497766 #25498506
:   :   <---- t0 OS9$81 <I$Detach> {8300} #25495684 #25498856
:   {proc=2} t0: OS9$31 <F$Ret64> {block_num=2 address=8200} #25498886
:   :   <---- t0 OS9$31 <F$Ret64> {block_num=2 address=8200} #25498886 #25499664
:   <---- t1 OS9$8f <I$Close> {path=3} #25488074 #25500840
{proc=2"Shell"} t1: OS9$1c <F$SUser> {} #25505638
:   <---- t1 OS9$1c <F$SUser> {} #25505638 #25507892
{proc=2"Shell"} t1: OS9$82 <I$Dup> {$0} #25511790
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #25513192
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #25513192 #25513862
:   <---- t1 OS9$82 <I$Dup> {$0} #25511790 -> path $3 #25515118
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=0} #25515168
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #25516480
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #25516480 #25517150
:   <---- t1 OS9$8f <I$Close> {path=0} #25515168 #25519210
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25519594
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #25520922
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #25520922 #25521592
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=d511 destPtr=00b6 size=0020} #25522866
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=d511 destPtr=00b6 size=0020} #25522866 #25524852
:   <---- t1 OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25519594 #25526046
{proc=2"Shell"} t1: OS9$84 <I$Open> {00b5='/Term'} #25526078
:   {proc=2} t0: OS9$30 <F$All64> {table=8200} #25527408
:   :   <---- t0 OS9$30 <F$All64> {table=8200} #25527408 -> base $8200 blocknum $2 addr $8280 #25529798
:   {proc=2} t0: OS9$49 <F$LDABX> {} #25529872
:   :   <---- t0 OS9$49 <F$LDABX> {} #25529872 #25530638
:   {proc=2} t0: OS9$10 <F$PrsNam> {path=''} #25530716
:   :   <---- t0 OS9$10 <F$PrsNam> {path=''} #25530716 #25534396
:   {proc=2} t0: OS9$80 <I$Attach> {00b6='X'} #25534438
:   :   {proc=1} t0: OS9$34 <F$SLink> {"Term" type f0 name@ b6 dat@ 7a40} #25535654
:   :   :   <---- t0 OS9$34 <F$SLink> {"Term" type f0 name@ b6 dat@ 7a40} #25535654 #25575126
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='VTIO'} #25575258
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='VTIO'} #25575258 -> addr $b8dd entry $be5c #25629118
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='SCF'} #25629204
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='SCF'} #25629204 -> addr $b165 entry $b410 #25685482
:   :   <---- t0 OS9$80 <I$Attach> {00b6='X'} #25534438 #25687784
:   {proc=2} t0: OS9$10 <F$PrsNam> {path=''} #25689138
:   :   <---- t0 OS9$10 <F$PrsNam> {path=''} #25689138 #25691236
:   {proc=2} t0: OS9$28 <F$SRqMem> {size=100} #25691298
:   :   <---- t0 OS9$28 <F$SRqMem> {size=100} #25691298 -> size $100 addr $7900 #25701422
:   {proc=1} t0: OS9$80 <I$Attach> {d511='Term'} #25707530
:   :   {proc=1} t0: OS9$34 <F$SLink> {"Term" type f0 name@ d511 dat@ 640} #25708746
:   :   :   <---- t0 OS9$34 <F$SLink> {"Term" type f0 name@ d511 dat@ 640} #25708746 #25756572
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=e0 module/file='VTIO'} #25756704
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=e0 module/file='VTIO'} #25756704 -> addr $b8dd entry $be5c #25810564
:   :   {proc=1} t0: OS9$00 <F$Link> {type/lang=d0 module/file='SCF'} #25810650
:   :   :   <---- t0 OS9$00 <F$Link> {type/lang=d0 module/file='SCF'} #25810650 -> addr $b165 entry $b410 #25866928
:   :   <---- t0 OS9$80 <I$Attach> {d511='Term'} #25707530 #25869230
:   <---- t1 OS9$84 <I$Open> {00b5='/Term'} #25526078 -> path $0 #25872932
{proc=2"Shell"} t1: OS9$82 <I$Dup> {$1} #25873106
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #25874542
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #25874542 #25875212
:   <---- t1 OS9$82 <I$Dup> {$1} #25873106 -> path $4 #25876468
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=1} #25876518
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #25877830
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #25877830 #25878500
:   <---- t1 OS9$8f <I$Close> {path=1} #25876518 #25880620
{proc=2"Shell"} t1: OS9$82 <I$Dup> {$0} #25880652
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #25881986
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #25881986 #25882656
:   <---- t1 OS9$82 <I$Dup> {$0} #25880652 -> path $1 #25883912
{proc=2"Shell"} t1: OS9$82 <I$Dup> {$2} #25884058
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #25885528
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #25885528 #25886198
:   <---- t1 OS9$82 <I$Dup> {$2} #25884058 -> path $5 #25887454
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=2} #25887504
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #25888816
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #25888816 #25889486
:   <---- t1 OS9$8f <I$Close> {path=2} #25887504 #25893300
{proc=2"Shell"} t1: OS9$82 <I$Dup> {$1} #25893332
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #25894700
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #25894700 #25895370
:   <---- t1 OS9$82 <I$Dup> {$1} #25893332 -> path $2 #25896626
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=0 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25896766
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #25898094
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #25898094 #25898764
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=d511 destPtr=006e size=0020} #25900038
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=d511 destPtr=006e size=0020} #25900038 #25902024
:   <---- t1 OS9$8d <I$GetStt> {path=0 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25896766 #25903218
{proc=2"Shell"} t1: OS9$8a <I$Write> {{0}} #25903300
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #25904568
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #25904568 #25905238
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=e03c destPtr=79ff size=0001} #25906382
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=e03c destPtr=79ff size=0001} #25906382 #25907994
:   <---- t1 OS9$8a <I$Write> {{0}} #25903300 #25913120
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25913654
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #25914982
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #25914982 #25915652
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=d511 destPtr=0214 size=0020} #25916926
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=d511 destPtr=0214 size=0020} #25916926 #25918912
:   <---- t1 OS9$8d <I$GetStt> {path=1 e==SS.DevNm  : Return Device name (32-bytes at [X])} #25913654 #25920106
{proc=2"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #25923798
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #25925044
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #25925044 #25925714
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #25923798 #25928186
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #25928228
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #25929556
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #25929556 #25930226
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=82a0 destPtr=0124 size=0020} #25931106
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=82a0 destPtr=0124 size=0020} #25931106 #25933092
:   <---- t1 OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #25928228 #25934286
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=3} #25934526
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #25935838
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #25935838 #25936508
:   <---- t1 OS9$8f <I$Close> {path=3} #25934526 #25938748
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=4} #25938882
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #25940194
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #25940194 #25940864
:   <---- t1 OS9$8f <I$Close> {path=4} #25938882 #25943164
{proc=2"Shell"} t1: OS9$8f <I$Close> {path=5} #25943298
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=1} #25944610
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=1} #25944610 #25945280
:   {proc=2} t0: OS9$37 <F$GProcP> {id=01} #25946478
:   :   <---- t0 OS9$37 <F$GProcP> {id=01} #25946478 #25947186
:   <---- t1 OS9$8f <I$Close> {path=5} #25943298 #25948822
{proc=2"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #25962084
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #25963330
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #25963330 #25964000
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #25962084 #25966472
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #25966514
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #25967842
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #25967842 #25968512
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=82a0 destPtr=0124 size=0020} #25969392
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=82a0 destPtr=0124 size=0020} #25969392 #25971378
:   <---- t1 OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #25966514 #25972572
{proc=2"Shell"} t1: OS9$8a <I$Write> {
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #25974002
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #25974002 #25974672
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=e02e destPtr=79f2 size=000e} #25975816
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=e02e destPtr=79f2 size=000e} #25975816 #25977752
:   <---- t1 OS9$8a <I$Write> {
{proc=2"Shell"} t1: OS9$15 <F$Time> {buf=2da} #25991154
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=0028 destPtr=02da size=0006} #25992270
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=0028 destPtr=02da size=0006} #25992270 #25993876
:   <---- t1 OS9$15 <F$Time> {buf=2da} #25991154 #25995044
{proc=2"Shell"} t1: OS9$8c <I$WritLn> {} #25998350
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #25999618
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #25999618 #26000288
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=02e0 destPtr=79e0 size=0020} #26001430
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=02e0 destPtr=79e0 size=0020} #26001430 #26003416
:   <---- t1 OS9$8c <I$WritLn> {} #25998350 #26020568
{proc=2"Shell"} t1: OS9$1c <F$SUser> {} #26020730
:   <---- t1 OS9$1c <F$SUser> {} #26020730 #26022984
{proc=2"Shell"} t1: OS9$8c <I$WritLn> {} #26023194
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #26024462
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #26024462 #26025132
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0212 destPtr=79f2 size=000e} #26026268
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0212 destPtr=79f2 size=000e} #26026268 #26028116
:   <---- t1 OS9$8c <I$WritLn> {} #26023194 #26043822
{proc=2"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #26043906
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #26045152
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #26045152 #26045822
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.Relea  : Release device} #26043906 #26048294
{proc=2"Shell"} t1: OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #26048344
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #26049672
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #26049672 #26050342
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=82a0 destPtr=0124 size=0020} #26051222
:   :   <---- t0 OS9$38 <F$Move> {srcTask=0 destTask=2 srcPtr=82a0 destPtr=0124 size=0020} #26051222 #26053208
:   <---- t1 OS9$8d <I$GetStt> {path=0 0==SS.Opt    : Read/Write PD Options} #26048344 #26054402
{proc=2"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.Opt    : Read/Write PD Options} #26054562
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #26055808
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #26055808 #26056478
:   {proc=2} t0: OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0124 destPtr=82a0 size=001a} #26057224
:   :   <---- t0 OS9$38 <F$Move> {srcTask=2 destTask=0 srcPtr=0124 destPtr=82a0 size=001a} #26057224 #26059164
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.Opt    : Read/Write PD Options} #26054562 #26062350
{proc=2"Shell"} t1: OS9$8e <I$SetStt> {path=0 SS.SSig   : Send signal on data ready} #26062394
:   {proc=2} t0: OS9$2f <F$Find64> {base=8200 id=2} #26063640
:   :   <---- t0 OS9$2f <F$Find64> {base=8200 id=2} #26063640 #26064310
:   <---- t1 OS9$8e <I$SetStt> {path=0 SS.SSig   : Send signal on data ready} #26062394 #26066850
{proc=2"Shell"} t1: OS9$0a <F$Sleep> {ticks=0000} #26066940
:   {proc=2} t0: OS9$40 <F$DelTsk> {proc_desc=7a00} #26068250
:   :   <---- t0 OS9$40 <F$DelTsk> {proc_desc=7a00} #26068250 #26068920
:   {proc=2} t0: OS9$2d <F$NProc> {} #26068974
```


## END
