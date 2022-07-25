import { configureStore, Action } from '@reduxjs/toolkit';
import { useDispatch } from 'react-redux';
import thunk, { ThunkAction } from 'redux-thunk';
import storage from 'redux-persist/lib/storage';
import { persistReducer, persistStore } from 'redux-persist';
import storageSession from 'reduxjs-toolkit-persist/lib/storage/session';
import rootReducer, { RootState } from './rootReducer';

const persistConfig = {
  key: 'root',
  storage: storageSession
};
// const userPersistConfig = {
//   key: 'user',
//   storage: storageSession
// };

const persistedReducer = persistReducer(persistConfig, rootReducer);

const store = configureStore({
  reducer: persistedReducer,
  middleware: [thunk]
});

export type AppDispatch = typeof store.dispatch;
export const useAppDispatch = () => useDispatch<AppDispatch>();
export type AppThunk = ThunkAction<void, RootState, unknown, Action<string>>;
export const persistor = persistStore(store);
export default store;
