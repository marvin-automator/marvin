import React from "react";
import PropTypes from 'prop-types';
import {Grid, Form, Input, Button} from "semantic-ui-react";
const CodeEditor = React.lazy(() => import('./CodeEditor'));

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
                    <React.Suspense fallback={<div>Loading...</div>}>
                    <CodeEditor width="100%" height="80vh" script={this.state.script} onChange={(v) => this.handleChange(null, {name: "code", value: v})} />
                    </React.Suspense>
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

