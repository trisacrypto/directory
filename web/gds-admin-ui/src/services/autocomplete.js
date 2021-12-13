const { APICore } = require("helpers/api/apiCore")

const api = new APICore()

const getAllAutocomplete = () => {
    return api.get('/autocomplete')
}

export default getAllAutocomplete