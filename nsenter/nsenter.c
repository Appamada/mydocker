#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <sched.h>
#include <unistd.h>

void __attribute__((constructor)) enter_namespace() {

    const char *docker_pid;
    const char *docker_command;

    if (getenv("docker_pid") == NULL) {
        printf("cannot get env docker_pid");
    } else {
        docker_pid = getenv("docker_pid");
        printf("docker_pid: %s", docker_pid);
    };


    if (getenv("docker_command") == NULL) {
        printf("cannot get env docker_command");
    } else {
        docker_pid = getenv("docker_command");
        printf("docker_command: %s", docker_command);
    };

    char nspath [256] = {0};
    char namespace[][4] = {"pid", "mnt", "ipc", "uts", "net"};

    for (int i=0; i<5; i++){
        sprintf(nspath,"proc/%d/ns/%s", atoi(docker_pid),namespace[i]);

        printf("namespace: %s", nspath);

        // setns - reassociate thread with a namespace
        // int setns(int fd, int nstype);
        // Given a file descriptor referring to a namespace, reassociate the calling thread with that namespace.

        // The  fd  argument  is a file descriptor referring to one of the namespace entries in a /proc/[pid]/ns/ directory; see namespaces(7) for further information on
        // /proc/[pid]/ns/. 

        // The nstype argument specifies which type of namespace the calling thread may be reassociated with.
       
        int fd = open(nspath, O_RDONLY);
        if (setns(fd, 0) == 0) {
            fprintf(stdout,"setns success on namespace: %s\n", namespace[i]);
        } else {
            fprintf(stderr,"setns failed on namespace: %s\n", namespace[i]);
        }
        close(fd);
    }

    system(docker_command);
    exit(0);
    return;
}