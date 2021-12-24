import { APICore } from "helpers/api/apiCore"

const api = new APICore()

const updateTrixoForm = (id, payload) => {
    return api.patch(`/vasps/${id}`, payload)
}

export default updateTrixoForm