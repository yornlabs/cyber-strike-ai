package multiagent

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/components/tool"
)

// injectToolNamesOnlyInstruction prepends a compact tool-name-only section into
// the system instruction so the model can reference current callable names.
func injectToolNamesOnlyInstruction(ctx context.Context, instruction string, tools []tool.BaseTool) string {
	names := collectToolNames(ctx, tools)
	if len(names) == 0 {
		return strings.TrimSpace(instruction)
	}
	hasToolSearch := false
	for _, n := range names {
		if strings.EqualFold(strings.TrimSpace(n), "tool_search") {
			hasToolSearch = true
			break
		}
	}

	var sb strings.Builder
	sb.WriteString("以下是当前会话中可调用的工具名称列表（仅名称，无参数定义）：\n")
	for _, name := range names {
		sb.WriteString("- ")
		sb.WriteString(name)
		sb.WriteByte('\n')
	}
	sb.WriteString("\n使用规则：\n")
	sb.WriteString("1) 上述仅为名称列表，不包含参数定义。\n")
	if hasToolSearch {
		sb.WriteString("2) 在调用具体工具前，应先使用 tool_search 查看工具详情与参数要求，再发起调用。\n")
	} else {
		sb.WriteString("2) 调用具体工具前，请先确认该工具的参数要求；不确定时先澄清再调用。\n")
	}
	sb.WriteString("3) 不要臆造不存在的工具名。\n\n")
	if s := strings.TrimSpace(instruction); s != "" {
		sb.WriteString(s)
	}
	return sb.String()
}

func collectToolNames(ctx context.Context, tools []tool.BaseTool) []string {
	if len(tools) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(tools))
	out := make([]string, 0, len(tools))
	for _, t := range tools {
		if t == nil {
			continue
		}
		info, err := t.Info(ctx)
		if err != nil || info == nil {
			continue
		}
		name := strings.TrimSpace(info.Name)
		if name == "" {
			continue
		}
		key := strings.ToLower(name)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, name)
	}
	return out
}

