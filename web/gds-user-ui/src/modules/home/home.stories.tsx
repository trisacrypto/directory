import { Meta, Story } from '@storybook/react';
import HomePage from '.';

export default {
  title: 'modules/Home',
  component: HomePage
} as Meta;

const Template: Story = (args) => <HomePage {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
