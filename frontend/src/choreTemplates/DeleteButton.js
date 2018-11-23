import React from "react";
import {Mutation} from "react-apollo";

import {DELETE_TEMPLATE, GET_CHORE_TEMPLATES, GET_TEMPLATE} from "./query";

const DeleteButton = ({id, component, ...rest}) => {
    return <Mutation mutation={DELETE_TEMPLATE} variables={{id}} update={updateCache.bind(null, id)}>
        {(deleteTemplate) => {
            let Element = component;
            return <Element {...rest} onClick={deleteTemplate} />
        }}
    </Mutation>
};

export default DeleteButton;

const updateCache = (id, cache, result) => {
    let ChoreTemplates;
    try {
        ChoreTemplates = cache.readQuery({query: GET_CHORE_TEMPLATES}).ChoreTemplates;
    } catch {
        // We don't want to set an initial state for the cache here, because that'd lead future queries to not hit the server.
        return
    }

    console.log(cache, result)
    cache.writeQuery({
        query: GET_CHORE_TEMPLATES,
        data: { ChoreTemplates: ChoreTemplates.filter((ct) => ct.id !== id)}
    });
};