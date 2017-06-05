package zclsrv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/sourcegraph/go-langserver/pkg/lsp"
	"github.com/sourcegraph/jsonrpc2"
)

// NewHandler creates a zcl language server handler.
func NewHandler() jsonrpc2.Handler {
	return jsonrpc2.HandlerWithError((&handler{}).Handle)
}

// handler is the main JSON-RPC handler for the language server.
type handler struct {
	mu         sync.Mutex
	RootFSPath string // root path of the project's files in the (possibly virtual) file system, without the "file://" prefix (typically /src/github.com/foo/bar)
	init       *lsp.InitializeParams
}

func (h *handler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	// Prevent any uncaught panics from taking the entire server down.
	defer func() {
		v := recover()
		if v != nil {
			err = fmt.Errorf("Panic! %#v", recover())
		}
	}()

	h.mu.Lock()
	if req.Method != "initialize" && h.init == nil {
		h.mu.Unlock()
		return nil, errors.New("server must be initialized")
	}
	h.mu.Unlock()
	/*if err := h.CheckReady(); err != nil {
		if req.Method == "exit" {
			err = nil
		}
		return nil, err
	}*/

	switch req.Method {
	case "initialize":
		if h.init != nil {
			return nil, errors.New("language server is already initialized")
		}
		if req.Params == nil {
			return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
		}
		if err := json.Unmarshal(*req.Params, &h.init); err != nil {
			return nil, err
		}

		//kind := lsp.TDSKFull
		return lsp.InitializeResult{
			Capabilities: lsp.ServerCapabilities{
			/*TextDocumentSync: lsp.TextDocumentSyncOptionsOrKind{
				Kind: &kind,
			},
			DefinitionProvider:           true,
			DocumentFormattingProvider:   true,
			DocumentSymbolProvider:       true,
			HoverProvider:                true,
			ReferencesProvider:           true,
			WorkspaceSymbolProvider:      true,
			XWorkspaceReferencesProvider: true,
			XDefinitionProvider:          true,
			XWorkspaceSymbolByProperties: true,
			SignatureHelpProvider:        &lsp.SignatureHelpOptions{TriggerCharacters: []string{"(", ","}},*/
			},
		}, nil

	case "initialized":
		// A notification that the client is ready to receive requests. Ignore
		return nil, nil
		/*
			case "shutdown":
				h.ShutDown()
				return nil, nil

			case "exit":
				if c, ok := conn.(*jsonrpc2.Conn); ok {
					c.Close()
				}
				return nil, nil

			case "$/cancelRequest":
				// notification, don't send back results/errors
				if req.Params == nil {
					return nil, nil
				}
				var params lsp.CancelParams
				if err := json.Unmarshal(*req.Params, &params); err != nil {
					return nil, nil
				}
				if cancelManager == nil {
					return nil, nil
				}
				cancelManager.Cancel(jsonrpc2.ID{
					Num:      params.ID.Num,
					Str:      params.ID.Str,
					IsString: params.ID.IsString,
				})
				return nil, nil

			case "textDocument/hover":
				if req.Params == nil {
					return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
				}
				var params lsp.TextDocumentPositionParams
				if err := json.Unmarshal(*req.Params, &params); err != nil {
					return nil, err
				}
				return h.handleHover(ctx, conn, req, params)

			case "textDocument/definition":
				if req.Params == nil {
					return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
				}
				var params lsp.TextDocumentPositionParams
				if err := json.Unmarshal(*req.Params, &params); err != nil {
					return nil, err
				}
				return h.handleDefinition(ctx, conn, req, params)

			case "textDocument/xdefinition":
				if req.Params == nil {
					return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
				}
				var params lsp.TextDocumentPositionParams
				if err := json.Unmarshal(*req.Params, &params); err != nil {
					return nil, err
				}
				return h.handleXDefinition(ctx, conn, req, params)

			case "textDocument/references":
				if req.Params == nil {
					return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
				}
				var params lsp.ReferenceParams
				if err := json.Unmarshal(*req.Params, &params); err != nil {
					return nil, err
				}
				return h.handleTextDocumentReferences(ctx, conn, req, params)

			case "textDocument/documentSymbol":
				if req.Params == nil {
					return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
				}
				var params lsp.DocumentSymbolParams
				if err := json.Unmarshal(*req.Params, &params); err != nil {
					return nil, err
				}
				return h.handleTextDocumentSymbol(ctx, conn, req, params)

			case "textDocument/signatureHelp":
				if req.Params == nil {
					return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
				}
				var params lsp.TextDocumentPositionParams
				if err := json.Unmarshal(*req.Params, &params); err != nil {
					return nil, err
				}
				return h.handleTextDocumentSignatureHelp(ctx, conn, req, params)

			case "textDocument/formatting":
				if req.Params == nil {
					return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
				}
				var params lsp.DocumentFormattingParams
				if err := json.Unmarshal(*req.Params, &params); err != nil {
					return nil, err
				}
				return h.handleTextDocumentFormatting(ctx, conn, req, params)

			case "workspace/symbol":
				if req.Params == nil {
					return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
				}
				var params lspext.WorkspaceSymbolParams
				if err := json.Unmarshal(*req.Params, &params); err != nil {
					return nil, err
				}
				return h.handleWorkspaceSymbol(ctx, conn, req, params)

			case "workspace/xreferences":
				if req.Params == nil {
					return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
				}
				var params lspext.WorkspaceReferencesParams
				if err := json.Unmarshal(*req.Params, &params); err != nil {
					return nil, err
				}
				return h.handleWorkspaceReferences(ctx, conn, req, params)*/

	default:
		/*if isFileSystemRequest(req.Method) {
			uri, fileChanged, err := h.handleFileSystemRequest(ctx, req)
			if fileChanged {
				// a file changed, so we must re-typecheck and re-enumerate symbols
				h.resetCaches(true)
			}
			if uri != "" {
				// a user is viewing this path, hint to add it to the cache
				go h.typecheck(ctx, conn, uri, lsp.Position{})
			}
			return nil, err
		}*/

		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	}
}
