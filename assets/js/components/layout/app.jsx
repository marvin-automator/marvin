import React from "react"
import {Route} from "react-router-dom"
import {Container} from "semantic-ui-react"

import Menu from "./menu.jsx"
import ChoresPage from "../chores/ChoresPage"

const App = () => {
    return <div>
        <Menu />
        <Container style={{marginTop: '7em'}}>
            <Route path="/chores" component={ChoresPage} />
        </Container>
    </div>
}

export default App