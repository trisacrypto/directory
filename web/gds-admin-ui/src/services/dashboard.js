import { APICore } from "../helpers/api/apiCore";

const api = new APICore();

function getSummary(params) {
    return api.get("/summary", params)
}

function getRegistrationsReviews(params) {
    return api.get("/reviews", params)
}


export { getSummary, getRegistrationsReviews }