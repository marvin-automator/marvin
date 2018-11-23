import React from "react";
import { Mutation } from "react-apollo";
import {Message} from "semantic-ui-react";
import { navigate } from "@reach/router"

import {CREATE_TEMPLATE, GET_CHORE_TEMPLATES} from "../choreTemplates/query";

import ChoreTemplateEditor from "../choreTemplates/ChoreTemplateEditor";

const CreateChoreTemplate = () => {
    return <Mutation mutation={CREATE_TEMPLATE} update={updateCache} onCompleted={({createChoreTemplate: {id}}) => navigate(`/templates/${id}`)}>
        {(createChoreTemplate, {error}) => {
            let errMsg = error ? <Message error>{error.toString()}</Message> : null;
            return <React.Fragment>
                {errMsg}
                <ChoreTemplateEditor onSave={(template) => createChoreTemplate({variables: template})}/>
            </React.Fragment>
        }}
    </Mutation>
};

const updateCache = (cache, { data: { createChoreTemplate } }) => {
    let ChoreTemplates;
    try {
        ChoreTemplates = cache.readQuery({query: GET_CHORE_TEMPLATES}).ChoreTemplates;
    } catch {
        // We don't want to set an initial state for the cache here, because that'd lead future queries to not hit the server.
        return
    }
    cache.writeQuery({
        query: GET_CHORE_TEMPLATES,
        data: { ChoreTemplates: ChoreTemplates.concat([createChoreTemplate]) }
    });
};

export default CreateChoreTemplate;