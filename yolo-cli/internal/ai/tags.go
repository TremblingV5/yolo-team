package ai

import (
	"regexp"
	"strings"
)

// TagType represents the type of a render tag.
type TagType string

const (
	TagProjectList    TagType = "project_list"
	TagProjectCard    TagType = "project_card"
	TagIssueList      TagType = "issue_list"
	TagIssueCard      TagType = "issue_card"
	TagTaskList       TagType = "task_list"
	TagExecutorList   TagType = "executor_list"
	TagKanbanBoard    TagType = "kanban_board"
	TagDocumentList   TagType = "document_list"
	TagStatusBadge    TagType = "status_badge"
	TagPriorityBadge  TagType = "priority_badge"
)

// Tag represents a parsed render tag.
type Tag struct {
	Type  TagType           `json:"type"`
	Attrs map[string]string `json:"attrs,omitempty"`
	Raw   string            `json:"raw"`
}

var tagPattern = regexp.MustCompile(`<\w+\s*[^>]*\/\s*>`)

// ParseTags extracts all render tags from text.
func ParseTags(text string) []Tag {
	matches := tagPattern.FindAllString(text, -1)
	if matches == nil {
		return nil
	}

	tags := make([]Tag, 0, len(matches))
	for _, m := range matches {
		tag := parseSingleTag(m)
		if tag != nil {
			tags = append(tags, *tag)
		}
	}

	return tags
}

// RemoveTags strips all tags from text, returning plain text.
func RemoveTags(text string) string {
	return tagPattern.ReplaceAllString(text, "")
}

// RenderTag renders a single tag as a text representation (backend fallback).
func RenderTag(tag Tag) string {
	var sb strings.Builder
	sb.WriteString("[")

	switch tag.Type {
	case TagProjectList:
		sb.WriteString("项目列表")
	case TagProjectCard:
		sb.WriteString("项目卡片")
		if key := tag.Attrs["key"]; key != "" {
			sb.WriteString(": ")
			sb.WriteString(key)
		}
	case TagIssueList:
		sb.WriteString("任务列表")
		if project := tag.Attrs["project"]; project != "" {
			sb.WriteString(" (项目: ")
			sb.WriteString(project)
			sb.WriteString(")")
		}
	case TagIssueCard:
		sb.WriteString("任务卡片")
		if key := tag.Attrs["key"]; key != "" {
			sb.WriteString(": ")
			sb.WriteString(key)
		}
	case TagTaskList:
		sb.WriteString("子任务列表")
		if issue := tag.Attrs["issue"]; issue != "" {
			sb.WriteString(" (任务: ")
			sb.WriteString(issue)
			sb.WriteString(")")
		}
	case TagExecutorList:
		sb.WriteString("执行人列表")
	case TagKanbanBoard:
		sb.WriteString("看板")
	case TagDocumentList:
		sb.WriteString("文档列表")
		if project := tag.Attrs["project"]; project != "" {
			sb.WriteString(" (项目: ")
			sb.WriteString(project)
			sb.WriteString(")")
		}
	case TagStatusBadge:
		sb.WriteString("状态: ")
		sb.WriteString(tag.Attrs["status"])
	case TagPriorityBadge:
		sb.WriteString("优先级: ")
		sb.WriteString(tag.Attrs["priority"])
	default:
		sb.WriteString(string(tag.Type))
	}

	sb.WriteString("]")
	return sb.String()
}

// parseSingleTag parses a single tag string into a Tag struct.
func parseSingleTag(raw string) *Tag {
	// Strip < > and trailing />
	inner := strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(raw), "<"), "/>")
	inner = strings.TrimSpace(inner)

	parts := strings.SplitN(inner, " ", 2)
	if len(parts) == 0 || parts[0] == "" {
		return nil
	}

	name := parts[0]
	var tagType TagType
	switch name {
	case "project_list":
		tagType = TagProjectList
	case "project_card":
		tagType = TagProjectCard
	case "issue_list":
		tagType = TagIssueList
	case "issue_card":
		tagType = TagIssueCard
	case "task_list":
		tagType = TagTaskList
	case "executor_list":
		tagType = TagExecutorList
	case "kanban_board":
		tagType = TagKanbanBoard
	case "document_list":
		tagType = TagDocumentList
	case "status_badge":
		tagType = TagStatusBadge
	case "priority_badge":
		tagType = TagPriorityBadge
	default:
		return nil
	}

	tag := &Tag{
		Type:  tagType,
		Attrs: make(map[string]string),
		Raw:   raw,
	}

	if len(parts) > 1 {
		parseAttrs(tag.Attrs, parts[1])
	}

	return tag
}

// parseAttrs parses key="value" pairs from the remaining tag string.
func parseAttrs(attrs map[string]string, raw string) {
	re := regexp.MustCompile(`(\w+)\s*=\s*"([^"]*)"`)
	matches := re.FindAllStringSubmatch(raw, -1)
	for _, m := range matches {
		if len(m) >= 3 {
			attrs[m[1]] = m[2]
		}
	}
}
