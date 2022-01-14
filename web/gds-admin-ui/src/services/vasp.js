import { getCookie } from "utils";
import { APICore } from "../helpers/api/apiCore";

const api = new APICore();

function getVasp(id, params) {
    return api.get(`/vasps/${id}`, params)
}

function updateVasp(id, payload) {
    const csrfToken = getCookie('csrf_token')
    return api.patch(`/vasps/${id}`, payload, {
        headers: {
            'X-CSRF-TOKEN': csrfToken
        }
    })
}

function getAdminVerificationToken(id) {
    return api.get(`/vasps/${id}/review`)
}

function reviewVasp(id, payload, params) {
    const csrfToken = getCookie('csrf_token')
    return api.create(`/vasps/${id}/review`, payload, {
        headers: {
            'X-CSRF-TOKEN': csrfToken
        }
    })
}

export { getVasp, updateVasp, reviewVasp, getAdminVerificationToken };