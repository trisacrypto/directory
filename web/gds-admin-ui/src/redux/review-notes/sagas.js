import NProgress from 'nprogress';
import toast from 'react-hot-toast';
import { call, put, takeEvery } from 'redux-saga/effects';
import { deleteReviewNote as deleteNote, getReviewNotes } from 'services/review-notes';

import { fetchReviewNotesApiResponseError, fetchReviewNotesApiResponseSuccess } from './actions';
import { FetchReviewNotesActionTypes } from './constants';

import { DeleteReviewNotesActionTypes } from '.';

function* fetchReviewNotes({ payload: { id } }) {
  NProgress.start();
  try {
    const response = yield call(getReviewNotes, id);
    const { notes } = response?.data;
    const sortedData = Array.isArray(notes)
      ? notes.sort((a, b) => {
        const date1 = a?.modified ? new Date(a?.modified) : new Date(a?.created);
        const date2 = b?.modified ? new Date(b?.modified) : new Date(b?.created);

        return date2 - date1;
        })
      : [];

    yield put(fetchReviewNotesApiResponseSuccess(sortedData));
    NProgress.done();
  } catch (error) {
    toast.error(error);
    yield put(fetchReviewNotesApiResponseError(error.message));
    NProgress.done();
  }
}

function* deleteReviewNote({ payload: { noteId, vaspId } }) {
  try {
    const response = yield call(deleteNote, noteId, vaspId);
    if (response) {
      NProgress.done();
    }
  } catch (error) {
    console.error(error);
  }
}

export function* reviewNotesSaga() {
  yield takeEvery(FetchReviewNotesActionTypes.FETCH_REVIEW_NOTES, fetchReviewNotes);
}

export function* deleteReviewNoteSaga() {
  yield takeEvery(DeleteReviewNotesActionTypes.DELETE_REVIEW_NOTES, deleteReviewNote);
}
