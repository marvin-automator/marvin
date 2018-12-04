import React from "react";
import {Button, Dropdown, Menu} from "semantic-ui-react";
import {Link} from "@reach/router";

import DeleteButton from "./DeleteButton"

const ActionsButton = ({id}) => {
    return <Dropdown button
                     icon="ellipsis vertical"
                     size="small"
                     aria-label="actions"
                     className="compact right floated basic"
    >
        <Dropdown.Menu>
            <Dropdown.Item icon="add" content="Create a chore from this template" as={Link} to={`/chores/new/${id}`} />
            <DeleteButton component={Dropdown.Item} icon="trash" content="Delete" id={id}/>
        </Dropdown.Menu>
    </Dropdown>

};

export default ActionsButton;
