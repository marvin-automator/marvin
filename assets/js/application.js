require("antd/dist/antd.css")

import React from "react"
import ReactDOM from "react-dom"
import {BrowserRouter, Route} from "react-router-dom"
import {ApolloProvider} from 'react-apollo'

import App from "./components/layout/app.jsx"
import {client} from "./graphql"

const Main = () => {
    return <ApolloProvider client={client}>
        <BrowserRouter basename="/app">
            <Route path="/" component={App} />
        </BrowserRouter>
    </ApolloProvider>
}

ReactDOM.render(<Main />, document.getElementById("root"));