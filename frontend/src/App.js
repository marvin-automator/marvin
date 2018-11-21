import React, { Component } from 'react';
import ApolloClient from 'apollo-boost';
import { ApolloProvider } from 'react-apollo';

import MainLayout from "./layout/MainLayout";
import Routes from "./pages"

const client = new ApolloClient();

const App = () => {
    return <ApolloProvider client={client}>
        <MainLayout>
            <Routes/>
        </MainLayout>
    </ApolloProvider>
}


export default App;
