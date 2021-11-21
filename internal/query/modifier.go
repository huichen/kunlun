package query

func GetModifierName(qtype QueryType) string {
	switch qtype {
	case LanguageQuery:
		return "lang"
	case RepoQuery:
		return "repo"
	case FileQuery:
		return "file"
	case CaseQuery:
		return "case"

	}
	return ""
}
