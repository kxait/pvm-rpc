#ifndef _PVM_WRAPPER_H
#define _PVM_WRAPPER_H

#include <pvm3.h>

typedef struct pvmhostinfo pvmhostinfo;

int pvm_catchout_stdout();

void *ptr_at(void **ptr, int idx);

int pvm_packf_string(char *fmt, char *arg);

int pvm_unpackf_string(char *fmt, char *arg);

pvmhostinfo* hostinfo_ptr();

char* str_ptr();

void unwrap_hostinfo(pvmhostinfo* hostinfo, int* tid, char* name, char* arch, int* speed);

#endif