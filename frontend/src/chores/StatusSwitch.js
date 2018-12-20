import React from "react";
import {Checkbox} from "semantic-ui-react";
import { Mutation } from 'react-apollo';

import {SET_CHORE_ACTIVE, GET_CHORE_BY_ID, GET_CHORES} from "./queries";

const StatusSwitch = ({active, id}) => {
    return <Mutation mutation={SET_CHORE_ACTIVE} variables={{choreId: id, active: !active}} update={updateCache}>
        {(switchStatus, {loading}) => {
            return <span>
                <Checkbox toggle checked={active} disabled={loading} onChange={() => switchStatus()}
                          label={active ? "active" : "inactive"} />
            </span>
        }}
    </Mutation>
};

const updateCache = (cache, { data: { setChoreActive } }) => {
    let chores;
    try {
        chores = cache.readQuery({query: GET_CHORES}).chores;
    } catch {
        // We don't want to set an initial state for the cache here, because that'd lead future queries to not hit the server.
        return
    }
    cache.writeQuery({
        query: GET_CHORES,
        data: { chores: chores.map((c) => c.id === setChoreActive.id ? setChoreActive : c)}
    });
    cache.writeQuery({
        query: GET_CHORE_BY_ID,
        data: {choreById: setChoreActive},
        variables: {id: setChoreActive.id}
    })
};

export default StatusSwitch;