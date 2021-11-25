import config from '../config';
import { ENVIRONMENT, Status, VERIFIED_CONTACT_STATUS } from 'constants/index';
import { DIRECTORY } from 'constants/index';
import TrisatestLogo from 'assets/images/gds-trisatest-logo.png';
import VaspDirectoryLogo from 'assets/images/gds-vaspdirectory-logo.png';
import dayjs from 'dayjs';
import crypto from 'crypto'

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
    return Object.fromEntries(Object.entries(data).map(([k, v]) => [k, (v / total).toFixed(2)]));
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

function verifiedContactStatus({ data, type = '', verifiedContact }) {

    // perform verified
    if (Object.keys(verifiedContact).includes(type.toLowerCase())) {
        return VERIFIED_CONTACT_STATUS.VERIFIED
    }

    // perform alternate verified
    if (Object.values(verifiedContact).includes(data?.email)) {
        return VERIFIED_CONTACT_STATUS.ALTERNATE_VERIFIED
    }

    return VERIFIED_CONTACT_STATUS.UNVERIFIED
}

const formatDate = (date) => date ? dayjs(date).format('DD-MM-YYYY') : 'N/A';

/**
 * 
 * @param {string} data string to hash
 */
function generateMd5(data = '') {
    return crypto.createHash('md5').update(data).digest("hex");
}

function currencyFormatter({ style = 'currency', currency = "USD" }) {
    return new Intl.NumberFormat('en-US', {
        style,
        currency,

    })
}

export { currencyFormatter as intlFormatter, verifiedContactStatus, generateMd5, formatDate, isValidHttpUrl, getDirectoryLogo, isTestNet, getDirectoryName, getDirectoryURL, getStatusClassName, formatDisplayedData, defaultEndpointPrefix, apiHost, getRatios, capitalizeFirstLetter, getCookie }
