require("antd/dist/antd.css")

import React from "react"
import ReactDOM from "react-dom"
import {BrowserRouter, Route} from "react-router-dom"

import App from "./components/layout/app.jsx"

const Main = () => {
    return <BrowserRouter basename="/app">
        <Route path="/" component={App} />
    </BrowserRouter>
}

ReactDOM.render(<Main />, document.getElementById("root"));