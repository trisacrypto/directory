import { Meta, Story } from '@storybook/react';
import Contacts from '.';

export default {
  title: 'components/Contacts',
  component: Contacts
} as Meta;

const Template: Story = (args) => <Contacts {...args} />;

export const Default = Template.bind({});
Default.args = {};
