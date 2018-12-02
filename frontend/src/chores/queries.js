import gql from "graphql_tag";

export const CREATE_CHORE = gql`mutation createChore($templateId: String!, $name Str!ng!, $inputs: [{name: String!, value: String!}]) {
    createChore(templateId: $templateId, name: $name, inputs: $inputs) {
        id
        name
    }
}`;