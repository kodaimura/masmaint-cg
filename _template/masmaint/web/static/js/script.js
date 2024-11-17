export const getJaTime = () => {
    const date = new Date();
    const utcOffset = date.getTimezoneOffset() * 60000;
    const jaTime = new Date(date.getTime() + utcOffset + 9 * 3600000);
    return jaTime;
};

export const nullToEmpty = (s) => {
    return (s == null) ? '' : s;
}

export const emptyToNull = (s) => {
    return (s == '') ? null : s;
}

export const parseIntOrReturnOriginal = (value) => {
    if (value === "") {
        return null;
    } 
    if (typeof value === "string" && /^-?\d+$/.test(value)) {
        return parseInt(value);
    }
    return value;
}

export const parseFloatOrReturnOriginal = (value) => {
    if (value === "") {
        return null;
    }
    if (typeof value === "string" && !isNaN(value) && value.trim() !== "") {
        return parseFloat(value);
    }
    return value;
}

export const parseBoolOrReturnOriginal = (value) => {
    if (value === "") {
        return null;
    }
    if (value === "true") {
        return true;
    }
    if (value === "false") {
        return false;
    }
    return value;
}