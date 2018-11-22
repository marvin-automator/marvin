import React from "react";
import { Router} from "@reach/router";

import ChoreTemplates from "./ChoreTemplates"
import CreateChoreTemplate from "./CreateChoreTemplate";

const Routes = () => {
    return <Router>
        <Home path="/" />
        <ChoreTemplates path="/templates" />
        <CreateChoreTemplate path="/templates/new"/>
    </Router>
};

const Home = () => {
    return <div path="/">
        Home
    </div>
}

export default Routes;
