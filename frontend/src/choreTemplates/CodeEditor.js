import React from "react";
import PropTypes from "prop-types";

import "codemirror/lib/codemirror.css"
import * as CM from "codemirror/lib/codemirror";
import "codemirror/mode/javascript/javascript";
import "codemirror/addon/dialog/dialog";
import "codemirror/addon/dialog/dialog.css"
import "codemirror/addon/fold/foldgutter.css";
import "codemirror/addon/fold/brace-fold";
import "codemirror/addon/fold/foldcode";
import "codemirror/addon/fold/foldgutter.css"

import "codemirror/addon/hint/javascript-hint";
import "codemirror/addon/hint/show-hint";
import"codemirror/addon/hint/show-hint.css";
import "codemirror/addon/tern/tern"
import "codemirror/addon/tern/tern.css"
import "codemirror/theme/material.css";

import {UnControlled as CodeMirror} from 'react-codemirror2'

// eslint-disable-next-line import/no-webpack-loader-syntax
let workerScript = require("file-loader!codemirror/addon/tern/worker");

class CodeEditor extends React.Component {
    render() {
        let options = {
            mode: "javascript",
            theme: 'material',
            lineNumbers: true,
        };

        return <CodeMirror
                value={this.props.code}
                options={options}
                onChange={(editor, data, value) => {this.props.onChange(value)}}
                editorDidMount={this.initializeEditor}
            />
    }

    initializeEditor = async (editor) => {
        console.log("initializing editor");
        let EcmaDefs = await fetchDefs();
        let server = new CM.TernServer({
            defs: [defs],
            useWorker: true, workerScript: workerScript,
            workerDeps: [
                "https://ternjs.net/node_modules/acorn/dist/acorn.js",
                "https://ternjs.net/node_modules/acorn-loose/dist/acorn-loose.js",
                "https://ternjs.net/node_modules/acorn-walk/dist/walk.js",
                "https://ternjs.net/lib/signal.js",
                "https://ternjs.net/lib/tern.js",
                "https://ternjs.net/lib/def.js",
                "https://ternjs.net/lib/comment.js",
                "https://ternjs.net/lib/infer.js",
                "https://ternjs.net/plugin/doc_comment.js",
            ]
        });
        editor.setOption("extraKeys", {
            "Ctrl-Space": function(cm) { server.complete(cm); },
            "Ctrl-I": function(cm) { server.showType(cm); },
            "Ctrl-O": function(cm) { server.showDocs(cm); },
            "Alt-.": function(cm) { server.jumpToDef(cm); },
            "Alt-,": function(cm) { server.jumpBack(cm); },
            "Ctrl-Q": function(cm) { server.rename(cm); },
            "Ctrl-.": function(cm) { server.selectName(cm); }
        });
        editor.on("cursorActivity", function(cm) { server.updateArgHints(cm); });
    }
}

CodeEditor.propTypes = {
    width: PropTypes.oneOfType([PropTypes.number, PropTypes.string]).isRequired,
    height: PropTypes.oneOfType([PropTypes.number, PropTypes.string]).isRequired,
    code: PropTypes.string,
    onChange: PropTypes.func,
};

CodeEditor.defaultProps = {
    onChange: () => {}
};

export default CodeEditor;

let defs = null;
async function fetchDefs() {
    if (defs !== null) {
        return defs;
    }

    let resp = await fetch("http://ternjs.net/defs/ecmascript.json");
    defs = await resp.json();
    return defs;
}
