/**
 * ImindLab
 *
 * Create by SongLi on 2023/12/28
 * Copyright © 2023 imind.tech All rights reserved.
 */

package util

import "time"

// 优化 https://github.com/uber-go/guide/blob/master/style.md Avoid Mutable Globals

type CacheTool struct {
	RandExpire func(int64) time.Duration
}

func NewCacheTool() CacheTool {
	return CacheTool{RandExpire: RandDuration}
}
