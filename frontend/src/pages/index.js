import React from "react";
import { Router} from "@reach/router";

import ChoreTemplates from "./ChoreTemplates"
import CreateChoreTemplate from "./CreateChoreTemplate";
import UpdateChoreTemplate from "./UpdateChoreTemplate";

import CreateChore from "./CreateChore"

const Routes = () => {
    return <Router>
        <Home path="/" />
        <ChoreTemplates path="/templates" />
        <CreateChoreTemplate path="/templates/new"/>
        <UpdateChoreTemplate path="/templates/:id" />

        <CreateChore path="/chores/new/:templateId" />
    </Router>
};

const Home = () => {
    return <div path="/">
        Home
    </div>
}

export default Routes;
