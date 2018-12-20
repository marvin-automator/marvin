import React from "react";
import withState from "recompose/withState"
import { Query, Mutation } from "react-apollo";
import { navigate } from "@reach/router"

import {Form, Input, Button} from "semantic-ui-react";

import {Message} from "semantic-ui-react";

import {GET_TEMPLATE} from "../choreTemplates/query";
import {CREATE_CHORE} from "../chores/queries";

import ChoreTemplateEditor from "../choreTemplates/ChoreTemplateEditor";

const enhancer = withState("tempChore", "setTempChore", {name: "",});
const CreateChore = enhancer(({templateId, tempChore, setTempChore}) => {
    return <Query query={GET_TEMPLATE} variables={{id: templateId}}>
        {({loading, error, data}) => {
            if(loading) return <p>Loading...</p>;
            if(error) return <p>There was a problem: {error}</p>;

            data = data.ChoreTemplateById;

            const handleChange = (e, {name, value}) => setTempChore({...tempChore, [name]: value});
            let inputValues = Object.keys(tempChore).filter((k) => k.startsWith("input-")).map((key) => {
                return {
                    name: key.split("-")[1],
                    value: tempChore[key]
                }
            });
            return <Mutation mutation={CREATE_CHORE} variables={{name: tempChore.name, inputs: inputValues, templateId: templateId}}
                             onCompleted={({createChore:{id}}) => navigate(`/chores/${id }`)}>
                {(createChore, {error}) => {
                    return <React.Fragment>
                        <h1>Chore from: {data.name}</h1>
                        <Form onSubmit={createChore}>
                            <Form.Field id="chore-name" name="name" label="Chore Name" control={Input} onChange={handleChange} />
                            <h2>Template Parameters</h2>
                            {data.templateSettings.inputs.map((input) => {
                                return <React.Fragment key={input.name}>
                                    <Form.Field id={input.name} name={`input-${input.name}`} label={input.name} control={Input}
                                                onChange={handleChange} />
                                    {input.description ? <p>{input.description}</p> : null}
                                </React.Fragment>
                            })}
                            <p><Button positive icon="add" type="submit" content="Create" /></p>
                        </Form>
                    </React.Fragment>
                }}
            </Mutation>
        }}
    </Query>
});

const updateCache = (cache, { data: { updateChoreTemplate } }) => {
    //TODO: Update the cache here.
};

export default CreateChore;