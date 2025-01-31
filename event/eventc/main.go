/*
 * This file is part of Golaxy Distributed Service Development Framework.
 *
 * Golaxy Distributed Service Development Framework is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * Golaxy Distributed Service Development Framework is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with Golaxy Distributed Service Development Framework. If not, see <http://www.gnu.org/licenses/>.
 *
 * Copyright (c) 2024 pangdogs.
 */

// Package eventc 使用 go:generate 功能，在编译前自动化生成代码
/*
	- 可以生成事件（Event）与事件表（Event Table）代码。
	- 用于生成事件代码时，在事件定义代码源文件（.go）头部，添加以下注释：
		//go:generate go run git.golaxy.org/core/event/eventc event
	- 用于生成事件表代码时，在事件定义代码源文件（.go）头部，添加以下注释：
		//go:generate go run git.golaxy.org/core/event/eventc eventtab --name={事件表名称}
	- 在 cmd 控制台中，进入事件定义代码源文件（.go）的目录，输入 go generate 指令即可生成代码，此外也可以使用 IDE 提供的 go generate 功能。
	- 编译本包并执行 eventc --help ，可以查看命令行参数，通过参数可以调整生成的代码。
*/
package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func main() {
	// 基础选项
	rootCmd := &cobra.Command{
		Short: "事件代码生成工具。",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
			{
				declFile := viper.GetString("decl_file")
				if declFile == "" {
					panic("[--decl_file]值不能为空")
				}
				stat, err := os.Stat(declFile)
				if err != nil {
					panic(fmt.Errorf("[--decl_file]文件错误，%s", err))
				}
				if stat.IsDir() {
					panic("[--decl_file]文件错误，不能为文件夹")
				}
			}
			{
				eventRegexp := viper.GetString("event_regexp")
				if eventRegexp == "" {
					panic("[--event_regexp]值不能为空")
				}
			}
			{
				packageEventAlias := viper.GetString("package_event_alias")
				if packageEventAlias == "" {
					panic("[--package_event_alias]值不能为空")
				}
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   true,
			DisableNoDescFlag:   true,
			DisableDescriptions: true,
		},
	}
	rootCmd.PersistentFlags().String("decl_file", os.Getenv("GOFILE"), "事件定义文件（.go）。")
	rootCmd.PersistentFlags().String("event_regexp", "^[eE]vent.+", "匹配事件定义时，使用的正则表达式。")
	rootCmd.PersistentFlags().String("package_event_alias", "event", fmt.Sprintf("导入Golaxy框架的（%s）包时使用的别名。", packageEventPath))
	rootCmd.PersistentFlags().Bool("copyright", true, "Golaxy分布式服务开发框架版权信息。")

	// 生成事件代码相关选项
	eventCmd := &cobra.Command{
		Use:   "event",
		Short: "根据定义的事件，生成事件代码。（支持定义选项：+event-gen:export=[0,1]&auto=[0,1]）",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
			loadDeclFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			genEvent()
		},
	}
	eventCmd.Flags().Bool("default_export", true, "生成事件代码时，发送事件的代码的可见性，事件定义选项+event-gen:export=[0,1]可以覆盖此配置。")
	eventCmd.Flags().Bool("default_auto", true, "生成事件代码时，是否生成简化绑定事件的代码，事件定义选项+event-gen:auto=[0,1]可以覆盖此配置。")

	// 生成事件表代码相关选项
	eventTabCmd := &cobra.Command{
		Use:   "eventtab",
		Short: "根据定义的事件，生成事件表代码。（支持定义选项：+event-tab-gen:recursion=[allow,disallow,discard,truncate,deepest]）",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
			{
				pkg := viper.GetString("package")
				if pkg == "" {
					panic("[eventtab --package]值不能为空")
				}
			}
			{
				dir := viper.GetString("dir")
				if dir == "" {
					panic("[eventtab --dir]值不能为空")
				}
			}
			{
				name := viper.GetString("name")
				if name == "" {
					panic("[eventtab --name]值不能为空")
				}
			}
			loadDeclFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			genEventTab()
		},
	}
	eventTabCmd.Flags().String("package", os.Getenv("GOPACKAGE"), "生成事件表代码，使用的包名。")
	eventTabCmd.Flags().String("dir", filepath.Dir(os.Getenv("GOFILE")), "生成事件表代码时，输出代码文件（.go）的目录。")
	eventTabCmd.Flags().String("name", defaultEventTab(), "生成的事件表名称。")

	rootCmd.AddCommand(eventCmd, eventTabCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
