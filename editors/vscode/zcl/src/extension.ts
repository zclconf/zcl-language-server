'use strict';

import * as path from 'path';

import { workspace, Disposable, ExtensionContext } from 'vscode';
import { LanguageClient, LanguageClientOptions, SettingMonitor, ServerOptions, TransportKind } from 'vscode-languageclient';

export function activate(context: ExtensionContext) {
    // FIXME: remove this once this stops being a prototype
    console.log('zcl extension activated');

    // FIXME: Need to deal with choosing the right binary depending on
    // what OS and architecture we're running on.
    let serverBin = context.asAbsolutePath(path.join('zcl-language-server'));

	let serverOptions: ServerOptions = {
		run: { command: serverBin },
		debug: { command: serverBin },
	}

	let clientOptions: LanguageClientOptions = {
		// Register the server for plain text documents
		documentSelector: ['zcl'],
		synchronize: {
			configurationSection: 'zcl',
		},
	}

    let disposable = new LanguageClient('zcl', 'zcl Language Server', serverOptions, clientOptions).start();

    context.subscriptions.push(disposable);
}

// this method is called when your extension is deactivated
export function deactivate() {
}