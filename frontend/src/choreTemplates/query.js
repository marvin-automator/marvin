import gql from "graphql-tag";

export const GET_CHORE_TEMPLATES = gql`query {
    ChoreTemplates {
        name
        id
        created
    }
}`;

export const CREATE_TEMPLATE = gql`
mutation create($name: String!, $script: String!) {
  createChoreTemplate(name: $name, script: $script) {
    id
    name
    script
    created
    templateSettings {
      inputs {name description}
      triggers{provider group action}
    }
  }
}
`;