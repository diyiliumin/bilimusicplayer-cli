// watchdog.c
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/wait.h>

int main(void)
{
    pid_t pid = fork();
    if (pid < 0) { perror("fork"); return 1; }

    if (pid == 0) {               /* 子进程 */
        execl("cmd/tui/mytui", "mytui", (char *)NULL);
        perror("execl mytui");    /* 失败才到这儿 */
        _exit(127);
    }

    wait(NULL);                   /* 等 mytui 结束 */
    system("pkill -x fake_hex");   // -x 表示精确匹配进程名（推荐！）
    usleep(200000);               /* 0.2 s 兜底延迟 */
//    system("pkill -f play");
    system("pkill -x ffplay");     // ffplay 是实际播放进程，比 -f play 更安全
    return 0;
}
