package api

type APIRequest struct {
	Messages  []APIMessage `json:"messages"`
	Stream    bool         `json:"stream"`
	Model     string       `json:"model"`
	PluginIDs []string     `json:"plugin_ids"`
}

type APIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
