package bootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

var flags = NewCommandFlags()

// NewRootCmd 创建根命令并绑定命令行参数和执行函数。
func NewRootCmd(f *CommandFlags, runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "A microservice server application",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			f.Init()
		},
		RunE: runE,
	}
	f.AddFlags(cmd)
	return cmd
}

// CommandFlags 命令传参
type CommandFlags struct {
	Conf       string // 引导配置文件路径，默认为：../../configs
	Env        string // 开发环境：dev、debug……
	ConfigHost string // 远程配置服务端地址
	ConfigType string // 远程配置服务端类型
	Daemon     bool   // 是否转为守护进程
}

func NewCommandFlags() *CommandFlags {
	return &CommandFlags{
		Conf:       "../../configs",
		Env:        "dev",
		ConfigHost: "127.0.0.1:8848",
		ConfigType: "nacos",
		Daemon:     false,
	}
}

func (f *CommandFlags) Init() {
	if f.Daemon {
		BeDaemon("-d")
	}
}

// AddFlags 将 flags 绑定到传入的 cobra.Command（通常是 root command）。
func (f *CommandFlags) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&f.Conf, "conf", "c", f.Conf, "config path, eg: -conf ../../configs")
	cmd.PersistentFlags().StringVarP(&f.Env, "env", "e", f.Env, "runtime environment, eg: -env dev")
	cmd.PersistentFlags().StringVarP(&f.ConfigHost, "chost", "s", f.ConfigHost, "config server host, eg: -chost 127.0.0.1:8500")
	cmd.PersistentFlags().StringVarP(&f.ConfigType, "ctype", "t", f.ConfigType, "config server type, eg: -ctype consul")
	cmd.PersistentFlags().BoolVarP(&f.Daemon, "daemon", "d", f.Daemon, "run app as a daemon with -d or --daemon")
}

// BeDaemon 将当前进程转为守护进程（尝试启动脱离的子进程并退出父进程）
func BeDaemon(arg string) {
	childArgs := stripSlice(os.Args, arg)
	cmd := subProcess(childArgs)
	if cmd == nil || cmd.Process == nil {
		// 启动失败，继续在当前进程运行
		return
	}
	fmt.Printf("[*] Daemon started in PID: %d\n", cmd.Process.Pid)
	os.Exit(0)
}

// stripSlice 从字符串切片中移除指定元素
func stripSlice(slice []string, element string) []string {
	for i := 0; i < len(slice); {
		if slice[i] == element && i != len(slice)-1 {
			slice = append(slice[:i], slice[i+1:]...)
		} else if slice[i] == element && i == len(slice)-1 {
			slice = slice[:i]
		} else {
			i++
		}
	}
	return slice
}

// Unix 平台的 subProcess 实现，使用 Setsid
func subProcess(args []string) *exec.Cmd {
	if len(args) == 0 {
		return nil
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err == nil {
		cmd.Stdin = devNull
		cmd.Stdout = devNull
		cmd.Stderr = devNull
	} else {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err = cmd.Start(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[-] Error starting daemon: %s\n", err)
		if devNull != nil {
			_ = devNull.Close()
		}
		return nil
	}
	if devNull != nil {
		_ = devNull.Close()
	}
	return cmd
}
