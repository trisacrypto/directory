import { Meta, Story } from '@storybook/react';
import { withRHF } from 'hoc/withRHF';
import TrixoQuestionnaireForm from '.';

export default {
  title: 'components/TrixoQuestionnaireForm',
  component: TrixoQuestionnaireForm,
  decorators: [withRHF(false)]
} as Meta;

const Template: Story = (args) => <TrixoQuestionnaireForm {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
