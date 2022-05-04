import { Story } from '@storybook/react';
import { UseFormRegister, UseFormSetValue, Control } from 'react-hook-form/dist/types';
import BasicDetailsForm from '.';

type BasicDetailsFormProps = {};

export default {
  title: 'components/BasicDetailsForm',
  component: BasicDetailsForm
};

const Template: Story<BasicDetailsFormProps> = (args) => <BasicDetailsForm {...args} />;

export const Default = Template.bind({});
Default.args = {};
