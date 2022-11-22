import { Meta, Story } from '@storybook/react';
import ChooseAnAccount, { ChooseAnAccountProps } from '.';

export default {
  title: 'components/ChooseAnAccount',
  component: ChooseAnAccount
} as Meta<ChooseAnAccountProps>;

const Template: Story<ChooseAnAccountProps> = (args) => <ChooseAnAccount {...args} />;

export const Standard = Template.bind({});
Template.args = {};
