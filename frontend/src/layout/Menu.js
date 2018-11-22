import React from 'react'
import {Icon, Menu, Segment, Sidebar } from 'semantic-ui-react'
import {Link} from "@reach/router";

const MarvinMenu = () => {
    return <Sidebar as={Menu} animation='overlay' icon='labeled' inverted vertical visible width='thin'>
        <Menu.Item as={Link} to="/">
            <Icon name='home' />
            Home
        </Menu.Item>
        <Menu.Item as={Link} to="/chores">
            <Icon name='tasks' />
            Chores
        </Menu.Item>
        <Menu.Item as={Link} to="/templates">
            <Icon name='file code outline' />
            Chore Templates
        </Menu.Item>
        <Menu.Item as={Link} to="/logs">
            <Icon name='list' />
            Logs
        </Menu.Item>
        <Menu.Item as={Link} to="/settings">
            <Icon name='cog' />
            Settings
        </Menu.Item>
        <Menu.Item as="a" href="/auth/logout" position="bottom">
            <Icon name='sign out' />
            Log out
        </Menu.Item>
    </Sidebar>
};

export default MarvinMenu;