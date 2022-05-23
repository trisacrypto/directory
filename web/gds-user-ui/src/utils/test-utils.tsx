import { ChakraProvider } from '@chakra-ui/react';
import { render as rtlRender, RenderOptions as RTLRenderOptions } from '@testing-library/react';
import React, { useState } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Provider } from 'react-redux';
import customTheme from './theme';
import rootReducer from 'application/store/rootReducer';
import { configureStore, Store } from '@reduxjs/toolkit';
import { I18nProvider } from '@lingui/react';
import { i18n } from '@lingui/core';
import { LanguageContext } from 'contexts/LanguageContext';

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
    const [language, setLanguage] = useState<string | null>('en');

    return (
      <I18nProvider i18n={i18n}>
        <LanguageContext.Provider value={['en', setLanguage]}>
          <Provider store={store}>
            <ChakraProvider theme={customTheme}>
              <FormProvider {...methods}>{children}</FormProvider>
            </ChakraProvider>
          </Provider>
        </LanguageContext.Provider>
      </I18nProvider>
    );
  }
  return rtlRender(ui, { wrapper: Wrapper, ...renderOptions });
}

export * from '@testing-library/react';
