export * from './array';


function formatDisplayedData(target) {
    if (typeof target === "boolean") {
        return target.toString()
    } else if (Array.isArray(target)) {
        return target.length ? target.toString() : "N/A"
    } else if (typeof target === "string") {
        return target ? target.trim() : "N/A"
    }

    return target ? target : "N/A"
}

export { formatDisplayedData }