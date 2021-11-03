import { createSelector } from 'reselect'

const reviewNotesState = state => state.ReviewNotes

const getAllReviewNotes = createSelector(reviewNotesState, (state) => state.data)
const getReviewNoteLoadingState = createSelector(reviewNotesState, state => state.loading)


export { getAllReviewNotes, getReviewNoteLoadingState }