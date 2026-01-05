// fake_hex.c
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <signal.h>
#include <sys/stat.h>
#include <sys/ioctl.h>  // <-- 新增：用于 ioctl 和 winsize

volatile sig_atomic_t keep_running = 1;

void signal_handler(int sig) {
    keep_running = 0;
}

int get_terminal_width() {
    struct winsize w;
    if (ioctl(STDOUT_FILENO, TIOCGWINSZ, &w) == 0 && w.ws_col > 0) {
        return w.ws_col;
    }
    return 80; // fallback
}

int main(int argc, char *argv[]) {
    if (argc != 3) {
        fprintf(stderr, "usage: %s <audio_file> <tab_name>\n", argv[0]); // <-- 修复：补 argv[0]
        exit(1);
    }

    const char *audio_file = argv[1];
    const char *tab_name = argv[2];

    // 注册信号处理器（注意：宏名大写！）
    signal(SIGTERM, signal_handler);
    signal(SIGINT,  signal_handler);

    // 检查文件是否存在
    struct stat st;
    if (stat(audio_file, &st) != 0) {
        perror("stat");
        exit(1);
    }

    // 构造 xxd 命令
    size_t cmd_len = strlen(audio_file) + 64;
    char *cmd = malloc(cmd_len);
    snprintf(cmd, cmd_len, "xxd -c 16 \"%s\"", audio_file);

    FILE *fp = popen(cmd, "r"); // <-- 修复：FILE 而不是 file
    if (!fp) {
        perror("popen xxd failed");
        free(cmd);
        exit(1);
    }

    char line[512];
    while (keep_running && fgets(line, sizeof(line), fp)) {
        // 清除行尾换行（如果存在）
        size_t len = strlen(line);
        if (len > 0 && line[len-1] == '\n') {
            line[len-1] = '\0';
        }

        int cols = get_terminal_width();

        if ((int)strlen(line) >= cols) {
            if (cols > 0) {
                line[cols] = '\0';  // 注意：这里不是 cols-1，因为后面要加 \n，但你仍想占满一行
            } else {
                line[0] = '\0';
            }
        }

        char status[256];

        snprintf(status, sizeof(status), "▶正在播放 %s       按q退出  p暂停/播放  x跳过此曲", tab_name);

        // 截断到终端宽度（防止换行）
        if ((int)strlen(status) >= cols) {
            if (cols > 1) {
                status[cols - 1] = '\0';
            } else {
                status[0] = '\0';
            }
        }

        // 打印 hex 行（带换行）
        printf("\r%s\n", line);
        // 打印状态栏（覆盖当前行，不换行，用空格清尾）
        printf("\r%-*s", cols, status);
        fflush(stdout);

        usleep(30000); // 30ms
    }

    pclose(fp);
    free(cmd);
    return 0;
}
