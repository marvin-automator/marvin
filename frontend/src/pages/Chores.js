import React from 'react';
import { Query } from 'react-apollo';
import { List, Grid, Header, Icon} from 'semantic-ui-react'
import {Link} from "@reach/router";

import {GET_CHORES} from "../chores/queries";

const Chores = () => {
    return <Query query={GET_CHORES}>
        {({ loading, error, data }) => {
            if (loading) return <div>Loading...</div>;
            if (error) return <div>Error: {error}</div>;

            return (
                <Grid>
                    <Grid.Row>
                        <Grid.Column width={16}>
                            <Header as="h1">Chores</Header>
                        </Grid.Column>
                    </Grid.Row>
                    <Grid.Row>
                        <Grid.Column width={16}>
                            <List divided relaxed>
                                {data.chores.map((chore) => {
                                    return <List.Item key={chore.id}>
                                        <List.Icon name="tasks" size='large' verticalAlign='middle' />
                                        <List.Content>
                                            <List.Header as={Link} to={`/chores/${chore.id}`}>{chore.name}</List.Header>
                                            <List.Description><Icon name="file code outline"/>{chore.template.name}</List.Description>
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

export default Chores;