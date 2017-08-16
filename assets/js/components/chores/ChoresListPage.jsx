import React from "react"
import {Button} from "semantic-ui-react"

const ChoresListPage = () => {
    return <div>
        <h1>
            Chores
            <span style={{float: "right"}}>
                <Button positive as="a" href="/chores/new" size="small" content="new" icon="plus"/>
            </span>
        </h1>
    </div>
}

export default ChoresListPage