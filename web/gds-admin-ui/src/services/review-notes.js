import { APICore } from "../helpers/api/apiCore";
import { getCookie } from "../utils";

const api = new APICore()

function getReviewNotes(id, params) {
    return api.get(`/vasps/${id}/notes`, params)
}

function postReviewNote(note, vaspId) {
    const payload = { text: note, note_id: '' }
    const csrfToken = getCookie('csrf_token')
    return api.create(`/vasps/${vaspId}/notes`, payload, {
        headers: {
            'X-CSRF-TOKEN': csrfToken
        }
    })
}

export { getReviewNotes, postReviewNote }