```bash

# build docker image
git clone https://github.com/vitalyisaev2/pprof_cgo_issue.git
docker build -t pprof_cgo_issue

# run docker container (omit --privileged if you already have core_pattern configured)
docker run -it --name=pprof_cgo_issue --privileged pprof_cgo_issue bash

# setup core_pattern if neccessary
echo '/tmp/core.%h.%e.%t' > /proc/sys/kernel/core_pattern
ulimit -c unlimited

# launch process
./pprof_cgo_issue
2019/04/30 12:14:40 profile: cpu profiling enabled, cpu.pprof
Segmentation fault (core dumped)

# tap into core dump
gdb ./pprof_cgo_issue /tmp/core.d685d2a145b4.pprof_cgo_issue.1556626480.20
(gdb) bt
#0  0x00007fb47640c246 in ?? () from /lib64/libgcc_s.so.1
#1  0x00007fb47640cefd in _Unwind_Backtrace () from /lib64/libgcc_s.so.1
#2  0x00000000007322bc in cgoTraceback (parg=0xc0002dd970, parg@entry=<error reading variable: value has been optimized out>) at traceback.c:82
#3  0x00000000007364e6 in x_cgo_callers (sig=27, info=0xc0002ddaf0, context=0xc0002dd9c0, cgoTraceback=<optimized out>, cgoCallers=<optimized out>, sigtramp=0x46a2f0 <runtime.sigtramp>) at gcc_traceback.c:22
#4  <signal handler called>
#5  0x00007fb476614232 in ?? () from /lib64/libcrypto.so.1.1
Backtrace stopped: Cannot access memory at address 0xa7382544026add44

```
