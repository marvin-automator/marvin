import React from "react"
import {Route} from "react-router-dom"

import ChoresListPage from "./ChoresListPage";

const ChoresPage = (props) => {
    return <div>
        <Route path={props.match.path} exact component={ChoresListPage}/>
    </div>
}

export default ChoresPage