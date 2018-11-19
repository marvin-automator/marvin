import React from 'react'
import {Icon, Menu, Segment, Sidebar } from 'semantic-ui-react'

const MarvinMenu = () => {
    return <Sidebar as={Menu} animation='overlay' icon='labeled' inverted vertical visible width='thin'>
        <Menu.Item as='a'>
            <Icon name='home' />
            Home
        </Menu.Item>
        <Menu.Item as='a'>
            <Icon name='tasks' />
            Chores
        </Menu.Item>
        <Menu.Item as='a'>
            <Icon name='file code outline' />
            Chore Templates
        </Menu.Item>
        <Menu.Item as='a'>
            <Icon name='list' />
            Logs
        </Menu.Item>
        <Menu.Item as='a'>
            <Icon name='cog' />
            Settings
        </Menu.Item>
        <Menu.Item as='a'>
            <Icon name='sign out' />
            Log out
        </Menu.Item>
    </Sidebar>
};

export default MarvinMenu;