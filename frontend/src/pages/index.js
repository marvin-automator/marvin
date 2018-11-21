import React from "react";
import { Router, Link } from "@reach/router";

import ChoreTemplates from "./ChoreTemplates"

const Routes = () => {
    return <Router>
        <Home path="/" />
        <ChoreTemplates path="/templates" />
    </Router>
};

const Home = () => {
    return <div path="/">
        Home
    </div>
}

export default Routes;
