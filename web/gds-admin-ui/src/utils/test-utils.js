import { Provider } from "react-redux";
import { render as rtlRender } from '@testing-library/react';
import { configureStore } from "../redux/store";
import { ModalProvider } from "contexts/modal";

export const render = (ui) => rtlRender(<Provider store={configureStore({})}>
    <ModalProvider>
        {ui}
    </ModalProvider>
</Provider>)

export * from '@testing-library/react'