import gql from "graphql-tag";

import {CHORE_TEMPLATE_ALL_FIELDS} from "../choreTemplates/query";

export const ALL_CHORE_FIELDS = gql`fragment allChoreFields on Chore {
    id
    name
    active
    template {
        ...AllTemplateFields
    }
    choreSettings {
        inputs {
            name
            value
        }
        triggers {
            provider
            group
            action
        }
    }
}
${CHORE_TEMPLATE_ALL_FIELDS}`;

export const GET_CHORES = gql`query choresList {
    chores {
        ...allChoreFields
    }
}

${ALL_CHORE_FIELDS}`;

export const CREATE_CHORE = gql`mutation createChore($templateId: String!, $name: String!, $inputs: [ChoreInput!]) {
    createChore(templateId: $templateId, name: $name, inputs: $inputs) {
        ...allChoreFields
    }
}

${ALL_CHORE_FIELDS}`;

export const GET_CHORE_BY_ID = gql`query getChore($choreId: String!) {
    choreById(id: $choreId) {
        ...allChoreFields
    }
}


${ALL_CHORE_FIELDS}`;

export const SET_CHORE_ACTIVE = gql`mutation setActive($choreId: String!, $active: Boolean!) {
    setChoreActive(id: $choreId, active: $active) {
        ...allChoreFields
    }
}

${ALL_CHORE_FIELDS}`;

export const GET_CHORE_LOGS = gql`query getLogs($id: String!, $upTo: String!, $count: Int!) {
  choreById(id: $id) {
    id
    logs(upTo: $upTo, count: $count) {
      id
      message
      type
      time
    }
  }
}`;

export const GET_LATEST_LOGS = gql`query getLatestLogs($id: String!, $count: Int!) {
  choreById(id: $id) {
    id
    logs(count: $count) {
      id
      message
      type
      time
    }
  }
}`;
