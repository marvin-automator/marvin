import React from 'react';
import { Query } from 'react-apollo';
import { List } from 'semantic-ui-react'

import {GET_CHORE_TEMPLATES} from "../choreTemplates/query";

const ChoreTemplates = () => {
    return <Query query={GET_CHORE_TEMPLATES}>
        {({ loading, error, data }) => {
            if (loading) return <div>Loading...</div>;
            if (error) return <div>Error: {error}</div>;

            console.log(data);

            return (
                <div>
                    <h1>Chore Templates</h1>
                    <List divided relaxed>
                        {data.ChoreTemplates.map((ct) => {
                            return <List.Item key={ct.id}>
                                <List.Icon name="file code outline" size='large' verticalAlign='middle' />
                                <List.Content>
                                    <List.Header>{ct.name}</List.Header>
                                    <List.Description>{ct.created}</List.Description>
                                </List.Content>
                            </List.Item>
                        })}
                    </List>
                </div>
            )
        }}
    </Query>
};

export default ChoreTemplates;