package vision

// ModelOptions represents additional model parameters listed in the documentation.
type ModelOptions struct {
	NumKeep          int      `yaml:"NumKeep,omitempty" json:"num_keep,omitempty"` // Ollama ↓
	Seed             int      `yaml:"Seed,omitempty" json:"seed,omitempty"`
	NumPredict       int      `yaml:"NumPredict,omitempty" json:"num_predict,omitempty"`
	Temperature      float64  `yaml:"Temperature,omitempty" json:"temperature,omitempty"`
	TopK             int      `yaml:"TopK,omitempty" json:"top_k,omitempty"`
	TopP             float64  `yaml:"TopP,omitempty" json:"top_p,omitempty"`
	MinP             float64  `yaml:"MinP,omitempty" json:"min_p,omitempty"`
	TypicalP         float64  `yaml:"TypicalP,omitempty" json:"typical_p,omitempty"`
	TfsZ             float64  `yaml:"TfsZ,omitempty" json:"tfs_z,omitempty"`
	RepeatLastN      int      `yaml:"RepeatLastN,omitempty" json:"repeat_last_n,omitempty"`
	RepeatPenalty    float64  `yaml:"RepeatPenalty,omitempty" json:"repeat_penalty,omitempty"`
	PresencePenalty  float64  `yaml:"PresencePenalty,omitempty" json:"presence_penalty,omitempty"`
	FrequencyPenalty float64  `yaml:"FrequencyPenalty,omitempty" json:"frequency_penalty,omitempty"`
	Mirostat         int      `yaml:"Mirostat,omitempty" json:"mirostat,omitempty"`
	MirostatTau      float64  `yaml:"MirostatTau,omitempty" json:"mirostat_tau,omitempty"`
	MirostatEta      float64  `yaml:"MirostatEta,omitempty" json:"mirostat_eta,omitempty"`
	PenalizeNewline  bool     `yaml:"PenalizeNewline,omitempty" json:"penalize_newline,omitempty"`
	Stop             []string `yaml:"Stop,omitempty" json:"stop,omitempty"`
	Numa             bool     `yaml:"Numa,omitempty" json:"numa,omitempty"`
	NumCtx           int      `yaml:"NumCtx,omitempty" json:"num_ctx,omitempty"`
	NumBatch         int      `yaml:"NumBatch,omitempty" json:"num_batch,omitempty"`
	NumGpu           int      `yaml:"NumGpu,omitempty" json:"num_gpu,omitempty"`
	MainGpu          int      `yaml:"MainGpu,omitempty" json:"main_gpu,omitempty"`
	LowVram          bool     `yaml:"LowVram,omitempty" json:"low_vram,omitempty"`
	VocabOnly        bool     `yaml:"VocabOnly,omitempty" json:"vocab_only,omitempty"`
	UseMmap          bool     `yaml:"UseMmap,omitempty" json:"use_mmap,omitempty"`
	UseMlock         bool     `yaml:"UseMlock,omitempty" json:"use_mlock,omitempty"`
	NumThread        int      `yaml:"NumThread,omitempty" json:"num_thread,omitempty"`
	MaxOutputTokens  int      `yaml:"MaxOutputTokens,omitempty" json:"max_output_tokens,omitempty"` // OpenAI ↓
	Detail           string   `yaml:"Detail,omitempty" json:"detail,omitempty"`
	ForceJson        bool     `yaml:"ForceJson,omitempty" json:"force_json,omitempty"`
	SchemaVersion    string   `yaml:"SchemaVersion,omitempty" json:"schema_version,omitempty"`
	CombineOutputs   string   `yaml:"CombineOutputs,omitempty" json:"combine_outputs,omitempty"`
}
