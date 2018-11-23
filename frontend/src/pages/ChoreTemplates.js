import React from 'react';
import { Query } from 'react-apollo';
import { List, Grid, Button, Header} from 'semantic-ui-react'
import {Link} from "@reach/router";

import {GET_CHORE_TEMPLATES} from "../choreTemplates/query";

const ChoreTemplates = () => {
    return <Query query={GET_CHORE_TEMPLATES}>
        {({ loading, error, data }) => {
            if (loading) return <div>Loading...</div>;
            if (error) return <div>Error: {error}</div>;

            console.log(data);

            return (
                <Grid>
                    <Grid.Row>
                        <Grid.Column width={12}>
                            <Header as="h1">Chore Templates</Header>
                        </Grid.Column>
                        <Grid.Column width={4}>
                            <Button icon="plus" content="New template" as={Link} to="/templates/new" color="green" />
                        </Grid.Column>
                    </Grid.Row>
                    <Grid.Row>
                        <Grid.Column width={16}>
                            <List divided relaxed>
                                {data.ChoreTemplates.map((ct) => {
                                    return <List.Item key={ct.id}>
                                        <List.Icon name="file code outline" size='large' verticalAlign='middle' />
                                        <List.Content>
                                            <Button basic icon="ellipsis vertical" size="small" aria-label="actions" floated="right" compact/>
                                            <List.Header as={Link} to={`/templates/${ct.id}`}>{ct.name}</List.Header>
                                            <List.Description as={Link} to={`/templates/${ct.id}`}>{ct.created}</List.Description>
                                        </List.Content>
                                    </List.Item>
                                })}
                            </List>
                        </Grid.Column>
                    </Grid.Row>
                </Grid>
            )
        }}
    </Query>
};

export default ChoreTemplates;