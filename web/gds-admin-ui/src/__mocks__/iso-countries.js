import { isoCountries } from "utils/country";

function getIsoCountry(country = '') {
    return isoCountries[country]
}

export default getIsoCountry