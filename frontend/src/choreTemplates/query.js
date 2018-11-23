import gql from "graphql-tag";

const CHORE_TEMPLATE_ALL_FIELDS = gql`fragment AllTemplateFields on ChoreTemplate {
    id
    name
    script
    created
    templateSettings {
        inputs {name description}
        triggers{provider group action}
    }
}`;

export const GET_CHORE_TEMPLATES = gql`query {
    ChoreTemplates {
        ...AllTemplateFields
    }
}
${CHORE_TEMPLATE_ALL_FIELDS}
`;

export const CREATE_TEMPLATE = gql`
mutation create($name: String!, $script: String!) {
  createChoreTemplate(name: $name, script: $script) {
      ...AllTemplateFields
  }
}
${CHORE_TEMPLATE_ALL_FIELDS}`;

export const UPDATE_TEMPLATE = gql`
mutation updateTemplate($id: String!, $name: String!, $script: String!) {
    updateChoreTemplate(id: $id, name: $name, script: $script) {
        ...AllTemplateFields
    }
}
${CHORE_TEMPLATE_ALL_FIELDS}
`;

export const GET_TEMPLATE = gql`
query getTemplate($id: String!) {
    ChoreTemplateById(id: $id) {
        ...AllTemplateFields
    }
}
${CHORE_TEMPLATE_ALL_FIELDS}
`