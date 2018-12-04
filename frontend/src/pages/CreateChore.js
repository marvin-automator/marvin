import React from "react";

import { Query, Mutation } from "react-apollo";

import {Form, Input, Divider, Button} from "semantic-ui-react";

import {Message} from "semantic-ui-react";

import {GET_TEMPLATE} from "../choreTemplates/query";
import {CREATE_CHORE} from "../chores/queries";

import ChoreTemplateEditor from "../choreTemplates/ChoreTemplateEditor";

const CreateChore = ({templateId}) => {
    return <Query query={GET_TEMPLATE} variables={{id: templateId}}>
        {({loading, error, data}) => {
            if(loading) return <p>Loading...</p>;
            if(error) return <p>There was a problem: {error}</p>;

            data = data.ChoreTemplateById;
            console.log(data);

            return <React.Fragment>
                <h1>Chore from: {data.name}</h1>
                <Form>
                    <Form.Field id="chore-name" name="name" label="Chore Name" control={Input} />
                    <h2>Template Parameters</h2>
                    {data.templateSettings.inputs.map((input) => {
                        return <React.Fragment>
                            <Form.Field id={input.name} name={input.name} label={input.name} control={Input} />
                            {input.description ? <span>{input.description}</span> : null}
                        </React.Fragment>
                    })}
                    <p><Button positive icon="add" type="submit" content="Create" /></p>
                </Form>
            </React.Fragment>
        }}
    </Query>
};

const updateCache = (cache, { data: { updateChoreTemplate } }) => {

};

export default CreateChore;