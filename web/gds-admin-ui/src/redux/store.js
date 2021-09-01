// @flow
import { createStore, applyMiddleware } from 'redux';
import createSagaMiddleware from 'redux-saga';
import reducers from './reducers';
import rootSaga from './sagas';

const sagaMiddleware = createSagaMiddleware();
const middlewares = [sagaMiddleware];

export function configureStore(initialState: {}): any {
    let store;

    if (window['__REDUX_DEVTOOLS_EXTENSION_COMPOSE__']) {
        store = createStore(reducers, initialState, window['__REDUX_DEVTOOLS_EXTENSION_COMPOSE__'](applyMiddleware(...middlewares)));
    } else {
        store = createStore(reducers, initialState, applyMiddleware(...middlewares));
    }
    sagaMiddleware.run(rootSaga);
    return store;
}
