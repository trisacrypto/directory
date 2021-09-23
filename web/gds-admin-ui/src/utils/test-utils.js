import { Provider } from "react-redux";
import { render as rtlRender } from '@testing-library/react';
import { configureStore } from "../redux/store";

export const render = (ui) => rtlRender(<Provider store={configureStore({})}>{ui}</Provider>)

export * from '@testing-library/react'