package llm

import (
	"context"
	"hololab-core-system/internal/ports"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type MockProvider struct{}

func (p MockProvider) Chat(ctx context.Context, systemPrompt string, history []ports.ChatMsg, userText string) (string, error) {
	_ = ctx
	_ = history

	u := strings.TrimSpace(userText)
	if u == "" {
		return formatList(
			"Tôi chưa nhận được nội dung câu hỏi.",
			"Bạn vui lòng nhập rõ yêu cầu để tôi hỗ trợ.",
		), nil
	}
	lu := strings.ToLower(u)

	if isInjectionAttempt(lu) {
		return formatList(
			"Tôi không thể thực hiện yêu cầu này.",
			"Bạn vui lòng đặt câu hỏi trong phạm vi thông tin đã được cấu hình.",
		), nil
	}

	name, job, bio := parseIdentity(systemPrompt)
	_ = parseStyle(systemPrompt)
	ak, akSpecified := parseAllowedKnowledge(systemPrompt)

	if containsAny(lu, "bạn là ai", "ban la ai", "bạn giúp được gì", "ban giup duoc gi", "who are you", "what can you do") {
		intro := "Tôi là trợ lý ảo."
		if name != "" {
			intro = "Tôi là " + name + "."
		}

		lines := []string{intro}

		if job != "" {
			lines = append(lines, "Vai trò hiện tại: "+job+".")
		}
		if bio != "" {
			lines = append(lines, "Thông tin nền: "+shorten(bio, 180))
		}
		lines = append(lines, "Tôi trả lời dựa trên thông tin đã được cấu hình cho trợ lý này.")
		return formatList(lines...), nil
	}

	if containsAny(lu, "short plan", "kế hoạch", "ke hoach", "plan related to your occupation", "related to your occupation", "roadmap", "plan") {
		title := "Dưới đây là một kế hoạch ngắn gọn để bạn bắt đầu."
		if job != "" {
			title = "Dưới đây là một kế hoạch ngắn gọn phù hợp với vai trò " + job + "."
		}

		contextLine := ""
		if strings.TrimSpace(bio) != "" {
			contextLine = shorten(bio, 140)
		}

		lines := []string{title}
		if contextLine != "" {
			lines = append(lines, "Bối cảnh hiện có: "+contextLine)
		}

		lines = append(lines,
			"Trước hết, chốt mục tiêu trong 1–2 tuần tới (kết quả đo được).",
			"Sau đó chọn 2–3 việc tác động trực tiếp đến mục tiêu, tránh dàn trải.",
			"Chia nhỏ theo ngày để đảm bảo mỗi ngày có tiến độ rõ ràng.",
			"Cuối cùng, đặt tiêu chí hoàn thành để dễ kiểm tra và điều chỉnh.",
			"Nếu bạn nói rõ mục tiêu và ràng buộc (thời gian/ngân sách), tôi sẽ đề xuất chi tiết hơn.",
		)
		return formatList(lines...), nil
	}

	if looksLikeFactualExact(lu) {
		if !akSpecified {
			return formatList(
				"Tôi không tìm thấy thông tin phù hợp trong phần dữ liệu đã được cấu hình để trả lời chính xác câu này.",
				"Bạn thêm 2–3 gạch đầu dòng về số liệu hoặc nguồn liên quan (facts/sources), rồi tôi trả lời lại ngay.",
			), nil
		}

		if hit, answer := akAnswer(u, ak, true); hit {
			return answer, nil
		}

		return formatList(
			"Tôi không tìm thấy dữ liệu phù hợp trong phần dữ liệu đã được cấu hình để trả lời chính xác câu này.",
			"Bạn thêm 2–3 gạch đầu dòng về số liệu hoặc nguồn liên quan (facts/sources), rồi tôi trả lời lại ngay.",
		), nil
	}

	if akSpecified {
		if hit, answer := akAnswer(u, ak, false); hit {
			return answer, nil
		}
		return formatList(
			"Tôi chưa thấy thông tin liên quan trong phần dữ liệu đã được cấu hình để trả lời chắc chắn.",
			"Bạn bổ sung thêm vài gạch đầu dòng về dữ kiện liên quan (facts/sources), tôi sẽ trả lời sát hơn.",
		), nil
	}

	return formatList(
		"Tôi hiểu câu hỏi của bạn.",
		"Hiện tại tôi có thể trả lời theo hướng dẫn chung, nhưng sẽ không tự suy diễn số liệu hoặc thông tin cụ thể.",
		"Bạn muốn tôi đưa ra checklist, kế hoạch, hay hướng dẫn từng bước?",
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
		t = cleanText(t)
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
		lines := []string{
			"Tôi tìm thấy các thông tin liên quan trong phần dữ liệu đã được cấu hình:",
			strings.Join(prefixDash(topics), "\n"),
			"Nếu bạn cần số liệu hoặc thông tin khác chưa có ở đây, hãy bổ sung thêm facts/sources để tôi trả lời chính xác hơn.",
		}
		return true, formatList(lines...)
	}

	return true, formatList(
		"Dựa trên thông tin đã được cấu hình, đây là những điểm trọng tâm:",
		strings.Join(prefixDash(topics), "\n"),
		"Nếu bạn nói rõ bối cảnh và mục tiêu, tôi sẽ đề xuất hướng thực thi cụ thể hơn.",
	)
}

func prefixDash(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, s := range lines {
		t := strings.TrimSpace(s)
		if t == "" {
			continue
		}
		out = append(out, "- "+t)
	}
	return out
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

	allowShort := map[string]bool{
		"mvp": true, "pmf": true, "gtm": true, "kpi": true, "okr": true,
		"ux": true, "ui": true, "api": true, "saas": true,
		"gdp": true, "vat": true,
	}

	out := make([]string, 0, len(words))
	seen := map[string]bool{}

	for _, w := range words {
		if w == "" {
			continue
		}
		if isStopWord(w) {
			continue
		}

		if len(w) < 4 && !allowShort[w] {
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

func formatList(lines ...string) string {
	var b strings.Builder
	for _, s := range lines {
		s = cleanText(s)
		if s == "" {
			continue
		}
		if isNumberedList(s) {
			b.WriteString(s)
		} else {
			b.WriteString(s)
		}
		b.WriteString("\n")
	}
	return strings.TrimSpace(b.String())
}

func cleanText(s string) string {
	s = strings.ReplaceAll(s, "*", "")
	s = strings.ReplaceAll(s, "#", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.TrimSpace(s)
	return s
}

func isNumberedList(s string) bool {
	if len(s) > 2 && unicode.IsDigit(rune(s[0])) && (s[1] == '.' || s[1] == ')') {
		return true
	}
	return false
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
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
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
