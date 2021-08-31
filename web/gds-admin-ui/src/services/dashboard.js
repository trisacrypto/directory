import { APICore } from "../helpers/api/apiCore";

const api = new APICore();

function getSummary(params) {
    return api.get("/summary", params)
}


export { getSummary }