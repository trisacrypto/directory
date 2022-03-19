/* eslint-disable react/display-name */
import { action } from '@storybook/addon-actions';
import { StoryFnReactReturnType } from '@storybook/react/dist/ts3.9/client/preview/types';
import { VFC, ReactNode, FC } from 'react';
import { FormProvider, useForm } from 'react-hook-form';

const StorybookFormProvider: VFC<{
  children: ReactNode;
  defaultValues?: {
    [x: string]: any;
  };
}> = ({ children, defaultValues }) => {
  const methods = useForm({
    defaultValues
  });

  return (
    <FormProvider {...methods}>
      <form onSubmit={methods.handleSubmit(action('[React Hooks Form] Submit'))}>{children}</form>
    </FormProvider>
  );
};

StorybookFormProvider.displayName = 'StorybookFormProvider';

export const withRHF =
  (showSubmitButton: boolean, defaultValues?: Record<string, any>) =>
  (Story: FC): StoryFnReactReturnType =>
    (
      <StorybookFormProvider defaultValues={defaultValues}>
        <Story />
        {showSubmitButton && <button type="submit">Submit</button>}
      </StorybookFormProvider>
    );
