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
			"Tôi không thể bỏ qua persona và rules đã cấu hình.",
			"Hãy hỏi trong phạm vi hồ sơ và *Allowed Knowledge* của bot.",
		), nil
	}

	name, job, bio := parseIdentity(systemPrompt)
	style := parseStyle(systemPrompt)
	ak, akSpecified := parseAllowedKnowledge(systemPrompt)

	_ = style

	if containsAny(lu, "bạn là ai", "ban la ai", "bạn giúp được gì", "ban giup duoc gi", "who are you", "what can you do") {
		head := "Mình là chatbot theo persona đã cấu hình."
		if name != "" {
			head = "Mình là **" + name + "** theo persona đã cấu hình."
		}
		out := []string{head}
		if job != "" {
			out = append(out, "Nghề nghiệp/Vai trò: **"+job+"**.")
		}
		if bio != "" {
			out = append(out, "Tiểu sử: "+shorten(bio, 180))
		}
		out = append(out,
			"Mình sẽ trả lời theo đúng phong cách và **không vượt quá Allowed Knowledge**.",
		)
		return bulletsVN(out...), nil
	}

	if containsAny(lu, "short plan", "kế hoạch", "ke hoach", "plan related to your occupation", "related to your occupation") {
		title := "Kế hoạch ngắn theo vai trò hiện tại"
		if job != "" {
			title = "Kế hoạch ngắn (" + job + ")"
		}

		context := "(chưa cấu hình tiểu sử)"
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
		if !akSpecified {
			return bulletsVN(
				"Mình **không có Allowed Knowledge cụ thể** cho bot này, nên không thể trả lời dạng **số liệu/nguồn chính xác**.",
				"Bạn hãy cập nhật *Allowed Knowledge* (facts + nguồn), rồi mình sẽ trả lời đúng phạm vi.",
			), nil
		}

		if hit, snippet := akMatch(u, ak); hit {
			return bulletsVN(
				"Mình trả lời theo *Allowed Knowledge* đã cấu hình:",
				snippet,
			), nil
		}

		return bulletsVN(
			"Mình **không thấy dữ liệu liên quan** trong *Allowed Knowledge* để trả lời chính xác cho câu hỏi này.",
			"Bạn hãy bổ sung facts/nguồn vào *Allowed Knowledge*, rồi mình sẽ trả lời theo đúng phong cách đã cấu hình.",
		), nil
	}

	if akSpecified {
		if hit, snippet := akMatch(u, ak); hit {
			return bulletsVN(
				"Mình trả lời theo *Allowed Knowledge* đã cấu hình:",
				snippet,
			), nil
		}
		return bulletsVN(
			"Mình chưa thấy thông tin liên quan trong *Allowed Knowledge* để trả lời chắc chắn.",
			"Bạn có thể bổ sung thêm dữ liệu vào *Allowed Knowledge* để mình bám đúng phạm vi.",
		), nil
	}

	return bulletsVN(
		"Mình hiểu câu hỏi của bạn: \""+u+"\".",
		"Mình có thể trả lời **hướng dẫn chung** theo persona, nhưng sẽ không bịa số liệu/nguồn.",
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

	reBioBlock := regexp.MustCompile(`(?is)^\s*-\s*Bio:\s*(.*?)\n\nSTYLE\s*/\s*TONE`)
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
	re := regexp.MustCompile(`(?is)STYLE\s*/\s*TONE\s*\(must follow\)\s*:\s*(.*?)\n\n`)
	m := re.FindStringSubmatch(systemPrompt)
	if len(m) == 2 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

func parseAllowedKnowledge(systemPrompt string) (knowledge string, specified bool) {
	reNot := regexp.MustCompile(`(?is)ALLOWED KNOWLEDGE\s*:\s*Not specified\.`)
	if reNot.FindStringIndex(systemPrompt) != nil {
		return "", false
	}

	re1 := regexp.MustCompile(`(?is)ALLOWED KNOWLEDGE\s*\(ONLY use these facts/assumptions\)\s*:\s*(.*)$`)
	m := re1.FindStringSubmatch(systemPrompt)
	if len(m) != 2 {
		return "", false
	}
	block := strings.TrimSpace(m[1])

	cutRe := regexp.MustCompile(`(?is)\n\nIf the user asks outside this knowledge.*$`)
	block = cutRe.ReplaceAllString(block, "")
	block = strings.TrimSpace(block)

	if block == "" {
		return "", false
	}
	return block, true
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
		out = append(out, "(Allowed Knowledge is empty)")
	}
	return out
}
