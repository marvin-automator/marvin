import React from "react"
import {Route} from "react-router-dom"

import ChoresListPage from "./ChoresListPage";
import NewChorePage from "./NewChorePage"

const ChoresPage = (props) => {
    return <div>
        <Route path={props.match.path} exact component={ChoresListPage}/>
        <Route path={props.match.path + "/new"} component={NewChorePage} />
    </div>
}

export default ChoresPage