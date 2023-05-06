#ifndef _PVM_WRAPPER_C
#define _PVM_WRAPPER_C

#include <stdio.h>
#include "pvm_wrapper.h"

int pvm_catchout_stdout()
{
    return pvm_catchout(stdout);
}

void *ptr_at(void **ptr, int idx)
{
    return ptr[idx];
}

pvmhostinfo* hostinfo_ptr() {
    struct pvmhostinfo *hostp = 0;
    return hostp;
}

char** str_ptr() {
    char* ptr = 0;
    return &ptr;
}


void unwrap_hostinfo(pvmhostinfo* hostinfo, int* tid, char** name, char** arch, int* speed) {
    *tid = hostinfo->hi_tid;
    *name = hostinfo->hi_name;
    *arch = hostinfo->hi_arch;
    *speed = hostinfo->hi_speed;

    printf("%s\n", hostinfo->hi_name);
}

// pvm_packf static bindings
int pvm_packf_string(char *fmt, char *arg)
{
    return pvm_packf(fmt, arg);
}

// pvm_unpackf static bindings
int pvm_unpackf_string(char *fmt, char *arg)
{
    return pvm_unpackf(fmt, arg);
}

#endif