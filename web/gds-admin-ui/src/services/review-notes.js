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

function deleteReviewNote(noteId, vaspId, params) {
    const csrfToken = getCookie('csrf_token')
    return api.delete(`/vasps/${vaspId}/notes/${noteId}`, {
        headers: {
            'X-CSRF-TOKEN': csrfToken
        }
    })
}

function updateReviewNote(note, noteID, vaspID) {
    const csrfToken = getCookie('csrf_token');
    const data = {
        text: note
    }
    return api.update(`/vasps/${vaspID}/notes/${noteID}`, data, {
        headers: {
            'X-CSRF-TOKEN': csrfToken
        }
    })
}

export { getReviewNotes, postReviewNote, deleteReviewNote, updateReviewNote }
