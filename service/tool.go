package service

import (
	"ai-prompt-shell/dao"
	"ai-prompt-shell/internal/utils"
	"net/http"
)

var tools = dao.NewToolCache()

/**
 * Get tool by ID from cache
 * @param toolId ID of the tool to retrieve
 * @return tool definition if found
 * @return os.ErrNotExist error if tool not found
 */
func GetTool(toolId string) (dao.Tool, error) {
	t, exists := tools.Get(toolId)
	if !exists {
		return dao.Tool{}, utils.NewHttpError(http.StatusNotFound, "Tool not found")
	}
	return t, nil
}

/**
 * Get all available tool IDs
 * @return slice of tool IDs
 * @return nil error for future compatibility (always succeeds currently)
 */
func ToolIDs() ([]string, error) {
	var results []string
	for k, _ := range tools.All() {
		results = append(results, k)
	}
	return results, nil
}
