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

function deleteVasp(id) {
    const csrfToken = getCookie('csrf_token')
    return api.delete(`vasps/${id}`, {
        headers: {
            'X-CSRF-TOKEN': csrfToken
        }
    })
}

function putContact(vaspId, kind, data) {
    const csrfToken = getCookie('csrf_token')
    return api.update(`/vasps/${vaspId}/contacts/${kind}`, data, {
        headers: {
            'X-CSRF-TOKEN': csrfToken
        }
    })
}

function removeContact(vaspId, kind) {
    const csrfToken = getCookie('csrf_token')
    return api.delete(`/vasps/${vaspId}/contacts/${kind}`, {
        headers: {
            'X-CSRF-TOKEN': csrfToken
        }
    })
}

export { removeContact, putContact, getVasp, updateVasp, reviewVasp, getAdminVerificationToken, deleteVasp };