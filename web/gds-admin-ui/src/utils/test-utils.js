import { Provider } from "react-redux";
import { render as rtlRender } from '@testing-library/react';
import { configureStore } from "../redux/store";
import { ModalProvider } from "contexts/modal";
import { createMemoryHistory } from 'history'
import { Router } from "react-router-dom";
import { Suspense } from "react";
import { QueryClientProvider, QueryClient } from '@tanstack/react-query';



const queryClient = new QueryClient();
export const render = (ui) => {
    const history = createMemoryHistory()

    return rtlRender(
        <Provider store={configureStore({})}>
            <QueryClientProvider client={queryClient}>
                <Suspense fallback="loading...">
                    <Router history={history}>
                        <ModalProvider>
                            {ui}
                        </ModalProvider>
                    </Router>
                </Suspense>
            </QueryClientProvider>
        </Provider>
    )
}

export * from '@testing-library/react'