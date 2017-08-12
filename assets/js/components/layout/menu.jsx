import React from "react"
import {Menu, Container, Image, Icon} from "semantic-ui-react"

const AppMenu = () => {
    return <Menu fixed="top" inverted>
        <Container icon="labeled">
            <Menu.Item as="a">
                <Image src="http://via.placeholder.com/64x64" size="mini" style={{marginRight: '1.5em'}}/>
            </Menu.Item>
            <Menu.Item as="a" href="/app/chores">Chores</Menu.Item>
            <Menu.Item as="a" href="/app/settings">Settings</Menu.Item>
            <Menu.Item as="a" href="/app/account" position="right"><Icon name="user" />{ACCOUNT_EMAIL}</Menu.Item>
            <Menu.Item as="a" href="/logout"><Icon name="log out"/>Log Out</Menu.Item>
        </Container>
    </Menu>
}

export default AppMenu
