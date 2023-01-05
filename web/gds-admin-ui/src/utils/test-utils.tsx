import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render as rtlRender } from '@testing-library/react';
import { createMemoryHistory } from 'history';
import { Suspense } from 'react';
import { Provider } from 'react-redux';
import { Router } from 'react-router-dom';

import { configureStore } from '../redux/store';
import { ModalProvider } from '../contexts/modal';
import React from 'react';

const queryClient = new QueryClient();
export const render = (ui: any) => {
    const history = createMemoryHistory();

    return rtlRender(
        <Provider store={configureStore({})}>
            <QueryClientProvider client={queryClient}>
                <Suspense fallback="loading...">
                    <Router history={history}>
                        <ModalProvider>{ui}</ModalProvider>
                    </Router>
                </Suspense>
            </QueryClientProvider>
        </Provider>
    );
};

export * from '@testing-library/react';
