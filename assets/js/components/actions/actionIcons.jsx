import React from "react";

let iconsFetched = false;
function ensureIconsLoaded(){
    if (iconsFetched) {
        return
    }

    iconsFetched = true;

    fetch(location.origin + "/api/icons", {
        method: "GET",
        credentials: "same-origin"
    }).then((resp) => {
        return resp.text()
    }).then((txt) => {
        document.body.innerHTML = txt + document.body.innerHTML;
    });
}

let SVGIcon = (props) => {
    ensureIconsLoaded();
    return <svg style={props.style}>
        <use xlinkHref={`#${props.prefix}_${props.id}`} />
    </svg>
}

export const ProviderIcon = (props) => {
    return <SVGIcon prefix="provider" id={props.provider} />
};

export const GroupIcon = (props) => {
    return <SVGIcon prefix="group" id={props.group} />
};

export const ActionIcon = (props) => {
    return <SVGIcon prefix="action" id={`${props.provider}_${props.action}`}/>
}