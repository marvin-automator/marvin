import React from "react"
import {Message} from "semantic-ui-react"

import {TriggerPicker} from "../actions/actionPicker"

class NewChoresPage extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            actions: []
        }
    }

    render() {
        return <div>
            <h1>New Chore</h1>
            {this.getScreen()}
        </div>
    }

    getScreen() {
        if (this.state.actions.length == 0) {
            return [
                <Message icon="rocket" header="Select a trigger" content="This will start of your chore when a certain event occurs." key="p1"/>,
                <TriggerPicker onSelect={console.log.bind(console)} key="p2" />
            ]
        }
    }
}

export default NewChoresPage