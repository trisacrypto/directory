import { APICore } from "helpers/api/apiCore"
import { getCookie } from "utils"

const api = new APICore()

const updateTrixoForm = (id, payload) => {
    const csrfToken = getCookie("csrf_token")
    return api.patch(`/vasps/${id}`, payload, {
        headers: {
            'X-CSRF-TOKEN': csrfToken
        }
    })
}

export default updateTrixoForm