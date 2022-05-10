import { ChakraProvider } from '@chakra-ui/react';
import { render as rtlRender, RenderOptions as RTLRenderOptions } from '@testing-library/react';
import React from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Provider } from 'react-redux';
import customTheme from './theme';
import rootReducer from 'application/store/rootReducer';
import { configureStore, Store } from '@reduxjs/toolkit';

type RenderOptions = {
  rhfDefaultValues?: { [x: string]: any };
  store?: Store;
  preloadedState?: StoreState;
  renderOptions?: Omit<RTLRenderOptions, 'wrapper'>;
};

type StoreState = {
  [x: string]: any;
};

export function render(
  ui: React.ReactElement<any, string | React.JSXElementConstructor<any>>,
  {
    rhfDefaultValues,
    preloadedState,
    store = configureStore({ reducer: rootReducer, preloadedState }),
    ...renderOptions
  }: RenderOptions = {}
) {
  function Wrapper({ children }: { children: React.ReactNode }) {
    const methods = useForm({ defaultValues: rhfDefaultValues });
    return (
      <Provider store={store}>
        <ChakraProvider theme={customTheme}>
          <FormProvider {...methods}>{children}</FormProvider>
        </ChakraProvider>
      </Provider>
    );
  }
  return rtlRender(ui, { wrapper: Wrapper, ...renderOptions });
}

export * from '@testing-library/react';
