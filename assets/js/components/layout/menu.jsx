import React from "react"
import Menu from "antd/lib/menu"
import Icon from "antd/lib/icon"

const AppMenu = () => {
    return <Menu theme="dark" mode="horizontal" style={{ lineHeight: '64px' }}>
        <Menu.Item>Chores</Menu.Item>
        <Menu.Item>Settings</Menu.Item>
        <Menu.Item style={{float:"right"}}><Icon type="logout" />Log Out</Menu.Item>
        <Menu.Item style={{float:"right"}}><Icon type="user" />{ACCOUNT_EMAIL}</Menu.Item>
    </Menu>
}

export default AppMenu
