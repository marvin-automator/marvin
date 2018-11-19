import React from 'react'
import {Segment, Sidebar, Container } from 'semantic-ui-react'

import Menu from "./Menu"

const MainLayout = ({children}) => {
    return <Sidebar.Pushable as={Segment} style={{height: "100vh"}}>
        <Menu/>

        <Sidebar.Pusher>
            <Segment basic>
                <Container>
                    {children}
                </Container>
            </Segment>
        </Sidebar.Pusher>
    </Sidebar.Pushable>
};

export default MainLayout;
