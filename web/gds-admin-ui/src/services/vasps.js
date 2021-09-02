import { APICore } from "../helpers/api/apiCore";

const api = new APICore();

function getVasp(id) {
    return api.get(`/vasps/${id}`)
}

export { getVasp };