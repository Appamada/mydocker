package nsenter

/*
#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <sched.h>
#include <unistd.h>

void __attribute__((constructor)) enter_namespace(void) {

    char *container_pid;
	container_pid=getenv("container_pid");

    if(container_pid){
        //fprintf(stdout, "got container_pid=%s\n",container_pid);
    }else{
        //fprintf(stdout, "missing container_pid env, skip nsenter\n",container_pid);
        return;
    }

	char *container_command;
	container_command=getenv("container_command");
    if(container_command){
        //fprintf(stdout, "got container_command=%s\n",container_command);
    }else{
        //fprintf(stdout, "missing container_command env, skip nsenter\n",container_command);
        return;
    }


    char nspath[256];
    char *namespace[] = {"pid", "mnt", "ipc", "uts", "net"};

    for (int i=0; i<5; i++){
        sprintf(nspath,"/proc/%s/ns/%s", container_pid,namespace[i]);

        // setns - reassociate thread with a namespace
        // int setns(int fd, int nstype);
        // Given a file descriptor referring to a namespace, reassociate the calling thread with that namespace.

        // The  fd  argument  is a file descriptor referring to one of the namespace entries in a /proc/[pid]/ns/ directory; see namespaces(7) for further information on
        // /proc/[pid]/ns/.

        // The nstype argument specifies which type of namespace the calling thread may be reassociated with.

        int fd = open(nspath, O_RDONLY);

		// if (fd < 0) {
		// 	fprintf(stderr,"open failed on namespace: %s\n", namespace[i]);
		// 	continue;
		// }

        if (setns(fd, 0) == 0) {
            // fprintf(stdout,"setns success on namespace: %s\n", namespace[i]);
        } else {
            // fprintf(stderr,"setns failed on namespace: %s\n", namespace[i]);
        }
        close(fd);
    }

    int res = system(container_command);
    exit(0);
    return;
}
*/
import "C"
