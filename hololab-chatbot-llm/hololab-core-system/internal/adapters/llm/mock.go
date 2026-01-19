package llm

import (
	"context"
	"hololab-core-system/internal/ports"
	"regexp"
	"strings"
)

type MockProvider struct{}

func (p MockProvider) Chat(ctx context.Context, systemPrompt string, history []ports.ChatMsg, userText string) (string, error) {
	_ = ctx
	_ = history

	u := strings.TrimSpace(userText)
	if u == "" {
		return bulletsVN(
			"Mình chưa nhận được nội dung câu hỏi.",
			"Bạn nhập rõ yêu cầu giúp mình nhé.",
		), nil
	}
	lu := strings.ToLower(u)

	if isInjectionAttempt(lu) {
		return bulletsVN(
			"Mình không thể đáp ứng yêu cầu đó.",
			"Bạn hãy hỏi theo đúng vai trò và mục tiêu của nhân vật nhé.",
		), nil
	}

	name, job, bio := parseIdentity(systemPrompt)
	_ = parseStyle(systemPrompt)
	ak, akSpecified := parseAllowedKnowledge(systemPrompt)

	if containsAny(lu, "bạn là ai", "ban la ai", "bạn giúp được gì", "ban giup duoc gi", "who are you", "what can you do") {
		head := "Mình là chatbot theo persona."
		if name != "" {
			head = "Mình là **" + name + "**."
		}
		out := []string{head}
		if job != "" {
			out = append(out, "Nghề nghiệp/Vai trò: **"+job+"**.")
		}
		if strings.TrimSpace(bio) != "" {
			out = append(out, "Tiểu sử: "+shorten(bio, 180))
		}
		out = append(out, "Mình sẽ trả lời đúng theo vai trò và phong cách đã mô tả.")
		return bulletsVN(out...), nil
	}

	if containsAny(lu, "short plan", "kế hoạch", "ke hoach", "plan related to your occupation", "related to your occupation") {
		title := "Kế hoạch ngắn theo vai trò hiện tại"
		if job != "" {
			title = "Kế hoạch ngắn (" + job + ")"
		}

		context := "(chưa có thêm bối cảnh)"
		if strings.TrimSpace(bio) != "" {
			context = shorten(bio, 140)
		}

		return bulletsVN(
			"**"+title+":**",
			"Bối cảnh: "+context,
			"1) Xác định mục tiêu: bạn muốn đạt kết quả gì trong 1–2 tuần tới?",
			"2) Chia việc theo 3 nhóm: (a) việc cốt lõi theo vai trò, (b) việc hỗ trợ, (c) việc chuẩn hoá/kiểm soát chất lượng.",
			"3) Checklist đầu ra: mỗi ngày 1–3 đầu việc nhỏ có thể hoàn thành.",
			"4) Tiêu chí chất lượng: “đúng – đủ – dễ kiểm tra”.",
			"5) Review cuối tuần: cái gì hiệu quả / cái gì cần bỏ.",
			"Bạn cho mình 2 thông tin: (a) mục tiêu cụ thể là gì? (b) ràng buộc thời gian/ngân sách?",
		), nil
	}

	if looksLikeFactualExact(lu) {
		if akSpecified {
			if hit, snippet := akMatch(u, ak); hit {
				return bulletsVN(
					"Mình trả lời dựa trên thông tin bạn đã cung cấp cho nhân vật:",
					snippet,
				), nil
			}
			return bulletsVN(
				"Mình chưa có đủ thông tin để trả lời chính xác câu hỏi này.",
				"Bạn bổ sung thêm dữ liệu (facts/nguồn) hoặc cho mình thêm bối cảnh để mình trả lời cụ thể hơn nhé.",
			), nil
		}

		return bulletsVN(
			"Mình chưa có đủ thông tin để trả lời dạng số liệu/nguồn chính xác.",
			"Bạn bổ sung thêm dữ liệu (facts/nguồn) hoặc bối cảnh để mình trả lời cụ thể hơn nhé.",
		), nil
	}

	if akSpecified {
		if hit, snippet := akMatch(u, ak); hit {
			return bulletsVN(
				"Mình trả lời dựa trên thông tin bạn đã cung cấp cho nhân vật:",
				snippet,
			), nil
		}
		return bulletsVN(
			"Mình chưa có đủ thông tin để trả lời chắc chắn cho câu hỏi này.",
			"Bạn cho mình thêm 1–2 chi tiết (bối cảnh / mục tiêu / ràng buộc) để mình trả lời chính xác hơn nhé.",
		), nil
	}

	return bulletsVN(
		"Mình hiểu câu hỏi của bạn: \""+u+"\".",
		"Mình có thể trả lời hướng dẫn chung theo persona, và sẽ tránh bịa số liệu/nguồn.",
		"Bạn muốn đầu ra dạng gì? (checklist / kế hoạch / hướng dẫn từng bước)",
	), nil
}

