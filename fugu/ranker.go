package fugu

// 1x Img -> +1
// >10x Img -> +1
// 1x Css -> +1
// >10x Css -> +1
// 1x Script -> +100
// >10x Script -> +100
// Cookie -> +1000

func Rank(p SitePrivacy) int {
	rank := 0
	if p.ImgCount > 0 {
		rank += 1
	}
	if p.ImgCount > 10 {
		rank += 1
	}
	if p.ScriptCount > 0 {
		rank += 100
	}
	if p.ScriptCount > 10 {
		rank += 100
	}
	if p.Cookie {
		rank += 1000
	}
	return rank
}
