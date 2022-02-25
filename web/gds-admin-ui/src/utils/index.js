import config from '../config';
import { ENVIRONMENT, Status, VERIFIED_CONTACT_STATUS } from 'constants/index';
import { DIRECTORY } from 'constants/index';
import TrisatestLogo from 'assets/images/gds-trisatest-logo.png';
import VaspDirectoryLogo from 'assets/images/gds-vaspdirectory-logo.png';
import dayjs from 'dayjs';
import crypto from 'crypto'
import toast from 'react-hot-toast';
import { downloadFile, generateCSV } from 'helpers/api/utils';

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
    const total = Object.values(data).reduce((acc, x) => acc + x, 0);
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
    return config.IS_TESTNET
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

function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
}

function getBase64Size(str) {
    const buffer = Buffer.from(`${str}`, 'base64');
    return buffer.length
}

/**
 * Copy an element to the clipboard
 * @param {string} target item to copy to clipboard
 * @returns Promise<void>
 */
async function copyToClipboard(target = '') {
    try {
        await navigator.clipboard.writeText(target);
        toast.success('Copied to clipboard')
    } catch (err) {
        throw err
    }
}

function exportToCsv(rows) {
    const { verified_contacts, ...rest } = rows[0]

    let rowHeader = Object.keys(rest)

    const _rows = rows.map(row => {
        const { verified_contacts, ...rest } = row
        return Object.values(rest)
    })
    _rows.unshift(rowHeader)

    let csvFile = '';
    for (let i = 0; i < _rows.length; i++) {
        csvFile += generateCSV(_rows[i]);
    }
    const filename = `${dayjs().format("YYYY-MM-DD")}-directory.csv`
    downloadFile(csvFile, filename, 'text/csv;charset=utf-8;')
}

function isValidIvmsAddress(address) {
    if (address) {
        return !!(address.country && address.address_type)
    }
    return false;
}

function hasAddressField(address) {
    if (isValidIvmsAddress(address) && !hasAddressLine(address)) {
        return !!(address.street_name && (address.building_number || address.building_name))
    }
    return false
}

function hasAddressLine(address) {
    if (isValidIvmsAddress(address)) {
        return Array.isArray(address.address_line) && address.address_line.length > 0
    }
    return false;
}

function hasAddressFieldAndLine(address) {
    if (hasAddressField(address) && hasAddressLine(address)) {
        console.warn("cannot render address")
        return true
    }
    return false
}

const getMustComplyRegulations = (status) => status ? "must" : "must not"
const getConductsCustomerKYC = (status) => status ? "does" : "does not"
const getMustSafeguardPii = (status) => status ? "must" : "is not required to"
const getSafeguardPii = (status) => status ? "does" : "does not"

function isOptionAvailable(verificationStatus = "") {
    if (!verificationStatus) {
        return false;
    }
    return ["NO_VERIFICATION", "SUBMITTED", "EMAIL_VERIFIED", "PENDING_REVIEW", "ERRORED"].includes(verificationStatus)
}

export const validateIsoCode = (cc = '') => {
    if (typeof cc === 'string' && cc.length !== 2) {
        const matches = cc.match(/\b(\w)/g);
        const acronym = matches?.join('')
        return acronym?.length === 2 ? acronym : ''
    }

    return cc
}

export { isOptionAvailable, getMustComplyRegulations, getConductsCustomerKYC, getMustSafeguardPii, getSafeguardPii, isValidIvmsAddress, hasAddressField, hasAddressLine, hasAddressFieldAndLine, exportToCsv, copyToClipboard, getBase64Size, formatBytes, currencyFormatter as intlFormatter, verifiedContactStatus, generateMd5, formatDate, isValidHttpUrl, getDirectoryLogo, isTestNet, getDirectoryName, getDirectoryURL, getStatusClassName, formatDisplayedData, defaultEndpointPrefix, apiHost, getRatios, capitalizeFirstLetter, getCookie }
