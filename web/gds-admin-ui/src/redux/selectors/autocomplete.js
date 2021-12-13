import { createSelector } from "reselect"

const autocompleteState = state => state.Autocomplete

const fetchAllAutocomplete = createSelector(autocompleteState, state => state.data?.names)
const fetchAutocompleteLoadingState = createSelector(autocompleteState, state => state.loading)

export { fetchAllAutocomplete, fetchAutocompleteLoadingState }