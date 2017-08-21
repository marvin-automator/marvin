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
        let div = document.createElement("div");
        div.innerHTML = txt;
        document.body.insertBefore(div.firstChild, document.body.firstChild)
    });
}

let SVGIcon = (props) => {
    ensureIconsLoaded();
    return <span><svg style={props.style}>
        <use xlinkHref={`#${props.prefix}_${props.id}`} />
    </svg></span>
};

export const ProviderIcon = (props) => {
    return <SVGIcon prefix="provider" id={props.provider} />
};

export const GroupIcon = (props) => {
    return <SVGIcon prefix="group" id={props.group} style={props.style} />
};

export const ActionIcon = (props) => {
    return <SVGIcon prefix="action" id={`${props.provider}_${props.action}`} style={props.style}/>
}