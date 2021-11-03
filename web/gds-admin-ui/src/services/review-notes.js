import { APICore } from "../helpers/api/apiCore";

const api = new APICore()

function getReviewNotes(id, params) {
    return api.get(`/vasps/${id}/notes`, params)
}

export { getReviewNotes }