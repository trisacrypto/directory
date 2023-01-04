import { applyMiddleware, createStore } from 'redux';
import createSagaMiddleware from 'redux-saga';

import reducers from './reducers';
import rootSaga from './sagas';

const sagaMiddleware = createSagaMiddleware();
const middlewares = [sagaMiddleware];

export function configureStore(initialState: {}): any {
  let store;

  if ((window as any).__REDUX_DEVTOOLS_EXTENSION_COMPOSE__) {
    store = createStore(
      reducers,
      initialState,
      (window as any).__REDUX_DEVTOOLS_EXTENSION_COMPOSE__(applyMiddleware(...middlewares))
    );
  } else {
    store = createStore(reducers, initialState, applyMiddleware(...middlewares));
  }
  sagaMiddleware.run(rootSaga);
  return store;
}