func parseIdentity(systemPrompt string) (name, job, bio string) {
	reName := regexp.MustCompile(`(?m)^\s*-\s*Name:\s*(.*)\s*$`)
	reJob := regexp.MustCompile(`(?m)^\s*-\s*Occupation:\s*(.*)\s*$`)

	if m := reName.FindStringSubmatch(systemPrompt); len(m) == 2 {
		name = strings.TrimSpace(m[1])
	}
	if m := reJob.FindStringSubmatch(systemPrompt); len(m) == 2 {
		job = strings.TrimSpace(m[1])
	}

	reBioBlock := regexp.MustCompile(`(?is)^\s*-\s*Bio:\s*(.*?)\n\n(?:STYLE\s*/\s*TONE|BEHAVIOR\s+RULES|RULES)`)
	if m := reBioBlock.FindStringSubmatch(systemPrompt); len(m) == 2 {
		bio = strings.TrimSpace(m[1])
	} else {
		reBioLine := regexp.MustCompile(`(?m)^\s*-\s*Bio:\s*(.*)\s*$`)
		if m2 := reBioLine.FindStringSubmatch(systemPrompt); len(m2) == 2 {
			bio = strings.TrimSpace(m2[1])
		}
	}

	return
}

func parseStyle(systemPrompt string) string {
	re := regexp.MustCompile(`(?is)STYLE\s*/\s*TONE(?:\s*\(.*?\))?\s*:\s*(.*?)\n\n`)
	m := re.FindStringSubmatch(systemPrompt)
	if len(m) == 2 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

func parseAllowedKnowledge(systemPrompt string) (knowledge string, specified bool) {
	reNot := regexp.MustCompile(`(?is)ALLOWED\s+KNOWLEDGE\s*:\s*Not specified\.`)
	if reNot.FindStringIndex(systemPrompt) != nil {
		return "", false
	}

	re := regexp.MustCompile(`(?is)(ALLOWED\s+KNOWLEDGE(?:\s*\(.*?\))?\s*:\s*)(.*)$`)
	m := re.FindStringSubmatch(systemPrompt)
	if len(m) == 3 {
		block := strings.TrimSpace(m[2])

		cutters := []*regexp.Regexp{
			regexp.MustCompile(`(?is)\n\nBEHAVIOR\s+RULES\s*:.*$`),
			regexp.MustCompile(`(?is)\n\nRULES\s*\(.*?\)\s*:.*$`),
			regexp.MustCompile(`(?is)\n\nIf\s+the\s+user\s+asks.*$`),
		}
		for _, cut := range cutters {
			block = cut.ReplaceAllString(block, "")
			block = strings.TrimSpace(block)
		}

		if block != "" {
			return block, true
		}
	}

	re2 := regexp.MustCompile(`(?is)You\s+may\s+rely\s+on\s+the\s+following\s+knowledge\s*:\s*(.*)$`)
	m2 := re2.FindStringSubmatch(systemPrompt)
	if len(m2) == 2 {
		block := strings.TrimSpace(m2[1])
		if block != "" {
			return block, true
		}
	}

	return "", false
}

func akMatch(question string, ak string) (bool, string) {
	q := strings.ToLower(question)
	akLower := strings.ToLower(ak)

	words := splitWords(q)
	keywords := make([]string, 0, len(words))
	for _, w := range words {
		if len(w) < 4 {
			continue
		}
		if isStopWord(w) {
			continue
		}
		keywords = append(keywords, w)
	}

	for _, kw := range keywords {
		if strings.Contains(akLower, kw) {
			return true, "• " + strings.Join(firstNonEmptyLines(ak, 5), "\n• ")
		}
	}

	if containsAny(q, "vat") && strings.Contains(akLower, "vat") {
		return true, "• " + strings.Join(firstNonEmptyLines(ak, 5), "\n• ")
	}
	if containsAny(q, "deadline", "hạn", "han") && strings.Contains(akLower, "deadline") {
		return true, "• " + strings.Join(firstNonEmptyLines(ak, 5), "\n• ")
	}

	return false, ""
}

func isInjectionAttempt(lu string) bool {
	return containsAny(lu,
		"ignore", "override", "system prompt", "developer message", "jailbreak",
		"bỏ qua", "bo qua", "tiết lộ prompt", "tiet lo prompt",
	)
}

func looksLikeFactualExact(lu string) bool {
	return containsAny(lu,
		"exact", "chính xác", "chinh xac",
		"số liệu", "so lieu", "statistics", "stat",
		"dẫn nguồn", "dan nguon", "sources", "source",
		"bao nhiêu", "bao nhieu", "how much", "how many",
		"percent", "%", "rate",
		"gdp", "vat",
	)
}

func containsAny(s string, subs ...string) bool {
	ss := strings.ToLower(s)
	for _, sub := range subs {
		if sub == "" {
			continue
		}
		if strings.Contains(ss, strings.ToLower(sub)) {
			return true
		}
	}
	return false
}

func bulletsVN(lines ...string) string {
	var b strings.Builder
	for _, s := range lines {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		b.WriteString("- ")
		b.WriteString(s)
		b.WriteString("\n")
	}
	return strings.TrimSpace(b.String())
}

func shorten(s string, max int) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func splitWords(s string) []string {
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteRune(' ')
		}
	}
	return strings.Fields(b.String())
}

func isStopWord(w string) bool {
	switch w {
	case "this", "that", "with", "from", "give", "what", "when", "where", "plan", "short", "your", "about":
		return true
	case "cho", "toi", "cua", "bao", "nhieu", "chinh", "xac", "dan", "nguon", "hay", "lam", "the", "nao":
		return true
	default:
		return false
	}
}

func firstNonEmptyLines(s string, max int) []string {
	raw := strings.Split(s, "\n")
	out := make([]string, 0, len(raw))
	for _, r := range raw {
		t := strings.TrimSpace(strings.TrimPrefix(r, "-"))
		if t != "" {
			out = append(out, t)
		}
		if len(out) >= max {
			break
		}
	}
	if len(out) == 0 {
		out = append(out, "(empty)")
	}
	return out
}
