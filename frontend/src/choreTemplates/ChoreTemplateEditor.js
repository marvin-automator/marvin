import React from "react";
import PropTypes from 'prop-types';
import {Grid, Form, Input, Button} from "semantic-ui-react";

class ChoreTemplateEditor extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            ...props.template
        }
    }

    handleChange = (e, {name, value}) => {
        this.setState({[name]: value})
    };

    render() {
        return <Grid>
            <Grid.Row>
                <Grid.Column width={6}>
                    <Form.Field
                        id='title'
                        control={Input}
                        name="name"
                        label='Template Name'
                        placeholder='Enter a name...'
                        value={this.state.name || ""}
                        onChange={this.handleChange}
                    />
                </Grid.Column>
                <Grid.Column width={2}>
                    <Button icon="save" content="Save" />
                </Grid.Column>
            </Grid.Row>
            <Grid.Row>
                <Grid.Column width="12">
                    <div id="editor" style={{backgroundColor: "blue", width: "100%", height: "100vh"}}>
                        {this.state.script}
                    </div>
                </Grid.Column>
            </Grid.Row>
        </Grid>
    }
}

ChoreTemplateEditor.propTypes = {
    template: PropTypes.shape({
        name: PropTypes.string,
        script: PropTypes.string,
    })
};

export default ChoreTemplateEditor

