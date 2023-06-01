import React, { Dispatch, SetStateAction } from 'react';
import { Meta, Story } from '@storybook/react';
import { withRHF } from 'hoc/withRHF';
import TrixoQuestionnaireForm from '.';
interface Props {
  data: any;
  isLoading?: boolean;
  shouldResetForm?: boolean;
  onResetFormState?: Dispatch<SetStateAction<boolean>>;
}

export default {
  title: 'components/TrixoQuestionnaireForm',
  component: TrixoQuestionnaireForm,
  decorators: [withRHF(false)]
} as Meta;

const Template: Story<Props> = (args) => <TrixoQuestionnaireForm {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
