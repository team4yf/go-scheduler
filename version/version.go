package version

import "time"

//VERSION 项目的版本
const VERSION = "beta"

//BuildAt 编译打包的时间
var BuildAt = time.Now().Format("20060102.15.04")
