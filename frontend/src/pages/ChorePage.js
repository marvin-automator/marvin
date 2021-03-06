import React from "react";
import { Grid, List, Button, Header} from 'semantic-ui-react'

import {GET_CHORE_BY_ID} from "../chores/queries";
import { Query } from 'react-apollo';

import ChoreLogs from "../chores/ChoreLogs";
import StatusSwitch from "../chores/StatusSwitch";

const ChorePage = ({id}) => {
    return <Query query={GET_CHORE_BY_ID} variables={{choreId: id}}>
        {({loading, error, data: {choreById}}) => {
            if(loading) return <p>Loading...</p>

            return <Grid>
                <Grid.Row>
                    <Grid.Column width={12}>
                        <Header as="h1">{choreById.name}</Header>
                    </Grid.Column>
                    <Grid.Column width={4}>
                        <StatusSwitch id={id} active={choreById.active}/>
                    </Grid.Column>
                </Grid.Row>

                <Grid.Row>
                    <Grid.Column width={12}>
                        <ChoreLogs id={id}/>
                    </Grid.Column>
                </Grid.Row>
            </Grid>
        }}
    </Query>
};

export default ChorePage;