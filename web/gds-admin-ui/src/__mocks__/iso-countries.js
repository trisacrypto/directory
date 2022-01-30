import { isoCountries } from "utils/country";

function getIsoCountrie(country = '') {
    return isoCountries[country]
}

export default getIsoCountrie