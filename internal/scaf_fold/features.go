package scaf_fold

import (
	"fmt"
	"strings"
)

// SetFeatures 选择生成的组件；运行时启用状态由项目配置决定。
func (d *TemplateData) SetFeatures(value string) error {
	for part := range strings.SplitSeq(value, ",") {
		switch strings.ToLower(strings.TrimSpace(part)) {
		case "":
		case "redis":
			d.Redis = true
		case "cron":
			d.Cron = true
		case "lark":
			d.Lark = true
		case "prometheus-query":
			d.Prometheus = true
		case "jwt":
			d.JWT = true
		default:
			return fmt.Errorf("unknown component %q", part)
		}
	}
	return d.Validate()
}
