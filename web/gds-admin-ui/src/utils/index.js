import config from '../config';
import { ENVIRONMENT } from '../constants';

export * from './array';


function defaultEndpointPrefix() {
    if (config.GDS_API_URL) {
        return config.GDS_API_URL
    }

    switch (process.env.NODE_ENV) {
        case ENVIRONMENT.DEV:
            return "http://localhost:4434/v2"
        case ENVIRONMENT.PROD:
            if (config.IS_TESTNET) {
                return "https://api.admin.trisatest.net/v2";
            } else {
                return "https://api.admin.vaspdirectory.net/v2";
            }
        default:
            throw new Error("Could not identify the api prefix");
    }
}

function apiHost() {
    const url = new URL(config.GDS_API_URL)
    return url.hostname;
}

function formatDisplayedData(target) {
    if (typeof target === "boolean") {
        return target.toString()
    } else if (Array.isArray(target)) {
        return target.length ? target.join(', ') : "N/A"
    } else if (typeof target === "string") {
        return target ? target.trim() : "N/A"
    }

    return target ? target : "N/A"
}

const convertCountsToPercentages = (data) => {
    const total = Object.values(data).reduce((acc, x) => acc + x);
    return Object.fromEntries(Object.entries(data).map(([k, v]) => [k, v / total]));
}

function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1).toLowerCase();
}

function getCookie(name = '') {
    if (document.cookie && document.cookie !== '') {
        const cookies = document.cookie.split(';');
        for (let i = 0; i < cookies.length; i++) {
            const cookie = cookies[i].trim();
            if (cookie.substring(0, name.length + 1) === (name + '=')) {
                return decodeURIComponent(cookie.substring(name.length + 1));
            }
        }
    }
    return '';
}




export { formatDisplayedData, defaultEndpointPrefix, apiHost, convertCountsToPercentages, capitalizeFirstLetter, getCookie }