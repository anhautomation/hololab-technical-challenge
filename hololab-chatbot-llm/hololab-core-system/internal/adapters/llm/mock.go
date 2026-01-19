package llm

import (
	"context"
	"hololab-core-system/internal/ports"
	"regexp"
	"sort"
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
			"Mình không thể bỏ qua persona và rules đã cấu hình.",
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
			out = append(out, "Vai trò: **"+job+"**.")
		}
		if bio != "" {
			out = append(out, "Tiểu sử: "+shorten(bio, 180))
		}
		out = append(out, "Mình trả lời theo đúng phong cách và **không vượt quá Allowed Knowledge**.")
		return bulletsVN(out...), nil
	}

	if containsAny(lu, "short plan", "kế hoạch", "ke hoach", "plan related to your occupation", "related to your occupation") {
		title := "Kế hoạch ngắn theo vai trò hiện tại"
		if job != "" {
			title = "Kế hoạch ngắn (" + job + ")"
		}
		contextLine := "(chưa cấu hình tiểu sử)"
		if strings.TrimSpace(bio) != "" {
			contextLine = shorten(bio, 140)
		}
		return bulletsVN(
			"**"+title+":**",
			"Bối cảnh: "+contextLine,
			"1) Mục tiêu 1–2 tuần: kết quả đo được là gì?",
			"2) Chia việc: (a) cốt lõi theo vai trò, (b) hỗ trợ, (c) kiểm soát chất lượng.",
			"3) Mỗi ngày 1–3 đầu việc nhỏ hoàn thành được.",
			"4) Tiêu chí: đúng – đủ – dễ kiểm tra.",
			"5) Review cuối tuần: giữ cái hiệu quả, bỏ cái thừa.",
			"Bạn cho mình 2 thông tin: (a) mục tiêu cụ thể, (b) ràng buộc thời gian/ngân sách.",
		), nil
	}

	if looksLikeFactualExact(lu) {
		if !akSpecified {
			return bulletsVN(
				"Mình **chưa có Allowed Knowledge cụ thể** cho bot này nên không thể trả lời dạng **số liệu/nguồn chính xác**.",
				"Bạn cập nhật *Allowed Knowledge* (facts + nguồn) rồi mình trả lời đúng phạm vi.",
			), nil
		}

		if hit, answer := akAnswer(u, ak, true); hit {
			return answer, nil
		}

		return bulletsVN(
			"Mình **không thấy dữ liệu liên quan** trong *Allowed Knowledge* để trả lời chính xác.",
			"Bạn bổ sung facts/nguồn vào *Allowed Knowledge* rồi mình trả lời theo đúng phạm vi.",
		), nil
	}

	if akSpecified {
		if hit, answer := akAnswer(u, ak, false); hit {
			return answer, nil
		}
		return bulletsVN(
			"Mình chưa thấy thông tin liên quan trong *Allowed Knowledge* để trả lời chắc chắn.",
			"Bạn có thể bổ sung thêm dữ liệu vào *Allowed Knowledge* để mình bám đúng phạm vi.",
		), nil
	}

	return bulletsVN(
		"Mình hiểu câu hỏi của bạn: \""+u+"\".",
		"Mình có thể trả lời **hướng dẫn chung**, nhưng sẽ không bịa số liệu/nguồn.",
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

func akAnswer(question string, ak string, exactMode bool) (bool, string) {
	lines := normalizeAKLines(ak)
	if len(lines) == 0 {
		return false, ""
	}

	q := strings.ToLower(question)
	keywords := extractKeywords(q)
	if len(keywords) == 0 {
		keywords = []string{q}
	}

	type scored struct {
		line  string
		score int
	}
	scoredLines := make([]scored, 0, len(lines))

	for _, ln := range lines {
		lnLow := strings.ToLower(ln)
		s := 0
		for _, kw := range keywords {
			if kw == "" {
				continue
			}
			if strings.Contains(lnLow, kw) {
				s += 3
			} else if fuzzyHit(lnLow, kw) {
				s += 1
			}
		}
		if s > 0 {
			scoredLines = append(scoredLines, scored{line: ln, score: s})
		}
	}

	if len(scoredLines) == 0 {
		return false, ""
	}

	sort.SliceStable(scoredLines, func(i, j int) bool {
		return scoredLines[i].score > scoredLines[j].score
	})

	topics := make([]string, 0, 5)
	seen := map[string]bool{}
	for _, it := range scoredLines {
		if len(topics) >= 5 {
			break
		}
		t := strings.TrimSpace(it.line)
		if t == "" {
			continue
		}
		if !seen[t] {
			seen[t] = true
			topics = append(topics, t)
		}
	}
	if len(topics) == 0 {
		return false, ""
	}

	if exactMode {
		return true, bulletsVN(
			"Mình trả lời dựa trên *Allowed Knowledge* đã cấu hình:",
			"Phần liên quan mình tìm thấy:",
			"• "+strings.Join(topics, "\n• "),
			"Nếu bạn cần **số liệu/nguồn chính xác** mà chưa có trong *Allowed Knowledge*, hãy bổ sung thêm dữ liệu.",
		)
	}

	return true, bulletsVN(
		"Mình bám theo *Allowed Knowledge* đã cấu hình và trả lời theo hướng thực thi:",
		"Trọng tâm liên quan:",
		"• "+strings.Join(topics, "\n• "),
		"Gợi ý triển khai nhanh:",
		"1) Chốt mục tiêu và ràng buộc (deadline / ngân sách / mức rủi ro chấp nhận được).",
		"2) Lập checklist theo các mục trọng tâm ở trên, chọn 1–2 mục ưu tiên làm trước.",
		"3) Xác định đầu ra có thể kiểm tra (bảng theo dõi / rule / quy trình / báo cáo).",
		"Bạn cho mình 2 thông tin để cụ thể hoá: (a) bối cảnh hiện tại, (b) outcome bạn muốn đạt là gì?",
	)
}

func normalizeAKLines(ak string) []string {
	raw := strings.Split(ak, "\n")
	out := make([]string, 0, len(raw))
	for _, r := range raw {
		t := strings.TrimSpace(r)
		t = strings.TrimPrefix(t, "-")
		t = strings.TrimPrefix(t, "•")
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		out = append(out, t)
	}
	return out
}

func extractKeywords(q string) []string {
	words := splitWords(q)
	out := make([]string, 0, len(words))
	seen := map[string]bool{}
	for _, w := range words {
		if len(w) < 4 {
			continue
		}
		if isStopWord(w) {
			continue
		}
		if !seen[w] {
			seen[w] = true
			out = append(out, w)
		}
	}
	return out
}

func fuzzyHit(text string, kw string) bool {
	if kw == "" {
		return false
	}
	if strings.HasPrefix(kw, "risk") && strings.Contains(text, "risk") {
		return true
	}
	if strings.HasPrefix(kw, "cost") && strings.Contains(text, "cost") {
		return true
	}
	if strings.HasPrefix(kw, "budget") && strings.Contains(text, "budget") {
		return true
	}
	if strings.HasPrefix(kw, "audit") && strings.Contains(text, "audit") {
		return true
	}
	if strings.HasPrefix(kw, "vendor") && strings.Contains(text, "vendor") {
		return true
	}
	if strings.HasPrefix(kw, "contract") && strings.Contains(text, "contract") {
		return true
	}
	if strings.HasPrefix(kw, "compliance") && strings.Contains(text, "compliance") {
		return true
	}
	return false
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
