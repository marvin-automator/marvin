import React from "react";
import gql from "graphql-tag";
import {graphql} from "react-apollo"
import {Grid, Card, Accordion, Image, Loader, Segment, Dimmer} from "semantic-ui-react"
import withProps from "recompose/withProps";

import {ActionIcon, GroupIcon} from "./actionIcons";

const actionFragment = gql`
fragment ActionDetails on Action {
    key
    name
    description
}`;

const actionsQuery = gql`
query ActionsQuery {
    groups {
        name
        provider
        actions {
            ... ActionDetails
        }
    }
}
${actionFragment}`;

const triggersQuery = gql`
query ActionsQuery {
    groups {
        name
        provider
        triggers {
            ... ActionDetails
        }
    }
}
${actionFragment}`;

let ActionsGrid = (props) => {
    return <Grid columns={4}>
        <Grid.Row>
            {props.actions.map((action) => {
                let resultAction = {action: action.key, provider: props.provider}
                return <Grid.Column key={action.key}>
                    <Card onClick={() => props.onSelect(resultAction)}>
                        <Image>
                            <ActionIcon action={action.key} provider={props.provider} />
                        </Image>
                        <Card.Content>
                            <Card.Header>{action.name}</Card.Header>
                            <Card.Description>{action.description}</Card.Description>
                        </Card.Content>
                    </Card>
                </Grid.Column>
            })}
        </Grid.Row>
    </Grid>
};

let BaseActionSelect = (props) => {
    if (props.data.loading) {
        return <Segment>
            <Dimmer active inverted>
                <Loader inverted>Loading</Loader>
            </Dimmer>
            <div style={{display:"hidden"}}>hidden</div>
        </Segment>
    }

    console.log(props.data)

    let panels = props.data.groups.map((group) => {
        return {
            title: group.name,
            content: <ActionsGrid actions={group[props.field]} provider={group.provider} onSelect={props.onSelect} />,
            key: group.name,
        }
    });
    return <Accordion panels={panels} styled exclusive={false} fluid />
};

export const ActionPicker = withProps({field:"actions"})(
    graphql(actionsQuery)(BaseActionSelect)
);
export const TriggerPicker = withProps({field:"triggers"})(
    graphql(triggersQuery)(BaseActionSelect)
);