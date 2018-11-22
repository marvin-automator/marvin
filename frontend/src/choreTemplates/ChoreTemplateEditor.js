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

    handleSubmit = () => {
        this.props.onSave(this.state);
    };

    render() {
        return <Grid as={Form} onSubmit={this.handleSubmit}>
            <Grid.Row>
                <Grid.Column width={16}>
                    <Form.Field
                        id='title'
                        control={Input}
                        name="name"
                        label='Template Name'
                        placeholder='Enter a name...'
                        value={this.state.name || ""}
                        onChange={this.handleChange}
                        inline={false}
                    />
                </Grid.Column>
            </Grid.Row>
            <Grid.Row>
                <Grid.Column width="12">
                    <div id="editor" style={{backgroundColor: "blue", width: "100%", height: "80vh"}}>
                        {this.state.script}
                    </div>
                </Grid.Column>
            </Grid.Row>
            <Grid.Row>
                <Grid.Column><Button icon="save" content="Save" primary type="submit" /></Grid.Column>
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

