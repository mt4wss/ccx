//go:build linux

package updater

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

// platformInstall 替换当前 AppImage 文件，再 exec 新的。
//
// 流程：
//  1. 检测当前可执行路径（os.Executable）
//  2. 把下载的新 AppImage 复制到 <exe>.new
//  3. chmod +x
//  4. 启动 nohup 子脚本（完全脱离父进程）：
//     a. for 循环重试 mv（最多 30 次 × 0.5s = 15s）
//     b. mv 成功后 touch 信号文件
//     c. exec 新程序
//  5. 父进程轮询信号文件（200ms 间隔，20s 超时）
//  6. 确认 mv 完成后 os.Exit(0)
//
// 注意：仅支持 AppImage 形态；如果用户用 deb/rpm 安装则跳过自动替换。
func platformInstall(localPath string) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取当前可执行路径失败: %w", err)
	}

	// 简单启发式：可执行不在 /usr/* 才认为是 AppImage 直装模式
	if strings.HasPrefix(exePath, "/usr/") {
		return fmt.Errorf("检测到系统安装版本，请通过包管理器更新")
	}

	newPath := exePath + ".new"
	signalFile := exePath + ".updated"
	if err := copyFile(localPath, newPath); err != nil {
		return fmt.Errorf("复制新版本失败: %w", err)
	}
	if err := os.Chmod(newPath, 0o755); err != nil {
		return fmt.Errorf("设置可执行权限失败: %w", err)
	}

	// 清理旧信号文件
	os.Remove(signalFile)

	// 子脚本：重试 mv，成功后 touch 信号文件，然后 exec 新程序
	script := fmt.Sprintf(`
for i in $(seq 1 30); do
    mv -f %q %q 2>/dev/null && break
    sleep 0.5
done
touch %q
exec %q
`, newPath, exePath, signalFile, exePath)

	cmd := exec.Command("nohup", "sh", "-c", script)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动替换脚本失败: %w", err)
	}
	// 脱离子进程，避免 shell 退出时 SIGHUP 传导
	_ = cmd.Process.Release()

	// 父进程轮询信号文件，确认 mv 完成
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		if _, statErr := os.Stat(signalFile); statErr == nil {
			// mv 已完成，安全退出
			os.Exit(0)
		}
		time.Sleep(200 * time.Millisecond)
	}

	// 超时兜底：延迟退出让 mv 有机会完成
	time.Sleep(2 * time.Second)
	os.Exit(0)
	return nil // unreachable
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
