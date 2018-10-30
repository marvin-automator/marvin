let marvin = {
    _triggers: [],
    _inputs: [],
    _valueResolvers: {},
    isSetup: true
};

marvin.input = (name, description) => {
    marvin._inputs.push({name, description});
    return new Promise((resolve, reject) => {
        marvin._valueResolvers[name] = {resolve, reject};
    })
};

marvin._resolveInput = (name, value) => {
    marvin._valueResolvers[name].resolve(value);
    delete marvin._valueResolvers[name]
};

marvin._rejecUnknownInputs = () => {
    Object.keys(marvin._valueResolvers).forEach((name) => {
        marvin._valueResolvers[name].reject(`Configuration input ${name} is not defined.`)
    })
}