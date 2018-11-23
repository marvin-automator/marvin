import React from "react";
import { Query, Mutation } from "react-apollo";
import {Message} from "semantic-ui-react";

import {UPDATE_TEMPLATE, GET_TEMPLATE} from "../choreTemplates/query";

import ChoreTemplateEditor from "../choreTemplates/ChoreTemplateEditor";

const UpdateChoreTemplate = ({id}) => {
    return <Mutation mutation={UPDATE_TEMPLATE} variables={{id}} update={updateCache}>
        {(updateChoreTemplate, {error}) => {
            let errMsg = error ? <Message error>{error.toString()}</Message> : null;
            return <React.Fragment>
                {errMsg}
                <Query query={GET_TEMPLATE} variables={{id}}>
                    {({loading, error, data}) => {
                        if (loading) return <div>Loading...</div>;
                        if(error) {
                            return <Message error>
                                    {error.toString()}
                            </Message>
                        }
                        return <ChoreTemplateEditor
                            template={data.ChoreTemplateById}
                            onSave={(template) => updateChoreTemplate({variables: {...template, id: id}})}
                        />
                    }}
                </Query>
            </React.Fragment>
        }}
    </Mutation>
};

const updateCache = (cache, { data: { updateChoreTemplate } }) => {
    let ChoreTemplates;
    try {
        ChoreTemplates = cache.readQuery({query: GET_CHORE_TEMPLATES}).ChoreTemplates;
        cache.writeQuery({
            query: GET_CHORE_TEMPLATES,
            data: { ChoreTemplates: ChoreTemplates.mqp((tpl) => {
                    return tpl.id === updateChoreTemplate.id ? updateChoreTemplate : tpl;
                })
            }
        });
    } catch {
        // nothing to do here, the templates weren't fetched before, so nothing to update.
    }

    cache.writeQuery({
        query: GET_TEMPLATE,
        data: {ChoreTemplateById: updateChoreTemplate},
        variables: {id: updateChoreTemplate.id}
    })
};

export default UpdateChoreTemplate;