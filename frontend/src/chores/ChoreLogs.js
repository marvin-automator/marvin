import React from "react";
import {Table, Segment, Loader, Header, Icon, Button} from "semantic-ui-react";
import {Query, Mutation} from "react-apollo"

import {GET_CHORE_LOGS, GET_LATEST_LOGS, CLEAR_CHORE_LOGS} from "./queries";

function wait(n) {
    return new Promise((resolve) => {
        setTimeout(resolve, n)
    })
}

export default class ChoreLogs extends React.Component {
    constructor() {
        super();

        this.perPage = 20;
        this.state = {
            segments: [
                (new Date()).toISOString().split(".")[0] + "Z"
            ]
        };
    }

    componentWillUnmount() {
        clearInterval(this.intervalHandle);
    }

    render() {
        return <>
            <Segment attached="top">
                <Header as="h2">Logs</Header>
                <ClearBtn id={this.props.id}/>
            </Segment>
            <Table compact="very" attached="bottom" verticalAlign="top">
                <Table.Body>
                    {this.state.segments.map((time, index) => {
                        return this.makeSegment(time, index === this.state.segments.length - 1)
                    })}
                </Table.Body>
            </Table>
        </>
    }

    makeSegment(time, shouldRefresh){
        return <Query key={time}
                      query={shouldRefresh ? GET_LATEST_LOGS : GET_CHORE_LOGS}
                      variables={{id: this.props.id, count: this.perPage}}
                      pollInterval={shouldRefresh ? 5000 : 0}>
            {({data, loading, error}) => {
                if(loading && !data.choreById) {
                    return <Table.Row>
                        <Table.Cell><Loader inline size="tiny" /></Table.Cell>
                    </Table.Row>
                }

                let logs = data.choreById.logs;
                return logs.map((log) => {
                    return <Table.Row negative={log.type === "error"} key={log.id}>
                        <Table.Cell verticalAlign="top">
                            {{
                                "info": <Icon name="info circle" color="blue"/>,
                                "error": <Icon name="exclamation triangle" color="red"/>
                            }[log.type]}
                        </Table.Cell>
                        <Table.Cell verticalAlign="top">{new Date(log.time).toLocaleTimeString()}</Table.Cell>
                        <Table.Cell verticalAlign="top">{log.message}</Table.Cell>
                    </Table.Row>
                });
            }}
        </Query>
    }
};

const ClearBtn = ({id}) => {
    return <Mutation mutation={CLEAR_CHORE_LOGS} variables={{id}} >
        {(clearLogs) => {
            return <Button content="Clear" onClick={clearLogs} />
        }}
    </Mutation>
}
