import config from '../config';
import { ENVIRONMENT, Status } from '../constants';
import { DIRECTORY } from '../constants';
import TrisatestLogo from '../assets/images/gds-trisatest-logo.png';
import VaspDirectoryLogo from '../assets/images/gds-vaspdirectory-logo.png';

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

const getRatios = (data) => {
    const total = Object.values(data).reduce((acc, x) => acc + x);
    return Object.fromEntries(Object.entries(data).map(([k, v]) => [k, (v / total).toFixed(1)]));
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

function getStatusClassName(status = '') {
    switch (status) {
        case Status.VERIFIED:
            return 'bg-success'
        case Status.SUBMITTED:
            return 'bg-secondary'
        case Status.PENDING_REVIEW:
        case Status.EMAIL_VERIFIED:
            return 'bg-warning'
        case Status.ERRORED:
        case Status.REJECTED:
            return 'bg-danger'
        case Status.APPEALED:
            return 'bg-primary'
        case Status.REVIEWED:
        case Status.ISSUING_CERTIFICATE:
            return 'bg-info'
        default:
            return undefined
    }
}

function isTestNet() {
    return config.IS_TESNET
}

function getDirectoryName() {
    return isTestNet() ? DIRECTORY.VASP_DIRECTORY : DIRECTORY.TRISATEST
}

function getDirectoryURL() {
    return isTestNet() ? "https://admin.vaspdirectory.net" : "https://admin.trisatest.net"
}

const getDirectoryLogo = () => {
    return isTestNet() ? TrisatestLogo : VaspDirectoryLogo
}

function isValidHttpUrl(string) {
    let url;

    try {
        url = new URL(string);
    } catch (_) {
        return false;
    }

    return url.protocol === "http:" || url.protocol === "https:";
}

/**
 * 
 * @param {Object} contact contact you want to verify
 * @param {Object} verifiedContact list of verified contacts
 * @returns boolean
 */
function isVerifiedContact(contact, verifiedContact) {
    const verifiedContacts = typeof verifiedContact === 'object' ? Object.values(verifiedContact) : []
    return verifiedContacts.includes(contact.email)
}

export { isVerifiedContact, isValidHttpUrl, getDirectoryLogo, isTestNet, getDirectoryName, getDirectoryURL, getStatusClassName, formatDisplayedData, defaultEndpointPrefix, apiHost, getRatios, capitalizeFirstLetter, getCookie }