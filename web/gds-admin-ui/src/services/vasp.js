import { APICore } from "../helpers/api/apiCore";

const api = new APICore();

function getVasp(id, params) {
    return api.get(`/vasps/${id}`, params)
}

export { getVasp };