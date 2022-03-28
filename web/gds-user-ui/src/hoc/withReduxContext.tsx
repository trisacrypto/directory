import { StoryFnReactReturnType } from '@storybook/react/dist/ts3.9/client/preview/types';
import store from 'application/store';
import { FC } from 'react';
import { Provider } from 'react-redux';

export const withReduxContext =
  () =>
  // eslint-disable-next-line react/display-name
  (Component: FC): StoryFnReactReturnType =>
    (
      <Provider store={store}>
        <Component />
      </Provider>
    );
