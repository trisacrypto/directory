import { Meta, Story } from '@storybook/react';
import Collaborators from '.';

export default {
  title: 'modules/Collaborators',
  component: Collaborators
} as Meta;

const Template: Story = (args) => <Collaborators {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
