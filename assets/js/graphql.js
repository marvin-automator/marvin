import {execute, HttpLink} from 'apollo-link';
import ApolloClient, {createNetworkInterface} from 'apollo-client'


let uri = location.origin + "/api/graphql";
export const link = new HttpLink({uri});
export const fetcher = (operation) => execute(link, operation);

let netInt = createNetworkInterface({
    uri: uri,
    opts: {
        credentials: 'same-origin',
    }
});
export const client = new ApolloClient({
    networkInterface: netInt,
    connectToDevTools: true
})