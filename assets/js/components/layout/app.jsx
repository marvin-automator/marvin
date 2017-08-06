import React from "react"
import Layout from "antd/lib/layout/layout"

import Menu from "./menu.jsx"

const App = () => {
    return <Layout>
        <Layout.Header>
            <img src="http://via.placeholder.com/128x128" style={{width: "64px", height: "63px", float:"left"}} />
            <Menu />
        </Layout.Header>
        <Layout.Content>

        </Layout.Content>
    </Layout>
}

export default App