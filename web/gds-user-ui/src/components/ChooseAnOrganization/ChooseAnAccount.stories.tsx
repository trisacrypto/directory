import { Meta, Story } from '@storybook/react';
import ChooseAnOrganization from '.';

export default {
  title: 'components/ChooseAnAccount',
  component: ChooseAnOrganization
} as Meta;

const Template: Story = (args) => <ChooseAnOrganization {...args} />;

export const Standard = Template.bind({});
Template.args = {};
