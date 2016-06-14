#include <stdio.h>
#include <unistd.h>
#include <time.h>
#include <stdint.h>
#include <pthread.h>

int p[2];
int n = 100000;

void *
sender(void *arg)
{
        struct timespec s;
        struct timespec delay = {0, 10000};

        (void)arg;

        for (int i = 0; i < n; i++) {
                nanosleep(&delay, NULL);
                clock_gettime(CLOCK_MONOTONIC, &s);
                write(p[1], &s, sizeof(s));
        }
        return NULL;
}

int
main(void)
{
        pipe(p);

#ifdef THREAD
        pthread_t sendthread;

        pthread_create(&sendthread, NULL, sender, NULL);
#else
        if (fork() == 0) {
                sender(NULL);
                return 0;
        }
#endif

        int bucket[64] = {0};

        struct timespec s, s2;

        for (int i = 0; i < n; i++) {
                read(p[0], &s, sizeof(s));
                clock_gettime(CLOCK_MONOTONIC, &s2);
                int64_t nsec = (s2.tv_sec - s.tv_sec) * 1000000000 +
                        (s2.tv_nsec - s.tv_nsec);
                int order = 63 - __builtin_clzl(nsec);

                bucket[order] += 1;
        }

        for (int i = 0; i < 64; i++) {
                printf("%d\t%d\n", i, bucket[i]);
        }
}
