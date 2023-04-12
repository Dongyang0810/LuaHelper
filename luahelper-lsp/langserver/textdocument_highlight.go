package langserver

import (
	"context"
	"luahelper-lsp/langserver/check"
	"luahelper-lsp/langserver/check/common"
	"luahelper-lsp/langserver/log"
	"luahelper-lsp/langserver/lspcommon"
	lsp "luahelper-lsp/langserver/protocol"
)

// TextDocumentHighlight 对变量单击选中着色
func (l *LspServer) TextDocumentHighlight(ctx context.Context, vs lsp.TextDocumentPositionParams) (retVec []lsp.DocumentHighlight,
	err error) {
	l.requestMutex.Lock()
	defer l.requestMutex.Unlock()

	if !l.isCanHighlight() {
		log.Error("IsCanHighlight is false")
		return
	}

	//拿到这个变量所在文件的文件名、整个文件字符串，光标点击处在整体字符串的偏移等信息
	comResult := l.beginFileRequest(vs.TextDocument.URI, vs.Position)
	if !comResult.result {
		return
	}

	if len(comResult.contents) == 0 || comResult.offset >= len(comResult.contents) {
		return
	}

	project := l.getAllProject()

	//GetVarStruct根据文件字符串和位置偏移等信息解析出是：光标点击的具体单词，位置，前面是否是冒号，是否有括号等
	varStruct := check.GetVarStruct(comResult.contents, comResult.offset, comResult.pos.Line, comResult.pos.Character)
	if !varStruct.ValidFlag {
		log.Error("TextDocumentHighlight varStruct.ValidFlag not valid")
		return
	}

	// 去掉前缀后的名字
	referenVecs := project.FindReferences(comResult.strFile, &varStruct, common.CRSHighlight)
	retVec = make([]lsp.DocumentHighlight, 0, len(referenVecs))
	for _, referVarInfo := range referenVecs {
		retVec = append(retVec, lsp.DocumentHighlight{
			Range: lspcommon.LocToRange(&referVarInfo.Loc),
			Kind:  lsp.Write,
		})
	}
	return
}
