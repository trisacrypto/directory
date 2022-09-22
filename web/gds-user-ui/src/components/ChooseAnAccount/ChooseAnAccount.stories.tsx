import { Meta, Story } from '@storybook/react';
import ChooseAnAccount from '.';

type ComponentType = typeof ChooseAnAccount;

export default {
  title: 'components/ChooseAnAccount',
  component: ChooseAnAccount
} as Meta<ComponentType>;

const Template: Story<ComponentType> = (args) => <ChooseAnAccount />;

export const Standard = Template.bind({});
Template.args = {};
