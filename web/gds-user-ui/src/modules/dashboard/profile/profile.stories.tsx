import { Meta, Story } from '@storybook/react';
import Profile from '.';

export default {
  title: 'modules/Profile',
  component: Profile
} as Meta;

const Template: Story = (args) => <Profile {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
